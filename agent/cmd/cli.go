package cmd

import (
	"flag"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

// GetCommand returns the main cobra.Command of the CLI
func GetCommand() *cobra.Command {
	verbose := 0
	command := &cobra.Command{
		Use:              "reverse-shell-agent",
		Short:            "A Go agents listening for remote commands",
		Long:             `Starts the agent and listens for remote commands.`,
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.Set("logtostderr", "true")
			flag.Set("v", strconv.Itoa(verbose))
			flag.CommandLine.Parse([]string{})
		},
		DisableAutoGenTag: true,
	}

	command.PersistentFlags().IntVarP(&verbose, "verbose", "v", 0, "Be verbose on log output")

	agent := NewAgentCli()
	command.AddCommand(NewStdinListenerCommand(agent))
	command.AddCommand(NewTCPListenerCommand(agent))
	command.AddCommand(NewTCPdirectListenerCommand(agent))
	command.AddCommand(NewWebsocketListenerCommand(agent))
	return command
}

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
