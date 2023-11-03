package xdb

import "time"

type (
	CommandType string

	Command struct {
		CommandType  CommandType
		TimerCommand *TimerCommand
	}

	TimerCommand struct {
		DelayInSeconds int64
	}
)

const (
	CommandTypeTimer CommandType = "Timer"
)

func NewTimerCommand(duration time.Duration) Command {
	return Command{
		CommandType: CommandTypeTimer,
		TimerCommand: &TimerCommand{
			DelayInSeconds: int64(duration.Seconds()),
		},
	}
}
