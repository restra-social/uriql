package dictionary

import "udhvabon.com/neuron/uriql/models"

type Def struct {
}

func GetDef() *Def {
	return &Def{}
}

func (f *Def) MatchSearchParam(resource, match string) *models.SearchParam {
	switch resource {

	case "Patient" :
		sp := map[string]models.SearchParam{
			"active": models.SearchParam{
				Type:      "token",
				FieldType: "boolean",
				Path:      []string{"active"},
			},
			"identifier" : models.SearchParam{
				Type: "token",
				FieldType: "identifier",
				Path: []string{"[]identifier.value"},
			},
			"gender" : models.SearchParam{
				Type: "token",
				FieldType: "code",
				Path: []string{"gender"},
			},
			"name": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"[]name.[]family", "[]name.[]given"},
			},
			"address-state": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"[]address.state"},
			},
			"address-use" : models.SearchParam{
				Type: 		"token",
				FieldType: 	"string",
				Path: []string{"[]address.use"},
			},
			"language": models.SearchParam{
				Type:      "token",
				FieldType: "coding",
				Path:      []string{"[]communication.language"},
			},
			"general-practitioner": models.SearchParam{
				Type: "reference",
				FieldType: "string",
				Path: []string{"[]generalPractitioner.reference"},
			},
			"organization": models.SearchParam{
				Type: "reference",
				FieldType: "string",
				Path: []string{"managingOrganization.reference"},
			},
		}

		if val, ok := sp[match]; ok {
			return &val
		}
	case "Practitioner" :

	case "Encounter" :
		sp := map[string]models.SearchParam{
			"length" : models.SearchParam{
				Type: "number",
				FieldType: "number",
				Path: []string{"length"},
			},
		}

		if val, ok := sp[match]; ok {
			return &val
		}
	}

	return &models.SearchParam{}
}