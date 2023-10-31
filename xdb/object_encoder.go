package xdb

import (
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

type ObjectEncoder interface {
	// GetEncodingType returns the encoding info that it can handle
	GetEncodingType() string
	// Encode serialize an object
	Encode(obj interface{}) (*xdbapi.EncodedObject, error)
	// Decode deserialize an object
	Decode(encodedObj *xdbapi.EncodedObject, resultPtr interface{}) error
}
