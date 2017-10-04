package dictionary

import "github.com/kite-social/uriql/models"

// N1QLDictionary : Example Dictionary for N1QL Query Builder
func N1QLDictionary() map[string]map[string]models.SearchParam {

	dict := map[string]map[string]models.SearchParam{

		"Patient": map[string]models.SearchParam{

			/*"active": models.SearchParam{
				Type:      "token",
				FieldType: "boolean",
				Path:      []string{"active"},
			},*/
			"identifier": models.SearchParam{
				Type:      "token",
				FieldType: "identifier",
				Path:      []string{"[]identifier.system", "[]identifier.value"},
			},
			"language": models.SearchParam{
				Type:      "token",
				FieldType: "coding",
				Path:      []string{"[]communication.language.[]coding.system", "[]communication.language.[]coding.code"},
			},
		},

		/*		"Encounter": map[string]models.SearchParam{
				"length": models.SearchParam{
					Type:      "number",
					FieldType: "number",
					Path:      []string{"length"},
				},
			},*/
	}

	return dict
}
