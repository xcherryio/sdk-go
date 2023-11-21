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

type MultiTablesProcess struct {
	xc.ProcessDefaults
}

func (b MultiTablesProcess) GetPersistenceSchema() xc.PersistenceSchema {
	return xc.NewPersistenceSchemaWithOptions(
		xc.NewEmptyLocalAttributesSchema(),
		xc.NewGlobalAttributesSchema(
			xc.NewDBTableSchema(
				tblName, pk,
				xcapi.NO_LOCKING,
				xc.NewDBColumnDef(attrKeyInt, "create_timestamp", true),
				xc.NewDBColumnDef(attrKeyStr, "first_name", true)),
			xc.NewDBTableSchema(
				tblName2, pk2,
				xcapi.NO_LOCKING,
				xc.NewDBColumnDef(attrKeyInt2, "sequence", false),
				xc.NewDBColumnDef(attrKeyStr2, "item_name", true)),
		),
		xc.NewPersistenceSchemaOptions(
			xc.NewNamedPersistencePolicy(
				loadNothingPolicyName, nil,
				xc.NewTablePolicy(tblName, xcapi.NO_LOCKING),
				xc.NewTablePolicy(tblName2, xcapi.NO_LOCKING),
			),
			xc.NewNamedPersistencePolicy(
				loadSequencePolicyName, nil,
				xc.NewTablePolicy(tblName, xcapi.NO_LOCKING),
				xc.NewTablePolicy(tblName2, xcapi.NO_LOCKING, attrKeyInt2),
			),
			xc.NewNamedPersistencePolicy(
				loadAllPolicyName, nil,
				xc.NewTablePolicy(tblName, xcapi.NO_LOCKING, attrKeyInt, attrKeyStr),
				xc.NewTablePolicy(tblName2, xcapi.NO_LOCKING, attrKeyInt2, attrKeyStr2),
			),
		),
	)
}

func (b MultiTablesProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(
		&multiTableStateForInitialReadWrite{}, // read from initial global attributes and write to them
		&multiTableStateToVerifyGlobalAttrs{}, // verify the global attributes write from the prev state
		&multiTableStateForTestLoadSequence{}) // test loading load sequence policy
}

type multiTableStateForInitialReadWrite struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b multiTableStateForInitialReadWrite) GetStateOptions() *xc.AsyncStateOptions {
	return &xc.AsyncStateOptions{
		PersistencePolicyName: ptr.Any(loadAllPolicyName),
	}
}

func (b multiTableStateForInitialReadWrite) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {

	var i int
	persistence.GetGlobalAttribute(attrKeyInt, &i)
	var str string
	persistence.GetGlobalAttribute(attrKeyStr, &str)
	if i != 111 {
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str != "aaa" {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	persistence.GetGlobalAttribute(attrKeyInt2, &i)
	persistence.GetGlobalAttribute(attrKeyStr2, &str)
	if i != 222 {
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str != "bbb" {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	persistence.SetGlobalAttribute(attrKeyInt, 333)
	persistence.SetGlobalAttribute(attrKeyStr, "ccc")

	persistence.SetGlobalAttribute(attrKeyInt2, 444)
	persistence.SetGlobalAttribute(attrKeyStr2, "ddd")

	return xc.SingleNextState(multiTableStateToVerifyGlobalAttrs{}, nil), nil
}

type multiTableStateToVerifyGlobalAttrs struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b multiTableStateToVerifyGlobalAttrs) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var i int
	persistence.GetGlobalAttribute(attrKeyInt, &i)
	var str string
	persistence.GetGlobalAttribute(attrKeyStr, &str)
	if i != 333 {
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str != "ccc" {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	var i2 int
	var str2 string
	persistence.GetGlobalAttribute(attrKeyInt2, &i2)
	persistence.GetGlobalAttribute(attrKeyStr2, &str2)
	if i2 != 0 { // because the default policy won't load this attribute
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str2 != "ddd" {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	return xc.SingleNextState(multiTableStateForTestLoadSequence{}, nil), nil
}

type multiTableStateForTestLoadSequence struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b multiTableStateForTestLoadSequence) GetStateOptions() *xc.AsyncStateOptions {
	return &xc.AsyncStateOptions{
		PersistencePolicyName: ptr.Any(loadSequencePolicyName),
	}
}

func (b multiTableStateForTestLoadSequence) Execute(
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

	persistence.GetGlobalAttribute(attrKeyInt2, &i)
	persistence.GetGlobalAttribute(attrKeyStr2, &str)
	if i != 444 {
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str != "" {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	return xc.ForceCompletingProcess, nil
}

func TestGlobalAttributesWithMultiTables(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := MultiTablesProcess{}

	now64 := time.Now().UnixNano()

	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xcapi.RETURN_ERROR_ON_CONFLICT,
		&xc.ProcessStartOptions{
			GlobalAttributeOptions: xc.NewGlobalAttributeOptions(
				xc.DBTableConfig{
					TableName: tblName,
					PKValue:   prcId, // use processId as the primary key value(string)
					InitialAttributes: map[string]interface{}{
						attrKeyInt: 111,
						attrKeyStr: "aaa",
					},
					InitialWriteConflictMode: xcapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
				},
				xc.DBTableConfig{
					TableName: tblName2,
					PKValue:   now64,
					InitialAttributes: map[string]interface{}{
						attrKeyInt2: 222,
						attrKeyStr2: "bbb",
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
}
