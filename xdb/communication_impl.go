package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type communicationImpl struct {
	encoder                     ObjectEncoder
	localQueueMessagesToPublish []xdbapi.LocalQueueMessage
}

func NewCommunication(encoder ObjectEncoder) Communication {
	return &communicationImpl{}
}

func (c communicationImpl) PublishToLocalQueue(queueName string, payload interface{}) {
	pl, err := c.encoder.Encode(payload)
	if err != nil {
		panic(err)
	}
	msg := xdbapi.LocalQueueMessage{
		QueueName: queueName,
		Payload:   pl,
	}
	c.localQueueMessagesToPublish = append(c.localQueueMessagesToPublish, msg)
}

func (c communicationImpl) GetLocalQueueMessagesToPublish() []xdbapi.LocalQueueMessage {
	return c.localQueueMessagesToPublish
}
