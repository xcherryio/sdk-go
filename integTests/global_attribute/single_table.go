package global_attribute

import (
	"context"
	"fmt"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc/ptr"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/xc"
)

type SingleTableProcess struct {
	xc.ProcessDefaults
}

func (b SingleTableProcess) GetPersistenceSchema() xc.PersistenceSchema {
	return xc.NewPersistenceSchemaWithOptions(
		xc.NewEmptyLocalAttributesSchema(),
		xc.NewGlobalAttributesSchema(
			xc.NewDBTableSchema(
				tblName, pk,
				xcapi.NO_LOCKING,
				xc.NewDBColumnDef(attrKeyInt, "create_timestamp", true),
				xc.NewDBColumnDef(attrKeyStr, "first_name", true)),
		),
		xc.NewPersistenceSchemaOptions(
			xc.NewNamedPersistencePolicy(
				loadNothingPolicyName, nil,
				xc.NewTablePolicy(tblName, xcapi.NO_LOCKING)),
		),
	)
}

func (b SingleTableProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(
		&stateForInitialReadWrite{}, // read from initial global attributes and write to them
		&stateToVerifyGlobalAttrs{}, // verify the global attributes write from the prev state
		&stateForTestLoadNothing{})  // test loading nothing policy
}

type stateForInitialReadWrite struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b stateForInitialReadWrite) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var mode xcapi.AttributeWriteConflictMode
	input.Get(&mode)
	expectedI := 123
	expectedStr := "abc"
	if mode == xcapi.OVERRIDE_ON_CONFLICT {
		expectedI = 123456
		expectedStr = "abcdef"
	}
	if mode == xcapi.IGNORE_CONFLICT {
		// value from last execution
		expectedI = 456
		expectedStr = "def"
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

	return xc.SingleNextState(stateToVerifyGlobalAttrs{}, i+1), nil
}

type stateToVerifyGlobalAttrs struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b stateToVerifyGlobalAttrs) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
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

	return xc.SingleNextState(stateForTestLoadNothing{}, nil), nil
}

type stateForTestLoadNothing struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b stateForTestLoadNothing) GetStateOptions() *xc.AsyncStateOptions {
	return &xc.AsyncStateOptions{
		PersistencePolicyName: ptr.Any(loadNothingPolicyName),
	}
}

func (b stateForTestLoadNothing) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
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

	return xc.ForceCompletingProcess, nil
}

func TestGlobalAttributesWithSingleTable(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := SingleTableProcess{}

	runId1, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xcapi.RETURN_ERROR_ON_CONFLICT,
		&xc.ProcessStartOptions{
			GlobalAttributeOptions: xc.NewGlobalAttributeOptions(
				xc.DBTableConfig{
					TableName: tblName,
					PKValue:   prcId, // use processId as the primary key value(string)
					InitialAttributes: map[string]interface{}{
						attrKeyInt: 123,
						attrKeyStr: "abc",
					},
					InitialWriteConflictMode: xcapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
				},
			),
		})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())

	// failed when trying to start the same process again with conflicted global attributes
	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, xcapi.RETURN_ERROR_ON_CONFLICT,
		&xc.ProcessStartOptions{
			IdReusePolicy: xcapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: xc.NewGlobalAttributeOptions(
				xc.DBTableConfig{
					TableName: tblName,
					PKValue:   prcId, // use processId as the primary key value(string)
					InitialAttributes: map[string]interface{}{
						attrKeyInt: 123,
						attrKeyStr: "abc",
					},
					InitialWriteConflictMode: xcapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
				},
			),
		})
	assert.NotNil(t, err)
	assert.True(t, xc.IsGlobalAttributeWriteFailure(err))

	// failed when trying to start the same process when writing str to int
	_, err = client.StartProcessWithOptions(context.Background(), prc, prcId, xcapi.RETURN_ERROR_ON_CONFLICT,
		&xc.ProcessStartOptions{
			IdReusePolicy: xcapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: xc.NewGlobalAttributeOptions(
				xc.DBTableConfig{
					TableName: tblName,
					PKValue:   prcId, // use processId as the primary key value(string)
					InitialAttributes: map[string]interface{}{
						attrKeyInt: "abc",
						attrKeyStr: "123",
					},
					InitialWriteConflictMode: xcapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
				},
			),
		})
	assert.NotNil(t, err)
	assert.True(t, xc.IsGlobalAttributeWriteFailure(err))

	// succeeded when trying to start the same process with override
	runId2, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xcapi.OVERRIDE_ON_CONFLICT,
		&xc.ProcessStartOptions{
			IdReusePolicy: xcapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: xc.NewGlobalAttributeOptions(
				xc.DBTableConfig{
					TableName: tblName,
					PKValue:   prcId, // use processId as the primary key value(string)
					InitialAttributes: map[string]interface{}{
						attrKeyInt: 123456,
						attrKeyStr: "abcdef",
					},
					InitialWriteConflictMode: xcapi.OVERRIDE_ON_CONFLICT.Ptr(),
				},
			),
		})
	assert.Nil(t, err)
	assert.NotEqual(t, runId1, runId2)

	time.Sleep(time.Second * 3)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())

	// succeeded when trying to start the same process with ignore
	runId3, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xcapi.IGNORE_CONFLICT,
		&xc.ProcessStartOptions{
			IdReusePolicy: xcapi.ALLOW_IF_NO_RUNNING.Ptr(),
			GlobalAttributeOptions: xc.NewGlobalAttributeOptions(
				xc.DBTableConfig{
					TableName: tblName,
					PKValue:   prcId, // use processId as the primary key value(string)
					InitialAttributes: map[string]interface{}{
						attrKeyInt: 123456,
						attrKeyStr: "abcdef",
					},
					InitialWriteConflictMode: xcapi.IGNORE_CONFLICT.Ptr(),
				},
			),
		})
	assert.Nil(t, err)
	assert.NotEqual(t, runId2, runId3)

	time.Sleep(time.Second * 3)
	resp, err = client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}

// TODO Test with different locking types
// TODO test using a different loading policy on starting state
