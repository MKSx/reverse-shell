package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maxlaverse/reverse-shell/rendezvous/api"
	"github.com/spf13/cobra"
)

// NewListSessionCommand creates a new cobra.Command for `reverse-shell-client rendezvous session-list`
func NewListSessionCommand(agent Cli) *cobra.Command {
	var url string
	cmd := &cobra.Command{
		Use:              "session-list",
		Short:            "list all the sessions available on a rendez-vous",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf("List of sessions:\n")
			l, _ := listSessions(url)
			for _, v := range l {
				fmt.Printf(" * %s => agent: %s, masters: %s, state: %s\n", v.Name, v.Agent, v.Clients, v.State)
			}

		},
	}

	cmd.Flags().StringVarP(&url, "url", "", "", "Url of the rendez-vous point")

	return cmd
}

func listSessions(url string) ([]api.SessionListResponseAgent, error) {
	resp, err := http.Get(url + "/session/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	target := make([]api.SessionListResponseAgent, 0)
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil
}
