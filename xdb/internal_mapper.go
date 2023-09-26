package xdb

import (
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
)

func toApiCommandRequest(request *CommandRequest) (*xdbapi.CommandRequest, error) {
	// TODO
	return nil, nil
}

func toIdlDecision(decision *StateDecision, prcType string, registry Registry, encoder ObjectEncoder) (*xdbapi.StateDecision, error) {
	if decision.ThreadCloseType != nil && len(decision.NextStates) > 0 {
		return nil, NewProcessDefinitionError("cannot have both next state and closing in a single decision")
	}

	if decision.ThreadCloseType != nil {
		return &xdbapi.StateDecision{
			ThreadCloseDecision: &xdbapi.ThreadCloseDecision{
				CloseType: decision.ThreadCloseType,
			},
		}, nil
	}

	var mvs []xdbapi.StateMovement
	for _, fromMv := range decision.NextStates {
		input, err := encoder.Encode(fromMv.NextStateInput)
		if err != nil {
			return nil, err
		}
		var config *xdbapi.AsyncStateConfig
		stateDef := registry.getProcessState(prcType, fromMv.NextStateId)
		if ShouldSkipWaitUntilAPI(stateDef) {
			config = &xdbapi.AsyncStateConfig{
				SkipWaitUntil: ptr.Any(true),
			}
		}
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
