package uriql

import (
	"strings"
	"udhvabon.com/neuron/uriql/models"
	"udhvabon.com/neuron/uriql/dictionary"
	"reflect"
)

type QueryDecoder struct {
	Def *dictionary.Def
}

func GetQueryDecoder(dict map[string]map[string]models.SearchParam) *QueryDecoder {
	return &QueryDecoder{
		Def: dictionary.GetDef(dict),
	}
}

// Path could be []name.[]family , []address.state, active
func (f *QueryDecoder) getArrayAndFieldFromPath(str string) models.SearchField {
	var fv models.SearchField
	if strings.Contains(str, ".") {
		fi := strings.Split(str, ".")
		// if []name.[]family
		if strings.HasPrefix(fi[0], "[]") && strings.HasPrefix(fi[1], "[]") {
			fv.Array = fi[0][2:len(fi[0])]
			fv.Field.Array = true
			fv.Field.Field = fi[1][2:len(fi[1])]
		} else if strings.HasPrefix(fi[0], "[]"){
			// if []address.city
			fv.Array = fi[0][2:len(fi[0])]
			fv.Field.Array = false
			fv.Field.Field = fi[1]
		}else {
			// if managingOrganization.reference
			fv.Object = fi[0]					// assuming object is ony 2 level deep !
			fv.Field.Array = false
			fv.Field.Field = fi[1]
		}
	} else {
		// if active, gender
		fv.Field.Array = false
		fv.Field.Field = str
	}
	return fv
}

