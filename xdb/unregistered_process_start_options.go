package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type UnregisteredProcessOptions struct {
	ProcessIdReusePolicy *xdbapi.ProcessIdReusePolicy
	StartStateOptions    *xdbapi.AsyncStateConfig
	TimeoutSeconds       int32
}
