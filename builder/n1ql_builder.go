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

	for l, param := range queryParam {

		field := param.DictionaryInfo.FieldsInfo

		arryLen := len(field.ArrayPath)

		// if Array Path Exists
		if arryLen > 0 {
			for i := 0; i <= arryLen; {
				// Construct the first array Parameter
				if i == 0 {
					queryString = append(queryString, fmt.Sprintf("ANY a%d IN r.%s SATISFIES ", i, field.ArrayPath[i]))
				}else if i == arryLen && arryLen > 1 {
					if field.ObjectPath == "" {
						o := fmt.Sprintf("a%d %s '%s' END", i-1, param.Condition, param.Value.Value)
						queryString = append(queryString, o)
					}
					//queryString = append(queryString, fmt.Sprintf(" END"))
				}else if arryLen == 1 || (arryLen > 1 && field.ObjectPath != "") {
						// Object part , the last JOb assign value to field
						o := fmt.Sprintf("a%d.%s %s '%s' END", i-1, field.ObjectPath, param.Condition, param.Value.Value)
						queryString = append(queryString, o)
				}else{
					q := fmt.Sprintf("ANY a%d IN a%d.%s SATISFIES ", i, i-1, field.ArrayPath[i])
					queryString = append(queryString, q)
				}
				i++
			}
		} else {
			// For object parameter
			queryString = append(queryString, fmt.Sprintf("r.%s = '%s'", param.FieldsInfo.ObjectPath, param.Value.Value))
		}

		if len(queryParam) > 1 && l < len(queryParam)-1 {
			queryString = append(queryString, " OR ")
		}

	}

	result := strings.Join(queryString, "")

	return result
}
