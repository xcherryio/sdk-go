package basic

import (
	"context"
	"testing"
	"time"

	"github.com/xcherryio/sdk-go/integTests/common"

	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/xc"
)

type IOProcess struct {
	xc.ProcessDefaults
}

func (b IOProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(&state1{}, &state2{})
}

type state1 struct {
	xc.AsyncStateDefaults
}

func (b state1) WaitUntil(ctx xc.Context, input xc.Object, communication xc.Communication) (*xc.CommandRequest, error) {
	return xc.EmptyCommandRequest(), nil
}

func (b state1) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)
	return xc.SingleNextState(state2{}, i+1), nil
}

type state2 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b state2) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)
	time.Sleep(time.Second * 1)
	return xc.ForceCompletingProcess, nil
}

func TestStartIOProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := IOProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 123)
	assert.Nil(t, err)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xc.DefaultWorkerUrl, resp.GetWorkerUrl())
	assert.Equal(t, xc.GetFinalProcessType(prc), resp.GetProcessType())
	assert.NotNil(t, resp.ProcessExecutionId)
	assert.Equal(t, xcapi.RUNNING, resp.GetStatus())

	time.Sleep(time.Second * 3)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}

func TestProcessIdReusePolicyDisallowReuse(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := IOProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, nil)
	assert.Nil(t, err)

	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		IdReusePolicy: xcapi.DISALLOW_REUSE.Ptr(),
	})
	assert.NotNil(t, err)
	assert.True(t, xc.IsProcessAlreadyStartedError(err))

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())

	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		IdReusePolicy: xcapi.DISALLOW_REUSE.Ptr(),
	})
	assert.NotNil(t, err)
}

func TestProcessIdReusePolicyAllowIfNoRunning(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := IOProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, nil)
	assert.Nil(t, err)
	// immediate start with the same id is not allowed
	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		IdReusePolicy: xcapi.ALLOW_IF_NO_RUNNING.Ptr(),
	})
	assert.NotNil(t, err)

	// after the previous process with the same id is completed, the new process can be started
	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())

	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		IdReusePolicy: xcapi.ALLOW_IF_NO_RUNNING.Ptr(),
	})
	assert.Nil(t, err)
}

func TestProcessIdReusePolicyTerminateIfRunning(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := IOProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, nil)
	assert.Nil(t, err)
	// immediate start with the same id
	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		IdReusePolicy: xcapi.TERMINATE_IF_RUNNING.Ptr(),
	})
	assert.Nil(t, err)
}

func TestProcessIdReusePolicyAllowIfPreviousExitAbnormallyCase1(t *testing.T, client xc.Client) {
	// 1st case, if previous run finished normally, then the new run is not allowed
	prcId := common.GenerateProcessId()
	prc := IOProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 124, nil)
	assert.Nil(t, err)
	// immediate start with the same id
	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		IdReusePolicy: xcapi.ALLOW_IF_PREVIOUS_EXIT_ABNORMALLY.Ptr(),
	})
	assert.NotNil(t, err)

	time.Sleep(time.Second * 5)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())

	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		IdReusePolicy: xcapi.ALLOW_IF_PREVIOUS_EXIT_ABNORMALLY.Ptr(),
	})
	assert.NotNil(t, err)
}

func TestProcessIdReusePolicyAllowIfPreviousExitAbnormallyCase2(t *testing.T, client xc.Client) {
	// 2nd case, if previous run finished abnormally, then the new run is allowed
	prcId := common.GenerateProcessId()
	prc := IOProcess{}
	runId1, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 124, nil)
	assert.Nil(t, err)
	err = client.StopProcess(context.Background(), prcId, xcapi.FAIL)
	assert.Nil(t, err)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.FAILED, resp.GetStatus())
	runId2, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		IdReusePolicy: xcapi.ALLOW_IF_PREVIOUS_EXIT_ABNORMALLY.Ptr(),
	})
	assert.Nil(t, err)
	assert.NotEqual(t, runId1, runId2)

	time.Sleep(time.Second * 5)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
