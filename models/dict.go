package models

type Dictionary struct {
	Model map[string]map[string]SearchParam
	Bucket string
}