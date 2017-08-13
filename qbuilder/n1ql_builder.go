package qbuilder

import (
	"github.com/kite-social/uriql/models"
	"fmt"
	"strings"
)

type N1QLQueryBuilder struct {
}

func GetN1QLQueryBuilder() *N1QLQueryBuilder {
	return &N1QLQueryBuilder{}
}

func (n *N1QLQueryBuilder) Build(allparam []models.QueryParam) string {

	var str string
	total := len(allparam)
	for i, model := range allparam {

		if i == 0 {
			// todo#fix fix condtion parameter
			// where ANY n IN %s satisfies (any name in n.`%s` SATISFIES name LIKE '%%%s%%' END) end;
			str += fmt.Sprintf("select * from `kite` as r where r.`type` = '%s' and ", model.Resource)
		}

		var conNVal string
		//#todo#fix token condition need to be fixed
		switch model.Condition {
		case "like":
			conNVal = fmt.Sprintf("%s '%%%s%%'", model.Condition, model.Value[0])
		case "=":
			conNVal = fmt.Sprintf("%s '%s'", model.Condition, model.Value[0])
		default:
			conNVal = fmt.Sprintf("%s '%s'", model.Condition, model.Value[0])
		}

		switch model.FHIRType {

		case "universal":
			str += "r."
			for _, k := range model.Field {
				if k.Array == false && k.Object == false {
					str += fmt.Sprintf("%s.", k.Field)
				} else {
					// not yet for _profile , _security and _tag
				}
			}
			str += fmt.Sprintf(" %s", conNVal)

		case "single":
			str = fmt.Sprintf("select * from `kite` as r where META(r).id = '%s::%s'", model.Resource, model.Value[0])
		case "number":
			// found just field so far
			str += fmt.Sprintf("r.%s %s", model.Field[0].Field, conNVal)
		case "string":
			// we now assume that all string type search field is within an array !!
			str += buildArrayQuery(model, conNVal, i, total)
		case "reference":
			for _, k := range model.Field {
				if k.Array == true {
					str += fmt.Sprintf("any ref in SPLIT(r.`%s`.`reference`,'/') satisfies ref %s end;", model.Field, conNVal)
				} else if k.Object == true {
					str += fmt.Sprintf("any ref in r.`%s` satisfies (any org in SPLIT(ref.`reference`, '/') satisfies org %s end) end;", model.Field, conNVal)
				}
			}
		case "token":
			switch model.FHIRFieldType {
			case "string":
				str += buildArrayQuery(model, conNVal, i, total)
			case "boolean":
				str += fmt.Sprintf("r.`%s` %s", model.Field[0].Field, conNVal)
			case "code":
				str += fmt.Sprintf("r.`%s` %s", model.Field[0].Field, conNVal)
			case "coding", "identifier":
				str += buildArrayQuery(model, conNVal, i, total)
			}
		}

	}

	return str
}

func buildArrayQuery(model models.QueryParam, conNVal string, loop, total int) (str string) {
	switch model.FHIRFieldType {
	case "coding":
		// select * from `default` as r where r.`resourceType` = 'Patient' and ANY n IN communication satisfies (any d in n.`language`.`coding` satisfies d.`display` = 'Dutch' and d.`system` = 'urn:ietf:bcp:47' end) end;
		if len(model.Value) > 1 {
			if model.Value[0] != "" {
				str += fmt.Sprintf("(any d in n.`%s`.`coding` satisfies d.`system` = '%s' and d.`code` = '%s' end)", model.Field[0].Field, model.Value[0], model.Value[1])
			} else {
				str += fmt.Sprintf("(any d in n.`%s`.`coding` satisfies  d.`code` = '%s' end)", model.Field[0].Field, model.Value[1])
			}
		}
		str += " end;"
	case "string":
		buildPath(&model)
		oldPath := filterPath(&model)
		model.Path = ""
		// all of them are object like address.city.name
		if model.ArrayCount == 0 {
			str += fmt.Sprintf("r.%s %s", oldPath, conNVal)
		} else {
			// condition might be address.[]city.name
			if model.ArrayCount == 1 {
				for i := 0; i < model.ArrayCount; i++ {
					// array does not exists so its a simple path like address.city.name
					buildPath(&model)
					str += fmt.Sprintf("(any n in r.%s satisfies n.%s %s end)", oldPath, trimDot(model.Path), conNVal)
				}
			} else {
				// multiple array found
				// condition might be address.[]city.name.[]room.whatever
				str += fmt.Sprintf("(any n in r.%s satisfies ", oldPath)
				for i := 0; i < model.ArrayCount; i++ {
					// array does not exists so its a simple path like address.city.name
					buildPath(&model)
					oldPath := filterPath(&model)
					model.Path = ""
					if i == 0 {
						str += fmt.Sprintf("(any d%d in n.%s satisfies ", i, oldPath)
					} else if i == model.ArrayCount-1 {
						// add the condition and value to the last nested subquery
						str += fmt.Sprintf("d%d.%s %s end)", i-1, oldPath, conNVal)
						break; // its done break the loop
					} else {
						str += fmt.Sprintf("(any d%d in d%d satisfies ", i, oldPath)
					}
				}
			}
		}

		if loop >= 0 && loop < total-1 {
			str += fmt.Sprintf(" or ")
		}
		// valid query of address.[]city.name.[]room.whatever will be
		/*
		select * from `default` as r where r.`type` = 'restaurant' and any n in r.`address`.`city` satisfies
		 (any d0 in n.`name`.`room` satisfies d0.`whatever` = 'dhaka' end) end;
		 */

	case "identifier":
		// select * from `default` as r where r.`resourceType` = 'Patient' and ANY n IN communication satisfies (any d in n.`language`.`coding` satisfies d.`display` = 'Dutch' and d.`system` = 'urn:ietf:bcp:47' end) end;
		if len(model.Value) > 1 {
			if model.Value[0] != "" {
				str += fmt.Sprintf("n.`system` = '%s' and n.`value` = '%s' end)", model.Value[0], model.Value[1])
			} else {
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
					str += "or "
				}
			}
		}
		str += " end;"
	}
	return str
}

func buildPath(model *models.QueryParam) {
	for i := 0; i < len(model.Field); i++ {
		model.Path += fmt.Sprintf("`%s`.", model.Field[i].Field)
		if model.Field[i].Array == true {
			model.Field = append(model.Field[:0], model.Field[1:]...)
			break;
		} else {
			model.Field = append(model.Field[:0], model.Field[1:]...)
			buildPath(model)
		}
		break;
	}
}

func filterPath(model *models.QueryParam) string {
	trim := strings.TrimSuffix(model.Path, ".")
	return fmt.Sprintf("%s", trim)
}

func trimDot(input string) string {
	trim := strings.TrimSuffix(input, ".")
	return trim
}
