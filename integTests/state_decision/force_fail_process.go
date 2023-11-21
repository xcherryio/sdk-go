package state_decision

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc"
)

type ForceFailProcess struct {
	xc.ProcessDefaults
}

func (b ForceFailProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(&forceFailState1{}, &forceFailState2{}, &forceFailState3{})
}

type forceFailState1 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b forceFailState1) GetStateId() string {
	return "state1"
}

func (b forceFailState1) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.MultiNextStates(forceFailState2{}, forceFailState3{}), nil
}

type forceFailState2 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b forceFailState2) GetStateId() string {
	return "state2"
}

func (b forceFailState2) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.ForceFailProcess, nil
}

type forceFailState3 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b forceFailState3) GetStateId() string {
	return "state3"
}

func (b forceFailState3) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	// TODO: add timer
	return xc.DeadEnd, nil
}

func TestForceFailProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := ForceFailProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, struct{}{})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.FAILED, resp.GetStatus())
}
