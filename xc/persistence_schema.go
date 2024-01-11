package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type PersistenceSchema struct {
	// LocalAttributeSchema is the schema for local attributes
	// LocalAttributes are attributes that are specific to a process execution
	LocalAttributeSchema *LocalAttributesSchema
	// AppDatabaseSchema is the schema for app database
	AppDatabaseSchema *AppDatabaseSchema
}

type AppDatabaseSchema struct {
	// Tables is table name to the table schema
	Tables map[string]AppDatabaseTableSchema
}

type AppDatabaseTableSchema struct {
	TableName      string
	PKColumns      []DatabasePKColumnDef
	OtherColumns   []DatabaseColumnDef
	DefaultLocking xcapi.DatabaseLockingType
}

type DatabasePKColumnDef struct {
	ColumnName string
	Hint       *DatabaseHint
}

type DatabaseColumnDef struct {
	ColumnName     string
	Hint           *DatabaseHint
	defaultLoading bool
}

// DatabaseHint is the hint for the DBConverter to convert database column to query value and vice versa
type DatabaseHint string

type LocalAttributesSchema struct {
	LocalAttributeKeys   map[string]bool
	DefaultLoadingPolicy LocalAttributeLoadingPolicy
}

type LocalAttributeLoadingPolicy struct {
	LocalAttributeKeysNoLock   map[string]bool
	LocalAttributeKeysWithLock map[string]bool
	LockingType                *xcapi.DatabaseLockingType
}

// ------------------ below are constructor/helpers for PersistenceSchema ------------------

func NewEmptyPersistenceSchema() PersistenceSchema {
	return NewPersistenceSchema(nil, nil)
}

// NewPersistenceSchema creates a new PersistenceSchema
func NewPersistenceSchema(
	localAttrSchema *LocalAttributesSchema,
	appDBSchema *AppDatabaseSchema,
) PersistenceSchema {
	return PersistenceSchema{
		LocalAttributeSchema: localAttrSchema,
		AppDatabaseSchema:    appDBSchema,
	}
}

func NewPersistenceSchemaWithOptions(
	localAttrSchema *LocalAttributesSchema,
	globalAttrSchema *AppDatabaseSchema,
	options PersistenceSchemaOptions,
) PersistenceSchema {
	return PersistenceSchema{
		AppDatabaseSchema:    globalAttrSchema,
		LocalAttributeSchema: localAttrSchema,
	}
}

type PersistenceSchemaOptions struct {
	// TODO in later PRs
}

type LocalAttributeLoadingType string

const (
	NotLoad      LocalAttributeLoadingType = "not load"
	LoadWithLock LocalAttributeLoadingType = "load with lock"
	LoadNoLock   LocalAttributeLoadingType = "load no lock"
)

type LocalAttributeDef struct {
	Key                string
	DefaultLoadingType LocalAttributeLoadingType
}

func NewLocalAttributeDef(key string, defaultLoadingType LocalAttributeLoadingType) LocalAttributeDef {
	return LocalAttributeDef{
		Key:                key,
		DefaultLoadingType: defaultLoadingType,
	}
}

func NewEmptyLocalAttributesSchema() *LocalAttributesSchema {
	return nil
}

func NewLocalAttributesSchema(
	LockingType *xcapi.DatabaseLockingType,
	// TODO: it's confusing. let's remove this, and let the default locking to be exclusive locking only
	localAttributesDef ...LocalAttributeDef,
) *LocalAttributesSchema {
	keys := map[string]bool{}
	keysWithLock := map[string]bool{}
	keysNoLock := map[string]bool{}
	for _, def := range localAttributesDef {
		keys[def.Key] = true
		switch def.DefaultLoadingType {
		case NotLoad:
		case LoadWithLock:
			keysWithLock[def.Key] = true
		case LoadNoLock:
			keysNoLock[def.Key] = true
		default:
			panic("unknown loading type " + def.DefaultLoadingType)
		}
	}

	return &LocalAttributesSchema{
		LocalAttributeKeys: keys,
		DefaultLoadingPolicy: LocalAttributeLoadingPolicy{
			LocalAttributeKeysNoLock:   keysNoLock,
			LocalAttributeKeysWithLock: keysWithLock,
			LockingType:                LockingType,
		},
	}
}

func NewAppDatabaseSchema(
	table ...AppDatabaseTableSchema,
) *AppDatabaseSchema {
	m := map[string]AppDatabaseTableSchema{}
	for _, t := range table {
		m[t.TableName] = t
	}
	return &AppDatabaseSchema{
		m,
	}
}

func NewAppDatabaseTableSchema(
	tableName string,
	pkColumns []DatabasePKColumnDef,
	otherColumns []DatabaseColumnDef,
	defaultLocking xcapi.DatabaseLockingType,
) AppDatabaseTableSchema {
	return AppDatabaseTableSchema{
		TableName:      tableName,
		PKColumns:      pkColumns,
		OtherColumns:   otherColumns,
		DefaultLocking: defaultLocking,
	}
}

func NewDatabasePKColumnDef(
	dbColumn string,
) DatabasePKColumnDef {
	return DatabasePKColumnDef{
		ColumnName: dbColumn,
	}
}

func NewDatabasePKColumnDefWithHint(
	dbColumn string, hint DatabaseHint,
) DatabaseColumnDef {
	return DatabaseColumnDef{
		ColumnName: dbColumn,
		Hint:       &hint,
	}
}

func NewDatabaseColumnDef(
	dbColumn string, defaultLoading bool,
) DatabaseColumnDef {
	return DatabaseColumnDef{
		ColumnName:     dbColumn,
		defaultLoading: defaultLoading,
	}
}

func NewDatabaseColumnDefWithHint(
	dbColumn string, defaultLoading bool, hint DatabaseHint,
) DatabaseColumnDef {
	return DatabaseColumnDef{
		ColumnName:     dbColumn,
		Hint:           &hint,
		defaultLoading: defaultLoading,
	}
}
