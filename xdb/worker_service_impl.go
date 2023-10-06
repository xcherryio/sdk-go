package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

type workerServiceImpl struct {
	registry Registry
	options  WorkerOptions
}

func (w *workerServiceImpl) HandleAsyncStateWaitUntil(ctx context.Context, request xdbapi.AsyncStateWaitUntilRequest) (resp *xdbapi.AsyncStateWaitUntilResponse, retErr error) {
	defer func() { captureStateExecutionError(recover(), &retErr) }()

	prcType := request.GetProcessType()
	stateDef := w.registry.getProcessState(prcType, request.GetStateId())
	input := NewObject(request.StateInput, w.options.ObjectEncoder)
	reqContext := request.GetContext()
	wfCtx := newXdbContext(reqContext)

	var comm Communication // TODO
	commandRequest, err := stateDef.WaitUntil(wfCtx, input, comm)
	if err != nil {
		return nil, err
	}

	idlCommandRequest, err := toApiCommandRequest(commandRequest)
	if err != nil {
		return nil, err
	}
	resp = &xdbapi.AsyncStateWaitUntilResponse{
		CommandRequest: *idlCommandRequest,
	}

	return resp, nil
}

func (w *workerServiceImpl) HandleAsyncStateExecute(ctx context.Context, request xdbapi.AsyncStateExecuteRequest) (resp *xdbapi.AsyncStateExecuteResponse, retErr error) {
	defer func() { captureStateExecutionError(recover(), &retErr) }()

	prcType := request.GetProcessType()
	stateDef := w.registry.getProcessState(prcType, request.GetStateId())
	input := NewObject(request.StateInput, w.options.ObjectEncoder)
	reqContext := request.GetContext()
	wfCtx := newXdbContext(reqContext)

	commandResults, err := fromApiCommandResults(request.CommandResults, w.options.ObjectEncoder)
	if err != nil {
		return nil, err
	}
	var pers Persistence   // TODO
	var comm Communication // TODO
	decision, err := stateDef.Execute(wfCtx, input, commandResults, pers, comm)
	if err != nil {
		return nil, err
	}
	idlDecision, err := toIdlDecision(decision, prcType, w.registry, w.options.ObjectEncoder)
	if err != nil {
		return nil, err
	}
	resp = &xdbapi.AsyncStateExecuteResponse{
		StateDecision: *idlDecision,
	}
	return resp, nil
}
