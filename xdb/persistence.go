package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type Persistence interface {
	// GetGlobalAttribute returns the global attribute value
	GetGlobalAttribute(key string, resultPtr interface{})
	// SetGlobalAttribute sets the global attribute value
	SetGlobalAttribute(key string, value interface{})
	// GetGlobalAttributesToReturn returns the global attributes to update
	getGlobalAttributesToUpdate() []xdbapi.GlobalAttributeValue
}
