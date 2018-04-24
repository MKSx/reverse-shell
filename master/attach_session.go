package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/message"
)

func AttachSession(url string, sessionId string) {

	// Support Proxy
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)

	fmt.Printf("Attaching to %s\n", sessionId)

	conn, _, err := websocket.DefaultDialer.Dial("ws"+url[4:]+"/session/attach/"+sessionId, http.Header{"Origin": {url}})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", sessionId)

	var processId = sessionId

	go func() {
		defer conn.Close()
		//	defer close(done)
		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("ReadMessage error:", err)
				return
			}
			b := message.FromBinary(m)
			switch v := b.(type) {
			case *message.ProcessOutput:
				os.Stdout.Write(v.Data)
			case *message.ProcessCreated:
				fmt.Printf("New session is named: %s\n", v.Id)
				processId = v.Id
			case *message.ProcessTerminated:
				fmt.Printf("Session closed: %s\n", v.Id)
				os.Exit(0)
			default:
				fmt.Printf("Received an unknown message type: %v\n", v)
			}
		}
	}()

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
				m := message.ExecuteCommand{
					Id:      processId,
					Command: msg[0:size],
				}
				conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(m))
			}
		}
	}
}
