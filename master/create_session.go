package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/maxlaverse/reverse-shell/rendezvous/api"
)

func CreateSession(url string, agent string) string {
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
