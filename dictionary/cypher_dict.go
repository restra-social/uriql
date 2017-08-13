package dictionary

import "github.com/kite-social/uriql/models"

// CypherDictionary : Example Dictionary for Cypher Query builder
func CypherDictionary() map[string]map[string]models.SearchParam {

	dict := map[string]map[string]models.SearchParam{

		"friend": map[string]models.SearchParam{

			"_id": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"id"},
			},
			"status": models.SearchParam{
				Type:      "relation",
				FieldType: "string",
				Path:      []string{"status"},
			},
			"since": models.SearchParam{
				Type:      "relation",
				FieldType: "string",
				Path:      []string{"since"},
			},
			"user": models.SearchParam{
				Type:      "node",
				FieldType: "string",
				Path:      []string{"user"},
			},
		},
	}

	return dict
}
