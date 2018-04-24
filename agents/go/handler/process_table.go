package handler

// ProcessTable stores the processes currently started by this agent
type ProcessTable struct {
	processes map[string]*Process
}

func newProcessTable() *ProcessTable {
	return &ProcessTable{
		processes: make(map[string]*Process),
	}
}

// New adds a process to the local process table
func (s *ProcessTable) New(processName string, command string) *Process {
	newProcess := NewProcess(processName, command)
	s.processes[newProcess.Name] = newProcess
	return newProcess
}

// List all the processes started by this agent
func (s *ProcessTable) List() []string {
	keys := make([]string, 0, len(s.processes))
	for k := range s.processes {
		keys = append(keys, k)
	}
	return keys
}

// Get returns a process from the process table
func (s *ProcessTable) Get(processName string) *Process {
	return s.processes[processName]
}

// Remove deletes a process from the process table
func (s *ProcessTable) Remove(process *Process) {
	delete(s.processes, process.Name)
}
