package schemaground

type Comparator interface {
	Compare(schema string) (CompareResponse, error)
}
