package global_attribute

import (
	"context"
	"fmt"
	"github.com/xdblab/xdb-golang-sdk/integTests/common"
	"github.com/xdblab/xdb-golang-sdk/xdb/ptr"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb"
)

type IOProcess struct {
	xdb.ProcessDefaults
}

const (
	attrKeyInt            = "key1"
	attrKeyStr            = "key2"
	loadNothingPolicyName = "loadNothing"
)

func (b IOProcess) GetPersistenceSchema() xdb.PersistenceSchema {
	return xdb.NewPersistenceSchema(
		xdb.NewGlobalAttributesSchema(
			"sample-table-1",
			"sample-str-pk",
			xdb.NewGlobalAttributeDef(attrKeyInt, "sample-int-col"),
			xdb.NewGlobalAttributeDef(attrKeyStr, "sample-string-col"),
		),
		nil,
		xdb.NewPersistenceLoadingPolicy(
			xdb.NewGlobalAttributeLoadingPolicy(
				xdbapi.NO_LOCKING,
				attrKeyInt, attrKeyStr,
			),
			nil,
		),
		map[string]xdb.PersistenceLoadingPolicy{
			loadNothingPolicyName: xdb.NewPersistenceLoadingPolicy(
				xdb.NewGlobalAttributeLoadingPolicy(
					xdbapi.NO_LOCKING,
				),
				nil,
			),
		},
	)
}

func (b IOProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(&state1{}, &state2{}, &state3{})
}

type state1 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b state1) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var mode xdbapi.AttributeWriteConflictMode
	input.Get(&mode)
	expectedI := 123
	expectedStr := "abc"
	if mode == xdbapi.OVERRIDE_ON_CONFLICT || mode == xdbapi.IGNORE_CONFLICT {
		expectedI = 123456
		expectedStr = "abcdef"
	}

	var i int
	persistence.GetGlobalAttribute(attrKeyInt, &i)
	var str string
	persistence.GetGlobalAttribute(attrKeyStr, &str)
	if i != expectedI {
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str != expectedStr {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	persistence.SetGlobalAttribute(attrKeyInt, 456)
	persistence.SetGlobalAttribute(attrKeyStr, "def")

	return xdb.SingleNextState(state2{}, i+1), nil
}

type state2 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b state2) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var i int
	persistence.GetGlobalAttribute(attrKeyInt, &i)
	var str string
	persistence.GetGlobalAttribute(attrKeyStr, &str)
	if i != 456 {
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str != "def" {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	return xdb.SingleNextState(state3{}, nil), nil
}

type state3 struct {
	xdb.AsyncStateNoWaitUntil
}

func (b state3) GetStateOptions() *xdb.AsyncStateOptions {
	return &xdb.AsyncStateOptions{
		PersistenceLoadingPolicyName: ptr.Any(loadNothingPolicyName),
	}
}

func (b state3) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var i int
	persistence.GetGlobalAttribute(attrKeyInt, &i)
	var str string
	persistence.GetGlobalAttribute(attrKeyStr, &str)
	if i != 0 {
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str != "" {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	return xdb.ForceCompletingProcess, nil
}

func TestGlobalAttributesWithSingleTable(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := IOProcess{}

	runId1, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.RETURN_ERROR_ON_CONFLICT,
		&xdb.ProcessStartOptions{
			GlobalAttributeOptions: &xdb.GlobalAttributeOptions{
				PrimaryAttributeValue: prcId, // use processId as the primary key value(string)
				InitialAttributes: map[string]interface{}{
					attrKeyInt: 123,
					attrKeyStr: "abc",
				},
				InitialWriteConflictMode: xdbapi.RETURN_ERROR_ON_CONFLICT,
			},
		})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())

	// failed when trying to start the same process again with conflicted global attributes
	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.RETURN_ERROR_ON_CONFLICT,
		&xdb.ProcessStartOptions{
			IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: &xdb.GlobalAttributeOptions{
				PrimaryAttributeValue: prcId,
				InitialAttributes: map[string]interface{}{
					attrKeyInt: 123,
					attrKeyStr: "abc",
				},
				InitialWriteConflictMode: xdbapi.RETURN_ERROR_ON_CONFLICT,
			},
		})
	assert.NotNil(t, err)
	assert.True(t, xdb.IsGlobalAttributeWriteFailure(err))

	// failed when trying to start the same process when writing str to int
	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.RETURN_ERROR_ON_CONFLICT,
		&xdb.ProcessStartOptions{
			IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: &xdb.GlobalAttributeOptions{
				PrimaryAttributeValue: prcId,
				InitialAttributes: map[string]interface{}{
					attrKeyInt: "abc",
					attrKeyStr: 123,
				},
				InitialWriteConflictMode: xdbapi.OVERRIDE_ON_CONFLICT,
			},
		})
	assert.NotNil(t, err)
	assert.True(t, xdb.IsGlobalAttributeWriteFailure(err))

	// succeeded when trying to start the same process with override
	runId2, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.OVERRIDE_ON_CONFLICT,
		&xdb.ProcessStartOptions{
			IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: &xdb.GlobalAttributeOptions{
				PrimaryAttributeValue: prcId,
				InitialAttributes: map[string]interface{}{
					attrKeyInt: 123456,
					attrKeyStr: "abcdef",
				},
				InitialWriteConflictMode: xdbapi.OVERRIDE_ON_CONFLICT,
			},
		})
	assert.Nil(t, err)
	assert.NotEqual(t, runId1, runId2)

	time.Sleep(time.Second * 3)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())

	// succeeded when trying to start the same process with ignore
	runId3, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.IGNORE_CONFLICT,
		&xdb.ProcessStartOptions{
			IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: &xdb.GlobalAttributeOptions{
				PrimaryAttributeValue: prcId,
				InitialAttributes: map[string]interface{}{
					attrKeyInt: 123456789,   // it will be ignored because of conflict
					attrKeyStr: "abcdefefg", // it will be ignored because of conflict
				},
				InitialWriteConflictMode: xdbapi.IGNORE_CONFLICT,
			},
		})
	assert.Nil(t, err)
	assert.NotEqual(t, runId2, runId3)

	time.Sleep(time.Second * 3)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}

// TODO test with multiple/alternative tables
// TODO Test with different locking types
// TODO test using a different loading policy on starting state
