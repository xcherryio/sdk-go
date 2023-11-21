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

type DeadEndProcess struct {
	xc.ProcessDefaults
}

func (b DeadEndProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(&deadEndState1{}, &deadEndState2{}, &deadEndState3{})
}

type deadEndState1 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b deadEndState1) GetStateId() string {
	return "state1"
}

func (b deadEndState1) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.MultiNextStates(deadEndState2{}, deadEndState3{}), nil
}

type deadEndState2 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b deadEndState2) GetStateId() string {
	return "state2"
}

func (b deadEndState2) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.DeadEnd, nil
}

type deadEndState3 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b deadEndState3) GetStateId() string {
	return "state3"
}

func (b deadEndState3) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.DeadEnd, nil
}

func TestDeadEndProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := DeadEndProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, struct{}{})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.RUNNING, resp.GetStatus())
}
