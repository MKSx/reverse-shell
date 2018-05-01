package cmd

import (
	"encoding/json"
	"net/http"

	"github.com/maxlaverse/reverse-shell/rendezvous/api"
)

func (h onAgentList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := []api.AgentListResponseAgent{}
	for _, c := range agentTable.agents {
		p = append(p, api.AgentListResponseAgent{Name: c.RemoteAddr().String()})
	}
	json.NewEncoder(w).Encode(p)
}

type onAgentList struct{}
