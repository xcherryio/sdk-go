package failure_recovery

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/integTests/common"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
)

type StateFailureRecoveryTestExecuteNoWaitUntilProcess struct {
	xdb.ProcessDefaults
}

func (b StateFailureRecoveryTestExecuteNoWaitUntilProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.WithStartingState(
		&executeNoWaitUntilInitState{},
		&executeNoWaitUntilFailState{},
		&executeNoWaitUntilRecoverState{})
}

type executeNoWaitUntilInitState struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

// TODO: investigate the issue of starting state options being applied to all states
// TODO: change the options to state2
func (d executeNoWaitUntilInitState) GetStateOptions() *xdb.AsyncStateOptions {
	stateOptions := &xdb.AsyncStateOptions{
		ExecuteTimeoutSeconds:   1,
		WaitUntilTimeoutSeconds: 1,
		WaitUntilRetryPolicy: &xdbapi.RetryPolicy{
			BackoffCoefficient:             ptr.Any(float32(1.0)),
			InitialIntervalSeconds:         ptr.Any(int32(1)),
			MaximumIntervalSeconds:         ptr.Any(int32(1)),
			MaximumAttemptsDurationSeconds: ptr.Any(int32(1)),
			MaximumAttempts:                ptr.Any(int32(1)),
		},
		ExecuteRetryPolicy: &xdbapi.RetryPolicy{
			BackoffCoefficient:             ptr.Any(float32(1.0)),
			InitialIntervalSeconds:         ptr.Any(int32(1)),
			MaximumIntervalSeconds:         ptr.Any(int32(1)),
			MaximumAttemptsDurationSeconds: ptr.Any(int32(1)),
			MaximumAttempts:                ptr.Any(int32(1)),
		},
	}

	stateOptions.SetFailureRecoveryOption(
		&executeNoWaitUntilRecoverState{},
		&xdb.AsyncStateOptions{})

	return stateOptions
}

func (b executeNoWaitUntilInitState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)
	return xdb.SingleNextState(&executeNoWaitUntilFailState{}, i+1), nil
}

type executeNoWaitUntilFailState struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b executeNoWaitUntilFailState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)

	return xdb.SingleNextState(&executeRecoverState{}, i+2), fmt.Errorf("error for test")
}

type executeNoWaitUntilRecoverState struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b executeNoWaitUntilRecoverState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	if ctx.GetRecoverFromStateApi() == nil || *(ctx.GetRecoverFromStateApi()) != xdbapi.EXECUTE_API {
		panic("should recover from execute api")
	}

	if ctx.GetRecoverFromStateExecutionId() == nil || *(ctx.GetRecoverFromStateExecutionId()) != "failure_recovery.executeNoWaitUntilFailState" {
		panic("should recover from state failure_recovery.executeNoWaitUntilFailState")
	}

	var i int
	input.Get(&i)

	if i == 2 {
		return xdb.GracefulCompletingProcess, nil
	}

	return xdb.ForceFailProcess, nil
}

func TestStateFailureRecoveryTestExecuteNoWaitUntilProcess(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := StateFailureRecoveryTestExecuteNoWaitUntilProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 1)
	require.NoError(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}
