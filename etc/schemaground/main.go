package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

// ----------------- MAIN -----------------

func main() {
	cfgA := loadOneConfigs("A")
	cfgB := loadOneConfigs("B")
	cfgC := loadOneConfigs("C")

	dbA, err := sql.Open("postgres", cfgA.URL)
	if err != nil {
		log.Fatalf("connect A error: %v", err)
	}
	defer dbA.Close()

	dbB, err := sql.Open("postgres", cfgB.URL)
	if err != nil {
		log.Fatalf("connect B error: %v", err)
	}
	defer dbB.Close()

	dbC, err := sql.Open("postgres", cfgC.URL)
	if err != nil {
		log.Fatalf("connect C error: %v", err)
	}
	defer dbB.Close()

	// ตั้งค่า schema ที่อยากเทียบ (ส่วนใหญ่คือ "public")
	schemaName := "public"

	infoA, err := loadDBSchema(dbA, cfgA.Name, schemaName)
	if err != nil {
		log.Fatalf("load schema A error: %v", err)
	}

	infoB, err := loadDBSchema(dbB, cfgB.Name, schemaName)
	if err != nil {
		log.Fatalf("load schema B error: %v", err)
	}

	infoC, err := loadDBSchema(dbC, cfgC.Name, schemaName)
	if err != nil {
		log.Fatalf("load schema C error: %v", err)
	}

	// A -> B
	resultAtoB := compareSchemas(infoA, infoB)
	// A -> C
	resultAtoC := compareSchemas(infoA, infoC)

	// schemaNameA := fmt.Sprintf("schema_%s.json", cfgA.Name)
	// if err := writeJSONFile(schemaNameA, resultAtoB.DBA); err != nil {
	// 	log.Fatalf("write json error: %v", err)
	// }
	// schemaNameB := fmt.Sprintf("schema_%s.json", cfgB.Name)
	// if err := writeJSONFile(schemaNameB, resultAtoB.DBB); err != nil {
	// 	log.Fatalf("write json error: %v", err)
	// }
	// schemaNameC := fmt.Sprintf("schema_%s.json", cfgC.Name)
	// if err := writeJSONFile(schemaNameC, resultAtoC.DBA); err != nil {
	// 	log.Fatalf("write json error: %v", err)
	// }

	// Schema between A and B
	// schemaDiffAToBName := fmt.Sprintf("schema_diff_[%s]_to_[%s].json", cfgA.Name, cfgB.Name)
	// if err := writeJSONFile(schemaDiffAToBName, resultAtoB.TableDiffs); err != nil {
	// 	log.Fatalf("write json error: %v", err)
	// }
	// Schema between A and C
	schemaDiffAToCName := fmt.Sprintf("schema_diff_[%s]_to_[%s].json", cfgA.Name, cfgC.Name)
	if err := writeJSONFile(schemaDiffAToCName, resultAtoC.TableDiffs); err != nil {
		log.Fatalf("write json error: %v", err)
	}
	fmt.Println("✅ Done. Output -> schema_diff.json")

	sort.Slice(resultAtoB.TableDiffs, func(i, j int) bool {
		return resultAtoB.TableDiffs[i].Table < resultAtoB.TableDiffs[j].Table
	})
	excelSchemaDiffAToBName := fmt.Sprintf("schema_diff_[%s]_to_[%s].xlsx", cfgA.Name, cfgB.Name)
	if err := writeExcelDiff(excelSchemaDiffAToBName, resultAtoB, cfgA, cfgB); err != nil {
		log.Fatalf("write excel error: %v", err)
	}
	sort.Slice(resultAtoC.TableDiffs, func(i, j int) bool {
		return resultAtoC.TableDiffs[i].Table < resultAtoC.TableDiffs[j].Table
	})
	excelSchemaDiffAToCName := fmt.Sprintf("schema_diff_[%s]_to_[%s].xlsx", cfgA.Name, cfgC.Name)
	if err := writeExcelDiff(excelSchemaDiffAToCName, resultAtoC, cfgA, cfgC); err != nil {
		log.Fatalf("write excel error: %v", err)
	}
	fmt.Println("✅ Excel diff -> schema_diff.xlsx")

}

