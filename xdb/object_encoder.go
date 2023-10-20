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
	// FromGlobalAttributeToDbValue converts global attribute value to database query value
	FromGlobalAttributeToDbValue(val interface{}, hint *string) (dbValue string, err error)
	// FromDbValueToGlobalAttribute converts database query value to global attribute value
	FromDbValueToGlobalAttribute(dbQueryValue string, hint *string, resultPtr interface{}) error
}
