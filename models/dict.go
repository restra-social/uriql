package models


// Dictionary model
type Dictionary struct {
	Model  map[string]map[string]SearchParam
	Bucket string
}
