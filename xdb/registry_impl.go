package xdb

type registryImpl struct {
	processStore  map[string]Process
	startingState map[string]AsyncState
	stateStore    map[string]map[string]AsyncState
}

func (r *registryImpl) AddProcess(processDef Process) error {
	if err := r.registerProcessType(processDef); err != nil {
		return err
	}
	if err := r.registerProcessState(processDef); err != nil {
		return err
	}
	return nil
}

func (r *registryImpl) AddProcesses(processDefs ...Process) error {
	for _, wf := range processDefs {
		err := r.AddProcess(wf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *registryImpl) GetAllRegisteredProcessTypes() []string {
	var res []string
	for wfType := range r.processStore {
		res = append(res, wfType)
	}
	return res
}

func (r *registryImpl) getProcess(prcType string) Process {
	return r.processStore[prcType]
}

func (r *registryImpl) getProcessStartingState(prcType string) AsyncState {
	return r.startingState[prcType]
}

func (r *registryImpl) getProcessState(prcType string, stateId string) AsyncState {
	return r.stateStore[prcType][stateId]
}

func (r *registryImpl) registerProcessType(prc Process) error {
	wfType := GetFinalProcessType(prc)
	_, ok := r.processStore[wfType]
	if ok {
		return NewProcessDefinitionError("Process type conflict: " + wfType)
	}
	r.processStore[wfType] = prc
	return nil
}

func (r *registryImpl) registerProcessState(prc Process) error {
	wfType := GetFinalProcessType(prc)
	stateMap := map[string]AsyncState{}
	for _, state := range prc.GetAsyncStateSchema().AllStates {
		stateId := GetFinalStateId(state)
		_, ok := stateMap[stateId]
		if ok {
			return NewProcessDefinitionError("Process %v cannot have duplicate stateId %v ", wfType, stateId)
		}
		stateMap[stateId] = state
	}
	r.stateStore[wfType] = stateMap
	if prc.GetAsyncStateSchema().StartingState != nil {
		r.startingState[wfType] = prc.GetAsyncStateSchema().StartingState
	}
	
	return nil
}
