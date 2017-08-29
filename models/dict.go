package models

// SearchParam : Search Parameter Information
type SearchParam struct {
	Type      string   //
	FieldType string   //
	Path      []string //
	Target    []string
}

// Dictionary model
type Dictionary struct {
	Model  map[string]map[string]SearchParam
	Bucket string
}
