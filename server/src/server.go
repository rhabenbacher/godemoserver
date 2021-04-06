package main

import (
	"log"
	"net/http"
	"os"
)

type serverLogs struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func setupServerLogs() serverLogs {
	return serverLogs{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (logs serverLogs) logStartupInfo(port string) {
	logs.infoLog.Printf("Start server on port %v", port)
	logs.infoLog.Println("Docker environment detected:", isRunningInDockerContainer())
	logs.infoLog.Println("Kubernetes environment (Pod) detected:", isRunningInKubernetesPod())
}

func (c *Config) startServer(startupMessages []string) {

	logs := setupServerLogs()

	server := &http.Server{
		Addr:     ":" + c.port,
		ErrorLog: logs.errorLog,
		Handler:  c.handlerFunc(logs),
	}

	logs.logStartupInfo(c.port)
	for _, message := range startupMessages {
		logs.infoLog.Println(message)
	}

	err := server.ListenAndServe()
	logs.errorLog.Fatal(err)
}
