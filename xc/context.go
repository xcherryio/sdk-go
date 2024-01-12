package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type Context interface {
	GetAttempt() int
	GetProcessId() string
	GetRecoverFromStateExecutionId() *string
	GetRecoverFromStateApi() *xcapi.WorkerApiType
}

func newContext(ctx xcapi.Context) Context {
	return &contextImpl{ctx: ctx}
}

type contextImpl struct {
	ctx xcapi.Context
}

func (c contextImpl) GetProcessId() string {
	return c.ctx.GetProcessId()
}

func (c contextImpl) GetAttempt() int {
	return int(c.ctx.GetAttempt())
}

func (c contextImpl) GetRecoverFromStateExecutionId() *string {
	return c.ctx.RecoverFromStateExecutionId
}

func (c contextImpl) GetRecoverFromStateApi() *xcapi.WorkerApiType {
	return c.ctx.RecoverFromApi
}
