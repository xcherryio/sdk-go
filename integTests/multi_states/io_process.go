package multi_states

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"strconv"
	"testing"
	"time"
)

const INPUT = 1

type MultiStatesProcess struct {
	xdb.ProcessDefaults
}

func (b MultiStatesProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.WithStartingState(&state1{}, &state2{}, &state3{})
}

type state1 struct {
	xdb.AsyncStateDefaults
}

func (b state1) WaitUntil(ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication) (*xdb.CommandRequest, error) {
	var i int
	input.Get(&i)

	if i != INPUT {
		panic("state1 WaitUntil: input is not expected. Expected: " + fmt.Sprint(INPUT) + ", actual: " + fmt.Sprint(i))
	}

	return nil, nil
}

func (b state1) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)

	if i != INPUT {
		panic("state1 Execute: input is not expected. Expected: " + fmt.Sprint(INPUT) + ", actual: " + fmt.Sprint(i))
	}

	return xdb.MultiNextStatesWithInput(xdb.NewStateMovement(state2{}, i+2), xdb.NewStateMovement(state3{}, i+3)), nil
}

type state2 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b state2) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)

	if i != INPUT+2 {
		panic("state2 Execute: input is not expected. Expected: " + fmt.Sprint(INPUT+2) + ", actual: " + fmt.Sprint(i))
	}

	return xdb.DeadEnd, nil
}

type state3 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b state3) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	var i int
	input.Get(&i)

	if i != INPUT+3 {
		panic("state3 Execute: input is not expected. Expected: " + fmt.Sprint(INPUT+3) + ", actual: " + fmt.Sprint(i))
	}

	return xdb.DeadEnd, nil
}

func TestTerminateMultiStatesProcess(t *testing.T, client xdb.Client) {
	prcId := "TestTerminateMultiStatesProcess-" + strconv.Itoa(int(time.Now().Unix()))
	prc := MultiStatesProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, INPUT, nil)
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdb.DefaultWorkerUrl, resp.GetWorkerUrl())
	assert.Equal(t, xdb.GetFinalProcessType(prc), resp.GetProcessType())
	assert.NotNil(t, resp.ProcessExecutionId)
	assert.Equal(t, xdbapi.RUNNING, resp.GetStatus())

	err = client.StopProcess(context.Background(), prcId, xdbapi.TERMINATE)
	assert.Nil(t, err)

	resp, err = client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.TERMINATED, resp.GetStatus())
}

func TestFailMultiStatesProcess(t *testing.T, client xdb.Client) {
	prcId := "TestFailMultiStatesProcess-" + strconv.Itoa(int(time.Now().Unix()))
	prc := MultiStatesProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, INPUT, nil)
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdb.DefaultWorkerUrl, resp.GetWorkerUrl())
	assert.Equal(t, xdb.GetFinalProcessType(prc), resp.GetProcessType())
	assert.NotNil(t, resp.ProcessExecutionId)
	assert.Equal(t, xdbapi.RUNNING, resp.GetStatus())

	err = client.StopProcess(context.Background(), prcId, xdbapi.FAIL)
	assert.Nil(t, err)

	resp, err = client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.FAILED, resp.GetStatus())
}
