package state_decision

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"strconv"
	"testing"
	"time"
)

type GracefulCompleteProcess struct {
	xdb.ProcessDefaults
}

func (b GracefulCompleteProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.WithStartingState(&gracefulCompleteState1{}, &gracefulCompleteState2{}, &gracefulCompleteState3{})
}

type gracefulCompleteState1 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b gracefulCompleteState1) GetStateId() string {
	return "state1"
}

func (b gracefulCompleteState1) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	return xdb.MultiNextStates(gracefulCompleteState2{}, gracefulCompleteState3{}), nil
}

type gracefulCompleteState2 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b gracefulCompleteState2) GetStateId() string {
	return "state2"
}

func (b gracefulCompleteState2) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	return xdb.GracefulCompletingProcess, nil
}

type gracefulCompleteState3 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b gracefulCompleteState3) GetStateId() string {
	return "state3"
}

func (b gracefulCompleteState3) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	// TODO: add timer
	return xdb.DeadEnd, nil
}

func TestGracefulCompleteProcess(t *testing.T, client xdb.Client) {
	prcId := "TestGracefulCompleteProcess-" + strconv.Itoa(int(time.Now().Unix()))
	prc := GracefulCompleteProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, struct{}{}, nil)
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}