package builder

import (
	"fmt"
	"github.com/kite-social/uriql/helper"
	"github.com/kite-social/uriql/models"
	"strings"
)

type QueryBuilder struct {
	Query []string
}

func (q *QueryBuilder) BuildQueryIndex(bucket string, resource string, dict map[string]models.SearchParam) (idx []string) {

	for _, param := range dict {
		var idx []string

		for _, path := range param.Path {

			fieldStack := helper.GetFieldInfoFromPath(path)
			arryLen := len(fieldStack.ArrayPath)


			bucketQuery := fmt.Sprintf("CREATE INDEX `%s` ON `%s`",  fieldStack.Name, bucket)
			idx = append(idx, bucketQuery)

			if arryLen > 0 {
				for i := 0; i <= arryLen; {

					// Construct the first array Parameter
					if i == 0 {
						idx = append(idx, fmt.Sprintf("(DISTINCT ARRAY "))
					} else if i == arryLen {
						// Construct the last array Parameter
						num := arryLen - i
						// If multiple array found then syntax will be difference
						if arryLen < 2 {
							field := fieldStack.ArrayPath[arryLen-i]
							if fieldStack.ObjectPath != "" {
								// Condition for []array.obj
								idx = append(idx, fmt.Sprintf("a%d.%s FOR a%d IN %s END, %s)", num, fieldStack.ObjectPath, num, field, field))
							} else {
								// For Covered Array Indexing , the Last parameter suppose to be the array
								covered := strings.Split(field, ".")
								// Condition for []array.[]array or []array.obj.[]array
								idx = append(idx, fmt.Sprintf("a%d FOR a%d IN %s END, %s)", num, num, field, covered[len(covered)-1]))
							}
						}else{
							idx = append(idx, fmt.Sprintf("FOR a%d IN %s END)", num, fieldStack.ArrayPath[arryLen-i]))
						}

					} else {
						// Everything in between the first array Parameter
						if fieldStack.ObjectPath != "" {
							idx = append(idx, fmt.Sprintf("(DISTINCT ARRAY a%d.%s FOR a%d IN a%d.%s END) ", i, fieldStack.ObjectPath, i, i-1, fieldStack.ArrayPath[i]))
						} else {
							idx = append(idx, fmt.Sprintf("(DISTINCT ARRAY a%d FOR a%d IN a%d.%s END) ", i, i, i-1, fieldStack.ArrayPath[i]))
						}
					}
					i++
				}
			} else {
				// For object parameter
				idx = append(idx, fmt.Sprintf("(%s)", fieldStack.ObjectPath))
			}

			endQuery := fmt.Sprintf(" WHERE resourceType = '%s' ", resource)
			idx = append(idx, endQuery)

			fmt.Println(strings.Join(idx, ""))
			idx = []string{}


		}
	}

	return idx
}
