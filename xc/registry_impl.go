package xc

type registryImpl struct {
	processStore           map[string]Process
	persistenceSchemaStore map[string]PersistenceSchema
	localAttrKeys          map[string]map[string]bool
	startingState          map[string]AsyncState
	stateStore             map[string]map[string]AsyncState
}

func (r *registryImpl) AddProcess(processDef Process) error {
	if err := r.registerProcessType(processDef); err != nil {
		return err
	}
	if err := r.registerProcessState(processDef); err != nil {
		return err
	}
	if err := r.registerPersistenceSchema(processDef); err != nil {
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
	for prcType := range r.processStore {
		res = append(res, prcType)
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

func (r *registryImpl) getPersistenceSchema(prcType string) PersistenceSchema {
	return r.persistenceSchemaStore[prcType]
}

func (r *registryImpl) registerProcessType(prc Process) error {
	prcType := GetFinalProcessType(prc)
	_, ok := r.processStore[prcType]
	if ok {
		return NewProcessDefinitionError("Process type conflict: " + prcType)
	}
	r.processStore[prcType] = prc
	return nil
}

func (r *registryImpl) registerProcessState(prc Process) error {
	prcType := GetFinalProcessType(prc)
	stateMap := map[string]AsyncState{}
	for _, state := range prc.GetAsyncStateSchema().AllStates {
		stateId := GetFinalStateId(state)
		_, ok := stateMap[stateId]
		if ok {
			return NewProcessDefinitionError("Process %v cannot have duplicate stateId %v ", prcType, stateId)
		}
		stateMap[stateId] = state
	}
	r.stateStore[prcType] = stateMap
	if prc.GetAsyncStateSchema().StartingState != nil {
		r.startingState[prcType] = prc.GetAsyncStateSchema().StartingState
	}

	return nil
}

func (r *registryImpl) registerPersistenceSchema(prc Process) error {
	prcType := GetFinalProcessType(prc)
	ps := prc.GetPersistenceSchema()
	r.persistenceSchemaStore[prcType] = ps

	localAttrKeys, err := ps.ValidateLocalAttributeForRegistry()
	if err != nil {
		return err
	}
	r.localAttrKeys[prcType] = localAttrKeys

	return nil
}

func (r *registryImpl) getLocalAttributeKeys(prcType string) map[string]bool {
	return r.localAttrKeys[prcType]
}
