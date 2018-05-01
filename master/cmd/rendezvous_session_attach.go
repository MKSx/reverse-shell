package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/spf13/cobra"
)

// NewAttachCommand creates a new cobra.Command for `reverse-shell-master rendezvous session-attach`
func NewAttachCommand(agent Cli) *cobra.Command {
	var url string
	cmd := &cobra.Command{
		Use:              "attach",
		Short:            "attach to an existing session",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			AttachSession(url, args[0])
		},
	}

	cmd.Flags().StringVarP(&url, "url", "", "", "Url of the rendez-vous point")

	return cmd
}

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
