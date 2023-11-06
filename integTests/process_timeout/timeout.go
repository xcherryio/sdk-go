package process_timeout

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/integTests/common"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
)

type TimeoutProcess struct {
	xdb.ProcessDefaults
}

func (b TimeoutProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(&timeoutState1{})
}

type timeoutState1 struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b timeoutState1) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	time.Sleep(time.Second * 10)
	return xdb.GracefulCompletingProcess, nil
}

func TestStartTimeoutProcessCase1(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := TimeoutProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xdb.ProcessStartOptions{
		TimeoutSeconds: ptr.Any(int32(2)),
		IdReusePolicy:  xdbapi.DISALLOW_REUSE.Ptr(),
	})
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
	assert.Equal(t, xdbapi.TIMEOUT, resp.GetStatus())
}

func TestStartTimeoutProcessCase2(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := TimeoutProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xdb.ProcessStartOptions{
		TimeoutSeconds: ptr.Any(int32(2)),
		IdReusePolicy:  xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
	})
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
	assert.Equal(t, xdbapi.TIMEOUT, resp.GetStatus())
}

func TestStartTimeoutProcessCase3(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := TimeoutProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xdb.ProcessStartOptions{
		TimeoutSeconds: ptr.Any(int32(2)),
		IdReusePolicy:  xdbapi.ALLOW_IF_PREVIOUS_EXIT_ABNORMALLY.Ptr(),
	})
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
	assert.Equal(t, xdbapi.TIMEOUT, resp.GetStatus())
}

func TestStartTimeoutProcessCase4(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := TimeoutProcess{}
	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, 123, &xdb.ProcessStartOptions{
		TimeoutSeconds: ptr.Any(int32(2)),
		IdReusePolicy:  xdbapi.TERMINATE_IF_RUNNING.Ptr(),
	})
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
	assert.Equal(t, xdbapi.TIMEOUT, resp.GetStatus())
}
