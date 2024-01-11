package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type ProcessStartOptions struct {
	// TimeoutSeconds is the timeout for the process execution.
	// Default is nil, which will use the value defined in process definition.
	// If provided, it will override the timeout defined in process definition.
	TimeoutSeconds *int32
	// IdReusePolicy is the policy for reusing process id.
	// Default: xcapi.ALLOW_IF_NO_RUNNING when set as nil.
	// This will override the IdReusePolicy defined in process definition.
	IdReusePolicy *xcapi.ProcessIdReusePolicy
	// InitialLocalAttribute is the initial local attributes to be set when starting the process execution
	// Not required.
	InitialLocalAttribute map[string]interface{}
	// AppDatabaseOptions is the options for App Database feature
	// Required if using App Database feature. Each table defined in AppDatabaseSchema must be included.
	AppDatabaseOptions *AppDatabaseOptions
}

type AppDatabaseOptions struct {
	// AppDatabaseTableConfig is the config for App Database feature
	// all tables defined in AppDatabaseSchema must be included
	tables map[string]AppDatabaseTableOptions
}

func NewAppDatabaseOptions(
	tableOptions ...AppDatabaseTableOptions,
) *AppDatabaseOptions {
	dbTableOpts := map[string]AppDatabaseTableOptions{}
	for _, tblOpts := range tableOptions {
		dbTableOpts[tblOpts.TableName] = tblOpts
	}
	return &AppDatabaseOptions{
		tables: dbTableOpts,
	}
}

type AppDatabaseTableOptions struct {
	TableName string
	Rows      []AppDatabaseRowOptions
}

type AppDatabaseRowOptions struct {
	// PrimaryKeyValues is the primary key values for selecting the row
	// all primary key columns defined in AppDatabaseSchema must be included
	// Key is the column name, value is the column value
	PrimaryKeyValues map[string]interface{}
	// InitialWriteColumns is the initial columns to be set for the row, when starting the process execution
	// Key is the column name, value is the column value
	// optional. Nil/empty map will be ignored as no initial write.
	InitialWriteColumns map[string]interface{}
	// InitialWriteConflictMode is for how to resolve the write conflict when setting the initialWrite
	// Required if InitialWriteColumns is not empty
	InitialWriteConflictMode *xcapi.WriteConflictMode
}
