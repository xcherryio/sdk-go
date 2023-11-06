package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type BasicClientProcessOptions struct {
	ProcessIdReusePolicy *xdbapi.ProcessIdReusePolicy
	StartStateOptions    *xdbapi.AsyncStateConfig
	// Default: 10 seconds when set as 0
	TimeoutSeconds        int32
	GlobalAttributeConfig *xdbapi.GlobalAttributeConfig
}
