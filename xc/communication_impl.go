package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type communicationImpl struct {
	encoder                     ObjectEncoder
	localQueueMessagesToPublish []xcapi.LocalQueueMessage
}

func NewCommunication(encoder ObjectEncoder) Communication {
	return &communicationImpl{
		encoder:                     encoder,
		localQueueMessagesToPublish: nil,
	}
}

func (c *communicationImpl) PublishToLocalQueue(queueName string, payload interface{}) {
	pl, err := c.encoder.Encode(payload)
	if err != nil {
		panic(err)
	}
	msg := xcapi.LocalQueueMessage{
		QueueName: queueName,
		Payload:   pl,
	}
	c.localQueueMessagesToPublish = append(c.localQueueMessagesToPublish, msg)
}

func (c *communicationImpl) GetLocalQueueMessagesToPublish() []xcapi.LocalQueueMessage {
	return c.localQueueMessagesToPublish
}
