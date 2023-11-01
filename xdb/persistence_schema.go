package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type PersistenceSchema struct {
	// GlobalAttributeSchema is the schema for global attributes
	// GlobalAttributes are attributes that are shared across all process executions
	// They are directly mapped to a table in the database
	GlobalAttributeSchema *GlobalAttributesSchema
	// LocalAttributeSchema is the schema for local attributes
	// LocalAttributes are attributes that are specific to a process execution
	LocalAttributeSchema *LocalAttributesSchema
	// OverrideLoadingPolicies is the loading policy with a name, which can be used as an override to the default
	// loading policy for global and local attribute schemas
	OverrideLoadingPolicies map[string]NamedPersistenceLoadingPolicy
}

type GlobalAttributesSchema struct {
	// Tables is table name to the table schema
	Tables map[string]DBTableSchema
}

type DBTableSchema struct {
	TableName string
	PK        string
	Columns   []DBColumnDef
	// DefaultTableLoadingPolicy is the default loading policy for this table
	DefaultTableLoadingPolicy TableLoadingPolicy
}

type DBColumnDef struct {
	GlobalAttributeKey string
	ColumnName         string
	Hint               *DBHint
	defaultLoading     bool
}

// DBHint is the hint for the DBConverter to convert database column to query value and vice versa
type DBHint string

type TableLoadingPolicy struct {
	TableName string
	// LoadingKeys are the attribute keys that will be loaded from the database
	LoadingKeys []string
	// TableLockingTypeDefault is the locking type for all the loaded attributes
	LockingType xdbapi.TableReadLockingPolicy
}

type NamedPersistenceLoadingPolicy struct {
	Name string
	// GlobalAttributeLoadingPolicy is the loading policy for global attributes
	// key is the table name
	GlobalAttributeTableLoadingPolicy map[string]TableLoadingPolicy
	LocalAttributeLoadingPolicy       *LocalAttributeLoadingPolicy
}

type LocalAttributesSchema struct {
	// TODO
}

type LocalAttributeLoadingPolicy struct {
	// TODO
}

func NewEmptyPersistenceSchema() PersistenceSchema {
	return NewPersistenceSchema(nil, nil)
}

// NewPersistenceSchema creates a new PersistenceSchema
// globalAttrSchema is the schema for global attributes
// localAttrSchema is the schema for local attributes
func NewPersistenceSchema(
	globalAttrSchema *GlobalAttributesSchema,
	localAttrSchema *LocalAttributesSchema,
) PersistenceSchema {
	return PersistenceSchema{
		GlobalAttributeSchema: globalAttrSchema,
		LocalAttributeSchema:  localAttrSchema,
	}
}

func NewPersistenceSchemaWithOptions(
	globalAttrSchema *GlobalAttributesSchema,
	localAttrSchema *LocalAttributesSchema,
	options PersistenceSchemaOptions,
) PersistenceSchema {
	return PersistenceSchema{
		GlobalAttributeSchema:   globalAttrSchema,
		LocalAttributeSchema:    localAttrSchema,
		OverrideLoadingPolicies: options.NameToLoadingPolicies,
	}
}

type PersistenceSchemaOptions struct {
	// NameToLoadingPolicies is the loading policy with a name, which can be used as an override to the default loading policy
	NameToLoadingPolicies map[string]NamedPersistenceLoadingPolicy
}

func NewGlobalAttributesSchema(
	table ...DBTableSchema,
) *GlobalAttributesSchema {
	m := map[string]DBTableSchema{}
	for _, t := range table {
		m[t.TableName] = t
	}
	return &GlobalAttributesSchema{
		m,
	}
}

func NewDBTableSchema(
	tableName string,
	pk string,
	defaultReadLocking xdbapi.TableReadLockingPolicy,
	columns ...DBColumnDef,
) DBTableSchema {

	var loadingKeys []string
	for _, col := range columns {
		if col.defaultLoading {
			loadingKeys = append(loadingKeys, col.GlobalAttributeKey)
		}
	}

	defaultLoadingPolicy := NewTableLoadingPolicy(tableName, defaultReadLocking, loadingKeys...)

	return DBTableSchema{
		TableName:                 tableName,
		PK:                        pk,
		Columns:                   columns,
		DefaultTableLoadingPolicy: defaultLoadingPolicy,
	}
}

func NewDBColumnDef(
	key string, dbColumn string, defaultLoading bool,
) DBColumnDef {
	return DBColumnDef{
		GlobalAttributeKey: key,
		ColumnName:         dbColumn,
		defaultLoading:     defaultLoading,
	}
}

func NewDBColumnDefWithHint(
	key string, dbColumn string, defaultLoading bool, hint DBHint,
) DBColumnDef {
	return DBColumnDef{
		GlobalAttributeKey: key,
		ColumnName:         dbColumn,
		Hint:               &hint,
		defaultLoading:     defaultLoading,
	}
}

func NewTableLoadingPolicy(
	tableName string,
	lockingType xdbapi.TableReadLockingPolicy,
	loadingKeys ...string,
) TableLoadingPolicy {
	return TableLoadingPolicy{
		TableName:   tableName,
		LoadingKeys: loadingKeys,
		LockingType: lockingType,
	}
}

func NewNamedPersistenceLoadingPolicy(
	name string,
	localAttributesLoadingPolicy *LocalAttributeLoadingPolicy,
	globalAttributeTableLoadingPolicy ...TableLoadingPolicy,
) NamedPersistenceLoadingPolicy {
	tblToPolicy := map[string]TableLoadingPolicy{}
	for _, p := range globalAttributeTableLoadingPolicy {
		tblToPolicy[p.TableName] = p
	}
	return NamedPersistenceLoadingPolicy{
		Name:                              name,
		GlobalAttributeTableLoadingPolicy: tblToPolicy,
		LocalAttributeLoadingPolicy:       localAttributesLoadingPolicy,
	}
}

func NewPersistenceSchemaOptions(
	namedPolicies ...NamedPersistenceLoadingPolicy,
) PersistenceSchemaOptions {
	nameToPolicy := map[string]NamedPersistenceLoadingPolicy{}
	for _, p := range namedPolicies {
		nameToPolicy[p.Name] = p
	}
	return PersistenceSchemaOptions{
		NameToLoadingPolicies: nameToPolicy,
	}
}
