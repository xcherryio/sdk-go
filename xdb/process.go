package xdb

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

	// GetProcessType defines the processType of this process definition.
	// See GetFinalProcessType for default value when return empty string.
	// It's the package + struct name of the process instance and ignores the import paths and aliases.
	// e.g. if the process is from myStruct{} under mywf package, the simple name is just "mywf.myStruct". Underneath, it's from reflect.TypeOf(wf).String().
	GetProcessType() string
}

// GetFinalProcessType returns the process type that will be registered
// if the process is from &myStruct{} or myStruct{} under mywf package, the method returns "mywf.myStruct"
func GetFinalProcessType(wf Process) string {
	wfType := wf.GetProcessType()
	if wfType == "" {
		simpleType := getSimpleTypeNameFromReflect(wf)
		return simpleType
	}
	return wfType
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
// Then myPcImpl doesn't have to implement GetProcessType or GetAsyncStateSchema
type ProcessDefaults struct {
}

func (d ProcessDefaults) GetProcessType() string {
	return ""
}

func (d ProcessDefaults) GetAsyncStateSchema() StateSchema {
	return StateSchema{}
}
