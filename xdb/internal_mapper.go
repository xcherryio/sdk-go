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
	for _, t := range request.Commands {
		if t.CommandType == CommandTypeTimer {
			timerCmd := xdbapi.TimerCommand{
				CommandId:                  t.CommandId,
				FiringUnixTimestampSeconds: t.TimerCommand.FiringUnixTimestampSeconds,
			}
			timerCmds = append(timerCmds, timerCmd)
		}
	}
	return &xdbapi.CommandRequest{
			WaitingType:   request.CommandWaitingType,
			TimerCommands: timerCmds,
		},
		nil
}

func fromApiCommandResults(results *xdbapi.CommandResults, _ ObjectEncoder) (CommandResults, error) {
	if results == nil {
		return CommandResults{}, nil
	}
	var timerResults []TimerCommandResult
	for _, t := range results.TimerResults {
		timerResult := TimerCommandResult{
			CommandId: t.CommandId,
			Status:    t.TimerStatus,
		}
		timerResults = append(timerResults, timerResult)
	}

	return CommandResults{
		Timers: timerResults,
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
	preferredPersistencePolicyName, stateCfg := fromAsyncStateOptionsToAsyncStateConfig(state.GetStateOptions())
	if ShouldSkipWaitUntilAPI(state) {
		stateCfg.SkipWaitUntil = ptr.Any(true)
	}

	stateCfg.LoadGlobalAttributesRequest = createLoadGlobalAttributesRequest(registry, prcType, preferredPersistencePolicyName)

	return stateCfg
}

func fromAsyncStateOptionsToAsyncStateConfig(
	stateOptions *AsyncStateOptions,
) (*string, *xdbapi.AsyncStateConfig) {
	stateCfg := &xdbapi.AsyncStateConfig{}
	if stateOptions == nil {
		return nil, stateCfg
	}

	stateCfg.WaitUntilApiTimeoutSeconds = &stateOptions.WaitUntilTimeoutSeconds
	stateCfg.ExecuteApiTimeoutSeconds = &stateOptions.ExecuteTimeoutSeconds
	stateCfg.WaitUntilApiRetryPolicy = stateOptions.WaitUntilRetryPolicy
	stateCfg.ExecuteApiRetryPolicy = stateOptions.ExecuteRetryPolicy
	stateCfg.StateFailureRecoveryOptions = stateOptions.FailureRecoveryOptions
	return stateOptions.PersistenceLoadingPolicyName, stateCfg
}

func createLoadGlobalAttributesRequest(
	registry Registry, prcType string, preferredPersistencePolicyName *string,
) *xdbapi.LoadGlobalAttributesRequest {
	persistenceSchema := registry.getPersistenceSchema(prcType)
	if persistenceSchema.GlobalAttributeSchema != nil {
		keyToDefs := registry.getGlobalAttributeKeyToDefs(prcType)
		persistencePolicy := persistenceSchema.DefaultLoadingPolicy.GlobalAttributeLoadingPolicy
		if preferredPersistencePolicyName != nil {
			policy, ok := persistenceSchema.NamedLoadingPolicies[*preferredPersistencePolicyName]
			if !ok {
				panic("persistence loading policy not found " + *preferredPersistencePolicyName)
			}
			persistencePolicy = policy.GlobalAttributeLoadingPolicy
		}
		return convertGlobalAttributeLoadingPolicyToLoadingRequest(
			keyToDefs,
			persistencePolicy,
		)
	}
	return nil
}

func convertGlobalAttributeLoadingPolicyToLoadingRequest(
	keyToDefs map[string]internalGlobalAttrDef,
	policy *GlobalAttributeLoadingPolicy,
) *xdbapi.LoadGlobalAttributesRequest {
	var attrs []xdbapi.GlobalAttributeKey
	for _, k := range policy.LoadingKeys {
		def := keyToDefs[k]
		attr := xdbapi.GlobalAttributeKey{
			DbColumn:                         def.colName,
			AlternativeTable:                 def.altTableName,
			AlternativeTableForeignKeyColumn: def.altTableForeignKey,
		}
		attrs = append(attrs, attr)
	}

	var tableReadLockingPolicyOverrides []xdbapi.TableReadLockingPolicy
	for tbl, lockType := range policy.TableLockingTypeOverrides {
		tableReadLockingPolicyOverrides = append(tableReadLockingPolicyOverrides, xdbapi.TableReadLockingPolicy{
			TableName:       tbl,
			ReadLockingType: lockType,
		})
	}
	return &xdbapi.LoadGlobalAttributesRequest{
		Attributes:                      attrs,
		DefaultReadLockingType:          policy.TableLockingTypeDefault.Ptr(),
		TableReadLockingPolicyOverrides: tableReadLockingPolicyOverrides,
	}
}
