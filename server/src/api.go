package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type TimeResponse struct {
	Name string
	Time string
}

func getMuxforApi() HandlerFunc {
	return func(logs serverLogs) *http.ServeMux {
		mux := http.NewServeMux()
		mux.HandleFunc("/time", logs.getTime)
		return mux
	}
}

func (logs serverLogs) getTime(w http.ResponseWriter, r *http.Request) {

	logs.infoLog.Printf("/time received request from %s", r.Header.Get("User-Agent"))
	tnow := time.Now()
	response := TimeResponse{os.Getenv("HOSTNAME"), tnow.Format("02.01.2006 15:04:05.000")}
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
