package xdb

type StateSchema struct {
	startingState     AsyncState
	nonStartingStates []AsyncState
}

func WithStartingState(startingState AsyncState, nonStartingStates ...AsyncState) StateSchema {
	return StateSchema{
		startingState:     startingState,
		nonStartingStates: nonStartingStates,
	}
}

func NoStartingState(nonStartingStates ...AsyncState) StateSchema {
	return StateSchema{
		nonStartingStates: nonStartingStates,
	}
}
