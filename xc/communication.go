package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type Communication interface {
	// PublishToLocalQueue publishes a message to a local queue
	// the payload can be empty(nil)
	PublishToLocalQueue(queueName string, payload interface{})

	// below is for internal implementation
	communicationInternal
}

type communicationInternal interface {
	GetLocalQueueMessagesToPublish() []xcapi.LocalQueueMessage
}
