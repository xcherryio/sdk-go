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

type StateFailureRecoveryTestExecuteNoWaitUntilProcess struct {
	xc.ProcessDefaults
}

func (b StateFailureRecoveryTestExecuteNoWaitUntilProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(
		&executeNoWaitUntilInitState{},
		&executeNoWaitUntilFailState{},
		&executeNoWaitUntilRecoverState{})
}

type executeNoWaitUntilInitState struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b executeNoWaitUntilInitState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)
	return xc.SingleNextState(&executeNoWaitUntilFailState{}, i+1), nil
}

type executeNoWaitUntilFailState struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (d executeNoWaitUntilFailState) GetStateOptions() *xc.AsyncStateOptions {
	stateOptions := &xc.AsyncStateOptions{
		ExecuteTimeoutSeconds: 1,
		ExecuteRetryPolicy: &xcapi.RetryPolicy{
			BackoffCoefficient:             ptr.Any(float32(1.0)),
			InitialIntervalSeconds:         ptr.Any(int32(1)),
			MaximumIntervalSeconds:         ptr.Any(int32(1)),
			MaximumAttemptsDurationSeconds: ptr.Any(int32(1)),
			MaximumAttempts:                ptr.Any(int32(1)),
		},
	}

	stateOptions.SetFailureRecoveryOption(&executeNoWaitUntilRecoverState{})

	return stateOptions
}

func (b executeNoWaitUntilFailState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)

	return xc.SingleNextState(&executeNoWaitUntilRecoverState{}, i+2), fmt.Errorf("error for test")
}

type executeNoWaitUntilRecoverState struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b executeNoWaitUntilRecoverState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	if ctx.GetRecoverFromStateApi() == nil || *(ctx.GetRecoverFromStateApi()) != xcapi.EXECUTE_API {
		panic("should recover from execute api")
	}

	if ctx.GetRecoverFromStateExecutionId() == nil || *(ctx.GetRecoverFromStateExecutionId()) != "failure_recovery.executeNoWaitUntilFailState" {
		panic("should recover from state failure_recovery.executeNoWaitUntilFailState")
	}

	var i int
	input.Get(&i)

	if i == 2 {
		return xc.GracefulCompletingProcess, nil
	}

	return xc.ForceFailProcess, nil
}

func TestStateFailureRecoveryTestExecuteNoWaitUntilProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := StateFailureRecoveryTestExecuteNoWaitUntilProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 1)
	require.NoError(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
