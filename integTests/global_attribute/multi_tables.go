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

type MultiTablesProcess struct {
	xdb.ProcessDefaults
}

func (b MultiTablesProcess) GetPersistenceSchema() xdb.PersistenceSchema {
	return xdb.NewPersistenceSchemaWithOptions(
		xdb.NewEmptyLocalAttributesSchema(),
		xdb.NewGlobalAttributesSchema(
			xdb.NewDBTableSchema(
				tblName, pk,
				xdbapi.NO_LOCKING,
				xdb.NewDBColumnDef(attrKeyInt, "create_timestamp", true),
				xdb.NewDBColumnDef(attrKeyStr, "first_name", true)),
			xdb.NewDBTableSchema(
				tblName2, pk2,
				xdbapi.NO_LOCKING,
				xdb.NewDBColumnDef(attrKeyInt2, "sequence", false),
				xdb.NewDBColumnDef(attrKeyStr2, "item_name", true)),
		),
		xdb.NewPersistenceSchemaOptions(
			xdb.NewNamedPersistencePolicy(
				loadNothingPolicyName, nil,
				xdb.NewTablePolicy(tblName, xdbapi.NO_LOCKING),
				xdb.NewTablePolicy(tblName2, xdbapi.NO_LOCKING),
			),
			xdb.NewNamedPersistencePolicy(
				loadSequencePolicyName, nil,
				xdb.NewTablePolicy(tblName, xdbapi.NO_LOCKING),
				xdb.NewTablePolicy(tblName2, xdbapi.NO_LOCKING, attrKeyInt2),
			),
			xdb.NewNamedPersistencePolicy(
				loadAllPolicyName, nil,
				xdb.NewTablePolicy(tblName, xdbapi.NO_LOCKING, attrKeyInt, attrKeyStr),
				xdb.NewTablePolicy(tblName2, xdbapi.NO_LOCKING, attrKeyInt2, attrKeyStr2),
			),
		),
	)
}

func (b MultiTablesProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(
		&multiTableStateForInitialReadWrite{}, // read from initial global attributes and write to them
		&multiTableStateToVerifyGlobalAttrs{}, // verify the global attributes write from the prev state
		&multiTableStateForTestLoadSequence{}) // test loading load sequence policy
}

type multiTableStateForInitialReadWrite struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b multiTableStateForInitialReadWrite) GetStateOptions() *xdb.AsyncStateOptions {
	return &xdb.AsyncStateOptions{
		PersistenceLoadingPolicyName: ptr.Any(loadAllPolicyName),
	}
}

func (b multiTableStateForInitialReadWrite) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {

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

	return xdb.SingleNextState(multiTableStateToVerifyGlobalAttrs{}, nil), nil
}

type multiTableStateToVerifyGlobalAttrs struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b multiTableStateToVerifyGlobalAttrs) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
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

	return xdb.SingleNextState(multiTableStateForTestLoadSequence{}, nil), nil
}

type multiTableStateForTestLoadSequence struct {
	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (b multiTableStateForTestLoadSequence) GetStateOptions() *xdb.AsyncStateOptions {
	return &xdb.AsyncStateOptions{
		PersistenceLoadingPolicyName: ptr.Any(loadSequencePolicyName),
	}
}

func (b multiTableStateForTestLoadSequence) Execute(
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

	persistence.GetGlobalAttribute(attrKeyInt2, &i)
	persistence.GetGlobalAttribute(attrKeyStr2, &str)
	if i != 444 {
		panic(fmt.Sprintf("unexpected value %d", i))
	}
	if str != "" {
		panic(fmt.Sprintf("unexpected value %s", str))
	}

	return xdb.ForceCompletingProcess, nil
}

func TestGlobalAttributesWithMultiTables(t *testing.T, client xdb.Client) {
	prcId := common.GenerateProcessId()
	prc := MultiTablesProcess{}

	now64 := time.Now().UnixNano()

	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xdbapi.RETURN_ERROR_ON_CONFLICT,
		&xdb.ProcessStartOptions{
			GlobalAttributeOptions: xdb.NewGlobalAttributeOptions(
				xdb.DBTableConfig{
					TableName: tblName,
					PKValue:   prcId, // use processId as the primary key value(string)
					InitialAttributes: map[string]interface{}{
						attrKeyInt: 111,
						attrKeyStr: "aaa",
					},
					InitialWriteConflictMode: xdbapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
				},
				xdb.DBTableConfig{
					TableName: tblName2,
					PKValue:   now64,
					InitialAttributes: map[string]interface{}{
						attrKeyInt2: 222,
						attrKeyStr2: "bbb",
					},
					InitialWriteConflictMode: xdbapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
				},
			),
		})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())
}
