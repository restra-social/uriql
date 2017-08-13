package uriql

import (
	"strings"
	"github.com/kite-social/uriql/models"
	"github.com/kite-social/uriql/helper"
	"reflect"
)

type QueryDecoder struct {
	Def *helper.Def
}

func GetQueryDecoder(dict map[string]map[string]models.SearchParam) *QueryDecoder {
	return &QueryDecoder{
		Def: helper.GetDef(dict),
	}
}

// Path could be []name.[]family , []address.state, active
func (f *QueryDecoder) getFieldInfoFromPath(str string) (fieldInfo []models.FieldInfo, count int) {

	var fv models.FieldInfo
	//count = 0 // initial array count

	if strings.Contains(str, ".") {
		fi := strings.Split(str, ".")
		end := len(fi)

		// loop through the end of the path except for the last
		for i := 0; i < end-1; i++ {
			if strings.HasPrefix(fi[i], "[]") {
				// lets say if []address.city
				fv.Array = true
				fv.Object = false
				fv.Field = fi[i][2:len(fi[i])] // address
				count++
			} else {
				// lets say if managingOrganization.reference
				fv.Field = fi[i] // managingOrganization
				fv.Array = false
				fv.Object = true
			}

			fieldInfo = append(fieldInfo, fv)
		}
		// the last one is the field so
		// if []name.[]family
		if strings.HasPrefix(fi[end-1], "[]") {
			fv.Array = true
			fv.Object = false
			fv.Field = fi[end-1][2:len(fi[end-1])]
			count++
		} else {
			fv.Array = false
			fv.Object = false
			fv.Field = fi[end-1]
		}

		fieldInfo = append(fieldInfo, fv)
	} else {
		// if active, gender
		fv.Array = false
		fv.Object = false
		fv.Field = str
		fieldInfo = append(fieldInfo, fv)
	}

	return fieldInfo, count
}

// todo--add better exception handeling
func (f *QueryDecoder) DecodeQueryString(request models.RequestInfo) []models.QueryParam {
	var decodedParam []models.QueryParam

	var queryStruct models.QueryParam

	// Assign the Request Info to Query Struct
	queryStruct.RequestInfo = request

	uri := strings.Split(request.Query, "?") // Trim ? from the Query Parameter

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

		queryStruct.Resource = uri[0]
		decodedParam = append(decodedParam, queryStruct)
		return decodedParam

	case "_lastUpdated":
		queryStruct.SearchResult.Type = "bundle"
		queryStruct.Field = []models.FieldInfo{
			{
				Field: "lastUpdated",
				Array: false,
			},
		}
		queryStruct.FHIRType = "universal"
		getConditionNVal(&queryStruct, queryParam)

		queryStruct.Resource = uri[0]
		decodedParam = append(decodedParam, queryStruct)
		return decodedParam

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

			getConditionNVal(&queryStruct, queryParam)

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

			// Additional Cases for Graph Type [relation, node , both]

		case "relation", "node":
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
		}
	}
	queryStruct.Resource = uri[0]

	for _, path := range info.Path {

		fv, count := f.getFieldInfoFromPath(path)
		queryStruct.ArrayCount = count
		queryStruct.Field = []models.FieldInfo{} // reset fields
		queryStruct.Field = append(queryStruct.Field, fv...)

		decodedParam = append(decodedParam, queryStruct)
	}

	return decodedParam
}

func getConditionNVal(queryStruct *models.QueryParam, queryParam string) {
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
