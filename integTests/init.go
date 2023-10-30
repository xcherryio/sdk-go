package integTests

import (
	"github.com/xdblab/xdb-golang-sdk/integTests/basic"
	"github.com/xdblab/xdb-golang-sdk/integTests/failure_recovery"
	"github.com/xdblab/xdb-golang-sdk/integTests/multi_states"
	"github.com/xdblab/xdb-golang-sdk/integTests/state_decision"
	"github.com/xdblab/xdb-golang-sdk/integTests/stateretry"
	"github.com/xdblab/xdb-golang-sdk/xdb"
)

var registry = xdb.NewRegistry()
var client = xdb.NewClient(registry, nil)
var workerService = xdb.NewWorkerService(registry, nil)

func init() {
	err := registry.AddProcesses(
		&basic.IOProcess{},
		&failure_recovery.StateFailureRecoveryTestExecuteProcess{},
		&failure_recovery.StateFailureRecoveryTestWaitUntilProcess{},
		&failure_recovery.StateFailureRecoveryTestExecuteNoWaitUntilProcess{},
		&failure_recovery.StateFailureRecoveryTestExecuteFailedAtStartProcess{},
		&multi_states.MultiStatesProcess{},
		&state_decision.GracefulCompleteProcess{},
		&state_decision.ForceCompleteProcess{},
		&state_decision.ForceFailProcess{},
		&state_decision.DeadEndProcess{},
		&stateretry.BackoffProcess{},
	)
	if err != nil {
		panic(err)
	}
}
