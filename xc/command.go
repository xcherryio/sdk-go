package xc

import "time"

type (
	CommandType string

	Command struct {
		CommandType       CommandType
		TimerCommand      *TimerCommand
		LocalQueueCommand *LocalQueueCommand
	}

	TimerCommand struct {
		DelayInSeconds int64
	}

	LocalQueueCommand struct {
		QueueName string
		Count     int
	}
)

const (
	CommandTypeTimer      CommandType = "Timer"
	CommandTypeLocalQueue CommandType = "LocalQueue"
)

func NewTimerCommand(duration time.Duration) Command {
	return Command{
		CommandType: CommandTypeTimer,
		TimerCommand: &TimerCommand{
			DelayInSeconds: int64(duration.Seconds()),
		},
	}
}

func NewLocalQueueCommand(queueName string, count int) Command {
	return Command{
		CommandType: CommandTypeLocalQueue,
		LocalQueueCommand: &LocalQueueCommand{
			QueueName: queueName,
			Count:     count,
		},
	}
}
