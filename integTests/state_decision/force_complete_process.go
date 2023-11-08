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

type ForceCompleteProcess struct {
	xdb.ProcessDefaults
}

func (b ForceCompleteProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(&forceCompleteState1{}, &forceCompleteState2{}, &forceCompleteState3{})
}

type forceCompleteState1 struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b forceCompleteState1) GetStateId() string {
	return "state1"
}

func (b forceCompleteState1) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	return xdb.MultiNextStates(forceCompleteState2{}, forceCompleteState3{}), nil
}

type forceCompleteState2 struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b forceCompleteState2) GetStateId() string {
	return "state2"
}

func (b forceCompleteState2) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	return xdb.ForceCompletingProcess, nil
}

type forceCompleteState3 struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b forceCompleteState3) GetStateId() string {
	return "state3"
}

func (b forceCompleteState3) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	// TODO: add timer
	return xdb.DeadEnd, nil
}

func TestForceCompleteProcess(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := ForceCompleteProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, struct{}{})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	resp, err := client.DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}
