package helper

import "github.com/restra-social/uriql/models"

// Def :
type Def struct {
	Dictionary *models.Dictionary
}

// GetDef : Get the dictionary
func GetDef(dictionary *models.Dictionary) *Def {
	return &Def{
		Dictionary: dictionary,
	}
}

// MatchSearchParam : Finds query parameter information from dictionary
func (f *Def) MatchSearchParam(resource, match string) *models.SearchParam {
	if res, ok := f.Dictionary.Model[resource]; ok {
		if res, ok := res[match]; ok {
			return &res
		}
	}
	return &models.SearchParam{}
}
