package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type Communication interface {
	// PublishToLocalQueue publishes a message to a local queue
	// the payload can be empty(nil)
	PublishToLocalQueue(queueName string, payload interface{})

	// below is for internal implementation
	communicationInternal
}

type communicationInternal interface {
	GetLocalQueueMessagesToPublish() []xdbapi.LocalQueueMessage
}
