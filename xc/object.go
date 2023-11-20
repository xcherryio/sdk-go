package xc

import "github.com/xcherryio/apis/goapi/xcapi"

// Object is a representation of EncodedObject
type Object interface {
	Get(resultPtr interface{})
}

type objectImpl struct {
	encodedObject *xcapi.EncodedObject
	objectEncoder ObjectEncoder
}

func NewObject(EncodedObject *xcapi.EncodedObject, ObjectEncoder ObjectEncoder) Object {
	return objectImpl{
		encodedObject: EncodedObject,
		objectEncoder: ObjectEncoder,
	}
}

// Get retrieves the actual object
// It just panics on error but the error can still be accessible if really need to do some customized handling(mostly you don't need to):
// 1. capturing panic yourself
// 2. get the error from WorkerService API, because WorkerService will use captureStateExecutionError to capture the error
func (o objectImpl) Get(resultPtr interface{}) {
	err := o.objectEncoder.Decode(o.encodedObject, resultPtr)
	if err != nil {
		panic(err)
	}
}
