package uriql

import (
	"github.com/kite-social/uriql/builder"
	"github.com/kite-social/uriql/helper"
	"github.com/kite-social/uriql/models"
	"strings"
)

// QueryDecoder : Get Query Decoder Object
type QueryDecoder struct {
	Def *helper.Def
}

// GetQueryDecoder : Get Query Decoder Object based on Dictionary
func GetQueryDecoder(dict *models.Dictionary) *QueryDecoder {
	return &QueryDecoder{
		Def: helper.GetDef(dict),
	}
}

/*
DecodeQueryString : Decodes Request information into Query Parameter
todo--add better exception handeling
*/
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

	// Universal Resource Search Parameter
	switch queryBase {
	case "_id":
		queryStruct.SearchResult.Type = "bundle"

		queryStruct.FHIRType = "universal"
		queryStruct.Condition = "="
		queryStruct.Value.Value = queryParam

		queryStruct.Resource = uri[0]
		decodedParam = append(decodedParam, queryStruct)
		return decodedParam

	case "_lastUpdated":
		queryStruct.SearchResult.Type = "bundle"

		queryStruct.FHIRType = "universal"
		getConditionNVal(&queryStruct, queryParam)

		queryStruct.Resource = uri[0]
		decodedParam = append(decodedParam, queryStruct)
		return decodedParam

	default:
		queryStruct.SearchResult.Type = "bundle"

		if strings.Contains(queryBase, ":") {
			// name:contains
			condition = strings.Split(queryBase, ":")           // name:contains=Mr
			queryStruct.Value.Modifiers = condition[1]          // `contains`
			info = f.Def.MatchSearchParam(uri[0], condition[0]) // Get information about the query from Dict `name`
		} else {
			info = f.Def.MatchSearchParam(uri[0], queryBase) // `name`
		}

		// #todo#fix better exception handeling

		if info.Type == "" && info.FieldType == "" {
			panic("Definition of [" + uri[1] + "] not found in dictionary")
		}

		queryStruct.FHIRType = info.Type
		queryStruct.FHIRFieldType = info.FieldType

		// Conditions for String type Parameter that could contain :contains, :exact etc Ex - ?name:contains=Mr.
		switch info.Type {

		case "string":
			// if condition specified like : contains or exact
			queryStruct.Value.Value = queryParam
			if len(condition) > 0 {
				if queryStruct.Value.Modifiers == "contains" {
					queryStruct.Condition = "like"
				} else if queryStruct.Value.Modifiers == "exact" {
					queryStruct.Condition = "=" // #todo-add search lower and upper both
				}
			} else {
				queryStruct.Condition = "="
			}

			// Conditions for token parameter Ex : - language=https://some.com|FR
		case "token":
			switch info.FieldType {
			// todo#fix handle different modifiers for token like :not , :text
			// If the parameter contains only single parameter or boolean value Ex :- ?active=true
			case "boolean":
				queryStruct.Condition = "="
				queryStruct.Value.Value = queryParam

			case "string":
				queryStruct.Condition = "="          // Assuming all string under token needs to be exact matched
				queryStruct.Value.Value = queryParam // Mr.
			case "coding", "identifier":
				queryStruct.Condition = "=" // Assuming all string under token needs to be exact matched
				if strings.Contains(queryParam, "|") {
					code := strings.Split(queryParam, "|")
					queryStruct.Value.Codable.System = code[0] // https://some.com
					queryStruct.Value.Codable.Code = code[1]   // FR
				}else{
					queryStruct.Value.Codable.Code = queryParam
				}

			case "code", "codeableConcept":

				queryStruct.Value.Value = queryParam

				if len(condition) > 0 {
					switch queryStruct.Value.Modifiers {
					case "not":
						queryStruct.Condition = "!="
					case "above":
						queryStruct.Condition = ">=" // #todo-add search lower and upper both
					case "below":
						queryStruct.Condition = "=<"
					case "in":
						queryStruct.Condition = "in"
					case "not-in":
						queryStruct.Condition = "not in"
					}
				} else {
					queryStruct.Condition = "="
				}

			}

			// If the parameter contains number value Ex : ?length=gt204
		case "number":
			getConditionNVal(&queryStruct, queryParam)

		case "reference":
			queryStruct.Condition = "="
			if strings.Contains(queryParam, "/") {
				val := strings.Split(queryParam, "/")
				queryStruct.Value.Reference.Target = val[0]
				queryStruct.Value.Reference.Value = val[1]
			} else {
				// if the format is Patient?general-practitioner:Practitioner=23
				queryStruct.Value.Reference.Target = queryStruct.Value.Modifiers // because then :Practitioner part goes to modifiers like but its not
				queryStruct.Value.Reference.Value = queryParam
			}

			// Additional Cases for Graph Type [relation, node , both]

		case "relation", "node":
			queryStruct.Value.Value = queryParam // Mr.
			if len(condition) > 0 {
				switch queryStruct.Value.Modifiers {

				case "contains":
					queryStruct.Condition = "like"
				case "exact":
					queryStruct.Condition = "=" // #todo-add search lower and upper both
				}
			} else {
				queryStruct.Condition = "="
			}
		}
	}
	queryStruct.Resource = uri[0]

	for _, path := range info.Path {

		fv := helper.GetFieldInfoFromPath(path)
		queryStruct.FieldsInfo.ArrayPath = fv.ArrayPath
		queryStruct.FieldsInfo.ObjectPath = fv.ObjectPath

		decodedParam = append(decodedParam, queryStruct)
	}

	return decodedParam
}

func getConditionNVal(queryStruct *models.QueryParam, queryParam string) {
	con := queryParam[0:2]
	if con == "gt" {
		queryStruct.Condition = ">"
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else if con == "lt" {
		queryStruct.Condition = "<"
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else if con == "ge" {
		queryStruct.Condition = ">="
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else if con == "le" {
		queryStruct.Condition = "=<"
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else if con == "ne" {
		queryStruct.Condition = "!="
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else {
		queryStruct.Condition = "="
		queryStruct.Value.Value = queryParam
	}
}

type QueryIndex struct {
	Resource string
	Indexes  []string
}

/*
DecodeQueryIndex : Builds Query Index Query out of Dictionary
todo--add better exception handeling
*/
func (f *QueryDecoder) DecodeQueryIndex() []QueryIndex {

	var index []QueryIndex
	for resource, dict := range f.Def.Dictionary.Model {
		var idx QueryIndex
		idx.Resource = resource
		var indexBuilder builder.QueryBuilder
		idx.Indexes = indexBuilder.BuildQueryIndex(f.Def.Dictionary.Bucket, resource, dict)
		index = append(index, idx)
	}
	return index
}
