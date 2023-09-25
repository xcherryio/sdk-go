package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

// Client is a full-featured client
type Client interface {
	// StartProcess starts a process execution
	// definition is the definition of the process
	// processId is the required business identifier for the process execution(can be used with ProcessIdReusePolicy
	// input the optional input for the startingState
	// options is optional includes like ProcessIdReusePolicy.
	// return the processExecutionId
	StartProcess(ctx context.Context, definition Process, processId string, input interface{}, options *ProcessOptions) (string, error)
}

// UnregisteredClient is a client without process registry
// It's the internal implementation of Client.
// But it can be used directly if there is good reason -- let you invoke the APIs to xdb server without much type validation checks(process type, queue names, etc).
type UnregisteredClient interface {
	// StartProcess starts a process execution
	// processType is the process type
	// startStateId is the stateId of the startingState
	// processId is the required business identifier for the process execution(can be used with ProcessIdReusePolicy
	// input the optional input for the startingState
	// options is optional includes like ProcessIdReusePolicy.
	// return the processExecutionId
	StartProcess(ctx context.Context, processType string, startStateId, processId string, input interface{}, options *UnregisteredProcessOptions) (string, error)
}

// NewUnregisteredClient returns a UnregisteredClient
func NewUnregisteredClient(options *ClientOptions) UnregisteredClient {
	if options == nil {
		options = GetLocalDefaultClientOptions()
	}

	apiClient := xdbapi.NewAPIClient(&xdbapi.Configuration{
		Servers: []xdbapi.ServerConfiguration{
			{
				URL: options.ServerUrl,
			},
		},
	})

	return &unregisteredClientImpl{
		options:   options,
		apiClient: apiClient,
	}
}

// NewClient returns a Client
func NewClient(registry Registry, options *ClientOptions) Client {
	if registry == nil {
		panic("A registry is required")
	}
	return &clientImpl{
		UnregisteredClient: NewUnregisteredClient(options),
		registry:           registry,
		options:            options,
	}
}
