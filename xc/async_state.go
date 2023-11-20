package xc

import (
	"reflect"
)

type AsyncState interface {
	// GetStateOptions defines the optional configuration of this state definition.
	GetStateOptions() *AsyncStateOptions

	// WaitUntil is the method to set up commands set up to wait for, before `Execute` API is invoked.
	//           It's optional -- use xc.AsyncStateNoWaitUntil to skip this( then Execute will be invoked directly instead)
	//
	//  ctx              the context info of this API invocation, like process start time, processId, etc
	//  input            the state input
	//  Communication    TODO
	// @return the requested commands for this state
	///
	WaitUntil(ctx Context, input Object, communication Communication) (*CommandRequest, error)

	// Execute is the method to execute and decide what to do next.
	// It's invoked after commands from WaitUntil are completed, or if WaitUntil is skipped(not implemented).
	//
	//  ctx              the context info of this API invocation, like process start time, processId, etc
	//  input            the state input
	//  CommandResults   the results of the command that executed by WaitUntil
	//  Persistence      TODO
	//  Communication    TODO
	// @return the decision of what to do next(e.g. transition to next states or closing process)
	Execute(ctx Context, input Object, commandResults CommandResults, persistence Persistence, communication Communication) (*StateDecision, error)
}

// GetFinalStateId returns the stateId that will be registered and used
// if the asyncState is from myStruct{} under mywf package, the method returns "mywf.myStruct"
func GetFinalStateId(asyncState AsyncState) string {
	options := asyncState.GetStateOptions()
	if options == nil || options.StateId == "" {
		return getSimpleTypeNameFromReflect(asyncState)
	}
	return options.StateId
}

// AsyncStateDefaults is a convenient struct to put into your state implementation to save the boilerplate code of returning default values
// Example usage:
//
//	type myStateImpl struct{
//	    AsyncStateDefaults
//	}
//
// Then myStateImpl doesn't have to implement WaitUntil, Execute or GetStateOptions
type AsyncStateDefaults struct {
	defaultStateOptions
}

// AsyncStateDefaultsSkipWaitUntil is required to skip WaitUntil
// put into your state implementation to save the boilerplate code of returning default values
// Example usage:
//
//	type myStateImpl struct{
//	    AsyncStateDefaultsSkipWaitUntil
//	}
//
// Then myStateImpl will skip WaitUntil, and doesn't have to implement GetStateOptions
type AsyncStateDefaultsSkipWaitUntil struct {
	defaultStateOptions
	skipWaitUntil
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
		if field.Type.String() == "xc.AsyncStateDefaultsSkipWaitUntil" {
			return true
		}
	}

	return false
}

type defaultStateOptions struct{}

func (d defaultStateOptions) GetStateOptions() *AsyncStateOptions {
	return nil
}

type skipWaitUntil struct{}

func (d skipWaitUntil) WaitUntil(ctx Context, input Object, communication Communication) (*CommandRequest, error) {
	panic("this method is for skipping WaitUntil. It should never be called")
}
