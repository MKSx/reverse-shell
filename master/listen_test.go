package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/message"
)

//
// func TestAutomaticSessionRecovery(t *testing.T) {
// 	stdinChannel := make(chan []byte)
// 	stdoutChannel := make(chan []byte)
//
// 	srv := httptest.NewServer(http.Handler(onConnectMaster{stdinChannel: stdinChannel, stdoutChannel: stdoutChannel}))
// 	u, _ := url.Parse(srv.URL)
// 	u.Scheme = "ws"
// 	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		log.Fatalf("cannot make websocket connection: %v", err)
// 	}
//
// 	conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(message.SessionTable{Sessions: []string{"abc"}}))
//
// 	g := <-stdoutChannel
// 	if string(g) != "Recovering session named: abc\n" {
// 		log.Fatalf("wrong output log: %q\n", g)
// 	}
//
// 	conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(message.ProcessTerminated{}))
// }

func TestAutomaticSessionCreation(t *testing.T) {
	stdinChannel := make(chan []byte)
	stdoutChannel := make(chan []byte)

	srv := httptest.NewServer(http.Handler(onConnectMaster{stdinChannel: stdinChannel, stdoutChannel: stdoutChannel}))
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("cannot make websocket connection: %v", err)
	}

	conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(message.SessionTable{}))

	g := <-stdoutChannel
	if string(g) != "No session yet for incoming connection starting a new session" {
		log.Fatalf("wrong output log: %q\n", g)
	}

	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Fatalf("cannot read message: %v", err)
	}
	b := message.FromBinary(p).(*message.CreateProcess)
	if b.CommandLine != "bash --norc" {
		log.Fatalf("wrong command line return: %q\n", b)
	}

	srv.CloseClientConnections()
}

//
// func TestSendingCommandTooSoon(t *testing.T) {
// 	stdinChannel := make(chan []byte)
// 	stdoutChannel := make(chan []byte)
//
// 	handler := onConnectMaster{stdinChannel: stdinChannel, stdoutChannel: stdoutChannel}
// 	srv := httptest.NewServer(http.Handler(handler))
// 	u, _ := url.Parse(srv.URL)
// 	u.Scheme = "ws"
// 	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		log.Fatalf("cannot make websocket connection: %v", err)
// 	}
//
// 	g := <-stdoutChannel
// 	if string(g) != "New incoming connection: starting a new session" {
// 		log.Fatalf("wrong output log: %q\n", g)
// 	}
//
// 	_, p, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Fatalf("cannot read message: %v", err)
// 	}
// 	b := message.FromBinary(p).(*message.CreateProcess)
// 	if b.CommandLine != "bash --norc" {
// 		log.Fatalf("wrong command line sent: %q\n", b)
// 	}
//
// 	stdinChannel <- []byte("uptime")
//
// 	g = <-stdoutChannel
// 	if string(g) != "Session is not ready" {
// 		log.Fatalf("wrong output log: %q\n", g)
// 	}
// }
//
// func TestSendingCommand(t *testing.T) {
// 	stdinChannel := make(chan []byte)
// 	stdoutChannel := make(chan []byte)
//
// 	handler := onConnectMaster{stdinChannel: stdinChannel, stdoutChannel: stdoutChannel}
// 	srv := httptest.NewServer(http.Handler(handler))
// 	u, _ := url.Parse(srv.URL)
// 	u.Scheme = "ws"
// 	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		log.Fatalf("cannot make websocket connection: %v", err)
// 	}
//
// 	g := <-stdoutChannel
// 	if string(g) != "New incoming connection: starting a new session" {
// 		log.Fatalf("wrong output log: %q\n", g)
// 	}
//
// 	_, p, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Fatalf("cannot read message: %v", err)
// 	}
// 	b := message.FromBinary(p).(*message.CreateProcess)
// 	if b.CommandLine != "bash --norc" {
// 		log.Fatalf("wrong command line sent: %q\n", b)
// 	}
//
// 	conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(message.ProcessCreated{Id: "14", WantedId: "2"}))
//
// 	g = <-stdoutChannel
// 	if string(g) != "New session is named: 14\n" {
// 		log.Fatalf("wrong output log: %q\n", g)
// 	}
//
// 	stdinChannel <- []byte("uptime")
//
// 	_, p, err = conn.ReadMessage()
// 	if err != nil {
// 		log.Fatalf("cannot read message: %v", err)
// 	}
// 	c := message.FromBinary(p).(*message.ExecuteCommand)
// 	if string(c.Command) != "uptime" {
// 		log.Fatalf("wrong command line sent: %q\n", c)
// 	}
// 	if c.Id != "14" {
// 		log.Fatalf("wrong session id sent: %q\n", c)
// 	}
//}
