package integTests

import (
	"github.com/xcherryio/sdk-go/integTests/basic"
	"github.com/xcherryio/sdk-go/integTests/command_request"
	"github.com/xcherryio/sdk-go/integTests/failure_recovery"
	"github.com/xcherryio/sdk-go/integTests/local_attribute"
	"github.com/xcherryio/sdk-go/integTests/multi_states"
	"github.com/xcherryio/sdk-go/integTests/process_timeout"
	"github.com/xcherryio/sdk-go/integTests/state_decision"
	"github.com/xcherryio/sdk-go/integTests/stateretry"
	"github.com/xcherryio/sdk-go/xc"
)

var registry = xc.NewRegistry()
var client = xc.NewClient(registry, getTestClientOptions())
var workerService = xc.NewWorkerService(registry, nil)

func init() {
	err := registry.AddProcesses(
		&basic.IOProcess{},
		&failure_recovery.StateFailureRecoveryTestExecuteProcess{},
		&failure_recovery.StateFailureRecoveryTestWaitUntilProcess{},
		&failure_recovery.StateFailureRecoveryTestExecuteNoWaitUntilProcess{},
		&failure_recovery.StateFailureRecoveryTestExecuteFailedAtStartProcess{},
		&multi_states.MultiStatesProcess{},
		&local_attribute.LocalAttributeTestProcess{},
		&state_decision.GracefulCompleteProcess{},
		&state_decision.ForceCompleteProcess{},
		&state_decision.ForceFailProcess{},
		&state_decision.DeadEndProcess{},
		&stateretry.BackoffProcess{},
		&process_timeout.TimeoutProcess{},
		&command_request.AnyOfTimerLocalQProcess{},
		&command_request.AllOfTimerLocalQProcess{},
	)
	if err != nil {
		panic(err)
	}
}

func getTestClientOptions() *xc.ClientOptions {
	options := xc.GetLocalDefaultClientOptions()
	options.DefaultProcessTimeoutSecondsOverride = 10
	return options
}
