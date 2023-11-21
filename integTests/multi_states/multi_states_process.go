package multi_states

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc"
)

const INPUT = 1

type MultiStatesProcess struct {
	xc.ProcessDefaults
}

func (b MultiStatesProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(&state1{}, &state2{}, &state3{})
}

type state1 struct {
	xc.AsyncStateDefaults
}

func (b state1) WaitUntil(ctx xc.Context, input xc.Object, communication xc.Communication) (*xc.CommandRequest, error) {
	var i int
	input.Get(&i)

	if i != INPUT {
		panic("state1 WaitUntil: input is not expected. Expected: " + fmt.Sprint(INPUT) + ", actual: " + fmt.Sprint(i))
	}

	return xc.EmptyCommandRequest(), nil
}

func (b state1) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)

	if i != INPUT {
		panic("state1 Execute: input is not expected. Expected: " + fmt.Sprint(INPUT) + ", actual: " + fmt.Sprint(i))
	}

	return xc.MultiNextStatesWithInput(xc.NewStateMovement(state2{}, i+2), xc.NewStateMovement(state3{}, i+3)), nil
}

type state2 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b state2) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)

	if i != INPUT+2 {
		panic("state2 Execute: input is not expected. Expected: " + fmt.Sprint(INPUT+2) + ", actual: " + fmt.Sprint(i))
	}

	return xc.DeadEnd, nil
}

type state3 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b state3) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	input.Get(&i)

	if i != INPUT+3 {
		panic("state3 Execute: input is not expected. Expected: " + fmt.Sprint(INPUT+3) + ", actual: " + fmt.Sprint(i))
	}

	return xc.DeadEnd, nil
}

func TestTerminateMultiStatesProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := MultiStatesProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, INPUT)
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xc.DefaultWorkerUrl, resp.GetWorkerUrl())
	assert.Equal(t, xc.GetFinalProcessType(prc), resp.GetProcessType())
	assert.NotNil(t, resp.ProcessExecutionId)
	assert.Equal(t, xcapi.RUNNING, resp.GetStatus())

	err = client.StopProcess(context.Background(), prcId, xcapi.TERMINATE)
	assert.Nil(t, err)

	resp, err = client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.TERMINATED, resp.GetStatus())
}

func TestFailMultiStatesProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := MultiStatesProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, INPUT)
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xc.DefaultWorkerUrl, resp.GetWorkerUrl())
	assert.Equal(t, xc.GetFinalProcessType(prc), resp.GetProcessType())
	assert.NotNil(t, resp.ProcessExecutionId)
	assert.Equal(t, xcapi.RUNNING, resp.GetStatus())

	err = client.StopProcess(context.Background(), prcId, xcapi.FAIL)
	assert.Nil(t, err)

	resp, err = client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.FAILED, resp.GetStatus())
}
