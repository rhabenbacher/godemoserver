package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const usage = `
Usage:
%s [command]

Commands:
  standalone     Run as standalone server
  rest           Provide a rest api
  frontend       Show response from rest api
`

const usageStandalone = `
Usage
%s standalone [flag]

Flags:

`

const usageRest = `
Usage
%s rest [flag]

Endpoint /time is provided

Flags:
`

const usageFrontend = `
Usage
%s frontend [flag]

Server calls rest Api and shows result :-)
The Environment variables API_HOST and API_PORT have to be set 
`

type Config struct {
	port        string
	handlerFunc func(logs *serverLogs) *http.ServeMux
}

func (c *Config) setupMenu() *flag.FlagSet {
	menu := flag.NewFlagSet("menu", flag.ExitOnError)
	menu.Usage = func() {
		fmt.Printf(usage, os.Args[0])
		menu.PrintDefaults()
	}
	return menu
}

func (c *Config) setupSubMenu(flagSetName string, defaultPort string, usageInfo string) *flag.FlagSet {
	menu := flag.NewFlagSet(flagSetName, flag.ExitOnError)
	menu.StringVar(&c.port, "port", defaultPort, "The server port")
	menu.Usage = func() {
		fmt.Printf(usageInfo, os.Args[0])
		menu.PrintDefaults()
	}
	return menu
}

func (c *Config) setupStandaloneSubMenu() *flag.FlagSet {
	return c.setupSubMenu("standalone", "9000", usageStandalone)
}

func (c *Config) setupRestSubMenu() *flag.FlagSet {
	return c.setupSubMenu("rest", "3000", usageRest)
}

func (c *Config) setupFrontendSubMenu() *flag.FlagSet {
	return c.setupSubMenu("frontend", "8000", usageFrontend)
}

func parseFlags(flagSet *flag.FlagSet) func() {

	return func() {
		if err := flagSet.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error parsing params %s, error: %s", os.Args[2:], err)
			return
		}
	}
}

func isRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

func isRunningInKubernetesPod() bool {
	if _, err := os.Stat("/var/run/secrets/kubernetes.io"); err == nil {
		return true
	}
	return false
}

func main() {
	conf := &Config{}
	menu := conf.setupMenu()

	if len(os.Args) > 5 || len(os.Args) == 1 {
		menu.Usage()
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "standalone":
		parseFlags(conf.setupStandaloneSubMenu())()
		conf.runStandaloneServer()
	case "rest":
		parseFlags(conf.setupRestSubMenu())()
		conf.runApiServer()
	case "frontend":
		parseFlags(conf.setupFrontendSubMenu())()
		conf.runFrontendServer()

	default:
		fmt.Println("Invalid command!")
		menu.Usage()

	}
}
