package dictionary

import "udhvabon.com/neuron/uriql/models"

func RestaurantDictionary() map[string]map[string]models.SearchParam {

	dict := map[string]map[string]models.SearchParam{

		"restaurant": map[string]models.SearchParam{

			"_id" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"id"},
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
			"phone" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"social.phone"},
			},
			"text" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"description"},
			},
			"verified" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"verified"},
			},
			"status" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"status"},
			},
			"test" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"[]medium.[]cost.type"},
			},
		},
	}

	return dict
}
