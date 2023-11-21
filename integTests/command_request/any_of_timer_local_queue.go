package command_request

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc"
	"github.com/xcherryio/sdk-go/xc/str"
	"testing"
	"time"
)

type AnyOfTimerLocalQProcess struct {
	xc.ProcessDefaults
}

func (b AnyOfTimerLocalQProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(&anyOfTimerLocalQState{})
}

type anyOfTimerLocalQState struct {
	xc.AsyncStateDefaults
}

func (b anyOfTimerLocalQState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return xc.AnyOf(
		xc.NewTimerCommand(time.Second*5),
		xc.NewLocalQueueCommand(testQueueName1, 2),
	), nil
}

func (b anyOfTimerLocalQState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var expected string
	input.Get(&expected)
	if expected == "timer" {
		if commandResults.GetFirstTimerStatus() == xcapi.COMPLETED_COMMAND &&
			commandResults.GetFirstLocalQueueCommand().GetStatus() == xcapi.WAITING_COMMAND {
			return xc.GracefulCompletingProcess, nil
		} else {
			return xc.ForceFailProcess, nil
		}
	} else if expected == "localQueue" {
		if commandResults.GetFirstTimerStatus() == xcapi.WAITING_COMMAND &&
			commandResults.GetFirstLocalQueueCommand().GetStatus() == xcapi.COMPLETED_COMMAND {
			var msg string
			commandResults.GetFirstLocalQueueCommand().GetFirstMessage(&msg)
			if msg != expected {
				panic("unexpected message:" + msg)
			}
			var secondMsg MyMessage
			msgs := commandResults.GetFirstLocalQueueCommand().GetMessages()
			msgs[1].Get(&secondMsg)
			if secondMsg != testMyMsq {
				panic("unexpected second message:" + str.AnyToJson(secondMsg))
			}
			return xc.GracefulCompletingProcess, nil
		} else {
			panic("unexpected command results" + str.AnyToJson(commandResults))
		}
	} else {
		panic("unexpected input:" + expected)
	}
}

func TestAnyOfTimerLocalQueueWithTimerFired(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := AnyOfTimerLocalQProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, "timer")
	assert.Nil(t, err)

	time.Sleep(time.Second * 6)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}

func TestAnyOfTimerLocalQueueWithLocalQueueMessagesReceived(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := AnyOfTimerLocalQProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, "localQueue")
	assert.Nil(t, err)

	err = client.PublishToLocalQueue(context.Background(), prcId, testQueueName1, "localQueue", nil)
	assert.Nil(t, err)
	err = client.PublishToLocalQueue(context.Background(), prcId, testQueueName1, testMyMsq, nil)
	assert.Nil(t, err)

	time.Sleep(time.Second * 6)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
