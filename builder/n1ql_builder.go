package builder

import (
	"fmt"
	"github.com/restra-social/uriql/models"
	"strings"
)

type n1QLQueryBuilder struct {
	bucketName string
}

// GetN1QLQueryBuilder : Get N1QL Builder Object
func GetN1QLQueryBuilder(bucket string) *n1QLQueryBuilder {
	return &n1QLQueryBuilder{bucketName: bucket}
}

func (n *n1QLQueryBuilder) Build(queryParam []models.QueryParam) string {

	var queryString []string

	bucketQuery := fmt.Sprintf("SELECT * FROM `%s` as r WHERE ", n.bucketName) // #todo fix resource
	queryString = append(queryString, bucketQuery)

	for n, param := range queryParam {

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
			tempQuery = fmt.Sprintf("r.resourceType = '%s' and ", param.Resource)
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
				queryString = append(queryString, fmt.Sprintf("r.`%s` = '%s'", param.FieldsInfo.ObjectPath, param.Value.Value))
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
				queryString = append(queryString, fmt.Sprintf("r.`%s` = '%s'", param.FieldsInfo.ObjectPath, param.Value.Value))
			}

			if paramLength > 1 && param.FHIRType == "string" {
				if len(queryParam) > 1 && n < len(queryParam)-1 {
					tempQuery = fmt.Sprintf(" and r.resourceType = '%s') OR (r.resourceType = '%s' and ", param.Resource, param.Resource)
					queryString = append(queryString, tempQuery)
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
