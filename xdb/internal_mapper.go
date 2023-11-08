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
	var localQCmds []xdbapi.LocalQueueCommand
	for _, t := range request.Commands {
		switch t.CommandType {
		case CommandTypeTimer:
			timerCmd := xdbapi.TimerCommand{
				DelayInSeconds: t.TimerCommand.DelayInSeconds,
			}
			timerCmds = append(timerCmds, timerCmd)
		case CommandTypeLocalQueue:
			localQCommand := xdbapi.LocalQueueCommand{
				QueueName: t.LocalQueueCommand.QueueName,
				Count:     ptr.Any(int32(t.LocalQueueCommand.Count)),
			}
			localQCmds = append(localQCmds, localQCommand)
		default:
			panic("unknown command type " + t.CommandType)
		}

	}
	return &xdbapi.CommandRequest{
			WaitingType:        request.CommandWaitingType,
			TimerCommands:      timerCmds,
			LocalQueueCommands: localQCmds,
		},
		nil
}

func fromApiCommandResults(results *xdbapi.CommandResults, encoder ObjectEncoder) (CommandResults, error) {
	if results == nil {
		return CommandResults{}, nil
	}
	var timerResults []TimerResult
	var localQResults []LocalQueueCommandResult
	for _, t := range results.TimerResults {
		timerResult := TimerResult{
			Status: t.Status,
		}
		timerResults = append(timerResults, timerResult)
	}

	for _, t := range results.LocalQueueResults {
		localQResult := LocalQueueCommandResult{
			Result:  t,
			Encoder: encoder,
		}
		localQResults = append(localQResults, localQResult)
	}

	return CommandResults{
		Timers:            timerResults,
		LocalQueueResults: localQResults,
	}, nil
}

func toApiDecision(
	decision *StateDecision, prcType string, registry Registry, encoder ObjectEncoder,
) (*xdbapi.StateDecision, error) {
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
		config := fromStateToAsyncStateConfig(stateDef, prcType, registry)
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

func fromStateToAsyncStateConfig(
	state AsyncState, prcType string, registry Registry,
) *xdbapi.AsyncStateConfig {
	stateCfg := fromAsyncStateOptionsToBasicAsyncStateConfig(state.GetStateOptions())
	if ShouldSkipWaitUntilAPI(state) {
		stateCfg.SkipWaitUntil = ptr.Any(true)
	}

	var preferredPersistencePolicyName *string
	var recoverState AsyncState
	if state.GetStateOptions() != nil {
		preferredPersistencePolicyName = state.GetStateOptions().PersistencePolicyName
		recoverState = state.GetStateOptions().FailureRecoveryState
	}

	stateCfg.LoadGlobalAttributesRequest = createLoadGlobalAttributesRequestIfNeeded(registry, prcType, preferredPersistencePolicyName)
	stateCfg.StateFailureRecoveryOptions = createFailureRecoveryOptionsIfNeeded(recoverState, prcType, registry)
	return stateCfg
}

func createFailureRecoveryOptionsIfNeeded(
	state AsyncState, prcType string, registry Registry,
) *xdbapi.StateFailureRecoveryOptions {
	if state == nil {
		return nil
	}

	stateId := GetFinalStateId(state)
	//NOTE: prevent stack overflow if the state recovering in a loop, e.g. state1 -> state2 -> state1
	if state.GetStateOptions() != nil && state.GetStateOptions().FailureRecoveryState != nil {
		panic("FailureRecoveryState cannot have FailureRecoveryState")
	}
	stateCfg := fromStateToAsyncStateConfig(state, prcType, registry)

	options := &xdbapi.StateFailureRecoveryOptions{
		Policy:                         xdbapi.PROCEED_TO_CONFIGURED_STATE,
		StateFailureProceedStateId:     &stateId,
		StateFailureProceedStateConfig: stateCfg,
	}
	return options
}

func fromAsyncStateOptionsToBasicAsyncStateConfig(
	stateOptions *AsyncStateOptions,
) *xdbapi.AsyncStateConfig {
	stateCfg := &xdbapi.AsyncStateConfig{}
	if stateOptions == nil {
		return stateCfg
	}

	stateCfg.WaitUntilApiTimeoutSeconds = &stateOptions.WaitUntilTimeoutSeconds
	stateCfg.ExecuteApiTimeoutSeconds = &stateOptions.ExecuteTimeoutSeconds
	stateCfg.WaitUntilApiRetryPolicy = stateOptions.WaitUntilRetryPolicy
	stateCfg.ExecuteApiRetryPolicy = stateOptions.ExecuteRetryPolicy
	return stateCfg
}

func createLoadGlobalAttributesRequestIfNeeded(
	registry Registry, prcType string, preferredPersistencePolicyName *string,
) *xdbapi.LoadGlobalAttributesRequest {
	persistenceSchema := registry.getPersistenceSchema(prcType)

	var preferredPolicy *NamedPersistencePolicy
	if preferredPersistencePolicyName != nil {
		preferredPolicyS, ok := persistenceSchema.OverridePersistencePolicies[*preferredPersistencePolicyName]
		if !ok {
			panic("persistence loading policy not found " + *preferredPersistencePolicyName)
		}
		preferredPolicy = &preferredPolicyS
	}

	var tblReqs []xdbapi.TableReadRequest
	if persistenceSchema.GlobalAttributeSchema != nil {
		keyToDefs := registry.getGlobalAttributeKeyToDefs(prcType)

		for _, tblSchema := range persistenceSchema.GlobalAttributeSchema.Tables {
			tblPolicy := getFinalTablePolicy(tblSchema, preferredPolicy)

			var colsToRead []xdbapi.TableColumnDef
			for _, key := range tblPolicy.LoadingKeys {
				def := keyToDefs[key]
				colsToRead = append(colsToRead, xdbapi.TableColumnDef{
					DbColumn: def.colDef.ColumnName,
				})
			}

			tblReqs = append(tblReqs, xdbapi.TableReadRequest{
				TableName:     ptr.Any(tblSchema.TableName),
				Columns:       colsToRead,
				LockingPolicy: ptr.Any(tblPolicy.LockingType),
			})
		}
	}
	if len(tblReqs) == 0 {
		return nil
	}
	return &xdbapi.LoadGlobalAttributesRequest{
		TableRequests: tblReqs,
	}
}

func getFinalTablePolicy(schema DBTableSchema, policy *NamedPersistencePolicy) TablePolicy {
	if policy != nil && policy.GlobalAttributePolicy != nil {
		p, ok := policy.GlobalAttributePolicy[schema.TableName]
		if ok {
			return p
		}
	}
	return schema.DefaultTablePolicy
}
