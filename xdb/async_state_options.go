package xdb

import (
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

type AsyncStateOptions struct {
	// StateId is the unique identifier of the state.
	// It is being used for WorkerService to choose the right AsyncState to execute Start/Execute APIs
	// Default: the pkgName.structName of the state struct, see GetFinalStateId() for details when set as empty string
	StateId string
	// WaitUntilTimeoutSeconds is the timeout for the waitUntil API call.
	// Default: 10 seconds(configurable in server) when set as 0
	// It will be capped to 60 seconds by server (configurable in server)
	WaitUntilTimeoutSeconds int32
	// ExecuteTimeoutSeconds is the timeout for the execute API call.
	// Default: 10 seconds(configurable in server) when set as 0
	// It will be capped to 60 seconds by server (configurable in server)
	ExecuteTimeoutSeconds int32
	// WaitUntilRetryPolicy is the retry policy for the waitUntil API call.
	// Default: infinite retry with 1 second initial interval, 120 seconds max interval, and 2 backoff factor,
	// when set as nil
	WaitUntilRetryPolicy *xdbapi.RetryPolicy
	// ExecuteRetryPolicy is the retry policy for the execute API call.
	// Default: infinite retry with 1 second initial interval, 120 seconds max interval, and 2 backoff factor,
	// when set as nil
	ExecuteRetryPolicy *xdbapi.RetryPolicy
	// FailureRecoveryState is the state to recover after current state execution fails
	// Default: no recovery when set as nil
	FailureRecoveryState AsyncState
	// PersistencePolicyName is the name of loading policy for persistence if not using default policy
	PersistencePolicyName *string
}

func (o *AsyncStateOptions) SetFailureRecoveryOption(destState AsyncState) *AsyncStateOptions {
	if destState == nil {
		panic("destState is nil")
	}

	o.FailureRecoveryState = destState

	return o
}
