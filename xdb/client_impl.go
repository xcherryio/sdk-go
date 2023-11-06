package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
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
		if startOptions == nil || startOptions.GlobalAttributeOptions == nil ||
			len(startOptions.GlobalAttributeOptions.DBTableConfigs) == 0 {
			return "", NewInvalidArgumentError("GlobalAttributeConfig is required for process with GlobalAttributeSchema")
		}
		dbTableCfgs := startOptions.GlobalAttributeOptions.DBTableConfigs
		schema := persSchema.GlobalAttributeSchema
		keyDefs := c.registry.getGlobalAttributeKeyToDefs(prcType)

		tableConfigs, err := c.convertToTableConfig(dbTableCfgs, *schema, keyDefs)
		if err != nil {
			return "", err
		}

		unregOpt.GlobalAttributeConfig = &xdbapi.GlobalAttributeConfig{
			TableConfigs: tableConfigs,
		}
	}

	unregOpt.ProcessIdReusePolicy = prcOptions.IdReusePolicy
	unregOpt.TimeoutSeconds = prcOptions.TimeoutSeconds

	return c.BasicClient.StartProcess(ctx, prcType, startStateId, processId, input, unregOpt)
}

func (c *clientImpl) PublishToLocalQueue(
	ctx context.Context, processId string, queueName string, payload interface{}, dedupUUID string,
) error {
	var pl *xdbapi.EncodedObject
	var err error
	if payload != nil {
		pl, err = c.clientOptions.ObjectEncoder.Encode(payload)
		if err != nil {
			return err
		}
	}
	msg := xdbapi.LocalQueueMessage{
		QueueName: queueName,
		Payload:   pl,
	}
	if dedupUUID != "" {
		msg.DedupId = ptr.Any(dedupUUID)
	}

	return c.BasicClient.PublishMessagesToLocalQueue(ctx, processId, []xdbapi.LocalQueueMessage{msg})
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

func (c *clientImpl) convertToTableConfig(
	dbCfgs map[string]DBTableConfig, schema GlobalAttributesSchema, keyToDefs map[string]internalGlobalAttrDef,
) ([]xdbapi.GlobalAttributeTableConfig, error) {
	var res []xdbapi.GlobalAttributeTableConfig
	for _, tbl := range schema.Tables {
		tblName := tbl.TableName
		cfg, ok := dbCfgs[tblName]
		if !ok {
			return nil, NewInvalidArgumentError("GlobalAttributeConfig.DBTableConfigs missing table: " + tblName)
		}
		dbVal, err := c.dbConverter.ToDBValue(cfg.PKValue, cfg.PKHint)
		if err != nil {
			return nil, err
		}
		var initWrite []xdbapi.TableColumnValue
		for key, attr := range cfg.InitialAttributes {
			def, ok := keyToDefs[key]
			if !ok {
				return nil, NewInvalidArgumentError("invalid attribute key for global attribute schema: " + key)
			}
			dbVal, err := c.dbConverter.ToDBValue(attr, def.colDef.Hint)
			if err != nil {
				return nil, err
			}
			initWrite = append(initWrite, xdbapi.TableColumnValue{
				DbColumn:     def.colDef.ColumnName,
				DbQueryValue: dbVal,
			})
		}
		tblConfig := xdbapi.GlobalAttributeTableConfig{
			TableName: tblName,
			PrimaryKey: xdbapi.TableColumnValue{
				DbColumn:     tbl.PK,
				DbQueryValue: dbVal,
			},
			InitialWrite:     initWrite,
			InitialWriteMode: cfg.InitialWriteConflictMode,
		}
		res = append(res, tblConfig)
	}
	return res, nil
}
