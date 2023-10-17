package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type XdbContext interface {
	GetAttempt() int
	GetProcessId() string
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
