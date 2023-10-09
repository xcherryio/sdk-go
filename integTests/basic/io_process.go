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
	return nil, nil
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
	return xdb.ForceCompletingProcess, nil
}

func TestStartIOProcess(t *testing.T, client xdb.Client) {
	prcId := "TestProceedOnStateStartFailWorkflow" + strconv.Itoa(int(time.Now().Unix()))
	prc := IOProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, 123, &xdb.ProcessOptions{
		IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
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
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())

}
