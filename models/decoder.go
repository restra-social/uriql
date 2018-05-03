package models

// FieldInfo : Json field information matched with query parameter to search
type FieldInfo struct {
	Field  string // Object field information, bought from dictionary based on query parameter
	Array  bool   // If the Field is an Array or Not
	Object bool   // If the Field is an Object or Not
}

// QueryInfo stores the decoded information about the query its slice of slice because of multiple
// Query or Composite Query Parameter
// Param length is the number of length passed by query parameter Ex : patient?name=fahim&dob=15-25-69&status=true | Param length is 3
// The First slice represents the inner query path lets say if `name` parameter needs to searched in multiple json path then the number of length
// Of []QueryParameter is the number of path to search
type QueryInfo struct {
	Params [][]QueryParamInfo
	Filter *Filter
}

// QueryParam : Decoded information about the Query Parameter
type QueryParamInfo struct {
	RequestInfo
	DictionaryInfo
	Resource     string
	Condition    string
	Value        ValueType
	SearchResult SearchResult
	Path         string
}

// ValueType : Contains the value parameter from the query
type ValueType struct {
	Codable struct {
		System string
		Code   string
	}
	Value     string
	Modifiers string
}
