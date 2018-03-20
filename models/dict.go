package models

// SearchParam : Search Parameter Information
type SearchParam struct {
	Type      string   //
	FieldType string   //
	Path      []string //
	Target    []string
	Join      []string
	Select    []string
}

// Dictionary model
type Dictionary struct {
	Model              map[string]map[string]SearchParam
	Bucket             string
	ResourceIdentifier string
}

// DictionaryInfo : Contains elaborated information about the fields exists in the dictionary
type DictionaryInfo struct {
	FHIRFieldType   string
	FHIRType        string
	SelectStatement []string
	FieldsInfo      struct {
		ArrayPath  []string
		ObjectPath string
	}
}
