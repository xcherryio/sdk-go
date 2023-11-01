package state_decision

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/integTests/common"
	"github.com/xdblab/xdb-golang-sdk/xdb"
)

type ForceFailProcess struct {
	xdb.ProcessDefaults
}

func (b ForceFailProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(&forceFailState1{}, &forceFailState2{}, &forceFailState3{})
}

type forceFailState1 struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b forceFailState1) GetStateId() string {
	return "state1"
}

func (b forceFailState1) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	return xdb.MultiNextStates(forceFailState2{}, forceFailState3{}), nil
}

type forceFailState2 struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b forceFailState2) GetStateId() string {
	return "state2"
}

func (b forceFailState2) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	return xdb.ForceFailProcess, nil
}

type forceFailState3 struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b forceFailState3) GetStateId() string {
	return "state3"
}

func (b forceFailState3) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	// TODO: add timer
	return xdb.DeadEnd, nil
}

func TestForceFailProcess(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := ForceFailProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, struct{}{})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.FAILED, resp.GetStatus())
}
