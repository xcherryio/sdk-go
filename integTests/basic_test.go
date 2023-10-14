package integTests

import (
	"github.com/xdblab/xdb-golang-sdk/integTests/basic"
	"github.com/xdblab/xdb-golang-sdk/integTests/multi_states"
	"github.com/xdblab/xdb-golang-sdk/integTests/state_decision"
	"testing"
)

func TestIOProcess(t *testing.T) {
	basic.TestStartIOProcess(t, client)
}

func TestStopProcess(t *testing.T) {
	multi_states.TestTerminateMultiStatesProcess(t, client)
	multi_states.TestFailMultiStatesProcess(t, client)
}

func TestStateDecision(t *testing.T) {
	state_decision.TestGracefulCompleteProcess(t, client)
	state_decision.TestForceCompleteProcess(t, client)
	state_decision.TestForceFailProcess(t, client)
	state_decision.TestDeadEndProcess(t, client)
}
