package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
)

type clientImpl struct {
	BasicClient
	registry Registry
}

func (c *clientImpl) GetBasicClient() BasicClient {
	return c.BasicClient
}

func (c *clientImpl) StartProcess(ctx context.Context, definition Process, processId string, input interface{}, options *ProcessOptions) (string, error) {
	prcType := GetFinalProcessType(definition)
	prc := c.registry.getProcess(prcType)
	if prc == nil {
		return "", NewInvalidArgumentError("Process is not registered")
	}

	state := c.registry.getProcessStartingState(prcType)

	unregOpt := &BasicClientProcessOptions{}

	startStateId := ""
	if state != nil {
		startStateId = GetFinalStateId(state)
		startStateOpt := &xdbapi.AsyncStateConfig{}
		if ShouldSkipWaitUntilAPI(state) {
			startStateOpt.SkipWaitUntil = ptr.Any(true)
		}
		unregOpt.StartStateOptions = startStateOpt
	}

	if options != nil {
		unregOpt.ProcessIdReusePolicy = options.IdReusePolicy
		unregOpt.TimeoutSeconds = options.TimeoutSeconds
	}
	return c.BasicClient.StartProcess(ctx, prcType, startStateId, processId, input, unregOpt)
}

func (c *clientImpl) StopProcess(ctx context.Context, processId string, stopType xdbapi.ProcessExecutionStopType) error {
	return c.BasicClient.StopProcess(ctx, processId, stopType)
}

func (c *clientImpl) DescribeCurrentProcessExecution(ctx context.Context, processId string) (*xdbapi.ProcessExecutionDescribeResponse, error) {
	return c.BasicClient.DescribeCurrentProcessExecution(ctx, processId)
}
