package models

type SearchParam struct {
	Type      string
	FieldType string
	Path      []string
}

type FieldInfo struct {
	Field  string
	Array  bool
	Object bool
}

type QueryParam struct {
	Resource      string
	ArrayCount    int
	Field         []FieldInfo
	FHIRFieldType string
	FHIRType      string
	Condition     string
	Value         []string
	SearchResult  SearchResult
	Path          string
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
