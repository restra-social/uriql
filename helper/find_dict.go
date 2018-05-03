package helper

import (
	"errors"
	"fmt"
	"github.com/restra-social/uriql/models"
)

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
func (f *Def) MatchSearchParam(resource, match string) (*models.SearchParam, error) {
	if res, ok := f.Dictionary.Model[resource]; ok {
		if res, ok := res[match]; ok {
			return &res, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Could not find field [%s] in the %s dictionary", match, resource))
}
