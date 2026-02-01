package schemaground_testing

import (
	"testing"

	"proundmhee/internal/modules/schemaground"

	"github.com/stretchr/testify/require"
)

func TestCompareSchemas(t *testing.T) {
	a := schemaground.DBInfo{
		Name: "A",
		Tables: map[string]schemaground.TableInfo{
			"users": {
				Name: "users",
				Columns: map[string]schemaground.ColumnInfo{
					"id":   {Name: "id", DataType: "int", IsNullable: false},
					"name": {Name: "name", DataType: "text", IsNullable: true},
				},
			},
		},
	}

	b := schemaground.DBInfo{
		Name: "B",
		Tables: map[string]schemaground.TableInfo{
			"users": {
				Name: "users",
				Columns: map[string]schemaground.ColumnInfo{
					"id": {Name: "id", DataType: "int", IsNullable: false},
				},
			},
			"posts": {
				Name: "posts",
				Columns: map[string]schemaground.ColumnInfo{
					"id": {Name: "id", DataType: "int", IsNullable: false},
				},
			},
		},
	}

	res := schemaground.CompareSchemas(a, b)
	require.Equal(t, "A", res.DBA)
	require.Equal(t, "B", res.DBB)
	require.Len(t, res.TableDiffs, 2)

	var foundUsersDiff bool
	var foundPostsOnlyInB bool
	for _, td := range res.TableDiffs {
		if td.Table == "users" && len(td.ColumnDiffs) == 1 {
			foundUsersDiff = true
		}
		if td.Table == "posts" && td.OnlyInB {
			foundPostsOnlyInB = true
		}
	}

	require.True(t, foundUsersDiff)
	require.True(t, foundPostsOnlyInB)
}
