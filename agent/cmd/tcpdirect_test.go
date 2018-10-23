package cmd

import (
	"fmt"
	"net"
	"testing"

	"github.com/maxlaverse/reverse-shell/agent/cmd/test"
	"github.com/stretchr/testify/assert"
)

func TestExecutingCommandDirectTcp(t *testing.T) {
	// Initialize listener
	listener := newTCPDirectListener("127.0.0.1", 0)

	// Start the listener
	err := listener.Start()
	if err != nil {
		t.Fatalf("Error starting the listener: %v", err)
	}

	// Get port of the agent
	p, err := test.AddrToPort(listener.Addr())
	if err != nil {
		t.Fatalf("Error getting the port of the agent: %v", err)
	}

	// Watch the listener's exit
	doneCh := make(chan error)
	go test.WaitFor(listener.Listen, doneCh)

	// Connect to the agent and send test command
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", p))
	if err != nil {
		t.Fatalf("Error connecting to the qgent: %v", err)
	}
	conn.Write([]byte("echo 'THIS IS A TEST';exit\n"))

	// Wait for the listener to exit
	err = <-doneCh

	// Read the response from the agent
	var msg = make([]byte, 1024)
	_, err2 := conn.Read(msg)
	if err2 != nil {
		t.Fatalf("Error reading response from agent: %v", err2)
	}

	// Test the result
	assert.Equal(t, err.Error(), "Lost connection", "bad")
	assert.Contains(t, string(msg), "THIS IS A TEST\r\nexit\r\nProcess terminated", "bad")
}
