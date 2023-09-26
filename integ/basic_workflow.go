package integ

import "github.com/xdblab/xdb-golang-sdk/xdb"

type basicWorkflow struct {
	xdb.ProcessDefaults
}

func (b basicWorkflow) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.WithStartingState(&state1{}, &state2{})
}

type state1 struct {
	xdb.AsyncStateDefaults
}

func (b state1) WaitUntil(ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication) (*xdb.CommandRequest, error) {
	return nil, nil
}

func (b state1) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)
	return xdb.SingleNextState(state2{}, i+1), nil
}

type state2 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b state2) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)
	return xdb.ForceCompletingProcess, nil
}
