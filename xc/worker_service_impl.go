package xc

import (
	"context"

	"github.com/xcherryio/apis/goapi/xcapi"
)

type workerServiceImpl struct {
	registry Registry
	options  WorkerOptions
}

func (w *workerServiceImpl) HandleAsyncStateWaitUntil(
	ctx context.Context, request xcapi.AsyncStateWaitUntilRequest,
) (resp *xcapi.AsyncStateWaitUntilResponse, retErr error) {
	defer func() { captureStateExecutionError(recover(), &retErr) }()

	prcType := request.GetProcessType()
	stateDef := w.registry.getProcessState(prcType, request.GetStateId())
	input := NewObject(request.StateInput, w.options.ObjectEncoder)
	reqContext := request.GetContext()
	wfCtx := newContext(reqContext)

	comm := NewCommunication(w.options.ObjectEncoder)
	commandRequest, err := stateDef.WaitUntil(wfCtx, input, comm)

	if err != nil {
		return nil, err
	}

	idlCommandRequest, err := toApiCommandRequest(commandRequest)
	if err != nil {
		return nil, err
	}
	resp = &xcapi.AsyncStateWaitUntilResponse{
		CommandRequest:      *idlCommandRequest,
		PublishToLocalQueue: comm.GetLocalQueueMessagesToPublish(),
	}

	return resp, nil
}

func (w *workerServiceImpl) HandleAsyncStateExecute(
	ctx context.Context, request xcapi.AsyncStateExecuteRequest,
) (resp *xcapi.AsyncStateExecuteResponse, retErr error) {
	defer func() { captureStateExecutionError(recover(), &retErr) }()

	prcType := request.GetProcessType()
	stateDef := w.registry.getProcessState(prcType, request.GetStateId())
	input := NewObject(request.StateInput, w.options.ObjectEncoder)
	reqContext := request.GetContext()
	wfCtx := newContext(reqContext)

	commandResults, err := fromApiCommandResults(request.CommandResults, w.options.ObjectEncoder)
	if err != nil {
		return nil, err
	}

	pers := w.createPersistenceImpl(prcType, request.LoadedGlobalAttributes, request.LoadedLocalAttributes)

	comm := NewCommunication(w.options.ObjectEncoder)
	decision, err := stateDef.Execute(wfCtx, input, commandResults, pers, comm)

	if err != nil {
		return nil, err
	}
	idlDecision, err := toApiDecision(decision, prcType, w.registry, w.options.ObjectEncoder)
	if err != nil {
		return nil, err
	}
	resp = &xcapi.AsyncStateExecuteResponse{
		StateDecision:       *idlDecision,
		PublishToLocalQueue: comm.GetLocalQueueMessagesToPublish(),
	}
	if len(pers.getGlobalAttributesToUpdate()) > 0 {
		resp.WriteToGlobalAttributes = pers.getGlobalAttributesToUpdate()
	}
	if len(pers.getLocalAttributesToUpdate()) > 0 {
		resp.WriteToLocalAttributes = pers.getLocalAttributesToUpdate()
	}
	return resp, nil
}

func (w *workerServiceImpl) createPersistenceImpl(
	prcType string, currGlobalAttrs *xcapi.LoadGlobalAttributeResponse,
	currLocalAttrs *xcapi.LoadLocalAttributesResponse,
) Persistence {
	gloAttrDefs := w.registry.getGlobalAttributeKeyToDefs(prcType)
	gloTblColToKey := w.registry.getGlobalAttributeTableColumnToKey(prcType)
	localAttributeKeys := w.registry.getLocalAttributeKeys(prcType)
	return NewPersistenceImpl(
		w.options.DBConverter,
		gloAttrDefs,
		gloTblColToKey,
		currGlobalAttrs,
		localAttributeKeys,
		currLocalAttrs)
}
