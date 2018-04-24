package cmd

import (
	"net"
	"testing"

	"github.com/maxlaverse/reverse-shell/agents/go/cmd/test"
	"github.com/stretchr/testify/assert"
)

func TestExecutingCommandTcp(t *testing.T) {
	// Start a fake master
	tcpMaster, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Error start dummy tcp master: %v", err)
	}

	port, err := test.AddrToPort(tcpMaster.Addr())
	if err != nil {
		t.Fatalf("Error while getting dummy tcp master port from '%s': %v", tcpMaster.Addr(), err)
	}

	// Initialize listener
	listener := newTCPListener("0.0.0.0", port)

	// Watch the listener's exit
	doneCh := make(chan error)
	go test.WaitFor(listener.Listen, doneCh)

	// Accept the agent connection and send test command
	r, err := tcpMaster.Accept()
	if err != nil {
		t.Fatalf("Error accepting connection on dummy tcp master: %v", err)
	}
	r.Write([]byte("echo 'THIS IS A TEST';exit\n"))

	// Wait for the listener to exit
	err = <-doneCh

	// Read the response from the agent
	var msg = make([]byte, 1024)
	_, err2 := r.Read(msg)
	if err2 != nil {
		t.Fatalf("Error reading response from agent: %v", err2)
	}

	// Test the result
	assert.Equal(t, err.Error(), "Lost connection", "Bad terminaison message")
	assert.Contains(t, string(msg), "THIS IS A TEST\r\nexit\r\nProcess terminated", "bad")
}
