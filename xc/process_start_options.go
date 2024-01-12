package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type ProcessStartOptions struct {
	// TimeoutSeconds is the timeout for the process execution.
	// Default: 0, mean which means infinite timeout.
	// This will override the timeout defined in process definition
	TimeoutSeconds *int32
	// IdReusePolicy is the policy for reusing process id.
	// Default: xcapi.ALLOW_IF_NO_RUNNING when set as nil.
	// This will override the IdReusePolicy defined in process definition.
	IdReusePolicy *xcapi.ProcessIdReusePolicy
	// InitialLocalAttribute is the initial local attributes to be set when starting the process execution
	InitialLocalAttribute map[string]interface{}
}
