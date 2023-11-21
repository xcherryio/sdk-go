package xc

type Registry interface {
	// AddProcess registers a process
	AddProcess(processDef Process) error
	// AddProcesses registers multiple processes
	AddProcesses(processDefs ...Process) error
	// GetAllRegisteredProcessTypes returns all the process types that have been registered
	GetAllRegisteredProcessTypes() []string

	// below are all for internal implementation
	getProcess(prcType string) Process
	getProcessStartingState(prcType string) AsyncState
	getProcessState(prcType string, id string) AsyncState
	getPersistenceSchema(prcType string) PersistenceSchema
	getGlobalAttributeKeyToDefs(prcType string) map[string]internalGlobalAttrDef
	getGlobalAttributeTableColumnToKey(prcType string) map[string]string
}

func NewRegistry() Registry {
	return &registryImpl{
		processStore:                map[string]Process{},
		startingState:               map[string]AsyncState{},
		stateStore:                  map[string]map[string]AsyncState{},
		persistenceSchemaStore:      map[string]PersistenceSchema{},
		globalAttributeKeyToDef:     map[string]map[string]internalGlobalAttrDef{},
		globalAttrTableColNameToKey: map[string]map[string]string{},
	}
}
