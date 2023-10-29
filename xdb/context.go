package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type XdbContext interface {
	GetAttempt() int
	GetProcessId() string
	GetRecoverFromStateExecutionId() *string
	GetRecoverFromStateApi() *xdbapi.StateApiType
}

func newXdbContext(ctx xdbapi.Context) XdbContext {
	return &contextImpl{ctx: ctx}
}

type contextImpl struct {
	ctx xdbapi.Context
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

func (c contextImpl) GetRecoverFromStateApi() *xdbapi.StateApiType {
	return c.ctx.RecoverFromApi
}
