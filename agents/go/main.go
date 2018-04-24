package main

import (
	"flag"
	"strconv"

	"github.com/maxlaverse/reverse-shell/agents/go/cmd"
	"github.com/spf13/cobra"
)

func main() {
	verbose := 0
	command := &cobra.Command{
		Use:              "reverse-shell-agent",
		Short:            "Agents listening for remote commands",
		Long:             `Starts an agent listening for remote commands. The Agent can receive remote commands is various ways.`,
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.Set("logtostderr", "true")
			flag.Set("v", strconv.Itoa(verbose))
			flag.CommandLine.Parse([]string{})
		},
	}

	command.PersistentFlags().IntVarP(&verbose, "verbose", "v", 0, "Be verbose on log output")

	agent := cmd.NewAgentCli()
	command.AddCommand(cmd.NewStdinListenerCommand(agent))
	command.AddCommand(cmd.NewTCPListenerCommand(agent))
	command.AddCommand(cmd.NewTCPdirectListenerCommand(agent))
	command.AddCommand(cmd.NewWebsocketListenerCommand(agent))

	if err := command.Execute(); err != nil {
		panic(err)
	}
}
