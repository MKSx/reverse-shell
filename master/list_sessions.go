package main

import (
	"encoding/json"
	"net/http"

	"github.com/maxlaverse/reverse-shell/rendezvous/api"
)

func ListSessions(url string) ([]api.SessionListResponseAgent, error) {
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
