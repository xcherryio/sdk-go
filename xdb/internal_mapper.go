package xdb

import (
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
)

func toApiCommandRequest(request *CommandRequest) (*xdbapi.CommandRequest, error) {
	if request == nil {
		return nil, NewProcessDefinitionError("command request cannot be nil")
	}
	var timerCmds []xdbapi.TimerCommand
	for _, t := range request.Commands {
		if t.CommandType == CommandTypeTimer {
			timerCmd := xdbapi.TimerCommand{
				CommandId:                  t.CommandId,
				FiringUnixTimestampSeconds: t.TimerCommand.FiringUnixTimestampSeconds,
			}
			timerCmds = append(timerCmds, timerCmd)
		}
	}
	return &xdbapi.CommandRequest{
			WaitingType:   request.CommandWaitingType,
			TimerCommands: timerCmds,
		},
		nil
}

func fromApiCommandResults(results *xdbapi.CommandResults, _ ObjectEncoder) (CommandResults, error) {
	if results == nil {
		return CommandResults{}, nil
	}
	var timerResults []TimerCommandResult
	for _, t := range results.TimerResults {
		timerResult := TimerCommandResult{
			CommandId: t.CommandId,
			Status:    t.TimerStatus,
		}
		timerResults = append(timerResults, timerResult)
	}

	return CommandResults{
		Timers: timerResults,
	}, nil
}

func toApiDecision(decision *StateDecision, prcType string, registry Registry, encoder ObjectEncoder) (*xdbapi.StateDecision, error) {
	if decision == nil {
		return nil, NewProcessDefinitionError("StateDecision cannot be nil")
	}
	if decision.ThreadCloseType != nil && len(decision.NextStates) > 0 {
		return nil, NewProcessDefinitionError("cannot have both next state and closing in a single decision")
	}

	if decision.ThreadCloseType != nil {
		return &xdbapi.StateDecision{
			ThreadCloseDecision: &xdbapi.ThreadCloseDecision{
				CloseType: *decision.ThreadCloseType,
			},
		}, nil
	}

	var mvs []xdbapi.StateMovement
	for _, fromMv := range decision.NextStates {
		input, err := encoder.Encode(fromMv.NextStateInput)
		if err != nil {
			return nil, err
		}
		stateDef := registry.getProcessState(prcType, fromMv.NextStateId)
		config := fromStateToAsyncStateConfig(stateDef)
		mv := xdbapi.StateMovement{
			StateId:     fromMv.NextStateId,
			StateInput:  input,
			StateConfig: config,
		}
		mvs = append(mvs, mv)
	}
	return &xdbapi.StateDecision{
		NextStates: mvs,
	}, nil
}

func fromStateToAsyncStateConfig(state AsyncState) *xdbapi.AsyncStateConfig {
	stateCfg := &xdbapi.AsyncStateConfig{}
	if ShouldSkipWaitUntilAPI(state) {
		stateCfg.SkipWaitUntil = ptr.Any(true)
	}
	options := state.GetStateOptions()
	if options != nil {
		stateCfg.WaitUntilApiTimeoutSeconds = &options.WaitUntilTimeoutSeconds
		stateCfg.ExecuteApiTimeoutSeconds = &options.ExecuteTimeoutSeconds
		stateCfg.WaitUntilApiRetryPolicy = options.WaitUntilRetryPolicy
		stateCfg.ExecuteApiRetryPolicy = options.ExecuteRetryPolicy
	}
	return stateCfg
}
