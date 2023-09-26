package xdb

type StateMovement struct {
	// NextStateId is required
	NextStateId string
	// NextStateInput is optional
	NextStateInput interface{}
}

func NewStateMovement(st AsyncState, input interface{}) StateMovement {
	return StateMovement{
		NextStateId:    GetFinalStateId(st),
		NextStateInput: input,
	}
}
