package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type StateDecision struct {
	NextStates      []StateMovement
	ThreadCloseType *xcapi.ThreadCloseType
}

func SingleNextState(state AsyncState, input interface{}) *StateDecision {
	return &StateDecision{
		NextStates: []StateMovement{
			{
				NextStateId:    GetFinalStateId(state),
				NextStateInput: input,
			},
		},
	}
}

func MultiNextStates(states ...AsyncState) *StateDecision {
	var movements []StateMovement
	for _, st := range states {
		movements = append(movements, StateMovement{
			NextStateId: GetFinalStateId(st),
		})
	}
	return &StateDecision{
		NextStates: movements,
	}
}

func MultiNextStatesWithInput(movements ...StateMovement) *StateDecision {
	return &StateDecision{
		NextStates: movements,
	}
}

var ForceCompletingProcess = &StateDecision{
	ThreadCloseType: xcapi.FORCE_COMPLETE_PROCESS.Ptr(),
}

var GracefulCompletingProcess = &StateDecision{
	ThreadCloseType: xcapi.GRACEFUL_COMPLETE_PROCESS.Ptr(),
}

var DeadEnd = &StateDecision{
	ThreadCloseType: xcapi.DEAD_END.Ptr(),
}

var ForceFailProcess = &StateDecision{
	ThreadCloseType: xcapi.FORCE_FAIL_PROCESS.Ptr(),
}