func writeExcelDiff(filename string, res SchemaDiffResult, cfgA, cfgB DBConfig) error {
	f := excelize.NewFile()

	// ---------- Sheet 1: TableDiffs ----------
	const sheetTable = "TableDiffs"
	index, err := f.NewSheet(sheetTable)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// header
	_ = f.SetSheetRow(sheetTable, "A1", &[]interface{}{
		"Table",
		fmt.Sprintf("Only In %s", cfgA.Name),
		fmt.Sprintf("Only In %s", cfgB.Name),
		"HasColumnDiffs",
	})

	row := 2
	for _, td := range res.TableDiffs {
		hasColDiff := len(td.ColumnDiffs) > 0
		cell := fmt.Sprintf("A%d", row)
		_ = f.SetSheetRow(sheetTable, cell, &[]interface{}{
			td.Table,
			td.OnlyInA,
			td.OnlyInB,
			hasColDiff,
		})
		row++
	}

	// --------------------------
	// 🎨 1) Header Style (ฟ้าอ่อน)
	// --------------------------
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#004b8d"}, // ฟ้าอ่อน
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	rolSheetTable := row

	// ---------- Sheet 2: ColumnDiffs ----------
	const sheetCol = "ColumnDiffs"
	f.NewSheet(sheetCol)

	_ = f.SetSheetRow(sheetCol, "A1", &[]interface{}{
		"Table",
		"Column",
		fmt.Sprintf("In %s Data Type", cfgA.Name),
		fmt.Sprintf("In %s Is Nullable", cfgA.Name),
		fmt.Sprintf("In %s Default", cfgA.Name),
		fmt.Sprintf("In %s Data Type", cfgB.Name),
		fmt.Sprintf("In %s Is Nullable", cfgB.Name),
		fmt.Sprintf("In %s Default", cfgB.Name),
	})

	row = 2
	for _, td := range res.TableDiffs {
		for _, cd := range td.ColumnDiffs {
			var (
				aType, aNull, aDef interface{}
				bType, bNull, bDef interface{}
			)

			if cd.InA != nil {
				aType = cd.InA.DataType
				aNull = cd.InA.IsNullable
				if cd.InA.Default != nil {
					aDef = *cd.InA.Default
				}
			}
			if cd.InB != nil {
				bType = cd.InB.DataType
				bNull = cd.InB.IsNullable
				if cd.InB.Default != nil {
					bDef = *cd.InB.Default
				}
			}

			cell := fmt.Sprintf("A%d", row)
			_ = f.SetSheetRow(sheetCol, cell, &[]interface{}{
				td.Table,
				cd.Column,
				aType,
				aNull,
				aDef,
				bType,
				bNull,
				bDef,
			})
			row++
		}
	}
	rowSheetColumn := row

	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	// apply border to TableDiffs
	f.SetCellStyle(sheetTable, "A1", fmt.Sprintf("D%d", rolSheetTable-1), borderStyle)

	// apply border to ColumnDiffs
	f.SetCellStyle(sheetCol, "A1", fmt.Sprintf("H%d", rowSheetColumn-1), borderStyle)

	// --------------------------
	// 🎨 2) Style TRUE → สีแดง
	// --------------------------
	trueStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#FF0000", // red
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})

	// --------------------------
	// 🎨 3) Style FALSE → สีเทาจาง
	// --------------------------
	falseStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#A0A0A0", // light gray
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})

	// Apply style for TRUE/FALSE
	for r := 2; r < row; r++ {
		for col := 2; col <= 4; col++ { // columns B, C, D
			colRune := 'A' + rune(col-1)            // คำนวณตัวอักษร column
			cell := fmt.Sprintf("%c%d", colRune, r) // %c สำหรับ rune
			val, _ := f.GetCellValue(sheetTable, cell)
			switch val {
			case "TRUE":
				f.SetCellStyle(sheetTable, cell, cell, trueStyle)
			case "FALSE":
				f.SetCellValue(sheetTable, cell, "-")
				f.SetCellStyle(sheetTable, cell, cell, falseStyle)
			}
		}
	}

	colAStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#DDEBF7"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	// For TableDiffs
	lastRowSheetTable := rolSheetTable - 1
	if lastRowSheetTable >= 2 {
		_ = f.SetCellStyle(sheetTable, "A2", fmt.Sprintf("A%d", rolSheetTable-1), colAStyle)
	}
	f.SetCellStyle(sheetTable, "A1", "D1", headerStyle)

	// For ColumnDiffs
	lastRowSheetColumn := rowSheetColumn - 1
	if lastRowSheetColumn >= 2 {
		_ = f.SetCellStyle(sheetCol, "A2", fmt.Sprintf("A%d", rowSheetColumn-1), colAStyle)
	}
	f.SetCellStyle(sheetCol, "A1", "H1", headerStyle)

	// ----------------------------------------------
	// 🎨 ColumnDiffs สีสำหรับแต่ละ Environment
	// ----------------------------------------------

	// สี A (C–E)
	aHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#004610"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})

	// สี B (F–H)
	bHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#805800"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})

	// C:E คือ A
	f.SetCellStyle(sheetCol, "C1", "E1", aHeaderStyle)
	// F:H คือ B
	f.SetCellStyle(sheetCol, "F1", "H1", bHeaderStyle)

	// auto width คร่าว ๆ
	for _, sh := range []string{sheetTable, sheetCol} {
		_ = f.SetColWidth(sh, "A", "H", 20)
	}

	// =========================
	// Freeze Row 1 + Column A
	// =========================
	f.SetPanes(sheetTable, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      1,
		YSplit:      1,
		TopLeftCell: "B2",
		ActivePane:  "bottomRight",
	})

	f.SetPanes(sheetCol, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      1,
		YSplit:      1,
		TopLeftCell: "B2",
		ActivePane:  "bottomRight",
	})

	// =========================
	// AutoFilter ที่ Row 1
	// =========================
	// AutoFilter row 1 (TableDiffs: A–D)
	_ = f.AutoFilter(sheetTable, fmt.Sprintf("A1:D%d", rolSheetTable-1), nil)

	// AutoFilter row 1 (ColumnDiffs: A–H)
	_ = f.AutoFilter(sheetCol, fmt.Sprintf("A1:H%d", rowSheetColumn-1), nil)

	f.DeleteSheet("Sheet1")
	return f.SaveAs(filename)
}

