package models

// FieldInfo : Json field information matched with query parameter to search
type FieldInfo struct {
	Field  string // Object field information, bought from dictionary based on query parameter
	Array  bool   // If the Field is an Array or Not
	Object bool   // If the Field is an Object or Not
}

// QueryParam : Decoded information about the Query Parameter
type QueryParam struct {
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
