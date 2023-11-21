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

type ForceCompleteProcess struct {
	xc.ProcessDefaults
}

func (b ForceCompleteProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(&forceCompleteState1{}, &forceCompleteState2{}, &forceCompleteState3{})
}

type forceCompleteState1 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b forceCompleteState1) GetStateId() string {
	return "state1"
}

func (b forceCompleteState1) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.MultiNextStates(forceCompleteState2{}, forceCompleteState3{}), nil
}

type forceCompleteState2 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b forceCompleteState2) GetStateId() string {
	return "state2"
}

func (b forceCompleteState2) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.ForceCompletingProcess, nil
}

type forceCompleteState3 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b forceCompleteState3) GetStateId() string {
	return "state3"
}

func (b forceCompleteState3) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	// TODO: add timer
	return xc.DeadEnd, nil
}

func TestForceCompleteProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := ForceCompleteProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, struct{}{})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
