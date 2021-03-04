package main

import (
	"net/http"
	"os"
	"text/template"
	"time"
)

// HOSTNAME = containername bzw Pod Name

type OutputStruc struct {
	Time         string
	IsDocker     bool
	IsKubernetes bool
	Hostname     string
}

const templateText = `
Hello, I'm serving in standalone mode !!!

{{if .IsKubernetes}}I'm running in a Pod with name {{.Hostname}}{{else if .IsDocker}}I'm running in a Container with ID {{.Hostname}}{{else}}I'm running somewhere{{end}} 

{{.Time}}
`

func (logs *serverLogs) getMuxforStandAloneServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", logs.home)
	return mux
}

func (logs *serverLogs) home(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC822)
	out := OutputStruc{now, isRunningInDockerContainer(), isRunningInKubernetesPod(), os.Getenv("HOSTNAME")}

	tmpl, err := template.New("standalone").Parse(templateText)
	if err != nil {
		logs.errorLog.Fatalln(err)
		return
	}
	err = tmpl.Execute(w, out)
	if err != nil {
		logs.errorLog.Fatalln(err)
		return
	}

}

func (logs *serverLogs) logStartupInfo(mode string, port string) {
	logs.infoLog.Printf("Start %v server on port %v", mode, port)
	logs.infoLog.Println("Docker environment detected:", isRunningInDockerContainer())
	logs.infoLog.Println("Kubernetes environment (Pod) detected:", isRunningInKubernetesPod())
}

func (c *Config) runStandaloneServer() {

	logs := &serverLogs{}
	logs.setup()

	server := &http.Server{
		Addr:     ":" + c.port,
		ErrorLog: logs.errorLog,
		Handler:  logs.getMuxforStandAloneServer(),
	}
	logs.logStartupInfo("standalone", c.port)
	err := server.ListenAndServe()
	logs.errorLog.Fatal(err)
}
