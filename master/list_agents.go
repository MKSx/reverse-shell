package main

import (
	"encoding/json"
	"net/http"

	"github.com/maxlaverse/reverse-shell/rendezvous/api"
)

func ListAgents(url string) ([]api.AgentListResponseAgent, error) {
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
