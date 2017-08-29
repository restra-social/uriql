package dictionary

import "github.com/kite-social/uriql/models"

// N1QLDictionary : Example Dictionary for N1QL Query Builder
func N1QLDictionary() map[string]map[string]models.SearchParam {

	dict := map[string]map[string]models.SearchParam{

		"restaurant": map[string]models.SearchParam{

			"_id": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"id"},
			},
			"title": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"title"},
			},
			"address": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"address.city.name", "address.state.name", "address.postal"},
			},
			"address-street": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"address.street"},
			},
			"phone": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"social.phone"},
			},
			"text": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"description"},
			},
			"verified": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"verified"},
			},
			"status": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"status"},
			},
		},

		"Patient": map[string]models.SearchParam{

			"language": models.SearchParam{
				Type:      "token",
				FieldType: "coding",
				Path:      []string{"[]communication.language"},
			},
		},

		"Observation": map[string]models.SearchParam{

			"subject": models.SearchParam{
				Type:      "reference",
				FieldType: "string",
				Path:      []string{"subject.reference"},
			},
		},


	}

	return dict
}
