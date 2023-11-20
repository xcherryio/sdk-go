package xc

import (
	"context"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/xc/ptr"
)

const (
	ApiPathAsyncStateWaitUntil = "/api/v1/xcherry/worker/async-state/wait-until"
	ApiPathAsyncStateExecute   = "/api/v1/xcherry/worker/async-state/execute"
)

// WorkerService is for worker to handle task requests from xCherry server
// Typically put it behind a REST controller, using the above API paths
type WorkerService interface {
	HandleAsyncStateWaitUntil(ctx context.Context, request xcapi.AsyncStateWaitUntilRequest) (*xcapi.AsyncStateWaitUntilResponse, error)
	HandleAsyncStateExecute(ctx context.Context, request xcapi.AsyncStateExecuteRequest) (*xcapi.AsyncStateExecuteResponse, error)
}

func NewWorkerService(registry Registry, options *WorkerOptions) WorkerService {
	if options == nil {
		options = ptr.Any(GetDefaultWorkerOptions())
	}
	return &workerServiceImpl{
		registry: registry,
		options:  *options,
	}
}
