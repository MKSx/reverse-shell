package cmd

import (
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/maxlaverse/reverse-shell/util"
)

type onSessionAttach struct{}

func (h onSessionAttach) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedSession := r.URL.Path[16:]
	glog.V(2).Infof("New attachement request for '%s'", requestedSession)

	session := sessionTable.FindSession(requestedSession)
	if session == nil {
		glog.V(2).Infof("Session not found")
		w.Write([]byte("Session not found"))
		return
	}

	conn, err := util.WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	session.clientConn = append(session.clientConn, conn)

	go func() {
		defer conn.Close()

		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				glog.V(2).Infof("ReadMessage error on the clientChannel: %s", err)
				return
			}
			b := message.FromBinary(m)
			switch v := b.(type) {
			case *message.CreateProcess:
				session.agentConn.WriteMessage(websocket.BinaryMessage, m)
			case *message.ExecuteCommand:
				session.agentConn.WriteMessage(websocket.BinaryMessage, m)
			default:
				glog.V(2).Infof("Received Client an unknown message type: %v", v)
			}
		}
	}()
}
