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

type SingleTableProcess struct {
	xdb.ProcessDefaults
}

const (
	attrKeyInt            = "key1"
	attrKeyStr            = "key2"
	loadNothingPolicyName = "loadNothing"
	tblName               = "sample_user_table"
	pk                    = "user_id"
)

func (b SingleTableProcess) GetPersistenceSchema() xdb.PersistenceSchema {
	return xdb.NewPersistenceSchemaWithOptions(
		xdb.NewGlobalAttributesSchema(
			xdb.NewDBTableSchema(
				tblName, pk,
				xdbapi.NO_LOCKING,
				xdb.NewDBColumnDef(attrKeyInt, "create_timestamp", true),
				xdb.NewDBColumnDef(attrKeyStr, "first_name", true)),
		),
		nil,
		xdb.NewPersistenceSchemaOptions(
			xdb.NewNamedPersistenceLoadingPolicy(
				loadNothingPolicyName, nil,
				xdb.NewTableLoadingPolicy(tblName, xdbapi.NO_LOCKING)),
		),
	)
}

func (b SingleTableProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(
		&stateForInitialReadWrite{}, // read from initial global attributes and write to them
		&stateToVerifyGlobalAttrs{}, // verify the global attributes write from the prev state
		&stateForTestLoadNothing{})  // test loading nothing policy
}

type stateForInitialReadWrite struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b stateForInitialReadWrite) Execute(
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

	return xdb.SingleNextState(stateToVerifyGlobalAttrs{}, i+1), nil
}

type stateToVerifyGlobalAttrs struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b stateToVerifyGlobalAttrs) Execute(
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

	return xdb.SingleNextState(stateForTestLoadNothing{}, nil), nil
}

type stateForTestLoadNothing struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b stateForTestLoadNothing) GetStateOptions() *xdb.AsyncStateOptions {
	return &xdb.AsyncStateOptions{
		PersistenceLoadingPolicyName: ptr.Any(loadNothingPolicyName),
	}
}

func (b stateForTestLoadNothing) Execute(
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
	prc := SingleTableProcess{}

	runId1, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.RETURN_ERROR_ON_CONFLICT,
		&xdb.ProcessStartOptions{
			GlobalAttributeOptions: &xdb.GlobalAttributeOptions{
				DBTableConfigs: map[string]xdb.DBTableConfig{
					tblName: xdb.DBTableConfig{
						PKValue: prcId, // use processId as the primary key value(string)
						InitialAttributes: map[string]interface{}{
							attrKeyInt: 123,
							attrKeyStr: "abc",
						},
						InitialWriteConflictMode: xdbapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
					},
				},
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
				DBTableConfigs: map[string]xdb.DBTableConfig{
					tblName: xdb.DBTableConfig{
						PKValue: prcId, // use processId as the primary key value(string)
						InitialAttributes: map[string]interface{}{
							attrKeyInt: 123,
							attrKeyStr: "abc",
						},
						InitialWriteConflictMode: xdbapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
					},
				},
			},
		})
	assert.NotNil(t, err)
	assert.True(t, xdb.IsGlobalAttributeWriteFailure(err))

	// failed when trying to start the same process when writing str to int
	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.RETURN_ERROR_ON_CONFLICT,
		&xdb.ProcessStartOptions{
			IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: &xdb.GlobalAttributeOptions{
				DBTableConfigs: map[string]xdb.DBTableConfig{
					tblName: xdb.DBTableConfig{
						PKValue: prcId, // use processId as the primary key value(string)
						InitialAttributes: map[string]interface{}{
							attrKeyInt: "abc",
							attrKeyStr: "123",
						},
						InitialWriteConflictMode: xdbapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
					},
				},
			},
		})
	assert.NotNil(t, err)
	assert.True(t, xdb.IsGlobalAttributeWriteFailure(err))

	// succeeded when trying to start the same process with override
	runId2, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.OVERRIDE_ON_CONFLICT,
		&xdb.ProcessStartOptions{
			IdReusePolicy: xdbapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: &xdb.GlobalAttributeOptions{
				DBTableConfigs: map[string]xdb.DBTableConfig{
					tblName: xdb.DBTableConfig{
						PKValue: prcId, // use processId as the primary key value(string)
						InitialAttributes: map[string]interface{}{
							attrKeyInt: 123,
							attrKeyStr: "abc",
						},
						InitialWriteConflictMode: xdbapi.OVERRIDE_ON_CONFLICT.Ptr(),
					},
				},
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
				DBTableConfigs: map[string]xdb.DBTableConfig{
					tblName: xdb.DBTableConfig{
						PKValue: prcId, // use processId as the primary key value(string)
						InitialAttributes: map[string]interface{}{
							attrKeyInt: 123,
							attrKeyStr: "abc",
						},
						InitialWriteConflictMode: xdbapi.IGNORE_CONFLICT.Ptr(),
					},
				},
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
