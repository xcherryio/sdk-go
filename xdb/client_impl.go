package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

type clientImpl struct {
	BasicClient
	dbConverter   DBConverter
	clientOptions ClientOptions
	registry      Registry
}

func (c *clientImpl) GetBasicClient() BasicClient {
	return c.BasicClient
}

func (c *clientImpl) StartProcessWithOptions(
	ctx context.Context, definition Process, processId string, input interface{}, startOptions *ProcessStartOptions,
) (string, error) {
	prcType := GetFinalProcessType(definition)
	prc := c.registry.getProcess(prcType)
	if prc == nil {
		return "", NewInvalidArgumentError("Process is not registered")
	}
	persSchema := c.registry.getPersistenceSchema(prcType)

	state := c.registry.getProcessStartingState(prcType)

	unregOpt := &BasicClientProcessOptions{}

	startStateId := ""
	if state != nil {
		startStateId = GetFinalStateId(state)
		unregOpt.StartStateOptions = fromStateToAsyncStateConfig(state, prcType, c.registry)
	}

	prcOptions := prc.GetProcessOptions()
	if startOptions != nil {
		// override the process startOptions
		if startOptions.IdReusePolicy != nil {
			prcOptions.IdReusePolicy = startOptions.IdReusePolicy
		}
		if startOptions.TimeoutSeconds != nil {
			prcOptions.TimeoutSeconds = *startOptions.TimeoutSeconds
		}
	}
	if persSchema.GlobalAttributeSchema != nil {
		if startOptions == nil || startOptions.GlobalAttributeOptions == nil {
			return "", NewInvalidArgumentError("GlobalAttributeConfig is required for process with GlobalAttributeSchema")
		}
		gaOptions := startOptions.GlobalAttributeOptions
		schema := persSchema.GlobalAttributeSchema
		gaDefs := c.registry.getGlobalAttributeKeyToDefs(prcType)
		pkValue, err := c.dbConverter.ToDBValue(gaOptions.PrimaryAttributeValue, schema.DefaultTablePrimaryKeyHint)
		if err != nil {
			return "", err
		}

		gaVals, err := c.convertGlobalAttributeValues(gaOptions.InitialAttributes, gaDefs)
		if err != nil {
			return "", err
		}

		unregOpt.GlobalAttributeConfig = &xdbapi.GlobalAttributeConfig{
			DefaultDbTable: &schema.DefaultTableName,
			PrimaryGlobalAttribute: &xdbapi.GlobalAttributeValue{
				DbColumn:     schema.DefaultTablePrimaryKey,
				DbQueryValue: pkValue,
			},
			InitialGlobalAttributes:         gaVals,
			InitialGlobalAttributeWriteMode: &gaOptions.InitialWriteConflictMode,
		}
	}

	unregOpt.ProcessIdReusePolicy = prcOptions.IdReusePolicy
	unregOpt.TimeoutSeconds = prcOptions.TimeoutSeconds

	return c.BasicClient.StartProcess(ctx, prcType, startStateId, processId, input, unregOpt)
}

func (c *clientImpl) StartProcess(
	ctx context.Context, definition Process, processId string, input interface{},
) (string, error) {
	return c.StartProcessWithOptions(ctx, definition, processId, input, nil)
}

func (c *clientImpl) StopProcess(
	ctx context.Context, processId string, stopType xdbapi.ProcessExecutionStopType,
) error {
	return c.BasicClient.StopProcess(ctx, processId, stopType)
}

func (c *clientImpl) DescribeCurrentProcessExecution(
	ctx context.Context, processId string,
) (*xdbapi.ProcessExecutionDescribeResponse, error) {
	return c.BasicClient.DescribeCurrentProcessExecution(ctx, processId)
}

func (c *clientImpl) convertGlobalAttributeValues(
	rawAttributes map[string]interface{}, schema map[string]internalGlobalAttrDef,
) ([]xdbapi.GlobalAttributeValue, error) {
	var vals []xdbapi.GlobalAttributeValue
	for k, v := range rawAttributes {
		attr, ok := schema[k]
		if !ok {
			return nil, NewInvalidArgumentError("GlobalAttributeConfig.InitialAttributes has unknown key: " + k)
		}
		dbValue, err := c.dbConverter.ToDBValue(v, attr.hint)
		if err != nil {
			return nil, err
		}
		vals = append(vals, xdbapi.GlobalAttributeValue{
			DbColumn:                         attr.colName,
			DbQueryValue:                     dbValue,
			AlternativeTable:                 attr.altTableName,
			AlternativeTableForeignKeyColumn: attr.altTableForeignKey,
		})
	}
	return vals, nil
}
