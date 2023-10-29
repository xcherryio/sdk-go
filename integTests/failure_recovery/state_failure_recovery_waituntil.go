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

type StateFailureRecoveryTestWaitUntilProcess struct {
	xdb.ProcessDefaults
}

func (b StateFailureRecoveryTestWaitUntilProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.WithStartingState(
		&stateFailureRecoveryTestWaitUntilState1{},
		&stateFailureRecoveryTestWaitUntilState2{},
		&stateFailureRecoveryTestWaitUntilState3{})
}

type stateFailureRecoveryTestWaitUntilState1 struct {
	xdb.AsyncStateDefaults
}

func (d stateFailureRecoveryTestWaitUntilState1) GetStateOptions() *xdb.AsyncStateOptions {
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

	stateOptions.SetFailureRecoveryOption(&stateFailureRecoveryTestWaitUntilState3{}, &xdb.AsyncStateOptions{})

	return stateOptions
}

func (b stateFailureRecoveryTestWaitUntilState1) WaitUntil(ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication) (*xdb.CommandRequest, error) {
	return xdb.EmptyCommandRequest(), nil
}

func (b stateFailureRecoveryTestWaitUntilState1) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)
	return xdb.SingleNextState(stateFailureRecoveryTestWaitUntilState2{}, i+1), nil
}

type stateFailureRecoveryTestWaitUntilState2 struct {
	xdb.AsyncStateDefaults
}

func (b stateFailureRecoveryTestWaitUntilState2) WaitUntil(ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication) (*xdb.CommandRequest, error) {
	return nil, fmt.Errorf("error for testing")
}

func (b stateFailureRecoveryTestWaitUntilState2) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)

	return xdb.SingleNextState(&stateFailureRecoveryTestWaitUntilState3{}, i+2), nil
}

type stateFailureRecoveryTestWaitUntilState3 struct {
	xdb.AsyncStateDefaults
}

func (b stateFailureRecoveryTestWaitUntilState3) WaitUntil(ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication) (*xdb.CommandRequest, error) {
	return xdb.EmptyCommandRequest(), nil
}

func (b stateFailureRecoveryTestWaitUntilState3) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	if ctx.GetRecoverFromStateApi() == nil || *(ctx.GetRecoverFromStateApi()) != xdbapi.WAIT_UNTIL_API {
		panic("should recover from execute api")
	}

	if ctx.GetRecoverFromStateExecutionId() == nil || *(ctx.GetRecoverFromStateExecutionId()) != "failure_recovery.stateFailureRecoveryTestWaitUntilState2" {
		panic("should recover from state failure_recovery.stateFailureRecoveryTestWaitUntilState2")
	}

	var i int
	input.Get(&i)

	if i == 2 {
		return xdb.GracefulCompletingProcess, nil
	}

	return xdb.ForceFailProcess, nil
}

func TestStateFailureRecoveryTestWaitUntilProcess(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := StateFailureRecoveryTestWaitUntilProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 1)
	require.NoError(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}
