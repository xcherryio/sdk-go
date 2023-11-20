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

type GracefulCompleteProcess struct {
	xc.ProcessDefaults
}

func (b GracefulCompleteProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(&gracefulCompleteState1{}, &gracefulCompleteState2{}, &gracefulCompleteState3{})
}

type gracefulCompleteState1 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b gracefulCompleteState1) GetStateId() string {
	return "state1"
}

func (b gracefulCompleteState1) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.MultiNextStates(gracefulCompleteState2{}, gracefulCompleteState3{}), nil
}

type gracefulCompleteState2 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b gracefulCompleteState2) GetStateId() string {
	return "state2"
}

func (b gracefulCompleteState2) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.GracefulCompletingProcess, nil
}

type gracefulCompleteState3 struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b gracefulCompleteState3) GetStateId() string {
	return "state3"
}

func (b gracefulCompleteState3) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence, communication xc.Communication,
) (*xc.StateDecision, error) {
	// TODO: add timer
	return xc.DeadEnd, nil
}

func TestGracefulCompleteProcess(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := GracefulCompleteProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, struct{}{})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
