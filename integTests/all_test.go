package integTests

import (
	"github.com/xdblab/xdb-golang-sdk/integTests/command_request"
	"testing"

	"github.com/xdblab/xdb-golang-sdk/integTests/global_attribute"
	"github.com/xdblab/xdb-golang-sdk/integTests/process_timeout"

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

func TestStateFailureRecoveryExecuteProcess(t *testing.T) {
	failure_recovery.TestStateFailureRecoveryTestExecuteProcess(t, client)
}

func TestStateFailureRecoveryWaitUntilProcess(t *testing.T) {
	failure_recovery.TestStateFailureRecoveryTestWaitUntilProcess(t, client)
}

func TestStateFailureRecoveryExecuteNoWaitUntilProcess(t *testing.T) {
	failure_recovery.TestStateFailureRecoveryTestExecuteNoWaitUntilProcess(t, client)
}

func TestStateFailureRecoveryExecuteFailedAtStartProcess(t *testing.T) {
	failure_recovery.TestStateFailureRecoveryTestExecuteFailedAtStartProcess(t, client)
}

func TestGlobalAttributesWithSingleTable(t *testing.T) {
	global_attribute.TestGlobalAttributesWithSingleTable(t, client)
}

func TestGlobalAttributesWithMultiTables(t *testing.T) {
	global_attribute.TestGlobalAttributesWithMultiTables(t, client)
}

func TestStartTimeoutProcessCase1(t *testing.T) {
	process_timeout.TestStartTimeoutProcessCase1(t, client)
}

func TestStartTimeoutProcessCase2(t *testing.T) {
	process_timeout.TestStartTimeoutProcessCase2(t, client)
}

func TestStartTimeoutProcessCase3(t *testing.T) {
	process_timeout.TestStartTimeoutProcessCase3(t, client)
}

func TestStartTimeoutProcessCase4(t *testing.T) {
	process_timeout.TestStartTimeoutProcessCase4(t, client)
}

func TestAnyOfTimerLocalQueueWithTimerFired(t *testing.T) {
	command_request.TestAnyOfTimerLocalQueueWithTimerFired(t, client)
}

func TestAnyOfTimerLocalQueueWithLocalQueueMessagesReceived(t *testing.T) {
	command_request.TestAnyOfTimerLocalQueueWithLocalQueueMessagesReceived(t, client)
}

func TestAllOfTimerLocalQueue(t *testing.T) {
	command_request.TestAllOfTimerLocalQueue(t, client)
}
