package builder

import (
	"fmt"
	"github.com/kite-social/uriql/models"
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

	bucketQuery := fmt.Sprintf("SELECT * FROM `%s` as r WHERE  r.resourceType = 'Patient' and ", n.bucketName) // #todo fix resource
	queryString = append(queryString, bucketQuery)

	for n, param := range queryParam {

		field := param.DictionaryInfo.FieldsInfo
		arryLen := len(field.ArrayPath)

		// if Array Path Exists
		if arryLen > 0 {
			for i := 0; i <= arryLen; {
				// Construct the First array Parameter
				if i == 0 {
					queryString = append(queryString, fmt.Sprintf("ANY a%d IN r.%s SATISFIES ", i, field.ArrayPath[i]))
				}else if i == arryLen && arryLen >= 1 {
					// The End Query
					if field.ObjectPath == "" {
						o := fmt.Sprintf("a%d %s '%s' END", i-1, param.Condition, param.Value.Value)
						queryString = append(queryString, o)
					}else{
						o := fmt.Sprintf("a%d.%s %s '%s' END", i-1, field.ObjectPath, param.Condition, param.Value.Value)
						queryString = append(queryString, o)
					}
					// If multiple array then add the END
					if i > 1 {
						queryString = append(queryString, fmt.Sprintf(" END"))
					}
				} else if i == arryLen && arryLen == 1 {
					q := fmt.Sprintf("a%d.%s %s '%s'", i-1, field.ObjectPath, param.Condition, param.Value.Value)
					queryString = append(queryString, q)
				} else{
					q := fmt.Sprintf("ANY a%d IN a%d.%s SATISFIES ", i, i-1, field.ArrayPath[i])
					queryString = append(queryString, q)
				}
				i++
			}
		} else {
			// For object parameter
			queryString = append(queryString, fmt.Sprintf("r.%s = '%s'", param.FieldsInfo.ObjectPath, param.Value.Value))
		}

		if len(queryParam) > 1 && n < len(queryParam)-1 {
			queryString = append(queryString, " OR ")
		}

	}

	result := strings.Join(queryString, "")

	return result
}
