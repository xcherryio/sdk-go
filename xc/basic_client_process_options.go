package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type BasicClientProcessOptions struct {
	ProcessIdReusePolicy *xcapi.ProcessIdReusePolicy
	StartStateOptions    *xcapi.AsyncStateConfig
	// default is 0 which indicate no timeout
	TimeoutSeconds        int32
	GlobalAttributeConfig *xcapi.GlobalAttributeConfig
	LocalAttributeConfig  *xcapi.LocalAttributeConfig
}
