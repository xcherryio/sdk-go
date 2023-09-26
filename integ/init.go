package integ

import "github.com/xdblab/xdb-golang-sdk/xdb"

var registry = xdb.NewRegistry()
var client = xdb.NewClient(registry, nil)
var workerService = xdb.NewWorkerService(registry, nil)

func init() {
	err := registry.AddProcesses(
		&basicWorkflow{},
	)
	if err != nil {
		panic(err)
	}
}
