package cmd

import (
	"github.com/maxlaverse/reverse-shell/master/handler"
	"github.com/spf13/cobra"
)

// NewListenCommand creates a new cobra.Command for `reverse-shell-master rendezvous listen`
func NewListenCommand(agent Cli) *cobra.Command {
	var port int32
	cmd := &cobra.Command{
		Use:              "listen",
		Short:            "listen for agents to connect using websocket",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			handler.Listen(port)
		},
	}

	cmd.Flags().Int32VarP(&port, "port", "", 0, "port to listen to ")

	return cmd
}
