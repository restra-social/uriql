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
				Path:      []string{"[]identifier.value", "[]identifier.system"},
			},
			"name": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"[]name.family", "[]name.[]given"},
			},
			"language": models.SearchParam{
				Type:      "token",
				FieldType: "coding",
				Path:      []string{"[]communication.language.[]coding.code"},
			},
			"organization": models.SearchParam{
				Type:      "reference",
				FieldType: "string",
				Path:      []string{"managingOrganization.[]reference"},
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
