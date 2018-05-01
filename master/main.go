package main

import (
	"flag"
	"strconv"

	"github.com/maxlaverse/reverse-shell/master/cmd"
	"github.com/spf13/cobra"
)

func main() {
	verbose := 0
	command := &cobra.Command{
		Use:              "reverse-shell-master",
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.Set("logtostderr", "true")
			flag.Set("v", strconv.Itoa(verbose))
			flag.CommandLine.Parse([]string{})
		},
	}

	command.PersistentFlags().IntVarP(&verbose, "verbose", "v", 0, "Be verbose on log output")

	master := cmd.NewMasterCli()
	command.AddCommand(cmd.NewRendezVousCommand(master))
	command.AddCommand(cmd.NewListenCommand(master))

	if err := command.Execute(); err != nil {
		panic(err)
	}
}
