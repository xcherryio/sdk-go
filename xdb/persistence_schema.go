package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type PersistenceSchema struct {
	// GlobalAttributeSchema is the schema for global attributes
	// They are attributes that are shared across all process executions
	// They are directly mapped to a table in the database
	GlobalAttributeSchema *GlobalAttributesSchema
	// LocalAttributeSchema is the schema for local attributes
	// They are attributes that are specific to a process execution
	LocalAttributeSchema *LocalAttributesSchema
	// DefaultLoadingPolicy is the default loading policy for AsyncStates and RPCs
	// PersistenceLoadingPolicy defines how to what attributes will be read from, and how to lock them for read
	DefaultLoadingPolicy PersistenceLoadingPolicy
	// NamedLoadingPolicies is the loading policy with a name, which can be used as an override
	NamedLoadingPolicies map[string]PersistenceLoadingPolicy
}

// NewPersistenceSchema creates a new PersistenceSchema
// globalAttrSchema is the schema for global attributes
// localAttrSchema is the schema for local attributes
// defaultLoadingPolicy is the default loading policy for AsyncStates and RPCs
// namedLoadingPolicies is the loading policy with a name, which can be used as an override
func NewPersistenceSchema(
	globalAttrSchema *GlobalAttributesSchema,
	localAttrSchema *LocalAttributesSchema,
	defaultLoadingPolicy PersistenceLoadingPolicy,
	namedLoadingPolicies map[string]PersistenceLoadingPolicy,
) PersistenceSchema {
	return PersistenceSchema{
		GlobalAttributeSchema: globalAttrSchema,
		LocalAttributeSchema:  localAttrSchema,
		DefaultLoadingPolicy:  defaultLoadingPolicy,
		NamedLoadingPolicies:  namedLoadingPolicies,
	}
}

func NewEmptyPersistenceSchema() PersistenceSchema {
	return NewPersistenceSchema(nil, nil,
		NewPersistenceLoadingPolicy(nil, nil),
		nil)
}

type internalGlobalAttrDef struct {
	colName            string
	hint               *string
	altTableName       *string
	altTableForeignKey *string
}

type GlobalAttributesSchema struct {
	// DefaultTableName is the name of the default table that the global attributes will be mapped to
	// To map to a different table, use SecondaryTable instead
	DefaultTableName string
	// DefaultTablePrimaryKey is the PK(primary key) of the default table
	// All the attributes will be mapped to the row using the primary key value
	DefaultTablePrimaryKey string
	// DefaultTablePrimaryKeyHint is the hint for the primary key value for converting to db value
	DefaultTablePrimaryKeyHint *string
	// DefaultTableAttributeDefs is the global attribute definition for the default table
	DefaultTableAttributeDefs []GlobalAttributeDef
	// GASecondaryTableDefs is the global attribute definition for the secondary tables
	GASecondaryTableDefs []GASecondaryTableDef
}

type GlobalAttributeDef struct {
	Key      string
	DBColumn string
	Hint     *string
}

func NewGlobalAttributeDef(
	key string, dbColumn string,
) GlobalAttributeDef {
	return GlobalAttributeDef{
		Key:      key,
		DBColumn: dbColumn,
	}
}

func NewGlobalAttributeDefWithHint(
	key string, dbColumn string, hint string,
) GlobalAttributeDef {
	return GlobalAttributeDef{
		Key:      key,
		DBColumn: dbColumn,
		Hint:     &hint,
	}
}

func NewGlobalAttributesSchema(
	defaultDbTable string,
	defaultDbTablePK string,
	attrs ...GlobalAttributeDef,
) *GlobalAttributesSchema {
	return &GlobalAttributesSchema{
		DefaultTableName:          defaultDbTable,
		DefaultTablePrimaryKey:    defaultDbTablePK,
		DefaultTableAttributeDefs: attrs,
	}
}

func NewGlobalAttributesSchemaWithHint(
	defaultDbTable string,
	defaultDbTablePK string,
	defaultDbTablePKHint string,
	attrs ...GlobalAttributeDef,
) GlobalAttributesSchema {
	return GlobalAttributesSchema{
		DefaultTableName:           defaultDbTable,
		DefaultTablePrimaryKey:     defaultDbTablePK,
		DefaultTablePrimaryKeyHint: &defaultDbTablePKHint,
		DefaultTableAttributeDefs:  attrs,
	}
}

type GASecondaryTableDef struct {
	DBTable    string
	ForeignKey string
	Attributes []GlobalAttributeDef
}

func NewGASecondaryTableDef(
	dbTable string,
	foreignKey string,
	attributes ...GlobalAttributeDef,
) GASecondaryTableDef {
	return GASecondaryTableDef{
		DBTable:    dbTable,
		ForeignKey: foreignKey,
		Attributes: attributes,
	}
}

func NewGlobalAttributesSchemaWithSecondaries(
	defaultTable string,
	defaultTablePK string,
	defaultTableAttrs []GlobalAttributeDef,
	secondaryTables ...GASecondaryTableDef,
) GlobalAttributesSchema {

	return GlobalAttributesSchema{
		DefaultTableName:          defaultTable,
		DefaultTablePrimaryKey:    defaultTablePK,
		DefaultTableAttributeDefs: defaultTableAttrs,
		GASecondaryTableDefs:      secondaryTables,
	}
}

type PersistenceLoadingPolicy struct {
	GlobalAttributeLoadingPolicy *GlobalAttributeLoadingPolicy
	LocalAttributeLoadingPolicy  *LocalAttributeLoadingPolicy
}

func NewPersistenceLoadingPolicy(
	globalAttributes *GlobalAttributeLoadingPolicy,
	localAttributes *LocalAttributeLoadingPolicy,
) PersistenceLoadingPolicy {
	return PersistenceLoadingPolicy{
		GlobalAttributeLoadingPolicy: globalAttributes,
		LocalAttributeLoadingPolicy:  localAttributes,
	}
}

type GlobalAttributeLoadingPolicy struct {
	// LoadingKeys are the attribute keys that will be loaded from the database
	LoadingKeys []string
	// TableLockingTypeDefault is the default locking type for all the loaded attributes, or for all tables
	TableLockingTypeDefault xdbapi.AttributeReadLockingType
	// TableLockingTypeOverrides is the override for the locking type for the TableLockingTypeDefault
	TableLockingTypeOverrides map[string]xdbapi.AttributeReadLockingType
}

func NewGlobalAttributeLoadingPolicy(
	tableLockingTypeDefault xdbapi.AttributeReadLockingType,
	loadingKeys ...string,
) *GlobalAttributeLoadingPolicy {
	return &GlobalAttributeLoadingPolicy{
		LoadingKeys:             loadingKeys,
		TableLockingTypeDefault: tableLockingTypeDefault,
	}
}

func NewGlobalAttributeLoadingPolicyWithOverrides(
	tableLockingTypeDefault xdbapi.AttributeReadLockingType,
	tableLockingTypeOverrides map[string]xdbapi.AttributeReadLockingType,
	loadingKeys ...string,
) *GlobalAttributeLoadingPolicy {
	return &GlobalAttributeLoadingPolicy{
		LoadingKeys:               loadingKeys,
		TableLockingTypeDefault:   tableLockingTypeDefault,
		TableLockingTypeOverrides: tableLockingTypeOverrides,
	}
}

type LocalAttributeLoadingPolicy struct {
	// TODO
}

type LocalAttributesSchema struct {
	// TODO
}
