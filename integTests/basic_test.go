package integTests

import (
	"github.com/xdblab/xdb-golang-sdk/integTests/basic"
	"github.com/xdblab/xdb-golang-sdk/integTests/multi_states"
	"testing"
)

func TestIOProcess(t *testing.T) {
	basic.TestStartIOProcess(t, client)
}

func TestStopMultiStatesProcess(t *testing.T) {
	multi_states.TestTerminateMultiStatesProcess(t, client)
	multi_states.TestFailMultiStatesProcess(t, client)
}
