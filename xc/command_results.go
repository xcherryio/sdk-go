package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type CommandResults struct {
	TimerResults      []TimerResult
	LocalQueueResults []LocalQueueCommandResult
}

type TimerResult struct {
	Status xcapi.CommandStatus
}

type LocalQueueCommandResult struct {
	Result  xcapi.LocalQueueResult
	Encoder ObjectEncoder
}

func (c CommandResults) GetFirstTimerStatus() xcapi.CommandStatus {
	return c.GetTimerStatus(0)
}

func (c CommandResults) GetTimerStatus(index int) xcapi.CommandStatus {
	cmd := c.TimerResults[index]
	return cmd.Status
}

func (c CommandResults) GetFirstLocalQueueCommand() LocalQueueCommandResult {
	return c.GetLocalQueueCommand(0)
}

func (c CommandResults) GetLocalQueueCommand(index int) LocalQueueCommandResult {
	cmd := c.LocalQueueResults[index]
	return cmd
}

func (lc LocalQueueCommandResult) GetStatus() xcapi.CommandStatus {
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
