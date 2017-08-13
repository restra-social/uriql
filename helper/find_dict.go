package helper

import "github.com/kite-social/uriql/models"

type Def struct {
	Dictionary map[string]map[string]models.SearchParam
}

func GetDef(dictionary map[string]map[string]models.SearchParam) *Def {
	return &Def{
		Dictionary: dictionary,
	}
}

func (f *Def) MatchSearchParam(resource, match string) *models.SearchParam {
		if res, ok := f.Dictionary[resource]; ok {
			if res, ok := res[match]; ok {
				return &res
			}
		}
	return &models.SearchParam{}
}