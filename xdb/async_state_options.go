package xdb

import (
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
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
	// FailureRecoveryOptions is information needed for failure recovery
	// Default: FAIL_PROCESS_ON_STATE_FAILURE
	FailureRecoveryOptions *xdbapi.StateFailureRecoveryOptions
}

func (o *AsyncStateOptions) SetFailureRecoveryOption(
	destState AsyncState, destStateOptions *AsyncStateOptions) *AsyncStateOptions {
	if destState == nil {
		return o
	}

	o.FailureRecoveryOptions = &xdbapi.StateFailureRecoveryOptions{
		Policy:                         xdbapi.PROCEED_TO_CONFIGURED_STATE,
		StateFailureProceedStateId:     ptr.Any(GetFinalStateId(destState)),
		StateFailureProceedStateConfig: fromAsyncStateOptionsToAsyncStateConfg(destStateOptions),
	}

	return o
}
