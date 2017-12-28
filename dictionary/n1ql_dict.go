package dictionary

import "github.com/restra-social/uriql/models"

func N1QLDictionary() map[string]map[string]models.SearchParam {

	dict := map[string]map[string]models.SearchParam{

		"profile": map[string]models.SearchParam{

			"_id": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"id"},
			},
			"gender": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"gender"},
			},
			"name": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"name.prefix", "name.first_name", "name.last_name"},
			},
			"address": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"address.city.name", "address.state.name", "address.street"},
			},
			"hobbies": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"[]hobbies"},
			},
			"postal": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"address.postal"},
			},
		},

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
			"status": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"status"},
			},
		},
	}

	return dict
}
