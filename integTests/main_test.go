package integTests

import (
	"fmt"
	"github.com/xcherryio/sdk-go/integTests/worker"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("start running integ test")
	closeFn := worker.StartGinWorker(workerService)
	code := m.Run()
	closeFn()
	fmt.Println("finished running integ test with status code", code)
	os.Exit(code)
}
