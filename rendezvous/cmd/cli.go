package cmd

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

const (
	rendezvousExample = `Start the rendezvous and the agent:
# On the rendezvous (1.2.3.4)
$ rendezvous -P 7777

# On the agent (3.4.5.6)
$ agent websocket -U http://1.2.3.4:7777

Open a shell and send some commands
# List the agents
$ ./master list-agents -U http://1.2.3.4:7777
List of agents:
* 3.4.5.6:65000

# Create a session
$ master create -U http://1.2.3.4:7777 3.4.5.6:65000
Attaching to admiring_meitn
Connected to admiring_meitn
bash-3.2$
`
)

func GetCommand() *cobra.Command {
	var port int32
	verbose := 0
	command := &cobra.Command{
		Use:   "reverse-shell-rendezvous",
		Short: "An http server listening for agents and masters",
		Long: "The rendezvous is an http server listening for agents and masters.	It can run behind a reverse-proxy and that reverse-proxy could to SSL offloading.",
		Example: rendezvousExample,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.Set("logtostderr", "true")
			flag.Set("v", strconv.Itoa(verbose))
			flag.CommandLine.Parse([]string{})
		},
		Run: func(*cobra.Command, []string) {
			start(port)
		},
		DisableAutoGenTag: true,
	}

	command.Flags().IntVarP(&verbose, "verbose", "v", 0, "Be verbose on log output")
	command.Flags().Int32VarP(&port, "port", "", 8080, "remote port to connect to")
	return command
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
