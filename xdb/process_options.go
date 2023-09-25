package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type ProcessOptions struct {
	IdReusePolicy  *xdbapi.ProcessIdReusePolicy
	TimeoutSeconds int32
}
