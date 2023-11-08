package command_request

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/integTests/common"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
	"testing"
	"time"
)

type AnyOfTimerLocalQProcess struct {
	xdb.ProcessDefaults
}

func (b AnyOfTimerLocalQProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(&anyOfTimerLocalQState{})
}

type anyOfTimerLocalQState struct {
	xdb.AsyncStateDefaults
}

const (
	testQueueName = "test-queue"
)

var testMyMsq = MyMessage{
	Str: "test-message",
	Int: 123,
}

type MyMessage struct {
	Str string
	Int int
}

func (b anyOfTimerLocalQState) WaitUntil(
	ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication,
) (*xdb.CommandRequest, error) {
	return xdb.AnyOf(
		xdb.NewTimerCommand(time.Second*5),
		xdb.NewLocalQueueCommand(testQueueName, 2),
	), nil
}

func (b anyOfTimerLocalQState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var expected string
	input.Get(&expected)
	if expected == "timer" {
		if commandResults.GetFirstTimerStatus() == xdbapi.COMPLETED_COMMAND &&
			commandResults.GetFirstLocalQueueCommand().GetStatus() == xdbapi.WAITING_COMMAND {
			return xdb.GracefulCompletingProcess, nil
		} else {
			return xdb.ForceFailProcess, nil
		}
	} else if expected == "localQueue" {
		if commandResults.GetFirstTimerStatus() == xdbapi.WAITING_COMMAND &&
			commandResults.GetFirstLocalQueueCommand().GetStatus() == xdbapi.COMPLETED_COMMAND {
			var msg string
			commandResults.GetFirstLocalQueueCommand().GetFirstMessage(&msg)
			if msg == expected {
				return xdb.GracefulCompletingProcess, nil
			}
			var secondMsg MyMessage
			msgs := commandResults.GetFirstLocalQueueCommand().GetMessages()
			msgs[1].Get(&secondMsg)
			if secondMsg == testMyMsq {
				panic("unexpected second message:" + ptr.AnyToJson(secondMsg))
			}
			panic("unexpected message:" + msg)
		} else {
			panic("unexpected command results" + ptr.AnyToJson(commandResults))
		}
	} else {
		panic("unexpected input:" + expected)
	}
}

func TestAnyOfTimerLocalQueueWithTimerFired(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := AnyOfTimerLocalQProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, "timer")
	assert.Nil(t, err)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)

	time.Sleep(time.Second * 6)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}
