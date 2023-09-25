package xdb

type Registry interface {
	// AddProcess registers a process
	AddProcess(processDef Process) error
	// AddProcesses registers multiple processes
	AddProcesses(processDefs ...Process) error
	// GetAllRegisteredProcessTypes returns all the process types that have been registered
	GetAllRegisteredProcessTypes() []string

	// below are all for internal implementation
	getProcess(wfType string) Process
	getProcessStartingState(wfType string) AsyncState
	getProcessState(wfType string, id string) AsyncState
}

func NewRegistry() Registry {
	return &registryImpl{
		processStore:  map[string]Process{},
		startingState: map[string]AsyncState{},
		stateStore:    map[string]map[string]AsyncState{},
	}
}
