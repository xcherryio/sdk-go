package stateretry

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/integTests/common"
	"github.com/xdblab/xdb-golang-sdk/xdb"
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
	lastTimestampMill int64
	WaiUntilFail      bool // for testing
	ExecuteSuccess    bool
}

func (b *stateDefaultPolicy) WaitUntil(ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication) (*xdb.CommandRequest, error) {
	if ctx.GetProcessId() != currTestProcessId {
		// ignore stale data
		return xdb.EmptyCommandRequest(), nil
	}

	if ctx.GetAttempt() == 1 {
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else if ctx.GetAttempt() == 2 {
		elapsedMillis := getCurrentTimeMillis() - b.lastTimestampMill
		if elapsedMillis < 500 || elapsedMillis > 1900 {
			b.WaiUntilFail = true
			fmt.Println("backoff interval is not correct", elapsedMillis, "expected 500-1900")
		}
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else if ctx.GetAttempt() == 3 {
		elapsedMillis := getCurrentTimeMillis() - b.lastTimestampMill
		if elapsedMillis < 1500 || elapsedMillis > 2900 {
			b.WaiUntilFail = true
			fmt.Println("backoff interval is not correct", elapsedMillis, "expected 1500-2900")
		}
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else {
		elapsedMillis := getCurrentTimeMillis() - b.lastTimestampMill
		if elapsedMillis < 3500 || elapsedMillis > 4900 {
			b.WaiUntilFail = true
			fmt.Println("backoff interval is not correct", elapsedMillis, "expected 3500-4900")
		}
		return xdb.EmptyCommandRequest(), nil
	}
}

func getCurrentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func (b *stateDefaultPolicy) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	if ctx.GetProcessId() != currTestProcessId {
		// ignore stale data
		return xdb.ForceCompletingProcess, nil
	}

	if ctx.GetAttempt() == 1 {
		return nil, fmt.Errorf("error for testing backoff retry")
	}
	b.ExecuteSuccess = true
	return xdb.SingleNextState(&stateCustomizedPolicy{}, nil), nil
}

type stateCustomizedPolicy struct {
	xdb.AsyncStateNoWaitUntil
	Success bool
}

func (b *stateCustomizedPolicy) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence, communication xdb.Communication,
) (*xdb.StateDecision, error) {
	if ctx.GetProcessId() != currTestProcessId {
		// ignore stale data
		return xdb.ForceCompletingProcess, nil
	}

	if ctx.GetAttempt() == 1 {
		return nil, fmt.Errorf("error for testing backoff retry")
	}
	b.Success = true
	return xdb.ForceCompletingProcess, nil
}

var currTestProcessId string

func TestBackoff(t *testing.T, client xdb.Client) {
	currTestProcessId = common.GenerateProcessId()
	prc := BackoffProcess{}
	_, err := client.StartProcess(context.Background(), prc, currTestProcessId, nil)
	assert.Nil(t, err)

	time.Sleep(time.Second * 15) // （1+2+4）+1+1 = 9 seconds
	resp, err := client.DescribeCurrentProcessExecution(context.Background(), currTestProcessId)
	assert.Nil(t, err)
	assert.Equal(t, xdbapi.COMPLETED, resp.GetStatus())

	assert.False(t, defaultState.WaiUntilFail)
	assert.True(t, defaultState.ExecuteSuccess)
	assert.True(t, customizedState.Success)
}
