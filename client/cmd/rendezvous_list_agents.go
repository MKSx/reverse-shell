package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maxlaverse/reverse-shell/rendezvous/api"
	"github.com/spf13/cobra"
)

// NewListAgentCommand creates a new cobra.Command for `reverse-shell-client rendezvous list-agents`
func NewListAgentCommand(agent Cli) *cobra.Command {
	var url string
	cmd := &cobra.Command{
		Use:              "list-agents",
		Short:            "list all the agents available on a rendez-vous",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("List of agents:\n")
			l, _ := listAgents(url)
			for _, v := range l {
				fmt.Printf(" * %s\n", v.Name)
			}

		},
	}

	cmd.Flags().StringVarP(&url, "url", "", "", "Url of the rendez-vous point")

	return cmd
}

func listAgents(url string) ([]api.AgentListResponseAgent, error) {
	resp, err := http.Get(url + "/agent/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	target := make([]api.AgentListResponseAgent, 0)
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil
}
