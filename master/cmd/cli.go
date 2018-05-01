package cmd

import (
	"time"

	"github.com/golang/glog"
)

// MasterCli is the base of the command-line tool client
type MasterCli struct {
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

// NewMasterCli is the root of the client
func NewMasterCli() Cli {
	return &MasterCli{}
}

// SafeStart starts a listener and restart it if it crashes
func (c *MasterCli) SafeStart(l Listener) error {
	l.Start()
	for {
		l.Listen()
		time.Sleep(3 * time.Second)
		glog.V(0).Infof("Main loop exited. Restarting")
	}
}
