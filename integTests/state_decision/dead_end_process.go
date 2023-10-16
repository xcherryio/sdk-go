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

type DeadEndProcess struct {
	xdb.ProcessDefaults
}

func (b DeadEndProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.WithStartingState(&deadEndState1{}, &deadEndState2{}, &deadEndState3{})
}

type deadEndState1 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b deadEndState1) GetStateId() string {
	return "state1"
}

func (b deadEndState1) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	return xdb.MultiNextStates(deadEndState2{}, deadEndState3{}), nil
}

type deadEndState2 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b deadEndState2) GetStateId() string {
	return "state2"
}

func (b deadEndState2) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	return xdb.DeadEnd, nil
}

type deadEndState3 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b deadEndState3) GetStateId() string {
	return "state3"
}

func (b deadEndState3) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	return xdb.DeadEnd, nil
}

func TestDeadEndProcess(t *testing.T, client xdb.Client) {
	prcId := "TestDeadEndProcess-" + strconv.Itoa(int(time.Now().Unix()))
	prc := DeadEndProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, struct{}{})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.RUNNING, resp.GetStatus())
}
