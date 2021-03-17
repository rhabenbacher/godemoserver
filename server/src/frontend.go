package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getMuxforFrontend() func(logs *serverLogs) *http.ServeMux {

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

func (logs *serverLogs) showTime(w http.ResponseWriter, r *http.Request) {

	host, port, _ := getEnvForFrontend()

	response, err := http.Get(fmt.Sprintf("http://%s:%s/time", host, port))
	if err != nil {
		http.Error(w, "Cannot connect to API :-(", http.StatusInternalServerError)
		logs.errorLog.Println("Frontend: showTime - http Get Error", err.Error())
		return
	}

	defer response.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(response.Body)

	var respStruct Response
	err = json.Unmarshal(bodyBytes, &respStruct)
	if err != nil {
		http.Error(w, "Cannot connect to API :-(", http.StatusInternalServerError)
		logs.errorLog.Println("Frontend: showTime - json Unmarshal", err.Error())
		return
	}

	fmt.Fprintf(w, "Host %v sent time: %v", respStruct.Name, respStruct.Time)

}

func (c *Config) runFrontendServer() {

	logs := &serverLogs{}
	logs.setup()

	apiHost, apiPort, ok := getEnvForFrontend()
	if !ok {
		logs := &serverLogs{}
		logs.setup()
		logs.errorLog.Println("Frontend: API_HOST or API_PORT not set")
		return
	}

	c.handlerFunc = getMuxforFrontend()
	c.startServer([]string{"serving in frontend mode", "API_HOST:" + apiHost, "API_PORT:" + apiPort})

}
