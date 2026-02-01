package schemaground

type DBConfig struct {
	Name string
	URL  string
}

type ColumnInfo struct {
	Name       string  `json:"name"`
	DataType   string  `json:"data_type"`
	IsNullable bool    `json:"is_nullable"`
	Default    *string `json:"default,omitempty"`
}

type TableInfo struct {
	Name    string                `json:"name"`
	Columns map[string]ColumnInfo `json:"columns"`
}

type DBInfo struct {
	Name   string               `json:"name"`
	Tables map[string]TableInfo `json:"tables"`
}

type ColumnDiff struct {
	Column string      `json:"column"`
	InA    *ColumnInfo `json:"in_a,omitempty"`
	InB    *ColumnInfo `json:"in_b,omitempty"`
}

type TableDiff struct {
	Table       string       `json:"table"`
	OnlyInA     bool         `json:"only_in_a,omitempty"`
	OnlyInB     bool         `json:"only_in_b,omitempty"`
	ColumnDiffs []ColumnDiff `json:"column_diffs,omitempty"`
}

type SchemaDiffSummary struct {
	DBA        string      `json:"db_a"`
	DBB        string      `json:"db_b"`
	TableDiffs []TableDiff `json:"table_diffs"`
}

type CompareResponse struct {
	Schema string            `json:"schema"`
	AtoB   SchemaDiffSummary `json:"a_to_b"`
	AtoC   SchemaDiffSummary `json:"a_to_c"`
}
