package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	response, err := http.Get(fmt.Sprintf("http://%s:%s/time", host, port))
	if err != nil {
		http.Error(w, "Cannot connect to API :-(", http.StatusInternalServerError)
		logs.errorLog.Println("Frontend: showTime - http Get Error", err.Error())
		return nil, err
	}
	return response, nil
}

func decodeTimeResponse(res *http.Response) (timeResponse *TimeResponse) {
	dec := json.NewDecoder(res.Body)
	defer res.Body.Close()
	err := dec.Decode(&timeResponse)
	if err != nil {
		timeResponse = nil
	}
	return
}

func (logs serverLogs) getTimefromAPIServer(w http.ResponseWriter) *TimeResponse {

	response, err := logs.httpGetTime(w)
	if err != nil {
		return nil
	}

	timeStruct := decodeTimeResponse(response)
	if timeStruct == nil {
		http.Error(w, "Error decoding Timer Response", http.StatusInternalServerError)
		logs.errorLog.Println("Frontend: Error decoding Time Response")
	}
	return timeStruct
}

func (logs serverLogs) showTime(w http.ResponseWriter, r *http.Request) {

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
