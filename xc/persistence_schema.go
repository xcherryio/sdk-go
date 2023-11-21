package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type PersistenceSchema struct {
	// LocalAttributeSchema is the schema for local attributes
	// LocalAttributes are attributes that are specific to a process execution
	LocalAttributeSchema *LocalAttributesSchema
	// GlobalAttributeSchema is the schema for global attributes
	// GlobalAttributes are attributes that are shared across all process executions
	// They are directly mapped to a table in the database
	GlobalAttributeSchema *GlobalAttributesSchema
	// OverridePersistencePolicies is the persistence policy with a name, which can be used as an override to the default
	// policy for global and local attribute schemas
	OverridePersistencePolicies map[string]NamedPersistencePolicy
}

type GlobalAttributesSchema struct {
	// Tables is table name to the table schema
	Tables map[string]DBTableSchema
}

type DBTableSchema struct {
	TableName string
	PK        string
	Columns   []DBColumnDef
	// DefaultTablePolicy is the default loading policy for this table
	DefaultTablePolicy TablePolicy
}

type DBColumnDef struct {
	GlobalAttributeKey string
	ColumnName         string
	Hint               *DBHint
	defaultLoading     bool
}

// DBHint is the hint for the DBConverter to convert database column to query value and vice versa
type DBHint string

type TablePolicy struct {
	TableName string
	// LoadingKeys are the attribute keys that will be loaded from the database
	LoadingKeys []string
	// TableLockingTypeDefault is the locking type for all the loaded attributes
	LockingType xcapi.TableReadLockingPolicy
}

type NamedPersistencePolicy struct {
	Name string
	// LocalAttributePolicy is the policy for local attributes
	LocalAttributePolicy *LocalAttributePolicy
	// GlobalAttributePolicy is the policy for global attributes
	// key is the table name
	GlobalAttributePolicy map[string]TablePolicy
}

type LocalAttributesSchema struct {
	// TODO
}

type LocalAttributePolicy struct {
	// TODO
}

func NewEmptyPersistenceSchema() PersistenceSchema {
	return NewPersistenceSchema(nil, nil)
}

// NewPersistenceSchema creates a new PersistenceSchema
// globalAttrSchema is the schema for global attributes
// localAttrSchema is the schema for local attributes
func NewPersistenceSchema(
	localAttrSchema *LocalAttributesSchema,
	globalAttrSchema *GlobalAttributesSchema,
) PersistenceSchema {
	return PersistenceSchema{
		GlobalAttributeSchema: globalAttrSchema,
		LocalAttributeSchema:  localAttrSchema,
	}
}

func NewPersistenceSchemaWithOptions(
	localAttrSchema *LocalAttributesSchema,
	globalAttrSchema *GlobalAttributesSchema,
	options PersistenceSchemaOptions,
) PersistenceSchema {
	return PersistenceSchema{
		GlobalAttributeSchema:       globalAttrSchema,
		LocalAttributeSchema:        localAttrSchema,
		OverridePersistencePolicies: options.NameToPolicy,
	}
}

type PersistenceSchemaOptions struct {
	// NameToPolicy is the persistence policy with a name,
	// which can be used as an override to the default policy
	NameToPolicy map[string]NamedPersistencePolicy
}

func NewEmptyLocalAttributesSchema() *LocalAttributesSchema {
	return nil
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

func NewEmptyGlobalAttributesSchema() *GlobalAttributesSchema {
	return nil
}

func NewDBTableSchema(
	tableName string,
	pk string,
	defaultReadLocking xcapi.TableReadLockingPolicy,
	columns ...DBColumnDef,
) DBTableSchema {

	var loadingKeys []string
	for _, col := range columns {
		if col.defaultLoading {
			loadingKeys = append(loadingKeys, col.GlobalAttributeKey)
		}
	}

	defaultPolicy := NewTablePolicy(tableName, defaultReadLocking, loadingKeys...)

	return DBTableSchema{
		TableName:          tableName,
		PK:                 pk,
		Columns:            columns,
		DefaultTablePolicy: defaultPolicy,
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

func NewTablePolicy(
	tableName string,
	lockingType xcapi.TableReadLockingPolicy,
	loadingKeys ...string,
) TablePolicy {
	return TablePolicy{
		TableName:   tableName,
		LoadingKeys: loadingKeys,
		LockingType: lockingType,
	}
}

func NewNamedPersistencePolicy(
	name string,
	localAttributesPolicy *LocalAttributePolicy,
	globalAttributeTablePolicy ...TablePolicy,
) NamedPersistencePolicy {
	tblToPolicy := map[string]TablePolicy{}
	for _, p := range globalAttributeTablePolicy {
		tblToPolicy[p.TableName] = p
	}
	return NamedPersistencePolicy{
		Name:                  name,
		GlobalAttributePolicy: tblToPolicy,
		LocalAttributePolicy:  localAttributesPolicy,
	}
}

func NewPersistenceSchemaOptions(
	namedPolicies ...NamedPersistencePolicy,
) PersistenceSchemaOptions {
	nameToPolicy := map[string]NamedPersistencePolicy{}
	for _, p := range namedPolicies {
		nameToPolicy[p.Name] = p
	}
	return PersistenceSchemaOptions{
		NameToPolicy: nameToPolicy,
	}
}
