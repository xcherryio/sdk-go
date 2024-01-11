package xc

func (s PersistenceSchema) ValidateAppDatabaseForRegistry() (
	internal, err error,
) {
	keyToDef = map[string]internalColumnDef{}
	tableColNameToKey = map[string]string{}

	if s.AppDatabaseSchema != nil {
		for _, tableSchema := range s.AppDatabaseSchema.Tables {
			if tableSchema.TableName == "" {
				return nil, nil, NewProcessDefinitionError("AppDatabaseSchema.Tables.TableName is empty")
			}
			if len(tableSchema.PKColumns) == 0 {
				return nil, nil, NewProcessDefinitionError("AppDatabaseSchema.Tables.PKColumns is empty")
			}
			for _, colDef := range tableSchema.OtherColumns {
				if colDef.ColumnName == "" {
					return nil, nil, NewProcessDefinitionError("AppDatabaseSchema.Tables.OtherColumns.ColumnName is empty")
				}
				key := colDef.GlobalAttributeKey
				if key == "" {
					return nil, nil, NewProcessDefinitionError("AppDatabaseSchema.Tables.OtherColumns.GlobalAttributeKey is empty")
				}
				if _, ok := keyToDef[key]; ok {
					return nil, nil, NewProcessDefinitionError("AppDatabaseSchema.Tables.OtherColumns.GlobalAttributeKey is duplicated " + key)
				}
				keyToDef[key] = internalColumnDef{
					tableName: tableSchema.TableName,
					colDef:    colDef,
				}

				tblColName := getTableColumnName(tableSchema.TableName, colDef.ColumnName)
				if _, ok := tableColNameToKey[tblColName]; ok {
					return nil, nil, NewProcessDefinitionError("AppDatabaseSchema.Tables.OtherColumns.ColumnName is duplicated " + tblColName)
				}
				tableColNameToKey[tblColName] = key
			}
		}
	}
	return keyToDef, tableColNameToKey, nil
}

func (s PersistenceSchema) ValidateLocalAttributeForRegistry() (map[string]bool, error) {
	localAttributeKeys := map[string]bool{}
	if s.LocalAttributeSchema != nil {
		localAttributeKeys = s.LocalAttributeSchema.LocalAttributeKeys

		if len(s.LocalAttributeSchema.DefaultLoadingPolicy.LocalAttributeKeysWithLock) > 0 &&
			s.LocalAttributeSchema.DefaultLoadingPolicy.LockingType == nil {
			return nil, NewProcessDefinitionError(
				"DefaultLocalAttributeLoadingPolicy KeysWithLock is not empty but locking type is not specified")
		}

		for key := range s.LocalAttributeSchema.DefaultLoadingPolicy.LocalAttributeKeysWithLock {
			if _, ok := localAttributeKeys[key]; !ok {
				return nil, NewProcessDefinitionError(
					"DefaultLocalAttributeLoadingPolicy KeysWithLock contains invalid key " + key)
			}
		}

		for key := range s.LocalAttributeSchema.DefaultLoadingPolicy.LocalAttributeKeysNoLock {
			if _, ok := localAttributeKeys[key]; !ok {
				return nil, NewProcessDefinitionError(
					"DefaultLocalAttributeLoadingPolicy KeysNoLock contains invalid key " + key)
			}

			if _, ok := s.LocalAttributeSchema.DefaultLoadingPolicy.LocalAttributeKeysWithLock[key]; ok {
				return nil, NewProcessDefinitionError(
					"DefaultLocalAttributeLoadingPolicy KeysNoLock and KeysWithLock contains duplicated key " + key)
			}
		}
	}
	return localAttributeKeys, nil
}
