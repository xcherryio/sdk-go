package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type CommandResults struct {
	Timers            []TimerCommandResult
	LocalQueueResults []LocalQueueCommandResult
}

type TimerCommandResult struct {
	Status xdbapi.CommandStatus
}

type LocalQueueCommandResult struct {
	Result  xdbapi.LocalQueueResult
	Encoder ObjectEncoder
}

func (c CommandResults) GetFirstTimerStatus() xdbapi.CommandStatus {
	return c.GetTimerStatus(0)
}

func (c CommandResults) GetTimerStatus(index int) xdbapi.CommandStatus {
	cmd := c.Timers[index]
	return cmd.Status
}

func (c CommandResults) GetFirstLocalQueueCommand() LocalQueueCommandResult {
	return c.GetLocalQueueCommand(0)
}

func (c CommandResults) GetLocalQueueCommand(index int) LocalQueueCommandResult {
	cmd := c.LocalQueueResults[index]
	return cmd
}

func (lc LocalQueueCommandResult) GetStatus() xdbapi.CommandStatus {
	return lc.Result.Status
}

func (lc LocalQueueCommandResult) GetQueueName() string {
	return lc.Result.QueueName
}

func (lc LocalQueueCommandResult) GetFirstMessage(ptr interface{}) {
	msg := lc.Result.GetMessages()[0]
	err := lc.Encoder.Decode(msg.Payload, ptr)
	if err != nil {
		panic(err)
	}
}

func (lc LocalQueueCommandResult) GetMessages() []Object {
	msgs := lc.Result.GetMessages()
	ret := make([]Object, len(msgs))
	for i, msg := range msgs {
		ret[i] = NewObject(msg.Payload, lc.Encoder)
	}
	return ret
}
