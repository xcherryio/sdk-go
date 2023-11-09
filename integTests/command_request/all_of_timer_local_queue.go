package command_request

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/integTests/common"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
	"github.com/xdblab/xdb-golang-sdk/xdb/str"
	"testing"
	"time"
)

type AllOfTimerLocalQProcess struct {
	xdb.ProcessDefaults
}

func (b AllOfTimerLocalQProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(
		&initState{},
		&allOfTimerLocalQState{},
		&publishMessagesState{},
	)
}

type initState struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (i initState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	return xdb.MultiNextStates(
		allOfTimerLocalQState{},
		publishMessagesState{},
	), nil
}

type allOfTimerLocalQState struct {
	xdb.AsyncStateDefaults
}

func (b allOfTimerLocalQState) WaitUntil(
	ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication,
) (*xdb.CommandRequest, error) {
	return xdb.AllOf(
		xdb.NewTimerCommand(time.Second*5),
		xdb.NewLocalQueueCommand(testQueueName1, 2),
		xdb.NewLocalQueueCommand(testQueueName2, 1),
		xdb.NewLocalQueueCommand(testQueueName3, 4),
	), nil
}

func (b allOfTimerLocalQState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	secondLocalQ := commandResults.GetLocalQueueCommand(1)
	thirdLocalQ := commandResults.GetLocalQueueCommand(2)

	if commandResults.GetFirstTimerStatus() == xdbapi.COMPLETED_COMMAND &&
		commandResults.GetFirstLocalQueueCommand().GetStatus() == xdbapi.COMPLETED_COMMAND &&
		secondLocalQ.GetStatus() == xdbapi.COMPLETED_COMMAND &&
		thirdLocalQ.GetStatus() == xdbapi.COMPLETED_COMMAND {

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

		return xdb.GracefulCompletingProcess, nil
	} else {
		panic("unexpected command results" + str.AnyToJson(commandResults))
	}
}

type publishMessagesState struct {
	xdb.AsyncStateDefaults
}

func (p publishMessagesState) WaitUntil(
	ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication,
) (*xdb.CommandRequest, error) {
	communication.PublishToLocalQueue(testQueueName3, "publishMessagesState")
	communication.PublishToLocalQueue(testQueueName3, testMyMsq)
	return xdb.EmptyCommandRequest(), nil
}

func (p publishMessagesState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	communication.PublishToLocalQueue(testQueueName3, 123)
	communication.PublishToLocalQueue(testQueueName3, testMyMsq)
	return xdb.DeadEnd, nil
}

func TestAllOfTimerLocalQueue(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := AllOfTimerLocalQProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, nil)
	assert.Nil(t, err)

	tuid, err := uuid.NewUUID()
	assert.Nil(t, err)

	err = client.BatchPublishToLocalQueue(context.Background(), prcId,
		xdb.LocalQueuePublishMessage{
			QueueName: testQueueName1,
			Payload:   "testLocalQ1",
			DedupSeed: ptr.Any("testLocalQ1"),
		},
		xdb.LocalQueuePublishMessage{
			QueueName: testQueueName2,
		},
		xdb.LocalQueuePublishMessage{
			QueueName: testQueueName1,
			Payload:   testMyMsq,
			DedupUUID: ptr.Any(tuid.String()),
		})
	assert.Nil(t, err)

	time.Sleep(time.Second * 6)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}
