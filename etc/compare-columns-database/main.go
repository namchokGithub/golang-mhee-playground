package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	_ "github.com/lib/pq"
)

type Column struct {
	Name string
	Type string
}

func fetchColumns(db *sql.DB, tableName string) (map[string]string, error) {
	query := `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_name = $1
	`
	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make(map[string]string)
	for rows.Next() {
		var name, dataType string
		if err := rows.Scan(&name, &dataType); err != nil {
			return nil, err
		}
		columns[name] = dataType
	}
	return columns, nil
}

func compareColumns(columns1, columns2 map[string]string) {
	fmt.Println("🟦 Differences:")

	allKeys := make(map[string]bool)
	for k := range columns1 {
		allKeys[k] = true
	}
	for k := range columns2 {
		allKeys[k] = true
	}

	var keys []string
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		t1, ok1 := columns1[key]
		t2, ok2 := columns2[key]

		switch {
		case ok1 && !ok2:
			fmt.Printf("- Column '%s' is missing in DB2\n", key)
		case !ok1 && ok2:
			fmt.Printf("- Column '%s' is missing in DB1\n", key)
		case t1 != t2:
			fmt.Printf("- Column '%s' has different types: DB1='%s', DB2='%s'\n", key, t1, t2)
		default:
			// match
		}
	}
}

func main() {
	db1, err := sql.Open("postgres", "host=localhost port=5432 user=youruser password=yourpass dbname=db1 sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB1: %v", err)
	}
	defer db1.Close()

	db2, err := sql.Open("postgres", "host=localhost port=5432 user=youruser password=yourpass dbname=db2 sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB2: %v", err)
	}
	defer db2.Close()

	tableName := "your_table"

	columns1, err := fetchColumns(db1, tableName)
	if err != nil {
		log.Fatalf("DB1 column fetch error: %v", err)
	}

	columns2, err := fetchColumns(db2, tableName)
	if err != nil {
		log.Fatalf("DB2 column fetch error: %v", err)
	}

	compareColumns(columns1, columns2)
}
