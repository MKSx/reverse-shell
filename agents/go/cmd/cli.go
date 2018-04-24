package cmd

import (
	"time"

	"github.com/golang/glog"
)

// AgentCli is the base of the command-line tool client
type AgentCli struct {
}

// Cli is the root of the client
type Cli interface {
	SafeStart(l Listener) error
}

// Listener is the interface for any kind of listener
type Listener interface {
	Start() error
	Listen() error
}

// NewAgentCli is the root of the client
func NewAgentCli() Cli {
	return &AgentCli{}
}

// SafeStart starts a listener and restart it if it crashes
func (c *AgentCli) SafeStart(l Listener) error {
	l.Start()
	for {
		l.Listen()
		time.Sleep(3 * time.Second)
		glog.V(0).Infof("Main loop exited. Restarting")
	}
}
