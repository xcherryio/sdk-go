package xdb

type internalGlobalAttrDef struct {
	tableName string
	colDef    DBColumnDef
}

func (s PersistenceSchema) ValidateForRegistry() (map[string]internalGlobalAttrDef, map[string]string, error) {
	keyToDef := map[string]internalGlobalAttrDef{}
	tableColNameToKey := map[string]string{}

	if s.GlobalAttributeSchema != nil {
		for _, tableSchema := range s.GlobalAttributeSchema.Tables {
			if tableSchema.TableName == "" {
				return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.Tables.TableName is empty")
			}
			if tableSchema.PK == "" {
				return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.Tables.PK is empty")
			}
			for _, colDef := range tableSchema.Columns {
				if colDef.ColumnName == "" {
					return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.Tables.Columns.ColumnName is empty")
				}
				key := colDef.GlobalAttributeKey
				if key == "" {
					return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.Tables.Columns.GlobalAttributeKey is empty")
				}
				if _, ok := keyToDef[key]; ok {
					return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.Tables.Columns.GlobalAttributeKey is duplicated " + key)
				}
				keyToDef[key] = internalGlobalAttrDef{
					tableName: tableSchema.TableName,
					colDef:    colDef,
				}

				tblColName := getTableColumnName(tableSchema.TableName, colDef.ColumnName)
				if _, ok := tableColNameToKey[tblColName]; ok {
					return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.Tables.Columns.ColumnName is duplicated " + tblColName)
				}
				tableColNameToKey[tblColName] = key
			}
		}
	}
	return keyToDef, tableColNameToKey, nil
}

func getTableColumnName(name string, column string) string {
	return name + "#" + column
}
