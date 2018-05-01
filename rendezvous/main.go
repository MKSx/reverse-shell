package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var agentTable = NewAgentTable()
var sessionTable = NewSessionTable()
var responseTable map[string]chan string = make(map[string]chan string)

func main() {
	var port int32
	verbose := 0
	command := &cobra.Command{
		Use:   "reverse-shell-rendezvous",
		Short: "An http server listening for agents and masters",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.Set("logtostderr", "true")
			flag.Set("v", strconv.Itoa(verbose))
			flag.CommandLine.Parse([]string{})
		},
		Run: func(*cobra.Command, []string) {
			start(port)
		},
	}

	command.Flags().IntVarP(&verbose, "verbose", "v", 0, "Be verbose on log output")
	command.Flags().Int32VarP(&port, "port", "", 8080, "remote port to connect to")

	if err := command.Execute(); err != nil {
		panic(err)
	}
}

func start(port int32) {
	go http.Handle("/agent/listen", onAgentConnection{})
	go http.Handle("/agent/list", onAgentList{})
	go http.Handle("/session/list", onSessionList{})
	go http.Handle("/session/attach/", onSessionAttach{})
	go http.Handle("/session/create", onSessionCreate{})

	glog.V(0).Infof("Ready for incoming connections")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
