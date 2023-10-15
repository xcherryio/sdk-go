package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type StateDecision struct {
	NextStates      []StateMovement
	ThreadCloseType *xdbapi.ThreadCloseType
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
	ThreadCloseType: xdbapi.FORCE_COMPLETE_PROCESS.Ptr(),
}

var GracefulCompletingProcess = &StateDecision{
	ThreadCloseType: xdbapi.GRACEFUL_COMPLETE_PROCESS.Ptr(),
}

var DeadEnd = &StateDecision{
	ThreadCloseType: xdbapi.DEAD_END.Ptr(),
}

var ForceFailProcess = &StateDecision{
	ThreadCloseType: xdbapi.FORCE_FAIL_PROCESS.Ptr(),
}
