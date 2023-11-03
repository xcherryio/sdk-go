package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type CommandResults struct {
	Timers []TimerCommandResult
}

type TimerCommandResult struct {
	Status xdbapi.CommandStatus
}

func (c CommandResults) GetTimerStatus() xdbapi.CommandStatus {
	if len(c.Timers) != 1 {
		panic("GetTimerCommandResult must be used when there is exactly one timer command")
	}
	cmd := c.Timers[0]
	return cmd.Status
}
