package integTests

import (
	"github.com/xdblab/xdb-golang-sdk/integTests/basic"
	"github.com/xdblab/xdb-golang-sdk/integTests/multi_states"
	"github.com/xdblab/xdb-golang-sdk/integTests/state_decision"
	"github.com/xdblab/xdb-golang-sdk/xdb"
)

var registry = xdb.NewRegistry()
var client = xdb.NewClient(registry, nil)
var workerService = xdb.NewWorkerService(registry, nil)

func init() {
	err := registry.AddProcesses(
		&basic.IOProcess{},
		&multi_states.MultiStatesProcess{},
		&state_decision.GracefulCompleteProcess{},
		&state_decision.ForceCompleteProcess{},
		&state_decision.ForceFailProcess{},
		&state_decision.DeadEndProcess{},
	)
	if err != nil {
		panic(err)
	}
}
