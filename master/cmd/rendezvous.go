package cmd

import (
	"github.com/spf13/cobra"
)

// NewRendezVousCommand creates a new cobra.Command for `reverse-shell-master rendezvous`
func NewRendezVousCommand(agent Cli) *cobra.Command {
	command := &cobra.Command{
		Use:              "rendezvous",
		Short:            "Connects to a remote rendez-vous point",
		TraverseChildren: true,
	}

	command.AddCommand(NewCreateCommand(agent))
	command.AddCommand(NewAttachCommand(agent))
	command.AddCommand(NewListSessionCommand(agent))
	command.AddCommand(NewListAgentCommand(agent))
	return command
}
