package helper

import "github.com/kite-social/uriql/models"

type Def struct {
	Dictionary *models.Dictionary
}

func GetDef(dictionary *models.Dictionary) *Def {
	return &Def{
		Dictionary: dictionary,
	}
}

func (f *Def) MatchSearchParam(resource, match string) *models.SearchParam {
		if res, ok := f.Dictionary.Model[resource]; ok {
			if res, ok := res[match]; ok {
				return &res
			}
		}
	return &models.SearchParam{}
}