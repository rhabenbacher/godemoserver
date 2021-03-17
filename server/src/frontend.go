package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getMuxforFrontend() HandlerFunc {

	return func(logs *serverLogs) *http.ServeMux {
		mux := http.NewServeMux()
		mux.HandleFunc("/", logs.showTime)
		return mux
	}
}

func getEnvForFrontend() (string, string, bool) {
	host := os.Getenv("API_HOST")
	port := os.Getenv("API_PORT")
	if len(host) > 0 && len(port) > 0 {
		return host, port, true
	}
	return "", "", false
}

func (logs *serverLogs) httpGetTime(w http.ResponseWriter) (*http.Response, error) {
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

func (logs *serverLogs) getTimefromAPIServer(w http.ResponseWriter) *TimeResponse {

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

func (logs *serverLogs) showTime(w http.ResponseWriter, r *http.Request) {

	respStruct := logs.getTimefromAPIServer(w)
	if respStruct == nil {
		return
	}
	fmt.Fprintf(w, "Host %v sent time: %v", respStruct.Name, respStruct.Time)
}

func (c *Config) runFrontendServer() {

	logs := &serverLogs{}
	logs.setup()

	apiHost, apiPort, ok := getEnvForFrontend()
	if !ok {
		log.Fatalln("Frontend: API_HOST or API_PORT not set")
		return
	}

	c.handlerFunc = getMuxforFrontend()
	c.startServer([]string{"serving in frontend mode", "API_HOST:" + apiHost, "API_PORT:" + apiPort})

}
