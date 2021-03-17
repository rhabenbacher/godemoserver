package main

import (
	htmltemplate "html/template"
	"math/big"
	"net/http"
	"os"
	"strconv"
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

type OutputStrucHtml struct {
	N          int
	FibN       string
	FibNLength int
	DurationMs int64
}

const templateText = `
Hello, I'm serving in standalone mode !!!

{{if .IsKubernetes}}I'm running in a Pod with name {{.Hostname}}{{else if .IsDocker}}I'm running in a Container with ID {{.Hostname}}{{else}}I'm running somewhere{{end}} 

{{.Time}}
`

const templateHtml = `
<!doctype html>
<html lang='en'>
<head>
<meta charset='utf-8'>
<title>GoServer - Standalone</title>

</head>
<style>
    h1 {
        font-size: 3rem;
		color: #FFB700;
    }

    p {
        padding: 1rem;
        margin-top: 3rem;
        font-size: 1rem;
        background-color: #F4F4F4;
        word-wrap: break-word;
        line-height: 1.25;
        
    }

    body {
        font-family: Consolas,monaco,monospace;
        margin-top: 3rem;
        margin-left: 3rem;
        width: 80%;
    }
</style>    
<body>
<h1>Fib({{.N}}) is a number with {{.FibNLength}} digits. The calculation took {{.DurationMs}} ms</h1> 
<p>{{.FibN}}</p>
</body>
    
</html>

`

func getMuxforStandAloneServer() HandlerFunc {

	return func(logs *serverLogs) *http.ServeMux {
		mux := http.NewServeMux()
		mux.HandleFunc("/", logs.home)
		mux.HandleFunc("/fibonacci", logs.fibonacci)
		return mux
	}
}

func calculateFibonacci(n int) *big.Int {
	switch {
	case n == 0:
		return big.NewInt(0)
	case n <= 2:
		return big.NewInt(1)
	}

	x, y := big.NewInt(1), big.NewInt(2)
	for i := 0; i < n-2; i++ {
		x.Add(x, y)
		x, y = y, x
	}
	return x

}

func (logs *serverLogs) fibonacci(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil || n < 0 {
		http.Error(w, "Parmeter n is required", http.StatusInternalServerError)
		return

	}
	start := time.Now()
	fibN := calculateFibonacci(n).String()
	t := time.Now()
	elapsedMs := t.Sub(start).Milliseconds()
	fibNlength := len(fibN)
	out := OutputStrucHtml{n, fibN, fibNlength, elapsedMs}

	tmpl, err := htmltemplate.New("fibonacci").Parse(templateHtml)
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

func (c *Config) runStandaloneServer() {
	c.handlerFunc = getMuxforStandAloneServer()
	c.startServer([]string{"serving in standalone mode"})

}
