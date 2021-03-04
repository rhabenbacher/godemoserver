package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Name string
	Time string
}

func (logs *serverLogs) getMuxforApi() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/time", logs.getTime)
	return mux
}

func (logs *serverLogs) getTime(w http.ResponseWriter, r *http.Request) {

	logs.infoLog.Println("getTime: received request")
	response := Response{os.Getenv("HOSTNAME"), time.Now().Format("02 Jan 06 15:04:01.1234 MST")}
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.errorLog.Println("getTime: error creating json")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}

func (c *Config) runApiServer() {

	logs := &serverLogs{}
	logs.setup()

	server := &http.Server{
		Addr:     ":" + c.port,
		ErrorLog: logs.errorLog,
		Handler:  logs.getMuxforApi(),
	}
	logs.logStartupInfo("api", c.port)
	err := server.ListenAndServe()
	logs.errorLog.Fatal(err)
}
