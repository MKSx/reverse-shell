package main

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/maxlaverse/reverse-shell/rendezvous/api"
)

type onSessionCreate struct{}

func (h onSessionCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Infof("On session create")
	decoder := json.NewDecoder(r.Body)
	var m api.CreateSession
	err := decoder.Decode(&m)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	agent := agentTable.FindAgent(m.Agent)
	if agent == nil {
		glog.V(2).Infof("Agent not found %s", m.Agent)
		return
	}

	glog.V(2).Infof("Agent found %s, creating session", m.Agent)
	responseTable["generated-token"] = make(chan string)
	m2 := message.CreateProcess{
		CommandLine: m.Command,
		Id:          "generated-token",
	}
	agent.WriteMessage(websocket.BinaryMessage, message.ToBinary(m2))

	t := <-responseTable["generated-token"]
	glog.V(2).Infof("New session, answering %s", t)
	w.Write([]byte(t))
}
