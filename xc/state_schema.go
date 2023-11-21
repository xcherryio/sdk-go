package xc

type StateSchema struct {
	StartingState AsyncState
	AllStates     []AsyncState
}

func NewStateSchema(startingState AsyncState, nonStartingStates ...AsyncState) StateSchema {
	allStates := nonStartingStates
	allStates = append(allStates, startingState)
	return StateSchema{
		StartingState: startingState,
		AllStates:     allStates,
	}
}

func NewStateSchemaNoStartingState(nonStartingStates ...AsyncState) StateSchema {
	return StateSchema{
		AllStates: nonStartingStates,
	}
}
