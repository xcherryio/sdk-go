package xc

import (
	"reflect"
	"strings"
)

// Process is the interface to define a process definition.
// Process is a top level concept in XDB
type Process interface {
	// GetAsyncStateSchema defines the AsyncStates of the process.
	// If there is no startingState, the process will not start any state execution after process stated.
	// Application can still use RPC to invoke new state execution later.
	GetAsyncStateSchema() StateSchema
	// GetPersistenceSchema defines the persistence schema of the process
	GetPersistenceSchema() PersistenceSchema
	// GetProcessOptions defines the options for the process
	// Note that they can be overridden by the ProcessStartOptions when starting a process
	GetProcessOptions() ProcessOptions
}

// GetFinalProcessType returns the process type that will be registered
// if the process is from &myStruct{} or myStruct{} under mywf package, the method returns "mywf.myStruct"
func GetFinalProcessType(wf Process) string {
	options := wf.GetProcessOptions()
	if options.ProcessType == "" {
		simpleType := getSimpleTypeNameFromReflect(wf)
		return simpleType
	}
	return options.ProcessType
}

func getSimpleTypeNameFromReflect(obj interface{}) string {
	rt := reflect.TypeOf(obj)
	rtStr := strings.TrimLeft(rt.String(), "*")
	return rtStr
}

// ProcessDefaults is a convenient struct to put into your process implementation to save the boilerplate code of returning default values
// Example usage :
//
//	type myPcImpl struct{
//	    ProcessDefaults
//	}
//
// Then myPcImpl doesn't have to implement GetProcessOptions or GetAsyncStateSchema
type ProcessDefaults struct {
}

var _ Process = (*ProcessDefaults)(nil)

func (d ProcessDefaults) GetAsyncStateSchema() StateSchema {
	return StateSchema{}
}

func (d ProcessDefaults) GetPersistenceSchema() PersistenceSchema {
	return NewEmptyPersistenceSchema()
}

func (d ProcessDefaults) GetProcessOptions() ProcessOptions {
	return NewDefaultProcessOptions()
}
