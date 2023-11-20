package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type CommandRequest struct {
	Commands           []Command
	CommandWaitingType xcapi.CommandWaitingType
}

// EmptyCommandRequest will jump to decide stage immediately.
func EmptyCommandRequest() *CommandRequest {
	return &CommandRequest{
		CommandWaitingType: xcapi.EMPTY_COMMAND,
	}
}

// AnyOf will wait for any of the commands to complete
func AnyOf(commands ...Command) *CommandRequest {
	return &CommandRequest{
		Commands:           commands,
		CommandWaitingType: xcapi.ANY_OF_COMPLETION,
	}
}

// AllOf will wait for all the commands to complete
func AllOf(commands ...Command) *CommandRequest {
	return &CommandRequest{
		Commands:           commands,
		CommandWaitingType: xcapi.ALL_OF_COMPLETION,
	}
}
