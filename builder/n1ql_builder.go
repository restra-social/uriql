package builder

import (
	"fmt"
	"github.com/restra-social/uriql/models"
	"strings"
)

const (
	DefaultSelectStatement = "r*"
	SelectResourceAs       = "r"
)

type N1QLQueryBuilder struct {
	bucketName             string
	resourceIdentifierName string
	page                   int
	limit                  int
}

// GetN1QLQueryBuilder : Get N1QL Builder Object
func GetN1QLQueryBuilder(bucket string, resourceIdentifier string) *N1QLQueryBuilder {
	return &N1QLQueryBuilder{bucketName: bucket, resourceIdentifierName: resourceIdentifier}
}

func (builder *N1QLQueryBuilder) Build(queryInfo *models.QueryInfo) string {

	if queryInfo.Filter.Limit == 0 {
		builder.limit = 10
	} else {
		builder.limit = queryInfo.Filter.Limit
	}
	if queryInfo.Filter.Page == 0 {
		builder.page = 1
	} else {
		builder.page = queryInfo.Filter.Page
	}

	var queryString []string

	// Check if Select choise is provided
	selectQuery := fmt.Sprintf("SELECT ")
	queryString = append(queryString, selectQuery)

	len := len(queryInfo.Params)

	for i, queryParam := range queryInfo.Params {
		// wrap the whole query with brackets because of
		// https://developer.couchbase.com/documentation/server/current/n1ql/n1ql-language-reference/logicalops.html
		// but dont apppend to last query

		if i == 0 { // Build only Once
			queryString = append(queryString, builder.BuildSelectQueryString(queryParam[i]))
			queryString = append(queryString, "(")
		}
		queryString = append(queryString, builder.BuildQueryString(queryParam))
		// join each composite parameter as AND operation
		if i != len-1 {
			queryString = append(queryString, ") AND (")
		}

		// close the bracket only at the end
		if i == len-1 {
			queryString = append(queryString, ")")
		}
	}

	if builder.page > 1 {
		queryString = append(queryString, fmt.Sprintf(" LIMIT %d OFFSET %d", builder.limit, (builder.page-1)*builder.limit))
	} else {
		queryString = append(queryString, fmt.Sprintf(" LIMIT %d", builder.limit))
	}

	result := strings.Join(queryString, "")

	return result
}

// #todo nested array is not supported
func (builder *N1QLQueryBuilder) BuildSelectQueryString(queryInfo models.QueryParamInfo) string {

	var query []string
	for _, field := range queryInfo.DictionaryInfo.SelectStatement {
		if strings.Contains(field, "[]") { // that is an array so make a select query

			lastField := strings.Split(field, ".")

			query = append(query, fmt.Sprintf("ARRAY v FOR v IN %s WHEN v.%s %s '%s' END", strings.TrimPrefix(lastField[0], "[]"), lastField[1], queryInfo.Condition, queryInfo.Value.Value))
		} else {
			query = append(query, field)
		}
	}

	selectQuery := strings.Join(query, ", ")

	return fmt.Sprintf("%s FROM `%s` AS %s WHERE ", selectQuery, builder.bucketName, SelectResourceAs) // #todo fix resource)
}

