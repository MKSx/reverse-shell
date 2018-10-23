package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/golang/glog"
	"github.com/maxlaverse/reverse-shell/agents/go/handler"
	"github.com/spf13/cobra"
)

const (
	bashCommand = "bash --norc"
)

type stdinListener struct {
	processOutput     chan *handler.ProcessOutput
	processTerminated chan *handler.ProcessTerminated
	inputCommand      chan []byte
	interruptChannel  chan os.Signal
	stop              chan error
	handler           *handler.Handler
	stdIn             io.Reader
	stdOut            io.Writer
}

// NewStdinListenerCommand creates a new cobra.Command for `agent stdin`
func NewStdinListenerCommand(agent Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "stdin",
		Short:            "Agent that listens for command on stdin",
		Long:             "Absolutely useless. It's basically just piping *stdin* to a process on the same machine.",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			newStdinListener(os.Stdin, os.Stdout).Listen()
		},
	}

	return cmd
}

func newStdinListener(stdIn io.Reader, stdOut io.Writer) *stdinListener {
	processOutput := make(chan *handler.ProcessOutput)
	processTerminated := make(chan *handler.ProcessTerminated)

	return &stdinListener{
		processOutput:     processOutput,
		processTerminated: processTerminated,
		inputCommand:      make(chan []byte),
		stop:              make(chan error),
		interruptChannel:  make(chan os.Signal, 1),
		handler:           handler.NewHandler(processOutput, processTerminated),
		stdIn:             stdIn,
		stdOut:            stdOut,
	}
}

func (l *stdinListener) Start() error {
	glog.V(1).Info("Starting stdinListener - noop")
	return nil
}

func (l *stdinListener) Listen() error {
	go l.pipeFromProcessOutput()
	go l.pipeToProcessInput()

	glog.V(1).Info("Creating process")
	processID := l.handler.CreateProcess(bashCommand)

	glog.V(2).Info("Installing signal handlers")
	signal.Notify(l.interruptChannel, os.Interrupt)

	glog.V(1).Info("Starting main loop")
	for {
		select {
		case <-l.interruptChannel:
			glog.V(0).Info("Received an interrupt message")
			l.handler.ExecuteCommand(processID, []byte{'\u0003'})
		case msg := <-l.inputCommand:
			glog.V(2).Info("Received an input command")
			l.handler.ExecuteCommand(processID, msg)
		case err := <-l.stop:
			glog.V(2).Info("Stopping main loop")
			return err
		}
	}
}

func (l *stdinListener) pipeToProcessInput() {
	glog.V(2).Info("Starting pipeToProcessInput() routine")
	for {
		select {
		default:
			var msg = make([]byte, 1024)
			glog.V(2).Info("Waiting for message in pipeToProcessInput()")
			size, err := l.stdIn.Read(msg)
			glog.V(2).Infof("Got message: %v", msg)
			if err == io.EOF {
				return
			} else if err != nil {
				panic(err)
			} else {
				l.inputCommand <- msg[0:size]
			}
		}
	}
}

func (l *stdinListener) pipeFromProcessOutput() {
	glog.V(2).Info("Starting pipeFromProcessOutput() routine")
	for {
		select {
		case a := <-l.processOutput:
			fmt.Fprintf(l.stdOut, "%s", a.Result)
		case <-l.processTerminated:
			glog.V(2).Infof("Received a processTerminated signal!")
			l.stop <- fmt.Errorf("Process terminated")
		}
	}
}
