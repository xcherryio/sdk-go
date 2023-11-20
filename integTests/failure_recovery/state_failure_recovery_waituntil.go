package failure_recovery

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc"
	"github.com/xcherryio/sdk-go/xc/ptr"
)

type StateFailureRecoveryTestWaitUntilProcess struct {
	xc.ProcessDefaults
}

func (b StateFailureRecoveryTestWaitUntilProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(
		&waitUntilInitState{},
		&waitUntilFailedState{},
		&waitUntilRecoverState{})
}

type waitUntilInitState struct {
	xc.AsyncStateDefaults
}

func (b waitUntilInitState) WaitUntil(ctx xc.Context, input xc.Object, communication xc.Communication) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b waitUntilInitState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)
	return xc.SingleNextState(waitUntilFailedState{}, i+1), nil
}

type waitUntilFailedState struct {
	xc.AsyncStateDefaults
}

func (d waitUntilFailedState) GetStateOptions() *xc.AsyncStateOptions {
	stateOptions := &xc.AsyncStateOptions{
		ExecuteTimeoutSeconds:   1,
		WaitUntilTimeoutSeconds: 1,
		WaitUntilRetryPolicy: &xcapi.RetryPolicy{
			BackoffCoefficient:             ptr.Any(float32(1.0)),
			InitialIntervalSeconds:         ptr.Any(int32(1)),
			MaximumIntervalSeconds:         ptr.Any(int32(1)),
			MaximumAttemptsDurationSeconds: ptr.Any(int32(1)),
			MaximumAttempts:                ptr.Any(int32(1)),
		},
		ExecuteRetryPolicy: &xcapi.RetryPolicy{
			BackoffCoefficient:             ptr.Any(float32(1.0)),
			InitialIntervalSeconds:         ptr.Any(int32(1)),
			MaximumIntervalSeconds:         ptr.Any(int32(1)),
			MaximumAttemptsDurationSeconds: ptr.Any(int32(1)),
			MaximumAttempts:                ptr.Any(int32(1)),
		},
	}

	stateOptions.SetFailureRecoveryOption(&waitUntilRecoverState{})

	return stateOptions
}

func (b waitUntilFailedState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return nil, fmt.Errorf("error for testing")
}

func (b waitUntilFailedState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)

	return xc.SingleNextState(&waitUntilRecoverState{}, i+2), nil
}

type waitUntilRecoverState struct {
	xc.AsyncStateDefaults
}

func (b waitUntilRecoverState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b waitUntilRecoverState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	if ctx.GetRecoverFromStateApi() == nil || *(ctx.GetRecoverFromStateApi()) != xcapi.WAIT_UNTIL_API {
		panic("should recover from execute api")
	}

	if ctx.GetRecoverFromStateExecutionId() == nil || *(ctx.GetRecoverFromStateExecutionId()) != "failure_recovery.waitUntilFailedState" {
		panic("should recover from state failure_recovery.waitUntilFailedState")
	}

	var i int
	input.Get(&i)

	if i == 2 {
		return xc.GracefulCompletingProcess, nil
	}

	return xc.ForceFailProcess, nil
}

func TestStateFailureRecoveryTestWaitUntilProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := StateFailureRecoveryTestWaitUntilProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 1)
	require.NoError(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
