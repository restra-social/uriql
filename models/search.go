package models

type SearchField struct {
	Object string
	Array  string
	Field  FieldInfo
}

type SearchParam struct {
	Type      string
	FieldType string
	Path      []string
}

type FieldInfo struct {
	Field string
	Array bool
}

type QueryParam struct {
	Resource      string
	Object        []string
	Array         []string
	Field         []FieldInfo
	FHIRFieldType string
	FHIRType      string
	Condition     string
	Value         []string
	SearchResult  SearchResult
}
type SearchResult struct {
	Type       string
	Sorting    []Sort
	Count      int64
	Include    []string
	RevInclude []string
}

type Sort struct {
	Field string
	Type  string
}
