package xdb

import "github.com/xdblab/xdb-apis/goapi/xdbapi"

type persistenceImpl struct {
	dbConverter DBConverter

	// for global attributes
	globalAttrDefs              map[string]internalGlobalAttrDef
	globalAttrTableColNameToKey map[string]string
	currGlobalAttrs             map[string]xdbapi.GlobalAttributeValue
	currUpdatedGlobalAttrs      map[string]xdbapi.GlobalAttributeValue
}

func NewPersistenceImpl(
	dbConverter DBConverter, defaultTable string,
	globalAttrDefs map[string]internalGlobalAttrDef, globalAttrTableColNameToKey map[string]string,
	currGlobalAttrs []xdbapi.GlobalAttributeValue,
) Persistence {
	currGlobalAttrsMap := map[string]xdbapi.GlobalAttributeValue{}
	for _, attr := range currGlobalAttrs {
		tbl := defaultTable
		if attr.AlternativeTable == nil {
			tbl = *attr.AlternativeTable
		}
		key, ok := globalAttrTableColNameToKey[getTableColumnName(tbl, attr.DbColumn)]
		if !ok {
			panic("global attribute not found " + attr.DbColumn)
		}
		currGlobalAttrsMap[key] = attr
	}

	return &persistenceImpl{
		dbConverter:                 dbConverter,
		globalAttrDefs:              globalAttrDefs,
		globalAttrTableColNameToKey: globalAttrTableColNameToKey,
		currGlobalAttrs:             currGlobalAttrsMap,
	}
}

func (p persistenceImpl) GetGlobalAttribute(key string, resultPtr interface{}) {
	def, ok := p.globalAttrDefs[key]
	if !ok {
		panic("global attribute not found " + key)
	}
	curVal, ok := p.currGlobalAttrs[key]
	if !ok {
		return
	}
	err := p.dbConverter.FromDBValue(curVal.DbQueryValue, def.hint, resultPtr)
	if err != nil {
		panic(err)
	}
}

func (p persistenceImpl) SetGlobalAttribute(key string, value interface{}) {
	def, ok := p.globalAttrDefs[key]
	if !ok {
		panic("global attribute not found " + key)
	}
	dbQueryValue, err := p.dbConverter.ToDBValue(value, def.hint)
	if err != nil {
		panic(err)
	}
	val := xdbapi.GlobalAttributeValue{
		DbQueryValue:                     dbQueryValue,
		DbColumn:                         def.colName,
		AlternativeTable:                 def.altTableName,
		AlternativeTableForeignKeyColumn: def.altTableForeignKey,
	}
	p.currGlobalAttrs[key] = val
	p.currUpdatedGlobalAttrs[key] = val
}

func (p persistenceImpl) getGlobalAttributesToUpdate() []xdbapi.GlobalAttributeValue {
	var res []xdbapi.GlobalAttributeValue
	for _, v := range p.currUpdatedGlobalAttrs {
		res = append(res, v)
	}
	return res
}