func (builder *N1QLQueryBuilder) BuildQueryString(queryParam []models.QueryParamInfo) string {

	var queryString []string

	for n, param := range queryParam {

		//#todo assigning limit and offset again and again , not good , saperate filter from request info
		//builder.limit = param.Filter.Limit
		//builder.page = param.Filter.Page

		field := param.DictionaryInfo.FieldsInfo
		arryLen := len(field.ArrayPath)

		paramLength := len(queryParam)

		objectPath := strings.Replace(field.ObjectPath, ".", "`.`", -1)

		// Combine Condition and Value
		var conNVal string
		//#todo#fix token condition need to be fixed
		// #todo , Filter string then pass to the sql , very dangerous to N1QL attack
		switch param.Condition {
		case "like":
			conNVal = fmt.Sprintf("%s '%%%s%%'", param.Condition, param.Value.Value)
		case "=":
			conNVal = fmt.Sprintf("%s '%s'", param.Condition, param.Value.Value)
		default:
			conNVal = fmt.Sprintf("%s '%s'", param.Condition, param.Value.Value)
		}

		var tempQuery string
		// this means there are multiple fields to search
		if paramLength > 1 && param.FHIRType == "string" {
			if n == 0 {
				tempQuery = fmt.Sprintf("(")
			}
		} else {
			tempQuery = fmt.Sprintf("r.`%s` = '%s' and ", builder.resourceIdentifierName, param.Resource)
		}
		queryString = append(queryString, tempQuery)

		// Switch Statement is For Special Case
		switch param.FHIRFieldType {

		case "identifier", "coding":
			// if Array Path Exists
			if arryLen > 0 {
				for i := 0; i <= arryLen; {
					var tempQuery string
					// Construct the First array Parameter
					if i == 0 {
						queryString = append(queryString, fmt.Sprintf("ANY a%d IN r.%s SATISFIES ", i, field.ArrayPath[i]))
					} else if i == arryLen && arryLen >= 1 {
						// The End Query

						if objectPath == "" {
							tempQuery = fmt.Sprintf("a%d %s END", i-1, conNVal)
						} else {
							if objectPath == "system" {
								tempQuery = fmt.Sprintf("a%d.`%s` %s '%s' END", i-1, objectPath, param.Condition, param.Value.Codable.System)
							} else if objectPath == "value" {
								tempQuery = fmt.Sprintf("a%d.`%s` %s '%s' END", i-1, objectPath, param.Condition, param.Value.Codable.Code)
							} else if objectPath == "code" {
								tempQuery = fmt.Sprintf("a%d.`%s` %s '%s' END", i-1, objectPath, param.Condition, param.Value.Codable.Code)
							} else {
								tempQuery = fmt.Sprintf("a%d.`%s` %s END", i-1, objectPath, conNVal)
							}
						}
						queryString = append(queryString, tempQuery)
						// If multiple array then add the END
						if i > 1 {
							queryString = append(queryString, fmt.Sprintf(" END"))
						}
					} else if i == arryLen && arryLen == 1 {
						if objectPath == "system" {
							tempQuery = fmt.Sprintf("a%d.`%s` %s '%s'", i-1, objectPath, param.Condition, param.Value.Codable.System)
						} else if objectPath == "value" {
							tempQuery = fmt.Sprintf("a%d.`%s` %s '%s'", i-1, objectPath, param.Condition, param.Value.Codable.Code)
						} else if objectPath == "code" {
							tempQuery = fmt.Sprintf("a%d.`%s` %s '%s'", i-1, objectPath, param.Condition, param.Value.Codable.Code)
						} else {
							tempQuery = fmt.Sprintf("a%d.`%s` %s", i-1, objectPath, conNVal)
						}
						queryString = append(queryString, tempQuery)
					} else {
						// convert object.array into object`.`array
						arrayPath := strings.Replace(field.ArrayPath[i], ".", "`.`", -1)
						q := fmt.Sprintf("ANY a%d IN a%d.`%s` SATISFIES ", i, i-1, arrayPath)
						queryString = append(queryString, q)
					}
					i++
				}
			} else {
				// For object parameter
				queryString = append(queryString, fmt.Sprintf("r.`%s` %s", objectPath, conNVal))
			}

			if len(queryParam) > 1 && n < len(queryParam)-1 {
				queryString = append(queryString, " AND ")
			}
		default:
			// if Array Path Exists
			if arryLen > 0 {
				for i := 0; i <= arryLen; {
					var tempQuery string
					// Construct the First array Parameter
					if i == 0 {
						queryString = append(queryString, fmt.Sprintf("ANY a%d IN r.%s SATISFIES ", i, field.ArrayPath[i]))
					} else if i == arryLen && arryLen >= 1 {
						// The End Query
						if objectPath == "" {
							tempQuery = fmt.Sprintf("a%d %s END", i-1, conNVal)
						} else {
							tempQuery = fmt.Sprintf("a%d.`%s` %s END", i-1, objectPath, conNVal)
						}
						queryString = append(queryString, tempQuery)
						// If multiple array then add the END
						if i > 1 {
							queryString = append(queryString, fmt.Sprintf(" END"))
						}
					} else if i == arryLen && arryLen == 1 {
						tempQuery = fmt.Sprintf("a%d.`%s` %s", i-1, objectPath, conNVal)
						queryString = append(queryString, tempQuery)
					} else {
						tempQuery = fmt.Sprintf("ANY a%d IN a%d.`%s` SATISFIES ", i, i-1, field.ArrayPath[i])
						queryString = append(queryString, tempQuery)
					}
					i++
				}
			} else {
				// For object parameter
				queryString = append(queryString, fmt.Sprintf("r.`%s` %s", objectPath, conNVal))
			}

			if paramLength > 1 && param.FHIRType == "string" {
				if len(queryParam) > 1 && n < len(queryParam)-1 {
					tempQuery = fmt.Sprintf(" and r.`%s` = '%s') OR (r.`%s` = '%s' and ", builder.resourceIdentifierName, param.Resource, builder.resourceIdentifierName, param.Resource)
					queryString = append(queryString, tempQuery)
				}
				// append ) at then end of query if multiple found

				if n == paramLength-1 {
					queryString = append(queryString, ")")
				}

			} else {
				if len(queryParam) > 1 && n < len(queryParam)-1 {
					queryString = append(queryString, " OR ")
					queryString = append(queryString, tempQuery)
				}
			}

		}
	}

	result := strings.Join(queryString, "")

	return result
}
