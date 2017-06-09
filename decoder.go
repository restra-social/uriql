package uriql

import (
	"strings"
	"udhvabon.com/neuron/uriql/models"
	"udhvabon.com/neuron/uriql/dictionary"
	"reflect"
)

type QueryDecoder struct {
	queryStruct models.QueryParam
	Def *dictionary.Def
}

func GetQueryDecoder(dict map[string]map[string]models.SearchParam) *QueryDecoder {
	return &QueryDecoder{
		Def: dictionary.GetDef(dict),
	}
}

// Path could be []name.[]family , []address.state, active
func (f *QueryDecoder) getFieldInfoFromPath(str string) []models.FieldInfo {
	var fieldInfo []models.FieldInfo
	var fv models.FieldInfo
	f.queryStruct.ArrayCount = 0 // initial array count
	fv.Order = 0;
	if strings.Contains(str, ".") {
		fi := strings.Split(str, ".")
		end := len(fi)

		// loop through the end of the path except for the last
		for i :=0; i < end-1; i++ {
			if strings.HasPrefix(fi[i], "[]") {
				// lets say if []address.city
				fv.Array = true
				fv.Object = false
				fv.Field = fi[i][2:len(fi[i])] 	// address
				f.queryStruct.ArrayCount++
			}else{
				// lets say if managingOrganization.reference
				fv.Field = fi[i]	// managingOrganization
				fv.Array = false
				fv.Object = true
			}
			fv.Order++
			fieldInfo = append(fieldInfo, fv)
		}
		// the last one is the field so
		// if []name.[]family
		if strings.HasPrefix(fi[end-1], "[]") {
			fv.Array = true
			fv.Object = false
			fv.Field = fi[end-1][2:len(fi[end-1])]
			f.queryStruct.ArrayCount++
		}else {
			fv.Array = false
			fv.Object = false
			fv.Field = fi[end-1]
		}
		fv.Order++
		fieldInfo = append(fieldInfo, fv)
	} else {
		// if active, gender
		fv.Array = false
		fv.Object = false
		fv.Field = str
		fieldInfo = append(fieldInfo, fv)
	}

	return fieldInfo
}

// todo--add better exception handeling
func (f *QueryDecoder) DecodeQueryString(query string) models.QueryParam {

	uri := strings.Split(query, "?")         // Trim ? from the Query Parameter
	if len(uri) == 1 {
		// if the parameter is /Patient/1234789
		f.queryStruct .SearchResult.Type = "resource"
		v := strings.Split(uri[0], "/")
		f.queryStruct.Condition = "="
		f.queryStruct.FHIRType = "single"
		f.queryStruct.Resource = v[0]
		f.queryStruct.Value = []string{v[1]}
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
			f.queryStruct.SearchResult.Type = "bundle"
			f.queryStruct.Field = []models.FieldInfo{
				{
					Field: "id",
					Array: false,
				},
			}
			f.queryStruct.FHIRType = "universal"
			f.queryStruct.Condition = "="
			f.queryStruct.Value = []string{queryParam}

		case "_lastUpdated":
			f.queryStruct.SearchResult.Type = "bundle"
			f.queryStruct.Field = []models.FieldInfo{
				{
					Field: "lastUpdated",
					Array: false,
				},
			}
			f.queryStruct.FHIRType = "universal"
			con := queryParam[0:2]
			if con == "gt" {
				f.queryStruct.Condition = ">"
				f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
			} else if con == "lt" {
				f.queryStruct.Condition = "<"
				f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
			} else if con == "ge" {
				f.queryStruct.Condition = ">="
				f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
			} else if con == "le" {
				f.queryStruct.Condition = "=<"
				f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
			} else if con == "ne" {
				f.queryStruct.Condition = "!="
				f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
			} else {
				f.queryStruct.Condition = "="
				f.queryStruct.Value = append(f.queryStruct.Value, queryParam)
			}

		default:
			f.queryStruct.SearchResult.Type = "bundle"
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

			for _, path = range info.Path {

				fv := f.getFieldInfoFromPath(path)
				f.queryStruct.Field = append(f.queryStruct.Field, fv...)

			}
			// Conditions for String type Parameter that could contain :contains, :exact etc Ex - ?name:contains=Mr.
			switch info.Type {

			case "string":
				// if conditon specified like : contains or exact
				f.queryStruct.Value = append(f.queryStruct.Value, queryParam) // Mr.
				f.queryStruct.FHIRType = info.Type
				f.queryStruct.FHIRFieldType = info.FieldType
				if len(condition) > 0 {
					if modifier == "contains" {
						f.queryStruct.Condition = "like"
					} else if modifier == "exact" {
						f.queryStruct.Condition = "=" // #todo-add search lower and upper both
					}
				} else {
					f.queryStruct.Condition = "="
				}

				// Conditions for token parameter Ex : - language=https://some.com|FR
			case "token":
				switch info.FieldType {
				// todo#fix handle different modifiers for token like :not , :text
				// If the parameter contains only single paramter or boolean value Ex :- ?active=true
				case "boolean":
					f.queryStruct.Condition = "="
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam)
					f.queryStruct.FHIRType = info.Type
					f.queryStruct.FHIRFieldType = info.FieldType
				case "string", "code":
					f.queryStruct.Condition = "="                               // Assuming all string under token needs to be exact matched
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam) // Mr.
					f.queryStruct.FHIRType = info.Type
					f.queryStruct.FHIRFieldType = info.FieldType

				case "coding", "identifier", "codeableConcept":
					f.queryStruct.FHIRType = info.Type
					f.queryStruct.FHIRFieldType = info.FieldType
					f.queryStruct.Value = strings.Split(queryParam, "|")

					if len(condition) > 0 {
						if modifier == "not" {
							f.queryStruct.Condition = "!="
						} else if modifier == "above" {
							f.queryStruct.Condition = ">=" // #todo-add search lower and upper both
						} else if modifier == "below" {
							f.queryStruct.Condition = "=<"
						} else if modifier == "in" {
							f.queryStruct.Condition = "in"
						} else if modifier == "not-in" {
							f.queryStruct.Condition = "not in"
						}
					} else {
						f.queryStruct.Condition = "="
					}

				}

				// If the parameter contains number value Ex : ?length=gt204
			case "number":
				f.queryStruct.FHIRType = info.Type
				f.queryStruct.FHIRFieldType = info.FieldType
				con := queryParam[0:2]
				if con == "gt" {
					f.queryStruct.Condition = ">"
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
				} else if con == "lt" {
					f.queryStruct.Condition = "<"
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
				} else if con == "ge" {
					f.queryStruct.Condition = ">="
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
				} else if con == "le" {
					f.queryStruct.Condition = "=<"
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
				} else if con == "ne" {
					f.queryStruct.Condition = "!="
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam[2:len(queryParam)])
				} else {
					f.queryStruct.Condition = "="
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam)
				}

			case "reference":
				f.queryStruct.FHIRType = info.Type
				f.queryStruct.FHIRFieldType = info.FieldType
				f.queryStruct.Condition = "="
				if strings.Contains(queryParam, "/") {
					val := strings.Split(queryParam, "/")
					f.queryStruct.Value = []string{val[1]}
				} else {
					// if the format is Patient?general-practitioner:Practitioner=23
					f.queryStruct.Value = append(f.queryStruct.Value, queryParam)
				}
			}
		}
		f.queryStruct.Resource = uri[0]
	}

	return f.queryStruct
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
