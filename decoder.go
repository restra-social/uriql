package uriql

import (
	"errors"
	"fmt"
	"github.com/restra-social/uriql/builder"
	"github.com/restra-social/uriql/helper"
	"github.com/restra-social/uriql/models"
	"net/url"
	"strconv"
	"strings"
)

const (
	QueryParameterIdentifier = "?"
	QueryParameterSeparator  = "&"
	QueryValueSaperator      = "="

	// Field Modifier
	ModifierIdentifier             = ":"
	ModifierFieldForWildcardSearch = "contains"
	ModifierFieldForExactMatch     = "exact"
	ModifierForNotEqual = "not"
	ModifierForGreaterThanEqual = "above"
	ModifierForLessThanEqual = "below"
	ModifierForWithin = "in"
	ModifierForNotWithin = "not-in"


	// Value Modifier
	ValueModifierForGraterThan = "gt"
	ValueModifierForGraterThanAndEqual = "ge"
	ValueModifierForLessThan = "lt"
	ValueModifierForLessThanAndEqual = "le"
	ValueModifierForNotEqual = "nt"


	PaginationIdentifier = "&_"

	DefaultPaginationParameter     = "%s&_page=1&_size=10"
	DefaultPaginationSizeParameter = "_size"
	DefaultPaginationPageParameter = "_page"

	DefaultSearchCondition         = "="
	DefaultWildcardSearchCondition = "like"

	UniversalIdentifierField  = "_id"
	UniversalLastUpdatedField = "_lastUpdated"

	DefaultSearchType = "bundle"

	// Base Search Type
	BaseStringTypeIdentifier    = "string"
	BaseTokenTypeIdentifier     = "token"
	BaseNumberTypeIdentifier    = "number"
	BaseReferenceTypeIdentifier = "reference"
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

// DecodeQueryString takes the query parameter , resource type and a filter object
// Returns QueryInfo Object
func (f *QueryDecoder) DecodeQueryString(request models.RequestInfo, filter *models.Filter) (*models.QueryInfo, error) {
	var queryParams [][]models.QueryParamInfo

	uri := strings.Split(request.Query, QueryParameterIdentifier) // Trim ? from the Query Parameter
	resource := request.Type
	query := uri[1]

	// Storage for Query Parameter , means how many path to search for query
	var queryParam []models.QueryParamInfo

	if strings.Contains(query, QueryParameterSeparator) { // & this means multiple search parameter found

		// add default pagination parameter if not added
		if !strings.Contains(query, PaginationIdentifier) {
			query = fmt.Sprintf(DefaultPaginationParameter, query) // #todo remove hard coded Param
		}

		// decode multiple query
		decodedQuery, err := url.ParseQuery(query)
		if err != nil {
			msg := fmt.Sprintf("Invalid Query String Format : %s", err.Error())
			return nil, errors.New(msg)
		}

		// check for pagination
		var mainQuery string
		// split and saperate main query
		mainQuery = strings.Split(query, PaginationIdentifier)[0] //name=fahim&address=dhaka&_size=20

		var limit int
		var page int

		if strings.Contains(query, DefaultPaginationSizeParameter) && strings.Contains(query, "_page") {
			limit, err = strconv.Atoi(decodedQuery.Get(DefaultPaginationSizeParameter))
			page, err = strconv.Atoi(decodedQuery.Get(DefaultPaginationPageParameter))
			filter.Limit = limit
			filter.Page = page
		} else if strings.Contains(query, DefaultPaginationPageParameter) {
			page, err = strconv.Atoi(decodedQuery.Get(DefaultPaginationPageParameter))
			filter.Page = page
			filter.Limit = 10 // If no limit set then use default as 10
		} else if strings.Contains(query, DefaultPaginationSizeParameter) {
			limit, err = strconv.Atoi(decodedQuery.Get(DefaultPaginationSizeParameter))
			filter.Limit = limit
		}

		multiParam := strings.Split(mainQuery, QueryParameterSeparator)
		for _, param := range multiParam {
			queryParam, err = f.DecodeQuery(resource, param, request)
			if err != nil {
				msg := fmt.Sprintf("Failed to Decode %s", err.Error())
				return nil, errors.New(msg)
			}
			queryParams = append(queryParams, queryParam)
		}

	} else {
		queryParam, err := f.DecodeQuery(resource, query, request)
		if err != nil {
			msg := fmt.Sprintf("Failed to Decode %s", err.Error())
			return nil, errors.New(msg)
		}
		queryParams = append(queryParams, queryParam)
	}

	var queryInfo models.QueryInfo
	queryInfo.Params = queryParams
	queryInfo.Filter = filter

	return &queryInfo, nil
}

/*
DecodeQueryString : Decodes Request information into Query Parameter
todo--add better exception handeling
*/
func (f *QueryDecoder) DecodeQuery(resource, queryRequest string, request models.RequestInfo) ([]models.QueryParamInfo, error) {

	// For error you know that
	var err error

	// Stores information about the search parameter
	var decodedParam []models.QueryParamInfo

	var queryStruct models.QueryParamInfo

	// Assign the Request Info to Query Structure
	queryStruct.RequestInfo = request

	valueGet := strings.Split(queryRequest, QueryValueSaperator) // Split where it gets = sign
	queryBase := valueGet[0]
	queryParam := valueGet[1]

	// info is for Storing Search Param Information from Dictionary
	var info *models.SearchParam

	// Condition is a slice because it could contain Modifiers
	var condition []string

	// Universal Resource Search Parameter, Some param is reserved and has specific task like _id , _lastUpdated
	// These parameter does not have modifier like _id:contains or _id:exact #todo
	switch queryBase {
	case UniversalIdentifierField:
		queryStruct.SearchResult.Type = DefaultSearchType

		queryStruct.FHIRType = "universal"
		queryStruct.Condition = DefaultSearchCondition
		queryStruct.Value.Value = queryParam

		queryStruct.Resource = resource
		decodedParam = append(decodedParam, queryStruct)
		return decodedParam, nil

	case UniversalLastUpdatedField:
		queryStruct.SearchResult.Type = DefaultSearchType

		queryStruct.FHIRType = "universal"
		getConditionNVal(&queryStruct, queryParam)

		queryStruct.Resource = resource
		decodedParam = append(decodedParam, queryStruct)
		return decodedParam, nil

	default:
		// Default Handler which Means the Query Param is random or what specified in the dictionary
		queryStruct.SearchResult.Type = DefaultSearchType

		// Check if the Parameter has Modifiers
		if strings.Contains(queryBase, ModifierIdentifier) {
			// name:contains
			condition = strings.Split(queryBase, ModifierIdentifier) // name:contains=Mr
			queryStruct.Value.Modifiers = condition[1]               // `contains`
			info , err = f.Def.MatchSearchParam(resource, condition[0])    // Get information about the query from Dict `name`
			if err != nil {
				return nil, err
			}
		} else {
			info , err = f.Def.MatchSearchParam(resource, queryBase) // `name`
			if err != nil {
				return nil, err
			}
		}

		queryStruct.FHIRType = info.Type
		queryStruct.FHIRFieldType = info.FieldType

		// Pass down the which fields to select for result
		queryStruct.SelectStatement = info.Select

		// Conditions for String type Parameter that could contain :contains, :exact etc Ex - ?name:contains=Mr.
		switch info.Type {

		case BaseStringTypeIdentifier:
			// if condition specified like : contains or exact
			queryStruct.Value.Value = queryParam
			if len(condition) > 0 {
				if queryStruct.Value.Modifiers == "contains" {
					queryStruct.Condition = DefaultWildcardSearchCondition
				} else if queryStruct.Value.Modifiers == "exact" {
					queryStruct.Condition = DefaultSearchCondition // #todo add search lower and upper both
				}
			} else {
				queryStruct.Condition = DefaultSearchCondition
			}

			// Conditions for token parameter Ex : - language=https://some.com|FR
		case BaseTokenTypeIdentifier:
			switch info.FieldType {
			// todo#fix handle different modifiers for token like :not , :text
			// If the parameter contains only single parameter or boolean value Ex :- ?active=true
			case "boolean":
				queryStruct.Condition = DefaultSearchCondition
				queryStruct.Value.Value = queryParam

			case "string":
				queryStruct.Condition = DefaultSearchCondition // Assuming all string under token needs to be exact matched
				queryStruct.Value.Value = queryParam           // Mr.
			case "coding", "identifier":
				queryStruct.Condition = DefaultSearchCondition // Assuming all string under token needs to be exact matched
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
					case ModifierForNotEqual:
						queryStruct.Condition = "!="
					case ModifierForGreaterThanEqual:
						queryStruct.Condition = ">=" // #todo-add search lower and upper both
					case ModifierForLessThanEqual:
						queryStruct.Condition = "=<"
					case ModifierForWithin:
						queryStruct.Condition = "in"
					case ModifierForNotWithin:
						queryStruct.Condition = "not in"
					}
				} else {
					queryStruct.Condition = DefaultSearchCondition
				}

			}

			// If the parameter contains number value Ex : ?length=gt204
		case BaseNumberTypeIdentifier:
			getConditionNVal(&queryStruct, queryParam)

		case BaseReferenceTypeIdentifier:
			queryStruct.Condition = DefaultSearchCondition

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

				case ModifierFieldForWildcardSearch:
					queryStruct.Condition = DefaultWildcardSearchCondition
				case ModifierFieldForExactMatch:
					queryStruct.Condition = DefaultSearchCondition // #todo-add search lower and upper both
				}
			} else {
				queryStruct.Condition = DefaultSearchCondition
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

	return decodedParam, nil
}

func getConditionNVal(queryStruct *models.QueryParamInfo, queryParam string) {
	con := queryParam[0:2]
	if con == ValueModifierForGraterThan {
		queryStruct.Condition = ">"
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else if con == ValueModifierForLessThan {
		queryStruct.Condition = "<"
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else if con == ValueModifierForGraterThanAndEqual {
		queryStruct.Condition = ">="
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else if con == ValueModifierForLessThanAndEqual {
		queryStruct.Condition = "=<"
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else if con == ValueModifierForNotEqual {
		queryStruct.Condition = "!="
		queryStruct.Value.Value = queryParam[2:len(queryParam)]
	} else {
		queryStruct.Condition = DefaultSearchCondition
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
