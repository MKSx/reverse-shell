package cmd

import (
	"fmt"
	"testing"

	"github.com/maxlaverse/reverse-shell/agent/cmd/test"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/stretchr/testify/assert"
)

func TestMissingArgument(t *testing.T) {
	cli := NewAgentCli()
	websocketListener := NewWebsocketListenerCommand(cli)
	websocketListener.SetArgs([]string{})

	websocketListener.SilenceUsage = true
	websocketListener.SilenceErrors = true
	err := websocketListener.Execute()

	assert.Equal(t, "a url must be given", err.Error())
}

func TestStructComparison(t *testing.T) {
	// Initialize a fake WebSocket master
	dummyMaster := test.NewDummyMaster(t)

	// Initialize the listener
	listener := newWebsocketListener(fmt.Sprintf("http://127.0.0.1:%d", dummyMaster.Port()))

	// Start the listener
	err := listener.Start()
	if err != nil {
		t.Fatalf("Error starting the listener: %v", err)
	}
	go listener.Listen()

	// Should receive the session table
	c := dummyMaster.ReadRawMessage().(*message.SessionTable)
	assert.Equal(t, len(c.Sessions), 0, "The session table was not empty")

	// Send a process creation request
	dummyMaster.SendMessage(message.CreateProcess{
		CommandLine: "bash --norc",
	})

	// Should receive a process structure
	c2 := dummyMaster.ReadRawMessage().(*message.ProcessCreated)

	// Sends the test command
	dummyMaster.SendMessage(message.ExecuteCommand{
		Command: []byte("echo 'THIS IS A TEST';exit\n"),
		Id:      c2.Id,
	})

	// Should receive the process output
	response := dummyMaster.ReadMessageUntilTerminated()

	assert.Contains(t, response, "THIS IS A TEST\r\nexit\r\n", "The response didn't contain the expected message")
}
