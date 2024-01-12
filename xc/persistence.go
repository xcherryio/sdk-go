package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type Persistence interface {
	// GetLocalAttribute returns the local attribute value
	GetLocalAttribute(key string, resultPtr interface{})
	// SetLocalAttribute sets the local attribute value
	SetLocalAttribute(key string, value interface{})

	// getLocalAttributesToReturn returns the local attributes to update
	getLocalAttributesToUpdate() []xcapi.KeyValue
}
