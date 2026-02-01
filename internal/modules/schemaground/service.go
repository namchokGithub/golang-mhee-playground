package schemaground

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Compare(schema string) (CompareResponse, error) {
	if schema == "" {
		schema = "public"
	}

	loadEnvFile()

	cfgA, err := loadOneConfig("A")
	if err != nil {
		return CompareResponse{}, err
	}
	cfgB, err := loadOneConfig("B")
	if err != nil {
		return CompareResponse{}, err
	}
	cfgC, err := loadOneConfig("C")
	if err != nil {
		return CompareResponse{}, err
	}

	dbA, err := sql.Open("postgres", cfgA.URL)
	if err != nil {
		return CompareResponse{}, fmt.Errorf("connect A error: %w", err)
	}
	defer dbA.Close()

	dbB, err := sql.Open("postgres", cfgB.URL)
	if err != nil {
		return CompareResponse{}, fmt.Errorf("connect B error: %w", err)
	}
	defer dbB.Close()

	dbC, err := sql.Open("postgres", cfgC.URL)
	if err != nil {
		return CompareResponse{}, fmt.Errorf("connect C error: %w", err)
	}
	defer dbC.Close()

	infoA, err := loadDBSchema(dbA, cfgA.Name, schema)
	if err != nil {
		return CompareResponse{}, fmt.Errorf("load schema A error: %w", err)
	}
	infoB, err := loadDBSchema(dbB, cfgB.Name, schema)
	if err != nil {
		return CompareResponse{}, fmt.Errorf("load schema B error: %w", err)
	}
	infoC, err := loadDBSchema(dbC, cfgC.Name, schema)
	if err != nil {
		return CompareResponse{}, fmt.Errorf("load schema C error: %w", err)
	}

	resultAtoB := CompareSchemas(infoA, infoB)
	resultAtoC := CompareSchemas(infoA, infoC)

	return CompareResponse{
		Schema: schema,
		AtoB: SchemaDiffSummary{
			DBA:        cfgA.Name,
			DBB:        cfgB.Name,
			TableDiffs: resultAtoB.TableDiffs,
		},
		AtoC: SchemaDiffSummary{
			DBA:        cfgA.Name,
			DBB:        cfgC.Name,
			TableDiffs: resultAtoC.TableDiffs,
		},
	}, nil
}

