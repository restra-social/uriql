package models

// RequestInfo : Query Request Information
type RequestInfo struct {
	UserID string // The user who is requesting a query E.g. 123456
	Type   string // The type of resource the user is requesting E.g. people
	Query  string // The query String E.g. people?name:contains=Jon
}

// FieldInfo : Json field information matched with query parameter to search
type FieldInfo struct {
	Field  string // Object field information, bought from dictionary based on query parameter
	Array  bool   // If the Field is an Array or Not
	Object bool   // If the Field is an Object or Not
}

// QueryParam : Decoded information about the Query Parameter
type QueryParam struct {
	RequestInfo
	Resource      string
	ArrayCount    int
	Field         []FieldInfo
	FHIRFieldType string
	FHIRType      string
	Condition     string
	Value         ValueType
	SearchResult  SearchResult
	Path          string
}

// ValueType : Contains the value parameter from the query
type ValueType struct {
	Reference struct {
		Target string
		Value  string
	}
	Codable struct {
		System string
		Code   string
	}
	Value     string
	Modifiers string
}

// SearchResult : Parameter to store search filter information
type SearchResult struct {
	Type       string
	Sorting    []Sort
	Count      int64
	Include    []string
	RevInclude []string
}

// Sort : Sorting information for Search Result
type Sort struct {
	Field string
	Type  string
}
