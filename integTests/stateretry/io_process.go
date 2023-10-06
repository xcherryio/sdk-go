package stateretry

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"strconv"
	"testing"
	"time"
)

type BackoffProcess struct {
	xdb.ProcessDefaults
}

var defaultState = &stateDefaultPolicy{}
var customizedState = &stateCustomizedPolicy{}

func (b BackoffProcess) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.WithStartingState(defaultState, customizedState)
}

type stateDefaultPolicy struct {
	xdb.AsyncStateDefaults
	WaitUntilCounter  int
	ExecuteCounter    int
	lastTimestampMill int64
	WaiUntilFail      bool // for testing
	ExecuteSuccess    bool
}

func (b *stateDefaultPolicy) WaitUntil(ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication) (*xdb.CommandRequest, error) {
	b.WaitUntilCounter++

	if b.WaitUntilCounter == 1 {
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else if b.WaitUntilCounter == 2 {
		currTimestampMills := getCurrentTimeMillis()
		elapsedMillis := currTimestampMills - b.lastTimestampMill
		if elapsedMillis < 500 || elapsedMillis > 1500 {
			// first backoff should be ~ 1 seconds (500ms ~ 1500ms)
			b.WaiUntilFail = true
		}
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else if b.WaitUntilCounter == 3 {
		currTimestampMills := getCurrentTimeMillis()
		elapsedMillis := currTimestampMills - b.lastTimestampMill
		if elapsedMillis < 1500 || elapsedMillis > 2500 {
			// first backoff should be ~ 2 seconds (1500ms ~ 2500ms)
			b.WaiUntilFail = true
		}
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else {
		currTimestampMills := getCurrentTimeMillis()
		elapsedMillis := currTimestampMills - b.lastTimestampMill
		if elapsedMillis < 3500 || elapsedMillis > 4500 {
			// first backoff should be ~ 4 seconds (3500ms ~ 4500ms)
			b.WaiUntilFail = true
		}
		return xdb.EmptyCommandRequest(), nil
	}
}

func getCurrentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func (b *stateDefaultPolicy) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	b.ExecuteCounter++
	if b.ExecuteCounter == 1 {
		return nil, fmt.Errorf("error for testing backoff retry")
	}
	b.ExecuteSuccess = true
	return xdb.SingleNextState(&stateCustomizedPolicy{}, nil), nil
}

type stateCustomizedPolicy struct {
	xdb.AsyncStateNoWaitUntil
	Counter int
	Success bool
}

func (b *stateCustomizedPolicy) Execute(ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication) (*xdb.StateDecision, error) {
	b.Counter++
	if b.Counter == 1 {
		return nil, fmt.Errorf("error for testing backoff retry")
	}
	b.Success = true
	return xdb.ForceCompletingProcess, nil
}

func TestBackoff(t *testing.T, client xdb.Client) {
	prcId := "TestBackoff" + strconv.Itoa(int(time.Now().Unix()))
	prc := BackoffProcess{}
	_, err := client.StartProcess(context.Background(), prc, prcId, nil)
	assert.Nil(t, err)

	time.Sleep(time.Second * 10) // （1+2+4）+1+1 = 89 seconds
	resp, err := client.GetBasicClient().DescribeCurrentProcessExecution(context.Background(), prcId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())

	assert.False(t, defaultState.WaiUntilFail)
	assert.True(t, defaultState.ExecuteSuccess)
	assert.True(t, customizedState.Success)
	assert.Equal(t, 4, defaultState.WaitUntilCounter)
	assert.Equal(t, 2, defaultState.ExecuteCounter)
	assert.Equal(t, 2, customizedState.Counter)
}
