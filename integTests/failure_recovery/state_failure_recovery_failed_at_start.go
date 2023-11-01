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

type StateFailureRecoveryTestExecuteFailedAtStartProcess struct {
	xdb.ProcessDefaults
}

func (b StateFailureRecoveryTestExecuteFailedAtStartProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(
		&executeFailedAtStartInitState{},
		&executeFailedAtStartSkippedState{},
		&executeFailedAtStartRecoverState{})
}

type executeFailedAtStartInitState struct {
	xdb.AsyncStateDefaults
}

func (d executeFailedAtStartInitState) GetStateOptions() *xdb.AsyncStateOptions {
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

	stateOptions.SetFailureRecoveryOption(&executeFailedAtStartRecoverState{})

	return stateOptions
}

func (b executeFailedAtStartInitState) WaitUntil(
	ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication,
) (*xdb.CommandRequest, error) {
	return xdb.EmptyCommandRequest(), nil
}

func (b executeFailedAtStartInitState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)
	return xdb.SingleNextState(&executeFailedAtStartSkippedState{}, i+1), fmt.Errorf("error for test")
}

type executeFailedAtStartSkippedState struct {
	xdb.AsyncStateDefaults
}

func (b executeFailedAtStartSkippedState) WaitUntil(
	ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication,
) (*xdb.CommandRequest, error) {
	return xdb.EmptyCommandRequest(), nil
}

func (b executeFailedAtStartSkippedState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)

	return xdb.SingleNextState(&executeFailedAtStartRecoverState{}, i+2), nil
}

type executeFailedAtStartRecoverState struct {
	xdb.AsyncStateDefaults
}

func (b executeFailedAtStartRecoverState) WaitUntil(
	ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication,
) (*xdb.CommandRequest, error) {
	return xdb.EmptyCommandRequest(), nil
}

func (b executeFailedAtStartRecoverState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	if ctx.GetRecoverFromStateApi() == nil || *(ctx.GetRecoverFromStateApi()) != xdbapi.EXECUTE_API {
		panic("should recover from execute api")
	}

	if ctx.GetRecoverFromStateExecutionId() == nil || *(ctx.GetRecoverFromStateExecutionId()) != "failure_recovery.executeFailedAtStartInitState" {
		panic("should recover from state failure_recovery.executeFailedAtStartInitState")
	}

	var i int
	input.Get(&i)

	if i == 1 {
		return xdb.GracefulCompletingProcess, nil
	}

	return xdb.ForceFailProcess, nil
}

func TestStateFailureRecoveryTestExecuteFailedAtStartProcess(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := StateFailureRecoveryTestExecuteFailedAtStartProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 1)
	require.NoError(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}
