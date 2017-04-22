package uriql

import (
	"udhvabon.com/neuron/uriql/models"
	"fmt"
)

type N1QLQueryBuilder struct {
}

func GetN1QLBuilder() *N1QLQueryBuilder {
	return &N1QLQueryBuilder{}
}

func (n *N1QLQueryBuilder) Build(model models.QueryParam) string {

	var str string
	// todo#fix fix condtion parameter
	// where ANY n IN %s satisfies (any name in n.`%s` SATISFIES name LIKE '%%%s%%' END) end;
	str += fmt.Sprintf("select * from `neuron` as r where r.`resourceType` = '%s' and ", model.Resource)

	var conNVal string
	//#todo#fix token condition need to be fixed
	switch model.Condition {
	case "like":
		conNVal = fmt.Sprintf("%s %%%s%%", model.Condition, model.Value[0])
	case "=":
		conNVal = fmt.Sprintf("%s '%s'", model.Condition, model.Value[0])
	default :
		conNVal = fmt.Sprintf("%s '%s'", model.Condition, model.Value[0])
	}

	switch model.FHIRType {

	case "universal":
		if len(model.Object) > 0 {
			if model.Field[0].Array == false {
				str += fmt.Sprintf("r.%s.%s %s", model.Object[0], model.Field[0].Field, conNVal)

			} else {
				// not yet for _profile , _security and _tag
			}
		} else {
			str += fmt.Sprintf("r.%s %s", model.Field[0].Field, conNVal)
		}

	case "single":
		str += fmt.Sprintf("META(r).id %s", conNVal)
	case "number":
		// found just field so far
		str += fmt.Sprintf("r.%s %s", model.Field[0].Field, conNVal)
	case "string":
		str += buildArrayQuery(model, conNVal)
	case "reference":
		if len(model.Array) > 0 {
			str += fmt.Sprintf("any ref in SPLIT(r.`%s`.`reference`,'/') satisfies ref %s end;", model.Array[0], conNVal)
		}else if len(model.Object) > 0 {
			str += fmt.Sprintf("any ref in r.`%s` satisfies (any org in SPLIT(ref.`reference`, '/') satisfies org %s end) end;", model.Object[0], conNVal)
		}
	case "token":
		switch model.FHIRFieldType {
		case "string":
			str += buildArrayQuery(model, conNVal)
		case "boolean":
			str += fmt.Sprintf("r.`%s` %s", model.Field[0].Field, conNVal)
		case "code":
			str += fmt.Sprintf("r.`%s` %s", model.Field[0].Field, conNVal)
		case "coding", "identifier" :
			str += buildArrayQuery(model, conNVal)
		}
	}

	return str
}

func buildArrayQuery(model models.QueryParam, conNVal string) (str string) {

	if len(model.Array) > 0 {
		str += fmt.Sprintf("ANY n IN %s satisfies ", model.Array[0])
	} else {
		str += "where "
	}
	switch model.FHIRFieldType {
	case "coding" :
		// select * from `default` as r where r.`resourceType` = 'Patient' and ANY n IN communication satisfies (any d in n.`language`.`coding` satisfies d.`display` = 'Dutch' and d.`system` = 'urn:ietf:bcp:47' end) end;
		if len(model.Value) > 1 {
			if model.Value[0] != "" {
				str += fmt.Sprintf("(any d in n.`%s`.`coding` satisfies d.`system` = '%s' and d.`code` = '%s' end)", model.Field[0].Field, model.Value[0], model.Value[1])
			}else{
				str += fmt.Sprintf("(any d in n.`%s`.`coding` satisfies  d.`code` = '%s' end)", model.Field[0].Field, model.Value[1])
			}
		}
		str += " end;"

	case "identifier" :
		// select * from `default` as r where r.`resourceType` = 'Patient' and ANY n IN communication satisfies (any d in n.`language`.`coding` satisfies d.`display` = 'Dutch' and d.`system` = 'urn:ietf:bcp:47' end) end;
		if len(model.Value) > 1 {
			if model.Value[0] != "" {
				str += fmt.Sprintf("n.`system` = '%s' and n.`value` = '%s' end)", model.Value[0], model.Value[1])
			}else{
				str += fmt.Sprintf("v.`value` = '%s' end)", model.Value[1])
			}
		}
		str += " end;"
	default:
		if len(model.Field) >= 1 {
			for i, field := range model.Field {
				if field.Array {
					str += fmt.Sprintf("(any %s in n.`%s` satisfies %s %s end) ", model.Field[i].Field, field.Field, model.Field[i].Field, conNVal)
				} else {
					str += fmt.Sprintf("n.`%s` %s", field.Field, conNVal)
				}
				// dont print and for the last field
				if i < len(model.Field)-1 {
					str += "and "
				}
			}
		}
		str += " end;"
	}
	return str
}
