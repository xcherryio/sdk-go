package xdb

import "time"

type (
	CommandType string

	Command struct {
		CommandId    *string
		CommandType  CommandType
		TimerCommand *TimerCommand
	}

	TimerCommand struct {
		FiringUnixTimestampSeconds int64
	}
)

const (
	CommandTypeTimer CommandType = "Timer"
)

func NewTimerCommand(duration time.Duration) Command {
	nowUnix := time.Now().Unix()
	return Command{
		CommandType: CommandTypeTimer,
		TimerCommand: &TimerCommand{
			FiringUnixTimestampSeconds: nowUnix + int64(duration.Seconds()),
		},
	}
}
