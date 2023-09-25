package xdb

import (
	"reflect"
)

type AsyncState interface {
	// GetStateId defines the StateId of this state definition.
	// the StateId is being used for WorkerService to choose the right AsyncState to execute Start/Execute APIs
	// It is a default value when return empty string.
	// Default is the package + struct name of the state definition and ignores the import paths and aliases.
	// e.g. if the process is from myStruct{} under mywf package, the simple name is just "mywf.myStruct". Underneath, it's from reflect.TypeOf(wf).String().
	GetStateId() string

	// WaitUntil is the method to set up commands set up to wait for, before `Execute` API is invoked.
	//           It's optional -- use xdb.AsyncStateNoWaitUntil to skip this( then Execute will be invoked directly instead)
	//
	//  ctx              the context info of this API invocation, like process start time, processId, etc
	//  input            the state input
	//  Communication    TODO
	// @return the requested commands for this state
	///
	WaitUntil(ctx XdbContext, input Object, communication Communication) (*CommandRequest, error)

	// Execute is the method to execute and decide what to do next.
	// It's invoked after commands from WaitUntil are completed, or if WaitUntil is skipped(not implemented).
	//
	//  ctx              the context info of this API invocation, like process start time, processId, etc
	//  input            the state input
	//  CommandResults   the results of the command that executed by WaitUntil
	//  Persistence      TODO
	//  Communication    TODO
	// @return the decision of what to do next(e.g. transition to next states or closing process)
	Execute(ctx XdbContext, input Object, commandResults CommandResults, persistence Persistence, communication Communication) (*StateDecision, error)
}

// GetFinalStateId returns the stateId that will be registered and used
// if the asyncState is from myStruct{} under mywf package, the method returns "mywf.myStruct"
func GetFinalStateId(asyncState AsyncState) string {
	sid := asyncState.GetStateId()
	if sid == "" {
		return getSimpleTypeNameFromReflect(asyncState)
	}
	return sid
}

// AsyncStateDefaults is a convenient struct to put into your state implementation to save the boilerplate code of returning default values
// Example usage:
//
//	type myStateImpl struct{
//	    AsyncStateDefaults
//	}
//
// Then myStateImpl doesn't have to implement GetStateId
type AsyncStateDefaults struct {
	defaultStateId
}

// AsyncStateNoWaitUntil is required to skip WaitUntil
// put into your state implementation to save the boilerplate code of returning default values
// Example usage:
//
//	type myStateImpl struct{
//	    AsyncStateNoWaitUntil
//	}
//
// Then myStateImpl will skip WaitUntil, and doesn't have to implement GetStateId
type AsyncStateNoWaitUntil struct {
	defaultStateId
	noWaitUntil
}

func ShouldSkipWaitUntilAPI(state AsyncState) bool {
	rt := reflect.TypeOf(state)
	var t reflect.Type
	if rt.Kind() == reflect.Pointer {
		t = rt.Elem()
	} else if rt.Kind() == reflect.Struct {
		t = rt
	} else {
		panic("a the state must be an pointer or a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.String() == "xdb.AsyncStateNoWaitUntil" {
			return true
		}
	}

	return false
}

type defaultStateId struct{}

func (d defaultStateId) GetStateId() string {
	return ""
}

type noWaitUntil struct{}

func (d noWaitUntil) WaitUntil(ctx XdbContext, input Object, persistence Persistence, communication Communication) (*CommandRequest, error) {
	panic("this method is for skipping WaitUntil. It should never be called")
}
