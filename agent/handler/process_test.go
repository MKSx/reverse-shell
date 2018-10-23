package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessExit(t *testing.T) {
	p := NewProcess("dummy", "true")
	p.WaitForExit()

	assert.Equal(t, PROCESS_EXITED, p.State, "Wrong process state")
}

func TestProcessKill(t *testing.T) {
	p := NewProcess("dummy", "sleep 86400")
	c := make(chan struct{})
	go func() {
		p.WaitForExit()
		c <- struct{}{}
	}()

	p.Kill()
	<-c
	assert.Equal(t, PROCESS_EXITED, p.State, "Wrong process state")
}

func TestProcessRunning(t *testing.T) {
	p := NewProcess("dummy", "sleep 86400")
	assert.Equal(t, PROCESS_RUNNING, p.State, "Wrong process state")
	p.Kill()
}

func TestProcessSend(t *testing.T) {
	p := NewProcess("dummy", "echo")

	processOutput := make(chan []byte)
	processTerminated := make(chan struct{})
	go p.Attach(processOutput, processTerminated)

	p.Send([]byte("hi this is a test"))
	p.Kill()
	p.WaitForExit()

	assert.Equal(t, PROCESS_EXITED, p.State, "Wrong process state")
	assert.Equal(t, "hi this is a test", string(<-processOutput), "Wrong process output")
}
