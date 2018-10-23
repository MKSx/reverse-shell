package handler

import (
	"errors"
	"fmt"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/golang/glog"
	"github.com/maxlaverse/reverse-shell/message"
)

type Handler struct {
	processOutput     chan *ProcessOutput
	processTerminated chan *ProcessTerminated
	processTable      *ProcessTable
}

func NewHandler(processOutput chan *ProcessOutput, processTerminated chan *ProcessTerminated) *Handler {
	return &Handler{
		processTable:      newProcessTable(),
		processOutput:     processOutput,
		processTerminated: processTerminated,
	}
}

func (h *Handler) Sessions() []string {
	return h.processTable.List()
}

func (h *Handler) CreateProcess(commandLine string) string {
	name := trunc(namesgenerator.GetRandomName(0), message.SessionIdentifierMaxLen-1)
	processOutput := make(chan []byte)
	processTerminated := make(chan struct{})

	glog.V(2).Infof("Creating process '%s' from '%s'", name, commandLine)
	newProcess := h.processTable.New(name, commandLine)
	go newProcess.Attach(processOutput, processTerminated)

	go h.pipeOutput(processOutput, processTerminated, newProcess)
	go func() {
		newProcess.WaitForExit()
		glog.V(2).Infof("Process exited")
		processTerminated <- struct{}{}
	}()
	return name
}

func (h *Handler) KillProcess(processName string) error {
	p := h.processTable.Get(processName)
	return p.Kill()
}

func (h *Handler) ExecuteCommand(id string, command []byte) error {
	p := h.processTable.Get(id)
	if p == nil {
		glog.V(2).Infof("Process '%s' not found in table. Current table contains:", id)
		for _, element := range h.processTable.List() {
			glog.V(2).Infof("* '%s'", element)
		}
		return errors.New("Process not found")
	}
	if p.State != PROCESS_RUNNING {
		glog.V(2).Infof("The process '%s' is not running anymore", id)
		return errors.New(fmt.Sprintf("The process '%s' is not running anymore", id))
	}
	return p.Send(command)
}

func (h *Handler) pipeOutput(processOutput chan []byte, processTerminated chan struct{}, p *Process) {
PipeLoop:
	for {
		select {

		case o := <-processOutput:
			glog.V(2).Infof("Received %d bytes in processOutput", len(o))
			h.processOutput <- &ProcessOutput{
				Process: p,
				Result:  o,
			}

		case <-processTerminated:
			glog.V(2).Infof("Received a processTerminated signal")
			h.processTerminated <- &ProcessTerminated{Process: p}
			break PipeLoop
		}
	}
}

func trunc(text string, size int) string {
	if len(text) > size {
		return text[0 : size-1]
	}
	return text
}
