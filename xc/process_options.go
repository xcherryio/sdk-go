package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type ProcessOptions struct {
	// TimeoutSeconds is the timeout for the process execution.
	// Default: 0, mean which means infinite timeout
	TimeoutSeconds int32
	// IdReusePolicy is the policy for reusing process id.
	// Default: xcapi.ALLOW_IF_NO_RUNNING when set as nil
	IdReusePolicy *xcapi.ProcessIdReusePolicy
	// GetProcessType defines the processType of this process definition.
	// GetFinalProcessType set the default value when return empty string --
	// It's the packageName.structName of the process instance and ignores the import paths and aliases.
	// e.g. if the process is from myStruct{} under mywf package, the simple name is just "mywf.myStruct". Underneath, it's from reflect.TypeOf(wf).String().
	ProcessType string
}

func NewDefaultProcessOptions() ProcessOptions {
	return ProcessOptions{
		TimeoutSeconds: 0,
		IdReusePolicy:  nil,
		ProcessType:    "",
	}
}
