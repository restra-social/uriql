package models

// RequestInfo : Query Request Information
type RequestInfo struct {
	UserID string // The user who is requesting a query E.g. 123456
	Type   string // The type of resource the user is requesting E.g. people
	Query  string // The query String E.g. people?name:contains=Jon
	Filter
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

// Filter represents Search Filter
// How many pages to show , whats the limit etc
type Filter struct {
	Page   int // Page Number , 1 or 2 etc
	Limit  int // Per Page Limit like 10 or 20
	OrderBy map[string]string
}
