package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type ProcessStartOptions struct {
	// TimeoutSeconds is the timeout for the process execution.
	// Default: 0, mean which means infinite timeout.
	// This will override the timeout defined in process definition
	TimeoutSeconds *int32
	// IdReusePolicy is the policy for reusing process id.
	// Default: xdbapi.ALLOW_IF_NO_RUNNING when set as nil.
	// This will override the IdReusePolicy defined in process definition.
	IdReusePolicy *xdbapi.ProcessIdReusePolicy
	// GlobalAttributeOptions is the options for global attributes
	// Required if using global attribute feature
	GlobalAttributeOptions *GlobalAttributeOptions
}

type GlobalAttributeOptions struct {
	// PrimaryKeyAttributeValue is the value of the primary key value for the default table
	// required if using global attributes
	PrimaryAttributeValue interface{}
	// InitialAttributes is the initial attributes to be set when starting the process execution
	// Key is the attribute key, value is the attribute value
	InitialAttributes map[string]interface{}
	// InitialWriteConflictMode is for how to resolve the write conflict when setting the initial attributes
	InitialWriteConflictMode xdbapi.AttributeWriteConflictMode
}
