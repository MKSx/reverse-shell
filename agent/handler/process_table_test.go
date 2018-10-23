package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessAdd(t *testing.T) {
	table := newProcessTable()
	table.New("dummy", "foo bar")

	p := table.processes["dummy"]
	assert.Equal(t, "dummy", p.Name, "Wrong process name")
	assert.Equal(t, PROCESS_FAILED_TO_START, p.State, "Wrong process state")
	assert.Equal(t, []string{"foo", "bar"}, p.cmd.Args, "The command had wrong argument")
}

func TestProcessList(t *testing.T) {
	table := newProcessTable()
	table.New("dummy", "foo bar")
	table.New("dummy2", "foo2 bar2")

	assert.Contains(t, table.List(), "dummy", "Missing dummy in the list")
	assert.Contains(t, table.List(), "dummy2", "Missing dummy2 in the list")
}

func TestProcessGet(t *testing.T) {
	table := newProcessTable()
	table.New("dummy", "foo bar")
	p := table.Get("dummy")
	assert.Equal(t, "dummy", p.Name, "Wrong process name")
	assert.Equal(t, PROCESS_FAILED_TO_START, p.State, "Wrong process state")
	assert.Equal(t, []string{"foo", "bar"}, p.cmd.Args, "The command had wrong argument")
}

func TestProcessRemove(t *testing.T) {
	table := newProcessTable()
	table.New("dummy", "foo bar")

	table.Remove(table.Get("dummy"))
	_, ok := table.processes["dummy"]

	assert.Equal(t, false, ok, "Process dummy still in the table")
}
