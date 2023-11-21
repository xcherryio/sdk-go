package xc

import (
	"github.com/xcherryio/apis/goapi/xcapi"
)

type ObjectEncoder interface {
	// GetEncodingType returns the encoding info that it can handle
	GetEncodingType() string
	// Encode serialize an object
	Encode(obj interface{}) (*xcapi.EncodedObject, error)
	// Decode deserialize an object
	Decode(encodedObj *xcapi.EncodedObject, resultPtr interface{}) error
}
