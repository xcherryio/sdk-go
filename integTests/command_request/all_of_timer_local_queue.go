package command_request

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc"
	"github.com/xcherryio/sdk-go/xc/ptr"
	"github.com/xcherryio/sdk-go/xc/str"
	"testing"
	"time"
)

type AllOfTimerLocalQProcess struct {
	xc.ProcessDefaults
}

func (b AllOfTimerLocalQProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(
		&initState{},
		&allOfTimerLocalQState{},
		&publishMessagesState{},
	)
}

type initState struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (i initState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	return xc.MultiNextStates(
		allOfTimerLocalQState{},
		publishMessagesState{},
	), nil
}

type allOfTimerLocalQState struct {
	xc.AsyncStateDefaults
}

func (b allOfTimerLocalQState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	return xc.AllOf(
		xc.NewTimerCommand(time.Second*5),
		xc.NewLocalQueueCommand(testQueueName1, 2),
		xc.NewLocalQueueCommand(testQueueName2, 1),
		xc.NewLocalQueueCommand(testQueueName3, 4),
	), nil
}

func (b allOfTimerLocalQState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	secondLocalQ := commandResults.GetLocalQueueCommand(1)
	thirdLocalQ := commandResults.GetLocalQueueCommand(2)

	if commandResults.GetFirstTimerStatus() == xcapi.COMPLETED_COMMAND &&
		commandResults.GetFirstLocalQueueCommand().GetStatus() == xcapi.COMPLETED_COMMAND &&
		secondLocalQ.GetStatus() == xcapi.COMPLETED_COMMAND &&
		thirdLocalQ.GetStatus() == xcapi.COMPLETED_COMMAND {

		// validate first queue
		var msg1 string
		commandResults.GetFirstLocalQueueCommand().GetFirstMessage(&msg1)
		if msg1 != "testLocalQ1" {
			panic("unexpected message:" + msg1)
		}
		var secondMsg MyMessage
		msgs := commandResults.GetFirstLocalQueueCommand().GetMessages()
		msgs[1].Get(&secondMsg)
		if secondMsg != testMyMsq {
			panic("unexpected second message:" + str.AnyToJson(secondMsg))
		}

		// validate second queue
		var msg2 string
		secondLocalQ.GetFirstMessage(&msg2)
		if msg2 != "" {
			panic("unexpected message:" + msg1)
		}

		// validate third queue
		msgs = thirdLocalQ.GetMessages()
		if len(msgs) != 4 {
			panic("unexpected message count:" + str.AnyToJson(msgs))
		}
		var s string
		var i int
		var myMsg2, myMsg3 MyMessage
		msgs[0].Get(&s)
		msgs[1].Get(&myMsg2)
		msgs[2].Get(&i)
		msgs[3].Get(&myMsg3)
		if s != "publishMessagesState" ||
			myMsg2 != testMyMsq ||
			i != 123 ||
			myMsg3 != testMyMsq {
			panic("unexpected messages:" + str.AnyToJson(msgs))
		}

		return xc.GracefulCompletingProcess, nil
	} else {
		panic("unexpected command results" + str.AnyToJson(commandResults))
	}
}

type publishMessagesState struct {
	xc.AsyncStateDefaults
}

func (p publishMessagesState) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	communication.PublishToLocalQueue(testQueueName3, "publishMessagesState")
	communication.PublishToLocalQueue(testQueueName3, testMyMsq)
	return xc.EmptyCommandRequest(), nil
}

func (p publishMessagesState) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	communication.PublishToLocalQueue(testQueueName3, 123)
	communication.PublishToLocalQueue(testQueueName3, testMyMsq)
	return xc.DeadEnd, nil
}

func TestAllOfTimerLocalQueue(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := AllOfTimerLocalQProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, nil)
	assert.Nil(t, err)

	tuid, err := uuid.NewUUID()
	assert.Nil(t, err)

	err = client.BatchPublishToLocalQueue(context.Background(), prcId,
		xc.LocalQueuePublishMessage{
			QueueName: testQueueName1,
			Payload:   "testLocalQ1",
			DedupSeed: ptr.Any("testLocalQ1"),
		},
		xc.LocalQueuePublishMessage{
			QueueName: testQueueName2,
		},
		xc.LocalQueuePublishMessage{
			QueueName: testQueueName1,
			Payload:   testMyMsq,
			DedupUUID: ptr.Any(tuid.String()),
		})
	assert.Nil(t, err)

	time.Sleep(time.Second * 6)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
