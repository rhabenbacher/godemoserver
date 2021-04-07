package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func getMuxforFrontend() HandlerFunc {

	return func(logs serverLogs) *http.ServeMux {
		mux := http.NewServeMux()
		mux.HandleFunc("/", logs.showTime)
		return mux
	}
}

func getEnvvar(name string) (string, error) {
	envvar := os.Getenv(name)
	if len(envvar) == 0 {
		return "", fmt.Errorf("%s not set", name)
	}
	return envvar, nil
}

func getEnvForFrontend() (string, string, error) {
	host, err := getEnvvar("API_HOST")
	if err != nil {
		return "", "", err
	}
	port, err := getEnvvar("API_PORT")
	if err != nil {
		return "", "", err
	}

	return host, port, nil
}

func (logs serverLogs) httpGetTime(w http.ResponseWriter) (*http.Response, error) {
	host, port, _ := getEnvForFrontend()
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, fmt.Sprintf("http://%s:%s/time", host, port), nil)
	if err != nil {
		panic(err)
	}

	return client.Do(req)

}

func decodeTimeResponse(res *http.Response) (timeResponse *TimeResponse) {
	dec := json.NewDecoder(res.Body)
	defer res.Body.Close()
	err := dec.Decode(&timeResponse)
	if err != nil {
		return nil
	}
	return timeResponse
}

func (logs serverLogs) errorNoApiConnection(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, "Cannot connect to API :-(", http.StatusInternalServerError)
		logs.errorLog.Println("Frontend: showTime - http Get Error", err.Error())
		return true
	}
	return false
}

func (logs serverLogs) errorApiWithBadStatus(w http.ResponseWriter, statusCode int) bool {
	if statusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("API responded with status %v :-(", statusCode), http.StatusInternalServerError)
		return true
	}
	return false
}

func (logs serverLogs) errorDecodingResponse(w http.ResponseWriter, isError bool) {
	if isError {
		http.Error(w, "Error decoding Timer Response", http.StatusInternalServerError)
		logs.errorLog.Println("Frontend: Error decoding Time Response")
	}
}

func (logs serverLogs) getTimefromAPIServer(w http.ResponseWriter) *TimeResponse {

	response, err := logs.httpGetTime(w)
	if logs.errorNoApiConnection(w, err) {
		return nil
	}

	if logs.errorApiWithBadStatus(w, response.StatusCode) {
		return nil
	}

	timeStruct := decodeTimeResponse(response)
	logs.errorDecodingResponse(w, timeStruct == nil)

	return timeStruct
}

func (logs serverLogs) showTime(w http.ResponseWriter, r *http.Request) {
	logs.infoLog.Printf("received request from %s", r.Header.Get("User-Agent"))
	respStruct := logs.getTimefromAPIServer(w)
	if respStruct == nil {
		return
	}
	fmt.Fprintf(w, "Host %v sent time: %v", respStruct.Name, respStruct.Time)
}

func (c *Config) runFrontendServer() {

	apiHost, apiPort, err := getEnvForFrontend()
	if err != nil {
		log.Println("API_HOST and API_PORT have to be set as environment variables")
		log.Fatalln(err)
		return
	}

	c.handlerFunc = getMuxforFrontend()
	c.startServer([]string{"serving in frontend mode", "API_HOST:" + apiHost, "API_PORT:" + apiPort})

}
