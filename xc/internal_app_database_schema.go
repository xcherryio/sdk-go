package xc

type internalAppDatabaseSchema struct {
	tables map[string]internalTableSchema
}

type internalTableSchema struct {
	columns map[string]internalColumnDef
}

type internalColumnDef struct {
	tableName string
	isPK      bool
	pkDef     DatabasePKColumnDef
	colDef    DatabaseColumnDef
}
