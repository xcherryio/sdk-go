package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
)

const (
	ApiPathAsyncStateWaitUntil = "/api/v1/xdb/worker/async-state/wait-until"
	ApiPathAsyncStateExecute   = "/api/v1/xdb/worker/async-state/execute"
)

// WorkerService is for worker to handle task requests from XDB server
// Typically put it behind a REST controller, using the above API paths
type WorkerService interface {
	HandleAsyncStateWaitUntil(ctx context.Context, request xdbapi.AsyncStateWaitUntilRequest) (*xdbapi.AsyncStateWaitUntilResponse, error)
	HandleAsyncStateExecute(ctx context.Context, request xdbapi.AsyncStateExecuteRequest) (*xdbapi.AsyncStateExecuteResponse, error)
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
