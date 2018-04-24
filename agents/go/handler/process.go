package handler

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	"github.com/kr/pty"
)

type Process struct {
	Name       string
	State      ProcessState
	cmd        *exec.Cmd
	descriptor *os.File
}

type ProcessState int

const (
	PROCESS_RUNNING ProcessState = 1 + iota
	PROCESS_FAILED_TO_START
	PROCESS_EXITED
	PROCESS_TIMEOUT
)

func NewProcess(id string, commandLine string) *Process {
	// Improve the split, with shellString
	args := strings.Split(commandLine, " ")

	// Use a CommandContext ?
	process := Process{
		Name: id,
		cmd:  exec.Command(args[0], args[1:]...),
	}

	glog.V(2).Infof("New process '%s' ready to be started", id)

	f, err := pty.Start(process.cmd)
	if err != nil {
		process.State = PROCESS_FAILED_TO_START
		return &process
	}
	process.descriptor = f
	process.State = PROCESS_RUNNING

	return &process
}

func (p *Process) Attach(outputChannel chan []byte, processCloseChannel chan struct{}) {
	for {
		select {
		default:
			var msg = make([]byte, 1024)
			size, err := p.descriptor.Read(msg)
			if err == io.EOF {
				glog.V(2).Infof("Process read returned EOF")
				processCloseChannel <- struct{}{}
				return
			} else if err != nil {
				glog.V(2).Infof("No idea what to do!An error has occured while reading: %s", err)
				panic(err)
			} else {
				glog.V(2).Infof("Received %d bytes from process", size)
				outputChannel <- msg[0:size]
			}
		}
	}
}

func (p *Process) WaitForExit() {
	p.cmd.Wait()
	p.State = PROCESS_EXITED
}

func (p *Process) Kill() error {
	p.descriptor.Sync()
	return p.cmd.Process.Kill()
}

func (p *Process) Send(data []byte) error {
	glog.V(2).Infof("Sending %d byte to process", len(data))
	n, err := p.descriptor.Write(data)
	glog.V(2).Infof("%d bytes written", n)
	return err
}
