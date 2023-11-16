package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type persistenceImpl struct {
	dbConverter DBConverter

	// for global attributes
	globalAttrDefs              map[string]internalGlobalAttrDef
	globalAttrTableColNameToKey map[string]string
	currGlobalAttrs             map[string]xcapi.TableColumnValue
	currUpdatedGlobalAttrs      map[string]xcapi.TableColumnValue
	// for local attributes
	localAttrKeys         map[string]bool
	currLocalAttrs        map[string]xcapi.EncodedObject
	currUpdatedLocalAttrs map[string]xcapi.EncodedObject
}

func NewPersistenceImpl(
	dbConverter DBConverter,
	globalAttrDefs map[string]internalGlobalAttrDef, globalAttrTableColNameToKey map[string]string,
	currGlobalAttrs *xcapi.LoadGlobalAttributeResponse,
	localAttrKeys map[string]bool,
	currLocalAttrs *xcapi.LoadLocalAttributesResponse,
) Persistence {
	currGlobalAttrsMap := map[string]xcapi.TableColumnValue{}

	if currGlobalAttrs != nil {
		for _, tblResp := range currGlobalAttrs.TableResponses {
			tblName := tblResp.GetTableName()
			for _, colVal := range tblResp.GetColumns() {
				key, ok := globalAttrTableColNameToKey[getTableColumnName(tblName, colVal.DbColumn)]
				if !ok {
					panic("global attribute not found " + colVal.DbColumn)
				}
				currGlobalAttrsMap[key] = colVal
			}
		}
	}

	currLocalAttrsMap := map[string]xcapi.EncodedObject{}
	if currLocalAttrs != nil {
		for _, kv := range currLocalAttrs.Attributes {
			if _, ok := localAttrKeys[kv.Key]; !ok {
				panic("local attribute not found " + kv.Key)
			}
			currLocalAttrsMap[kv.Key] = kv.Value
		}
	}

	return &persistenceImpl{
		dbConverter:                 dbConverter,
		globalAttrDefs:              globalAttrDefs,
		globalAttrTableColNameToKey: globalAttrTableColNameToKey,
		currGlobalAttrs:             currGlobalAttrsMap,
		// start from empty
		currUpdatedGlobalAttrs: map[string]xcapi.TableColumnValue{},

		localAttrKeys:         localAttrKeys,
		currLocalAttrs:        currLocalAttrsMap,
		currUpdatedLocalAttrs: map[string]xcapi.EncodedObject{},
	}
}

func (p *persistenceImpl) GetGlobalAttribute(key string, resultPtr interface{}) {
	def, ok := p.globalAttrDefs[key]
	if !ok {
		panic("global attribute not found " + key)
	}
	curVal, ok := p.currGlobalAttrs[key]
	if !ok {
		return
	}
	err := p.dbConverter.FromDBValue(curVal.DbQueryValue, def.colDef.Hint, resultPtr)
	if err != nil {
		panic(err)
	}
}

func (p *persistenceImpl) SetGlobalAttribute(key string, value interface{}) {
	def, ok := p.globalAttrDefs[key]
	if !ok {
		panic("global attribute not found " + key)
	}
	dbQueryValue, err := p.dbConverter.ToDBValue(value, def.colDef.Hint)
	if err != nil {
		panic(err)
	}
	val := xcapi.TableColumnValue{
		DbQueryValue: dbQueryValue,
		DbColumn:     def.colDef.ColumnName,
	}
	p.currGlobalAttrs[key] = val
	p.currUpdatedGlobalAttrs[key] = val
}

func (p *persistenceImpl) getGlobalAttributesToUpdate() []xcapi.GlobalAttributeTableRowUpdate {
	res := map[string]xcapi.GlobalAttributeTableRowUpdate{}
	for k, v := range p.currUpdatedGlobalAttrs {
		def := p.globalAttrDefs[k]
		tblName := def.tableName

		tblUpdate, ok := res[tblName]
		if !ok {
			tblUpdate = xcapi.GlobalAttributeTableRowUpdate{
				TableName: tblName,
			}
		}

		tblUpdate.UpdateColumns = append(tblUpdate.UpdateColumns, v)
		res[tblName] = tblUpdate
	}

	var res2 []xcapi.GlobalAttributeTableRowUpdate
	for _, v := range res {
		res2 = append(res2, v)
	}
	return res2
}

func (p *persistenceImpl) GetLocalAttribute(key string) *xcapi.EncodedObject {
	_, ok := p.localAttrKeys[key]
	if !ok {
		panic("local attribute not found " + key)
	}

	curVal, ok := p.currLocalAttrs[key]
	if !ok {
		return nil
	}

	return &curVal
}

func (p *persistenceImpl) SetLocalAttribute(key string, value xcapi.EncodedObject) {
	_, ok := p.localAttrKeys[key]
	if !ok {
		panic("local attribute not found " + key)
	}

	p.currLocalAttrs[key] = value
	p.currUpdatedLocalAttrs[key] = value
}

func (p *persistenceImpl) getLocalAttributesToUpdate() []xcapi.KeyValue {
	var res []xcapi.KeyValue
	for k, v := range p.currUpdatedLocalAttrs {
		res = append(res, xcapi.KeyValue{
			Key:   k,
			Value: v,
		})
	}
	return res
}
