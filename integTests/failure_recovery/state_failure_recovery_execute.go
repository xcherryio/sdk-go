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

type StateFailureRecoveryTestExecuteProcess struct {
	xc.ProcessDefaults
}

func (b StateFailureRecoveryTestExecuteProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(
		&executeInitState{},
		&executeFailState{},
		&executeRecoverState{})
}

type executeInitState struct {
	xc.AsyncStateDefaults
}

func (b executeInitState) WaitUntil(ctx xc.Context, input xc.Object, communication xc.Communication) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b executeInitState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)
	return xc.SingleNextState(&executeFailState{}, i+1), nil
}

type executeFailState struct {
	xc.AsyncStateDefaults
}

func (d executeFailState) GetStateOptions() *xc.AsyncStateOptions {
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

	stateOptions.SetFailureRecoveryOption(&executeRecoverState{})

	return stateOptions
}

func (b executeFailState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b executeFailState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)

	return xc.SingleNextState(&executeRecoverState{}, i+2), fmt.Errorf("error for test")
}

type executeRecoverState struct {
	xc.AsyncStateDefaults
}

func (b executeRecoverState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b executeRecoverState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	if ctx.GetRecoverFromStateApi() == nil || *(ctx.GetRecoverFromStateApi()) != xcapi.EXECUTE_API {
		panic("should recover from execute api")
	}

	if ctx.GetRecoverFromStateExecutionId() == nil || *(ctx.GetRecoverFromStateExecutionId()) != "failure_recovery.executeFailState" {
		panic("should recover from state failure_recovery.executeFailState")
	}

	var i int
	input.Get(&i)

	if i == 2 {
		return xc.GracefulCompletingProcess, nil
	}

	return xc.ForceFailProcess, nil
}

func TestStateFailureRecoveryTestExecuteProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := StateFailureRecoveryTestExecuteProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 1)
	require.NoError(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