// ----------------- LOAD SCHEMA -----------------

func loadDBSchema(db *sql.DB, name, schema string) (DBInfo, error) {
	// ดึง list table เฉพาะ BASE TABLE
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

	// ดึง columns ทั้งหมด
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
			// ถ้าเป็น view หรือ table ที่ไม่อยู่ใน BASE TABLE ที่เราดึงไว้ ก็ข้าม
			continue
		}

		var defPtr *string
		if defaultVal.Valid {
			// ทำให้ short หน่อย (ตัด newline / space)
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

// ----------------- COMPARE LOGIC -----------------

func compareSchemas(a, b DBInfo) SchemaDiffResult {
	var diffs []TableDiff

	// 1) วิ่งตาม tables ของ A
	for name, tblA := range a.Tables {
		tblB, ok := b.Tables[name]
		if !ok {
			// มีเฉพาะฝั่ง A
			diffs = append(diffs, TableDiff{
				Table:   name,
				OnlyInA: true,
			})
			continue
		}

		// มีทั้งสองฝั่ง -> เทียบ columns
		colDiffs := compareColumns(tblA.Columns, tblB.Columns)
		if len(colDiffs) > 0 {
			diffs = append(diffs, TableDiff{
				Table:       name,
				ColumnDiffs: colDiffs,
			})
		}
	}

	// 2) หา tables ที่มีเฉพาะใน B
	for name := range b.Tables {
		if _, ok := a.Tables[name]; !ok {
			diffs = append(diffs, TableDiff{
				Table:   name,
				OnlyInB: true,
			})
		}
	}

	return SchemaDiffResult{
		DBA:        a,
		DBB:        b,
		TableDiffs: diffs,
	}
}

func compareColumns(colsA, colsB map[string]ColumnInfo) []ColumnDiff {
	var diffs []ColumnDiff

	// เช็คจาก A เป็นหลัก
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

	// หา columns ที่มีเฉพาะใน B
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

	// เทียบ default แบบ pointer
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

// ----------------- JSON OUTPUT -----------------

func writeJSONFile(filename string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// ----------------- CONFIG ZONE -----------------

// ใช้ ENV แยกแต่ละ environment ไปเลย
//
//	export DB_A_URL="postgres://user:pass@host:5432/dbname?sslmode=disable"
//	export DB_B_URL="postgres://user:pass@host:5432/dbname?sslmode=disable"
//
// สามารถแก้ให้ไปอ่านจาก config package ของโปรเจกต์จริงได้
type DBConfig struct {
	Name string // label สั้นๆ เช่น "staging", "prod"
	URL  string // postgres connection URL
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("⚠️  No .env file found, using system env instead")
	}
}

func loadOneConfigs(prefix string) DBConfig {
	url := buildPostgresURL(prefix)
	cfgA := DBConfig{
		Name: os.Getenv(fmt.Sprintf("DB_%s_NAME", prefix)),
		URL:  url,
	}
	fmt.Println(fmt.Sprintf("DB_%s_URL =", prefix), cfgA.URL)

	if cfgA.URL == "" {
		log.Fatalf("❌ Missing DB config, please set POSTGRES_* variables for %s", prefix)
	}

	return cfgA
}

func buildPostgresURL(suffix string) string {
	// suffix = "A" หรือ "B"
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
		// ถ้า key สำคัญว่าง ให้ return "" ไว้เช็คทีหลัง
		return ""
	}

	url := makePgURL(host, port, user, pass, db, ssl)
	return url
}

func makePgURL2(host, port, user, pass, db, ssl string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, db, ssl,
	)
}

func makePgURL(host, port, user, pass, db, ssl string) string {
	u := &url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   db,
	}
	u.User = url.UserPassword(user, pass) // Go จะ encode ให้เอง

	q := u.Query()
	q.Set("sslmode", ssl)
	u.RawQuery = q.Encode()

	return u.String()
}

// ----------------- SCHEMA STRUCTS -----------------

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

// Diffs

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

type SchemaDiffResult struct {
	DBA        DBInfo      `json:"db_a"`
	DBB        DBInfo      `json:"db_b"`
	TableDiffs []TableDiff `json:"table_diffs"`
}
