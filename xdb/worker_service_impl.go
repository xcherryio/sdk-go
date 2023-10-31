package xdb

import (
	"context"

	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

type workerServiceImpl struct {
	registry Registry
	options  WorkerOptions
}

func (w *workerServiceImpl) HandleAsyncStateWaitUntil(
	ctx context.Context, request xdbapi.AsyncStateWaitUntilRequest,
) (resp *xdbapi.AsyncStateWaitUntilResponse, retErr error) {
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

func (w *workerServiceImpl) HandleAsyncStateExecute(
	ctx context.Context, request xdbapi.AsyncStateExecuteRequest,
) (resp *xdbapi.AsyncStateExecuteResponse, retErr error) {
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

	pers := w.createPersistenceImpl(prcType, request.LoadedGlobalAttributes)

	var comm Communication // TODO
	decision, err := stateDef.Execute(wfCtx, input, commandResults, pers, comm)

	if err != nil {
		return nil, err
	}
	idlDecision, err := toApiDecision(decision, prcType, w.registry, w.options.ObjectEncoder)
	if err != nil {
		return nil, err
	}
	resp = &xdbapi.AsyncStateExecuteResponse{
		StateDecision: *idlDecision,
	}
	if len(pers.getGlobalAttributesToUpdate()) > 0 {
		resp.UpsertGlobalAttributes = pers.getGlobalAttributesToUpdate()
	}
	return resp, nil
}

func (w *workerServiceImpl) createPersistenceImpl(
	prcType string, attributesResp *xdbapi.LoadGlobalAttributeResponse,
) Persistence {
	persSchema := w.registry.getPersistenceSchema(prcType)
	gloAttrDefs := w.registry.getGlobalAttributeKeyToDefs(prcType)
	gloTblColToKey := w.registry.getGlobalAttributeTableColumnToKey(prcType)
	var currGloAttrs []xdbapi.GlobalAttributeValue
	if attributesResp != nil {
		currGloAttrs = attributesResp.GetAttributes()
	}
	defaultTable := ""
	if persSchema.GlobalAttributeSchema != nil {
		defaultTable = persSchema.GlobalAttributeSchema.DefaultTableName
	}
	return NewPersistenceImpl(w.options.DBConverter, defaultTable, gloAttrDefs, gloTblColToKey, currGloAttrs)
}
