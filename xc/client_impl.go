package xc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/xc/ptr"
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
	if persSchema.AppDatabaseSchema != nil {
		if startOptions == nil || startOptions.GlobalAttributeOptions == nil ||
			len(startOptions.GlobalAttributeOptions.DBTableConfigs) == 0 {
			return "", NewInvalidArgumentError("GlobalAttributeConfig is required for process with AppDatabaseSchema")
		}
		dbTableCfgs := startOptions.GlobalAttributeOptions.DBTableConfigs
		schema := persSchema.AppDatabaseSchema
		keyDefs := c.registry.getGlobalAttributeKeyToDefs(prcType)

		tableConfigs, err := c.convertToTableConfig(dbTableCfgs, *schema, keyDefs)
		if err != nil {
			return "", err
		}

		unregOpt.GlobalAttributeConfig = &xcapi.GlobalAttributeConfig{
			TableConfigs: tableConfigs,
		}
	}
	if persSchema.LocalAttributeSchema != nil {
		if startOptions != nil && len(startOptions.InitialLocalAttribute) > 0 {
			var initialWrite []xcapi.KeyValue
			for key, attr := range startOptions.InitialLocalAttribute {
				if _, ok := persSchema.LocalAttributeSchema.LocalAttributeKeys[key]; !ok {
					return "", NewInvalidArgumentError("invalid attribute key for local attribute schema: " + key)
				}

				encodedValPtr, err := GetDefaultObjectEncoder().Encode(attr)
				if err != nil {
					return "", err
				}

				initialWrite = append(initialWrite, *xcapi.NewKeyValue(key, *encodedValPtr))
			}

			unregOpt.LocalAttributeConfig = &xcapi.LocalAttributeConfig{
				InitialWrite: initialWrite,
			}
		}
	}

	unregOpt.ProcessIdReusePolicy = prcOptions.IdReusePolicy
	unregOpt.TimeoutSeconds = prcOptions.TimeoutSeconds

	return c.BasicClient.StartProcess(ctx, prcType, startStateId, processId, input, unregOpt)
}

func (c *clientImpl) PublishToLocalQueue(
	ctx context.Context, processId string, queueName string, payload interface{}, options *LocalQueuePublishOptions,
) error {
	msg, err := c.convertToAPIMessage(queueName, payload, options)
	if err != nil {
		return err
	}
	return c.BasicClient.PublishToLocalQueue(ctx, processId, []xcapi.LocalQueueMessage{msg})
}

func (c *clientImpl) BatchPublishToLocalQueue(
	ctx context.Context, processId string, messages ...LocalQueuePublishMessage,
) error {
	var msgs []xcapi.LocalQueueMessage
	for _, m := range messages {
		msg, err := c.convertToAPIMessage(m.QueueName, m.Payload, &LocalQueuePublishOptions{
			DedupSeed: m.DedupSeed,
			DedupUUID: m.DedupUUID,
		})
		if err != nil {
			return err
		}
		msgs = append(msgs, msg)
	}
	return c.BasicClient.PublishToLocalQueue(ctx, processId, msgs)
}

func (c *clientImpl) convertToAPIMessage(
	queueName string, payload interface{}, options *LocalQueuePublishOptions,
) (xcapi.LocalQueueMessage, error) {
	var pl *xcapi.EncodedObject
	var err error
	if payload != nil {
		pl, err = c.clientOptions.ObjectEncoder.Encode(payload)
		if err != nil {
			return xcapi.LocalQueueMessage{}, err
		}
	}
	msg := xcapi.LocalQueueMessage{
		QueueName: queueName,
		Payload:   pl,
	}
	if options != nil {
		if options.DedupSeed != nil {
			guid := uuid.NewMD5(uuid.NameSpaceOID, []byte(*options.DedupSeed))
			msg.DedupId = ptr.Any(guid.String())
		}
		if options.DedupUUID != nil {
			_, err := uuid.Parse(*options.DedupUUID)
			if err != nil {
				return xcapi.LocalQueueMessage{}, fmt.Errorf("invalid DedupUUId %v , err: %w", *options.DedupUUID, err)
			}
			msg.DedupId = options.DedupUUID
		}
	}
	return msg, nil
}

func (c *clientImpl) StartProcess(
	ctx context.Context, definition Process, processId string, input interface{},
) (string, error) {
	return c.StartProcessWithOptions(ctx, definition, processId, input, nil)
}

func (c *clientImpl) StopProcess(
	ctx context.Context, processId string, stopType xcapi.ProcessExecutionStopType,
) error {
	return c.BasicClient.StopProcess(ctx, processId, stopType)
}

func (c *clientImpl) DescribeCurrentProcessExecution(
	ctx context.Context, processId string,
) (*xcapi.ProcessExecutionDescribeResponse, error) {
	return c.BasicClient.DescribeCurrentProcessExecution(ctx, processId)
}

func (c *clientImpl) convertToTableConfig(
	dbCfgs map[string]AppDatabaseTableOptions, schema AppDatabaseSchema, keyToDefs map[string]internalColumnDef,
) ([]xcapi.GlobalAttributeTableConfig, error) {
	var res []xcapi.GlobalAttributeTableConfig
	for _, tbl := range schema.Tables {
		tblName := tbl.TableName
		cfg, ok := dbCfgs[tblName]
		if !ok {
			return nil, NewInvalidArgumentError("GlobalAttributeConfig.tables missing table: " + tblName)
		}
		dbVal, err := c.dbConverter.ToDBValue(cfg.PKValue, cfg.PKHint)
		if err != nil {
			return nil, err
		}
		var initWrite []xcapi.TableColumnValue
		for key, attr := range cfg.InitialAttributes {
			def, ok := keyToDefs[key]
			if !ok {
				return nil, NewInvalidArgumentError("invalid attribute key for global attribute schema: " + key)
			}
			dbVal, err := c.dbConverter.ToDBValue(attr, def.colDef.Hint)
			if err != nil {
				return nil, err
			}
			initWrite = append(initWrite, xcapi.TableColumnValue{
				DbColumn:     def.colDef.ColumnName,
				DbQueryValue: dbVal,
			})
		}
		tblConfig := xcapi.GlobalAttributeTableConfig{
			TableName: tblName,
			PrimaryKey: xcapi.TableColumnValue{
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
