package local_attribute

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc"
	"testing"
	"time"
)

type LocalAttributeTestProcess struct {
	xc.ProcessDefaults
}

func (b LocalAttributeTestProcess) GetPersistenceSchema() xc.PersistenceSchema {
	keys := map[string]bool{}
	keys["localAttr1"] = true
	defaultPolicy := xc.LocalAttributePolicy{
		LocalAttributeKeysNoLock: keys,
	}
	return xc.NewPersistenceSchemaWithOptions(
		xc.NewLocalAttributesSchema(keys, defaultPolicy),
		nil,
		xc.NewPersistenceSchemaOptions(),
	)
}

func (b LocalAttributeTestProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(
		&stateForInitialReadWrite{}, // read from initial global attributes and write to them
		&stateToVerifyLocalAttrs{})  // verify the global attributes write from the prev state
}

type stateForInitialReadWrite struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b stateForInitialReadWrite) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var localAttr string
	persistence.GetLocalAttribute("localAttr1", &localAttr)
	if localAttr != "initial" {
		panic(fmt.Sprintf("unexpected value %s", localAttr))
	}

	persistence.SetLocalAttribute("localAttr1", "updated")

	return xc.SingleNextState(stateToVerifyLocalAttrs{}, 1), nil
}

type stateToVerifyLocalAttrs struct {
	xc.AsyncStateDefaultsSkipWaitUntil
}

func (b stateToVerifyLocalAttrs) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	var localAttr string
	persistence.GetLocalAttribute("localAttr1", &localAttr)
	if localAttr != "updated" {
		panic(fmt.Sprintf("unexpected value %s", localAttr))
	}
	return xc.GracefulCompletingProcess, nil
}

func TestLocalAttributes(t *testing.T, client xc.Client) {
	prcId := common.GenerateProcessId()
	prc := LocalAttributeTestProcess{}

	initialWrite := map[string]interface{}{}
	initialWrite["localAttr1"] = "initial"

	_, err := client.StartProcessWithOptions(context.Background(), prc, prcId, xcapi.RETURN_ERROR_ON_CONFLICT,
		&xc.ProcessStartOptions{
			InitialLocalAttribute: initialWrite,
		})
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())
}
