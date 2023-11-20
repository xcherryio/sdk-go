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

type StateFailureRecoveryTestExecuteFailedAtStartProcess struct {
	xc.ProcessDefaults
}

func (b StateFailureRecoveryTestExecuteFailedAtStartProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(
		&executeFailedAtStartInitState{},
		&executeFailedAtStartSkippedState{},
		&executeFailedAtStartRecoverState{})
}

type executeFailedAtStartInitState struct {
	xc.AsyncStateDefaults
}

func (d executeFailedAtStartInitState) GetStateOptions() *xc.AsyncStateOptions {
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

	stateOptions.SetFailureRecoveryOption(&executeFailedAtStartRecoverState{})

	return stateOptions
}

func (b executeFailedAtStartInitState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b executeFailedAtStartInitState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)
	return xc.SingleNextState(&executeFailedAtStartSkippedState{}, i+1), fmt.Errorf("error for test")
}

type executeFailedAtStartSkippedState struct {
	xc.AsyncStateDefaults
}

func (b executeFailedAtStartSkippedState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b executeFailedAtStartSkippedState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)

	return xc.SingleNextState(&executeFailedAtStartRecoverState{}, i+2), nil
}

type executeFailedAtStartRecoverState struct {
	xc.AsyncStateDefaults
}

func (b executeFailedAtStartRecoverState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b executeFailedAtStartRecoverState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	if ctx.GetRecoverFromStateApi() == nil || *(ctx.GetRecoverFromStateApi()) != xcapi.EXECUTE_API {
		panic("should recover from execute api")
	}

	if ctx.GetRecoverFromStateExecutionId() == nil || *(ctx.GetRecoverFromStateExecutionId()) != "failure_recovery.executeFailedAtStartInitState" {
		panic("should recover from state failure_recovery.executeFailedAtStartInitState")
	}

	var i int
	input.Get(&i)

	if i == 1 {
		return xc.GracefulCompletingProcess, nil
	}

	return xc.ForceFailProcess, nil
}

func TestStateFailureRecoveryTestExecuteFailedAtStartProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := StateFailureRecoveryTestExecuteFailedAtStartProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 1)
	require.NoError(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
