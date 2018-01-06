package uriql

import (
	"fmt"
	"github.com/restra-social/uriql/builder"
	"github.com/restra-social/uriql/helper"
	"github.com/restra-social/uriql/models"
	"net/url"
	"strconv"
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

func (f *QueryDecoder) DecodeQueryString(request models.RequestInfo) [][]models.QueryParam {
	var queryParams [][]models.QueryParam

	uri := strings.Split(request.Query, "?") // Trim ? from the Query Parameter
	resource := uri[0]
	query := uri[1]

	var queryParam []models.QueryParam

	if strings.Contains(query, "&") {

		// decode multiple query
		decodedQuery, err := url.ParseQuery(query)
		if err != nil {
			fmt.Println("Invalid Query String Format ", err.Error())
			return nil
		}

		// check for pagination
		var mainQuery string
		if strings.Contains(query, "_") {
			// split and saperate main query
			mainQuery = strings.Split(query, "&_")[0] //name=fahim&address=dhaka&_size=20

			var limit int
			var page int

			if strings.Contains(query, "_size") && strings.Contains(query, "_page") {
				limit, err = strconv.Atoi(decodedQuery.Get("_size"))
				page, err = strconv.Atoi(decodedQuery.Get("_page"))
				request.Limit = limit
				request.Page = page
			} else if strings.Contains(query, "_page") {
				page, err = strconv.Atoi(decodedQuery.Get("_page"))
				request.Page = page
				request.Limit = 10 // If no limit set then use default as 10
			} else if strings.Contains(query, "_size") {
				limit, err = strconv.Atoi(decodedQuery.Get("_size"))
				request.Limit = limit
			}
		}

		multiParam := strings.Split(mainQuery, "&")
		for _, param := range multiParam {
			queryParam = f.DecodeQuery(resource, param, request)
			queryParams = append(queryParams, queryParam)
		}
	} else {
		queryParam = f.DecodeQuery(resource, query, request)
		queryParams = append(queryParams, queryParam)
	}

	return queryParams
}

/*
DecodeQueryString : Decodes Request information into Query Parameter
todo--add better exception handeling
*/
func (f *QueryDecoder) DecodeQuery(resource, queryRequest string, request models.RequestInfo) []models.QueryParam {
	var decodedParam []models.QueryParam

	var queryStruct models.QueryParam

	// Assign the Request Info to Query Struct
	// #todo add request info to parameter for further debugging purpose
	queryStruct.RequestInfo = request

	valueGet := strings.Split(queryRequest, "=") // Split where it gets = sign
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

		queryStruct.Resource = resource
		decodedParam = append(decodedParam, queryStruct)
		return decodedParam

	case "_lastUpdated":
		queryStruct.SearchResult.Type = "bundle"

		queryStruct.FHIRType = "universal"
		getConditionNVal(&queryStruct, queryParam)

		queryStruct.Resource = resource
		decodedParam = append(decodedParam, queryStruct)
		return decodedParam

	default:
		queryStruct.SearchResult.Type = "bundle"

		if strings.Contains(queryBase, ":") {
			// name:contains
			condition = strings.Split(queryBase, ":")             // name:contains=Mr
			queryStruct.Value.Modifiers = condition[1]            // `contains`
			info = f.Def.MatchSearchParam(resource, condition[0]) // Get information about the query from Dict `name`
		} else {
			info = f.Def.MatchSearchParam(resource, queryBase) // `name`
		}

		// #todo#fix better exception handeling

		if info.Type == "" && info.FieldType == "" {
			panic("Definition of [" + queryRequest + "] not found in dictionary")
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
				} else {
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
				queryStruct.Value.Value = queryParam
			} else {
				// if the format is Patient?general-practitioner:Practitioner=23
				// because then :Practitioner part goes to modifiers like but its not
				queryStruct.Value.Value = fmt.Sprintf("%s/%s", queryStruct.Value.Modifiers, queryParam)
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
	queryStruct.Resource = resource

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

/*
DecodeQueryIndex : Builds Query Index Query out of Dictionary
todo--add better exception handeling
*/
func (f *QueryDecoder) DecodeQueryIndex(resourceType string) models.Migrations {

	var queryIndexes models.Migrations

	indexes := make(map[string]models.IndexInfo)

	for resource, dict := range f.Def.Dictionary.Model {

		indexes[resource] = builder.BuildQueryIndex(f.Def.Dictionary.Bucket, resource, dict, resourceType)
	}

	queryIndexes.Indexes.Migration = indexes
	return queryIndexes
}
