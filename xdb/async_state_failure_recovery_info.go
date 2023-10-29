package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type AsyncStateFailureRecoveryInfo struct {
	Policy                          xdbapi.StateFailureRecoveryPolicy
	StateFailureProceedStateId      *string
	StateFailureProceedStateOptions *AsyncStateOptions
}
