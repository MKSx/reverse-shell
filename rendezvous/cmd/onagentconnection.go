package cmd

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/maxlaverse/reverse-shell/util"
)

type onAgentConnection struct{}

func (h onAgentConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := util.WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.V(2).Infof("Err:%s", err)
		return
	}
	glog.V(2).Infof("New agent %s", conn.RemoteAddr().String())

	agentTable.AddAgent(conn)

	go func() {
		defer conn.Close()
		//	defer close(done)
		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				glog.V(2).Infof("Agent '%s' disconnected. Clearing sessions! Reason: %s", conn.RemoteAddr().String(), err)

				ses := sessionTable.FindSessionByAgent(conn)
				for _, masterConn := range ses {
					glog.V(2).Infof("Session '%s' was lost due to agent failure", masterConn.Id)
					a := message.ProcessTerminated{
						Id: masterConn.Id,
					}
					masterConn.State = SESSION_LOST
					if masterConn != nil {
						for _, c := range masterConn.masterConn {
							c.WriteMessage(websocket.BinaryMessage, message.ToBinary(a))
						}
					}
				}
				agentTable.RemoveAgent(conn)
				return
			}
			b := message.FromBinary(m)
			switch v := b.(type) {
			case *message.ProcessOutput:
				glog.V(2).Infof("New Agent ProcessOutput for: %s (%d)", v.Id, len(v.Id))
				masterConn := sessionTable.FindSession(v.Id)
				if masterConn == nil {
					glog.V(2).Infof("That's bad session was lost %s", v.Id)
				} else {
					for _, c := range masterConn.masterConn {
						c.WriteMessage(websocket.BinaryMessage, m)
					}
				}

			case *message.ProcessCreated:
				glog.V(2).Infof("New Agent ProcessCreated for: %s (%d), %s\n", v.Id, len(v.Id), v.WantedId)

				s := Session{
					Id:        v.Id,
					agentConn: conn,
					State:     SESSION_OPEN,
				}
				sessionTable.AddSession(&s)

				if responseTable[v.WantedId] != nil {
					responseTable[v.WantedId] <- v.Id
					close(responseTable[v.WantedId])
					responseTable[v.WantedId] = nil
				}
			case *message.ProcessTerminated:
				glog.V(2).Infof("Session ended for: %s (%d), %s\n", v.Id, len(v.Id))
				//Create just session, sent back id, wait for attachement

				masterConn := sessionTable.FindSession(v.Id)
				masterConn.State = SESSION_CLOSED
				if masterConn == nil {
					glog.V(2).Infof("That's bad session was lost %s", v.Id)
				} else {
					for _, c := range masterConn.masterConn {
						c.WriteMessage(websocket.BinaryMessage, m)
					}
				}
			case *message.SessionTable:
				glog.V(2).Infof("Restoring session table")
				//Create just session, sent back id, wait for attachement
				for _, v2 := range v.Sessions {
					s := Session{
						Id:        v2,
						agentConn: conn,
						State:     SESSION_OPEN,
					}
					glog.V(2).Infof("Adding session: %s", v2)
					sessionTable.AddSession(&s)
				}

			default:
				glog.V(2).Infof("Received Agent an unknown message type: %v", v)
			}
		}
	}()
}
