package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

// Client is a full-featured client
type Client interface {
	GetBasicClient() BasicClient
	// StartProcess starts a process execution
	// definition is the definition of the process
	// processId is the required business identifier for the process execution (can be used with ProcessIdReusePolicy)
	// input the optional input for the startingState
	// return the processExecutionId
	StartProcess(ctx context.Context, definition Process, processId string, input interface{}) (string, error)
	// StartProcessWithOptions starts a process execution with options, which will override the options defined in process definition
	StartProcessWithOptions(ctx context.Context, definition Process, processId string, input interface{}, options *ProcessStartOptions) (string, error)
	// StopProcess stops a process execution
	// processId is the required business identifier for the process execution
	StopProcess(ctx context.Context, processId string, stopType xdbapi.ProcessExecutionStopType) error
	// DescribeCurrentProcessExecution returns a process execution info
	// processId is the required business identifier for the process execution
	DescribeCurrentProcessExecution(ctx context.Context, processId string) (*xdbapi.ProcessExecutionDescribeResponse, error)
}

// BasicClient is a base client without process registry
// It's the internal implementation of Client.
// But it can be used directly if there is good reason -- let you invoke the APIs to xdb server without much type validation checks(process type, queue names, etc).
type BasicClient interface {
	// StartProcess starts a process execution
	// processType is the process type
	// startStateId is the stateId of the startingState
	// processId is the required business identifier for the process execution(can be used with ProcessIdReusePolicy
	// input the optional input for the startingState
	// options is optional includes like ProcessIdReusePolicy.
	// return the processExecutionId
	StartProcess(ctx context.Context, processType string, startStateId, processId string, input interface{}, options *BasicClientProcessOptions) (string, error)
	// StopProcess stops a process execution
	// processId is the required business identifier for the process execution
	StopProcess(ctx context.Context, processId string, stopType xdbapi.ProcessExecutionStopType) error
	// DescribeCurrentProcessExecution returns a process execution info
	// processId is the required business identifier for the process execution
	DescribeCurrentProcessExecution(ctx context.Context, processId string) (*xdbapi.ProcessExecutionDescribeResponse, error)
}

// NewClient returns a Client
func NewClient(registry Registry, options *ClientOptions) Client {
	if registry == nil {
		panic("A registry is required")
	}
	if options == nil {
		options = GetLocalDefaultClientOptions()
	}
	return &clientImpl{
		BasicClient:   NewBasicClient(*options),
		clientOptions: *options,
		registry:      registry,
	}
}

// NewBasicClient returns a BasicClient
func NewBasicClient(options ClientOptions) BasicClient {

	cfg := &xdbapi.Configuration{
		Servers: []xdbapi.ServerConfiguration{
			{
				URL: options.ServerUrl,
			},
		},
	}
	if options.EnabledDebugLogging {
		cfg.Debug = true
	}

	apiClient := xdbapi.NewAPIClient(cfg)

	return &basicClientImpl{
		options:   options,
		apiClient: apiClient,
	}
}
