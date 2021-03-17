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

func getMuxforApi() func(logs *serverLogs) *http.ServeMux {
	return func(logs *serverLogs) *http.ServeMux {
		mux := http.NewServeMux()
		mux.HandleFunc("/time", logs.getTime)
		return mux
	}
}

func (logs *serverLogs) getTime(w http.ResponseWriter, r *http.Request) {

	logs.infoLog.Println("getTime: received request")
	tnow := time.Now()
	response := Response{os.Getenv("HOSTNAME"), tnow.Format("02.01.2006 15:04:05.000")}
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
	c.handlerFunc = getMuxforApi()
	c.startServer([]string{"serving in api mode", "providing endpoint /time"})

}
