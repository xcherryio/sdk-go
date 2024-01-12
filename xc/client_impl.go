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
