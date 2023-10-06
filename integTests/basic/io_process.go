package basic

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb"
)

type IOProcess struct {
	xdb.ProcessDefaults
}

func (b IOProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.WithStartingState(&state1{}, &state2{})
}

type state1 struct {
	xdb.AsyncStateDefaults
}

func (b state1) WaitUntil(ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication) (*xdb.CommandRequest, error) {
	return xdb.EmptyCommandRequest(), nil
}

func (b state1) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)
	return xdb.SingleNextState(state2{}, i+1), nil
}

type state2 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b state2) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)
	time.Sleep(time.Second * 1)
	return xdb.ForceCompletingProcess, nil
}

func TestStartIOProcess(t *testing.T, client xdb.Client) {
	prcId := "TestProceedOnStateStartFailWorkflow" + strconv.Itoa(int(time.Now().Unix()))
	prc := IOProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 123)
	assert.Nil(t, err)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdb.DefaultWorkerUrl, resp.GetWorkerUrl())
	assert.Equal(t, xdb.GetFinalProcessType(prc), resp.GetProcessType())
	assert.NotNil(t, resp.ProcessExecutionId)
	assert.Equal(t, xdbapi.RUNNING, resp.GetStatus())

	time.Sleep(time.Second * 3)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}

func TestProcessIdReusePolicyDisallowReuse(t *testing.T, client xdb.Client) {
	prcId := "TestProcessIdReuseDisallowReuse" + strconv.Itoa(int(time.Now().Unix()))
	prc := IOProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 123, nil)
	assert.Nil(t, err)

	_, err = client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.DISALLOW_REUSE.Ptr(),
	})
	assert.NotNil(t, err)
	assert.True(t, xdb.IsProcessAlreadyStartedError(err))

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())

	_, err = client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.DISALLOW_REUSE.Ptr(),
	})
	assert.NotNil(t, err)
}

func TestProcessIdReusePolicyAllowIfNoRunning(t *testing.T, client xdb.Client) {
	prcId := "TestProcessIdReuseAllowIfNoRunning" + strconv.Itoa(int(time.Now().Unix()))
	prc := IOProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 123, nil)
	assert.Nil(t, err)
	// immediate start with the same id is not allowed
	_, err = client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
	})
	assert.NotNil(t, err)

	// after the previous process with the same id is completed, the new process can be started
	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())

	_, err = client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
	})
	assert.Nil(t, err)
}

func TestProcessIdReusePolicyTerminateIfRunning(t *testing.T, client xdb.Client) {
	prcId := "TestProcessIdReuseTerminateIfRunning" + strconv.Itoa(int(time.Now().Unix()))
	prc := IOProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 123, nil)
	assert.Nil(t, err)
	// immediate start with the same id
	_, err = client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.TERMINATE_IF_RUNNING.Ptr(),
	})
	assert.Nil(t, err)
}

func TestProcessIdReusePolicyAllowIfPreviousExitAbnormally(t *testing.T, client xdb.Client) {
	// 1st case, if previous run finished normally, then the new run is not allowed
	prcId := "TestProcessIdReusePolicyAllowIfPreviousExitAbnormally" + strconv.Itoa(int(time.Now().Unix()))
	prc := IOProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 124, nil)
	assert.Nil(t, err)
	// immediate start with the same id
	_, err = client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.ALLOW_IF_PREVIOUS_EXIT_ABNORMALLY.Ptr(),
	})
	assert.NotNil(t, err)

	time.Sleep(time.Second * 5)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())

	_, err = client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.ALLOW_IF_PREVIOUS_EXIT_ABNORMALLY.Ptr(),
	})
	assert.NotNil(t, err)

	// 2nd case, if previous run finished abnormally, then the new run is allowed
	prcId = "TestProcessIdReusePolicyAllowIfPreviousExitAbnormally" + strconv.Itoa(int(time.Now().Unix()))
	prc = IOProcess{}
	_, err = client.StartProcess(context.Background(), prc, prcId, 124, nil)
	assert.Nil(t, err)
	err = client.StopProcess(context.Background(), prcId, xdbapi.FAIL)
	assert.Nil(t, err)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.FAILED, resp.GetStatus())
	_, err = client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.ALLOW_IF_PREVIOUS_EXIT_ABNORMALLY.Ptr(),
	})
	assert.Nil(t, err)
}
