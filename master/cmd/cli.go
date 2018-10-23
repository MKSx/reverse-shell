package cmd

import (
	"flag"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	verbose := 0
	command := &cobra.Command{
		Use:              "reverse-shell-master",
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.Set("logtostderr", "true")
			flag.Set("v", strconv.Itoa(verbose))
			flag.CommandLine.Parse([]string{})
		},
		DisableAutoGenTag: true,
	}

	command.PersistentFlags().IntVarP(&verbose, "verbose", "v", 0, "Be verbose on log output")

	master := NewMasterCli()
	command.AddCommand(NewRendezVousCommand(master))
	command.AddCommand(NewListenCommand(master))
	return command
}

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
