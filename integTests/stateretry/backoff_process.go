package stateretry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/integTests/common"
	"github.com/xcherryio/sdk-go/xc"
	"github.com/xcherryio/sdk-go/xc/ptr"
)

type BackoffProcess struct {
	xc.ProcessDefaults
}

var defaultState = &stateDefaultPolicy{}
var customizedState = &stateCustomizedPolicy{}

func (b BackoffProcess) GetAsyncStateSchema() xc.StateSchema {
	return xc.NewStateSchema(defaultState, customizedState)
}

type stateDefaultPolicy struct {
	xc.AsyncStateDefaults
	lastTimestampMill int64
	WaiUntilFail      bool // for testing
	ExecuteSuccess    bool
}

func (d stateDefaultPolicy) GetStateOptions() *xc.AsyncStateOptions {
	return &xc.AsyncStateOptions{
		WaitUntilRetryPolicy: &xcapi.RetryPolicy{
			InitialIntervalSeconds: ptr.Any(int32(2)),
		},
	}
}

func (b *stateDefaultPolicy) WaitUntil(
	ctx xc.Context, input xc.Object, communication xc.Communication,
) (*xc.CommandRequest, error) {
	if ctx.GetProcessId() != currTestProcessId {
		// ignore stale data
		return xc.EmptyCommandRequest(), nil
	}

	if ctx.GetAttempt() == 1 {
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else if ctx.GetAttempt() == 2 {
		elapsedMillis := getCurrentTimeMillis() - b.lastTimestampMill
		if elapsedMillis < 500 || elapsedMillis > 3500 { // ~2s for 1.5 sec buffer
			b.WaiUntilFail = true
			fmt.Println("backoff interval is not correct", elapsedMillis, "expected 500-3500")
		}
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else if ctx.GetAttempt() == 3 {
		elapsedMillis := getCurrentTimeMillis() - b.lastTimestampMill
		if elapsedMillis < 2000 || elapsedMillis > 6000 { // ~4s for 2 sec buffer
			b.WaiUntilFail = true
			fmt.Println("backoff interval is not correct", elapsedMillis, "expected 2000-6000")
		}
		b.lastTimestampMill = getCurrentTimeMillis()
		return nil, fmt.Errorf("error for testing backoff retry")
	} else {
		elapsedMillis := getCurrentTimeMillis() - b.lastTimestampMill
		if elapsedMillis < 6000 || elapsedMillis > 10000 { // ~8s for 2 sec buffer
			b.WaiUntilFail = true
			fmt.Println("backoff interval is not correct", elapsedMillis, "expected 6000-10000")
		}
		return xc.EmptyCommandRequest(), nil
	}
}

func getCurrentTimeMillis() int64 {
	now := time.Now()
	currMs := now.UnixNano() / 1000000
	fmt.Println("current time millis", currMs, "currTime is ", now)
	return currMs
}

func (b *stateDefaultPolicy) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	if ctx.GetProcessId() != currTestProcessId {
		// ignore stale data
		return xc.ForceCompletingProcess, nil
	}

	if ctx.GetAttempt() == 1 {
		return nil, fmt.Errorf("error for testing backoff retry")
	}
	b.ExecuteSuccess = true
	return xc.SingleNextState(&stateCustomizedPolicy{}, nil), nil
}

type stateCustomizedPolicy struct {
	xc.AsyncStateDefaultsSkipWaitUntil
	Success bool
}

func (b *stateCustomizedPolicy) Execute(
	ctx xc.Context, input xc.Object, commandResults xc.CommandResults, persistence xc.Persistence,
	communication xc.Communication,
) (*xc.StateDecision, error) {
	if ctx.GetProcessId() != currTestProcessId {
		// ignore stale data
		return xc.ForceCompletingProcess, nil
	}

	if ctx.GetAttempt() == 1 {
		return nil, fmt.Errorf("error for testing backoff retry")
	}
	b.Success = true
	return xc.ForceCompletingProcess, nil
}

var currTestProcessId string

func TestBackoff(t *testing.T, client xc.Client) {
	currTestProcessId = common.GenerateProcessId()
	prc := BackoffProcess{}
	_, err := client.StartProcessWithOptions(
		context.Background(), prc, currTestProcessId, nil, &xc.ProcessStartOptions{
			TimeoutSeconds: ptr.Any(int32(30)),
		})
	assert.Nil(t, err)

	time.Sleep(time.Second * 20) // （2+4+8）+1+1 = 9 seconds
	resp, err := client.DescribeCurrentProcessExecution(context.Background(), currTestProcessId)
	assert.Nil(t, err)
	assert.Equal(t, xcapi.COMPLETED, resp.GetStatus())

	assert.False(t, defaultState.WaiUntilFail)
	assert.True(t, defaultState.ExecuteSuccess)
	assert.True(t, customizedState.Success)
}
