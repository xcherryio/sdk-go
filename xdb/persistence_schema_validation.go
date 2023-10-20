package xdb

func (s PersistenceSchema) Validate() (map[string]internalGlobalAttrDef, map[string]string, error) {
	keyToDef := map[string]internalGlobalAttrDef{}
	tableColNameToKey := map[string]string{}

	if s.GlobalAttributeSchema != nil {
		gas := *s.GlobalAttributeSchema
		if gas.DefaultTableName == "" {
			return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.DefaultTableName is empty")
		}
		if gas.DefaultTablePrimaryKey == "" {
			return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.DefaultTablePrimaryKey is empty")
		}
		if s.DefaultLoadingPolicy.GlobalAttributeLoadingPolicy == nil {
			return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.DefaultLoadingPolicy.GlobalAttributeLoadingPolicy is empty")
		}

		for _, attr := range gas.DefaultTableAttributeDefs {
			if attr.Key == "" {
				return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.DefaultTableAttributeDefs.Key is empty")
			}
			if attr.DBColumn == "" {
				return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.DefaultTableAttributeDefs.DBColumn is empty")
			}
			if _, ok := keyToDef[attr.Key]; ok {
				return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.DefaultTableAttributeDefs.Key is duplicated")
			}
			keyToDef[attr.Key] = internalGlobalAttrDef{
				colName: attr.DBColumn,
				hint:    attr.Hint,
			}
			tableColName := getTableColumnName(gas.DefaultTableName, attr.DBColumn)
			if _, ok := tableColNameToKey[tableColName]; ok {
				return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.DefaultTableAttributeDefs.DBColumn is duplicated")
			}
			tableColNameToKey[tableColName] = attr.Key
		}
		for _, secondaryTable := range gas.GASecondaryTableDefs {
			for _, attr := range secondaryTable.Attributes {
				if attr.Key == "" {
					return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.GASecondaryTableDefs.Attributes.Key is empty")
				}
				if attr.DBColumn == "" {
					return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.GASecondaryTableDefs.DBColumn is empty")
				}
				if _, ok := keyToDef[attr.Key]; ok {
					return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.GASecondaryTableDefs.Key is duplicated")
				}
				keyToDef[attr.Key] = internalGlobalAttrDef{
					colName:            attr.DBColumn,
					hint:               attr.Hint,
					altTableName:       &secondaryTable.DBTable,
					altTableForeignKey: &secondaryTable.ForeignKey,
				}
				tableColName := secondaryTable.DBTable + "#" + attr.DBColumn
				if _, ok := tableColNameToKey[tableColName]; ok {
					return nil, nil, NewProcessDefinitionError("GlobalAttributeSchema.GASecondaryTableDefs.DBColumn is duplicated")
				}
				tableColNameToKey[tableColName] = attr.Key
			}
		}
	}
	return keyToDef, tableColNameToKey, nil
}

func getTableColumnName(name string, column string) string {
	return name + "#" + column
}
