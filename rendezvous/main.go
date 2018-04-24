package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
	flags "github.com/jessevdk/go-flags"
)

var agentTable = NewAgentTable()
var sessionTable = NewSessionTable()
var responseTable map[string]chan string = make(map[string]chan string)

func Start(port int32) {
	go http.Handle("/agent/listen", onAgentConnection{})
	go http.Handle("/agent/list", onAgentList{})
	go http.Handle("/session/list", onSessionList{})
	go http.Handle("/session/attach/", onSessionAttach{})
	go http.Handle("/session/create", onSessionCreate{})

	glog.V(0).Infof("Ready for incoming connections")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

var options struct {
	Port int32 `short:"P" long:"port" env:"PORT" description:"Port" required:"true"`
}

func main() {
	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	Start(options.Port)
}
