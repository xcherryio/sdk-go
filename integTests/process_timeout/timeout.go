package process_timeout

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc"
	"github.com/xcherryio/sdk-go/xc/ptr"
)

type TimeoutProcess struct {
	xc.ProcessDefaults
}

func (b TimeoutProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(&timeoutState1{})
}

type timeoutState1 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b timeoutState1) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	time.Sleep(time.Second * 10)
	return xc.GracefulCompletingProcess, nil
}

func TestStartTimeoutProcessCase1(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := TimeoutProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		TimeoutSeconds: ptr.Any(int32(2)),
		IdReusePolicy:  xcapi.DISALLOW_REUSE.Ptr(),
	})
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
	assert.Equal(t, xcapi.TIMEOUT, resp.GetStatus())
}

func TestStartTimeoutProcessCase2(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := TimeoutProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		TimeoutSeconds: ptr.Any(int32(2)),
		IdReusePolicy:  xcapi.ALLOW_IF_NO_RUNNING.Ptr(),
	})
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
	assert.Equal(t, xcapi.TIMEOUT, resp.GetStatus())
}

func TestStartTimeoutProcessCase3(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := TimeoutProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		TimeoutSeconds: ptr.Any(int32(2)),
		IdReusePolicy:  xcapi.ALLOW_IF_PREVIOUS_EXIT_ABNORMALLY.Ptr(),
	})
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
	assert.Equal(t, xcapi.TIMEOUT, resp.GetStatus())
}

func TestStartTimeoutProcessCase4(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := TimeoutProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xc.ProcessStartOptions{
		TimeoutSeconds: ptr.Any(int32(2)),
		IdReusePolicy:  xcapi.TERMINATE_IF_RUNNING.Ptr(),
	})
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
	assert.Equal(t, xcapi.TIMEOUT, resp.GetStatus())
}
