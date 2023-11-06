package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type CommandRequest struct {
	Commands           []Command
	CommandWaitingType xdbapi.CommandWaitingType
}

// EmptyCommandRequest will jump to decide stage immediately.
func EmptyCommandRequest() *CommandRequest {
	return &CommandRequest{
		CommandWaitingType: xdbapi.EMPTY_COMMAND,
	}
}

// AnyOf will wait for any of the commands to complete
func AnyOf(commands ...Command) *CommandRequest {
	return &CommandRequest{
		Commands:           commands,
		CommandWaitingType: xdbapi.ANY_OF_COMPLETION,
	}
}

// AllOf will wait for all the commands to complete
func AllOf(commands ...Command) *CommandRequest {
	return &CommandRequest{
		Commands:           commands,
		CommandWaitingType: xdbapi.ALL_OF_COMPLETION,
	}
}
