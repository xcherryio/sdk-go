package integTests

import (
	"fmt"
	"os"
	"testing"

	"github.com/xdblab/xdb-golang-sdk/integTests/worker"
)

func TestMain(m *testing.M) {
	fmt.Println("start running integ test")
	closeFn := worker.StartGinWorker(workerService)
	code := m.Run()
	closeFn()
	fmt.Println("finished running integ test with status code", code)
	os.Exit(code)
}
