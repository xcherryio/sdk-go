package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

type clientImpl struct {
	BasicClient
	registry Registry
}

func (c *clientImpl) GetBasicClient() BasicClient {
	return c.BasicClient
}

func (c *clientImpl) StartProcessWithOptions(
	ctx context.Context, definition Process, processId string, input interface{}, optionsOverride *ProcessOptions,
) (string, error) {
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
		unregOpt.StartStateOptions = fromStateToAsyncStateConfig(state)
	}

	options := prc.GetProcessOptions()
	if optionsOverride != nil {
		options = optionsOverride
	}
	if options != nil {
		unregOpt.ProcessIdReusePolicy = options.IdReusePolicy
		unregOpt.TimeoutSeconds = options.TimeoutSeconds
	}
	return c.BasicClient.StartProcess(ctx, prcType, startStateId, processId, input, unregOpt)
}

func (c *clientImpl) StartProcess(ctx context.Context, definition Process, processId string, input interface{}) (string, error) {
	return c.StartProcessWithOptions(ctx, definition, processId, input, nil)
}

func (c *clientImpl) StopProcess(ctx context.Context, processId string, stopType xdbapi.ProcessExecutionStopType) error {
	return c.BasicClient.StopProcess(ctx, processId, stopType)
}

func (c *clientImpl) DescribeCurrentProcessExecution(ctx context.Context, processId string) (*xdbapi.ProcessExecutionDescribeResponse, error) {
	return c.BasicClient.DescribeCurrentProcessExecution(ctx, processId)
}