// todo--add better exception handeling
func (f *QueryDecoder) DecodeQueryString(query string) models.QueryParam {
	var queryStruct models.QueryParam

	var fv models.SearchField

	uri := strings.Split(query, "?")         // Trim ? from the Query Parameter
	if len(uri) == 1 {
		// if the parameter is /Patient/1234789
		queryStruct.SearchResult.Type = "resource"
		v := strings.Split(uri[0], "/")
		queryStruct.Condition = "="
		queryStruct.FHIRType = "single"
		queryStruct.Resource = v[0]
		queryStruct.Value = []string{v[1]}
	}else {
		valueGet := strings.Split(uri[1], "=") // Split where it gets = sign
		queryBase := valueGet[0]
		queryParam := valueGet[1]

		var info *models.SearchParam
		var condition []string
		var modifier string

		// Universal Resource Search Parameter
		switch queryBase {
		case "_id":
			queryStruct.SearchResult.Type = "bundle"
			queryStruct.Field = []models.FieldInfo{
				{
					Field: "id",
					Array: false,
				},
			}
			queryStruct.FHIRType = "universal"
			queryStruct.Condition = "="
			queryStruct.Value = []string{queryParam}

		case "_lastUpdated":
			queryStruct.SearchResult.Type = "bundle"
			queryStruct.Object = []string{"meta"}
			queryStruct.Field = []models.FieldInfo{
				{
					Field: "lastUpdated",
					Array: false,
				},
			}
			queryStruct.FHIRType = "universal"
			con := queryParam[0:2]
			if con == "gt" {
				queryStruct.Condition = ">"
				queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
			} else if con == "lt" {
				queryStruct.Condition = "<"
				queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
			} else if con == "ge" {
				queryStruct.Condition = ">="
				queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
			} else if con == "le" {
				queryStruct.Condition = "=<"
				queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
			} else if con == "ne" {
				queryStruct.Condition = "!="
				queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
			} else {
				queryStruct.Condition = "="
				queryStruct.Value = append(queryStruct.Value, queryParam)
			}

		default:
			queryStruct.SearchResult.Type = "bundle"
			if strings.Contains(queryBase, ":") {
				// name:contains
				condition = strings.Split(queryBase, ":") // name , containers
				modifier = condition[1]
				info = f.Def.MatchSearchParam(uri[0], condition[0]) // Get information about the query from Dict
			} else {
				info = f.Def.MatchSearchParam(uri[0], queryBase)
			}

			// #todo#fix better exception handeling

			if info.Type == "" && info.FieldType == "" {
				panic("Definaiton of [" + uri[1] + "] not found in dictionary")
			}

			var path string
			var fieldInfo models.FieldInfo

			for _, path = range info.Path {

				fv = f.getArrayAndFieldFromPath(path)
				if fv.Array != "" {
					if !in_slice(fv.Array, queryStruct.Array) {
						queryStruct.Array = append(queryStruct.Array, fv.Array)
					}
				}
				if fv.Object != "" {
					queryStruct.Object = append(queryStruct.Object, fv.Object)
				}
				if fv.Field.Array {
					fieldInfo.Array = true
					fieldInfo.Field = fv.Field.Field
					queryStruct.Field = append(queryStruct.Field, fieldInfo)
				} else {
					fieldInfo.Array = false
					fieldInfo.Field = fv.Field.Field
					queryStruct.Field = append(queryStruct.Field, fieldInfo)
				}
			}
			// Conditions for String type Parameter that could contain :contains, :exact etc Ex - ?name:contains=Mr.
			switch info.Type {

			case "string":
				// if conditon specified like : contains or exact
				queryStruct.Value = append(queryStruct.Value, queryParam) // Mr.
				queryStruct.FHIRType = info.Type
				queryStruct.FHIRFieldType = info.FieldType
				if len(condition) > 0 {
					if modifier == "contains" {
						queryStruct.Condition = "like"
					} else if modifier == "exact" {
						queryStruct.Condition = "=" // #todo-add search lower and upper both
					}
				} else {
					queryStruct.Condition = "="
				}

				// Conditions for token parameter Ex : - language=https://some.com|FR
			case "token":
				switch info.FieldType {
				// todo#fix handle different modifiers for token like :not , :text
				// If the parameter contains only single paramter or boolean value Ex :- ?active=true
				case "boolean":
					queryStruct.Condition = "="
					queryStruct.Value = append(queryStruct.Value, queryParam)
					queryStruct.FHIRType = info.Type
					queryStruct.FHIRFieldType = info.FieldType
				case "string", "code":
					queryStruct.Condition = "="                               // Assuming all string under token needs to be exact matched
					queryStruct.Value = append(queryStruct.Value, queryParam) // Mr.
					queryStruct.FHIRType = info.Type
					queryStruct.FHIRFieldType = info.FieldType

				case "coding", "identifier", "codeableConcept":
					queryStruct.FHIRType = info.Type
					queryStruct.FHIRFieldType = info.FieldType
					queryStruct.Value = strings.Split(queryParam, "|")

					if len(condition) > 0 {
						if modifier == "not" {
							queryStruct.Condition = "!="
						} else if modifier == "above" {
							queryStruct.Condition = ">=" // #todo-add search lower and upper both
						} else if modifier == "below" {
							queryStruct.Condition = "=<"
						} else if modifier == "in" {
							queryStruct.Condition = "in"
						} else if modifier == "not-in" {
							queryStruct.Condition = "not in"
						}
					} else {
						queryStruct.Condition = "="
					}

				}

				// If the parameter contains number value Ex : ?length=gt204
			case "number":
				queryStruct.FHIRType = info.Type
				queryStruct.FHIRFieldType = info.FieldType
				con := queryParam[0:2]
				if con == "gt" {
					queryStruct.Condition = ">"
					queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
				} else if con == "lt" {
					queryStruct.Condition = "<"
					queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
				} else if con == "ge" {
					queryStruct.Condition = ">="
					queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
				} else if con == "le" {
					queryStruct.Condition = "=<"
					queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
				} else if con == "ne" {
					queryStruct.Condition = "!="
					queryStruct.Value = append(queryStruct.Value, queryParam[2:len(queryParam)])
				} else {
					queryStruct.Condition = "="
					queryStruct.Value = append(queryStruct.Value, queryParam)
				}

			case "reference":
				queryStruct.FHIRType = info.Type
				queryStruct.FHIRFieldType = info.FieldType
				queryStruct.Condition = "="
				if strings.Contains(queryParam, "/") {
					val := strings.Split(queryParam, "/")
					queryStruct.Value = []string{val[1]}
				} else {
					// if the format is Patient?general-practitioner:Practitioner=23
					queryStruct.Value = append(queryStruct.Value, queryParam)
				}
			}
		}
		queryStruct.Resource = uri[0]
	}

	return queryStruct
}

func in_slice(v interface{}, in interface{}) (ok bool) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}