func loadDBSchema(db *sql.DB, name, schema string) (DBInfo, error) {
	tablesQuery := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = $1
		AND table_type = 'BASE TABLE'
		ORDER BY table_name;
	`

	rows, err := db.Query(tablesQuery, schema)
	if err != nil {
		return DBInfo{}, err
	}
	defer rows.Close()

	tables := make(map[string]TableInfo)

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return DBInfo{}, err
		}
		tables[tableName] = TableInfo{
			Name:    tableName,
			Columns: make(map[string]ColumnInfo),
		}
	}
	if err := rows.Err(); err != nil {
		return DBInfo{}, err
	}

	columnsQuery := `
		SELECT table_name,
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_schema = $1
		ORDER BY table_name, ordinal_position;
	`

	colRows, err := db.Query(columnsQuery, schema)
	if err != nil {
		return DBInfo{}, err
	}
	defer colRows.Close()

	for colRows.Next() {
		var tbl, col, dataType, isNullable string
		var defaultVal sql.NullString
		if err := colRows.Scan(&tbl, &col, &dataType, &isNullable, &defaultVal); err != nil {
			return DBInfo{}, err
		}

		t, ok := tables[tbl]
		if !ok {
			continue
		}

		var defPtr *string
		if defaultVal.Valid {
			s := strings.TrimSpace(defaultVal.String)
			defPtr = &s
		}

		t.Columns[col] = ColumnInfo{
			Name:       col,
			DataType:   dataType,
			IsNullable: strings.EqualFold(isNullable, "YES"),
			Default:    defPtr,
		}
		tables[tbl] = t
	}
	if err := colRows.Err(); err != nil {
		return DBInfo{}, err
	}

	return DBInfo{
		Name:   name,
		Tables: tables,
	}, nil
}

func CompareSchemas(a, b DBInfo) SchemaDiffSummary {
	var diffs []TableDiff

	for name, tblA := range a.Tables {
		tblB, ok := b.Tables[name]
		if !ok {
			diffs = append(diffs, TableDiff{
				Table:   name,
				OnlyInA: true,
			})
			continue
		}

		colDiffs := compareColumns(tblA.Columns, tblB.Columns)
		if len(colDiffs) > 0 {
			diffs = append(diffs, TableDiff{
				Table:       name,
				ColumnDiffs: colDiffs,
			})
		}
	}

	for name := range b.Tables {
		if _, ok := a.Tables[name]; !ok {
			diffs = append(diffs, TableDiff{
				Table:   name,
				OnlyInB: true,
			})
		}
	}

	return SchemaDiffSummary{
		DBA:        a.Name,
		DBB:        b.Name,
		TableDiffs: diffs,
	}
}

func compareColumns(colsA, colsB map[string]ColumnInfo) []ColumnDiff {
	var diffs []ColumnDiff

	for name, cA := range colsA {
		cB, ok := colsB[name]
		if !ok {
			diffs = append(diffs, ColumnDiff{
				Column: name,
				InA:    &cA,
				InB:    nil,
			})
			continue
		}

		if !equalColumn(cA, cB) {
			cA2 := cA
			cB2 := cB
			diffs = append(diffs, ColumnDiff{
				Column: name,
				InA:    &cA2,
				InB:    &cB2,
			})
		}
	}

	for name, cB := range colsB {
		if _, ok := colsA[name]; !ok {
			cB2 := cB
			diffs = append(diffs, ColumnDiff{
				Column: name,
				InA:    nil,
				InB:    &cB2,
			})
		}
	}

	return diffs
}

func equalColumn(a, b ColumnInfo) bool {
	if !strings.EqualFold(a.DataType, b.DataType) {
		return false
	}
	if a.IsNullable != b.IsNullable {
		return false
	}

	switch {
	case a.Default == nil && b.Default == nil:
		return true
	case a.Default == nil && b.Default != nil:
		return false
	case a.Default != nil && b.Default == nil:
		return false
	default:
		return strings.TrimSpace(*a.Default) == strings.TrimSpace(*b.Default)
	}
}

func loadEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found, using system env instead")
	}
}

func loadOneConfig(prefix string) (DBConfig, error) {
	cfg := DBConfig{
		Name: os.Getenv(fmt.Sprintf("DB_%s_NAME", prefix)),
		URL:  buildPostgresURL(prefix),
	}

	if cfg.URL == "" {
		return DBConfig{}, errors.New("missing db config for " + prefix)
	}

	return cfg, nil
}

func buildPostgresURL(suffix string) string {
	host := os.Getenv("POSTGRES_HOST_" + suffix)
	port := os.Getenv("POSTGRES_PORT_" + suffix)
	user := os.Getenv("POSTGRES_USER_" + suffix)
	pass := os.Getenv("POSTGRES_PASSWORD_" + suffix)
	db := os.Getenv("POSTGRES_DB_" + suffix)
	ssl := os.Getenv("POSTGRES_SSLMODE_" + suffix)

	if port == "" {
		port = "5432"
	}
	if ssl == "" {
		ssl = "disable"
	}

	if host == "" || user == "" || db == "" {
		return ""
	}

	return makePgURL(host, port, user, pass, db, ssl)
}

func makePgURL(host, port, user, pass, db, ssl string) string {
	u := &url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   db,
	}
	u.User = url.UserPassword(user, pass)

	q := u.Query()
	q.Set("sslmode", ssl)
	u.RawQuery = q.Encode()

	return u.String()
}
