package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type BasicClientProcessOptions struct {
	ProcessIdReusePolicy *xdbapi.ProcessIdReusePolicy
	StartStateOptions    *xdbapi.AsyncStateConfig
	// default is 0 which indicate no timeout
	TimeoutSeconds        int32
	GlobalAttributeConfig *xdbapi.GlobalAttributeConfig
}
