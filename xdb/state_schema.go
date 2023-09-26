package xdb

type StateSchema struct {
	StartingState AsyncState
	AllStates     []AsyncState
}

func WithStartingState(startingState AsyncState, nonStartingStates ...AsyncState) StateSchema {
	allStates := nonStartingStates
	allStates = append(allStates, startingState)
	return StateSchema{
		StartingState: startingState,
		AllStates:     allStates,
	}
}

func NoStartingState(nonStartingStates ...AsyncState) StateSchema {
	return StateSchema{
		AllStates: nonStartingStates,
	}
}
