package integTests

import (
	"testing"

	"github.com/xdblab/xdb-golang-sdk/integTests/basic"
	"github.com/xdblab/xdb-golang-sdk/integTests/multi_states"
	"github.com/xdblab/xdb-golang-sdk/integTests/state_decision"
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

func TestProcessIdReusePolicyDisallowReuse(t *testing.T) {
	basic.TestProcessIdReusePolicyDisallowReuse(t, client)
}

func TestProcessIdReusePolicyAllowIfNoRunning(t *testing.T) {
	basic.TestProcessIdReusePolicyAllowIfNoRunning(t, client)
}

func TestProcessIdReusePolicyTerminateIfRunning(t *testing.T) {
	basic.TestProcessIdReusePolicyTerminateIfRunning(t, client)
}

func TestProcessIdReusePolicyAllowIfPreviousExitAbnormally(t *testing.T) {
	basic.TestProcessIdReusePolicyAllowIfPreviousExitAbnormally(t, client)
}
