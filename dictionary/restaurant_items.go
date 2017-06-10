package dictionary

import "udhvabon.com/neuron/uriql/models"

func RestaurantItemsDictionary() map[string]map[string]models.SearchParam {

	dict := map[string]map[string]models.SearchParam{

		"restaurant-items": map[string]models.SearchParam{

			"_id" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"id"},
			},
			"_rid": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"r_id"},
			},
			"category": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"category.name"},
			},
			"category-id": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"category.id"},
			},
			"foods" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"[]foods.name"},
			},
			"foods-size" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"[]foods.size"},
			},
			"foods-price" : models.SearchParam{
				Type: "string",
				FieldType: "string",
				Path: []string{"[]foods.price"},
			},
		},
	}

	return dict
}
