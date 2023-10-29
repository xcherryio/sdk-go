package integTests

import (
	"testing"

	"github.com/xdblab/xdb-golang-sdk/integTests/failure_recovery"
	"github.com/xdblab/xdb-golang-sdk/integTests/stateretry"

	"github.com/xdblab/xdb-golang-sdk/integTests/basic"
	"github.com/xdblab/xdb-golang-sdk/integTests/multi_states"
	"github.com/xdblab/xdb-golang-sdk/integTests/state_decision"
)

func TestIOProcess(t *testing.T) {
	basic.TestStartIOProcess(t, client)
}

func TestStateBackoffRetry(t *testing.T) {
	stateretry.TestBackoff(t, client)
}

func TestTerminateProcess(t *testing.T) {
	multi_states.TestTerminateMultiStatesProcess(t, client)
}

func TestStopProcessByFail(t *testing.T) {
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

func TestProcessIdReusePolicyAllowIfPreviousExitAbnormallyCase1(t *testing.T) {
	basic.TestProcessIdReusePolicyAllowIfPreviousExitAbnormallyCase1(t, client)
}

func TestProcessIdReusePolicyAllowIfPreviousExitAbnormallyCase2(t *testing.T) {
	basic.TestProcessIdReusePolicyAllowIfPreviousExitAbnormallyCase2(t, client)
}

func TestStateFailureRecoveryProcess(t *testing.T) {
	failure_recovery.TestStateFailureRecoveryTestExecuteProcess(t, client)
}

func TestStateFailureRecoveryWaitUntilProcess(t *testing.T) {
	failure_recovery.TestStateFailureRecoveryTestWaitUntilProcess(t, client)
}
