package builder

import (
	"fmt"
	"github.com/kite-social/uriql/helper"
	"github.com/kite-social/uriql/models"
	"strings"
)

func BuildQueryIndex(bucket string, resource string, dict map[string]models.SearchParam) (index []string) {

	for _, param := range dict {
		var idx string

		for _, path := range param.Path {

			fieldStack := helper.GetFieldInfoFromPath(path)

			if fieldStack.Length == 1 {
				idx += fmt.Sprintf("CREATE INDEX `%s_%s` ON `%s`(%s) WHERE resourceType = `%s`", strings.ToLower(resource), fieldStack.Name, bucket, path, resource)
				idx += fmt.Sprintf(" ")
			} else {
				var objects []string
				var fields string
				var arrays []string

				for i := 0; i < fieldStack.Length; {

					fieldInfo := fieldStack.Fields.Pop().Value

					if fieldInfo.Object == true {
						objects = append([]string{fieldInfo.Field}, objects...)
						if i > 0 {
							fields = strings.Join(objects, ".")
						}else{
							// this object is a part of array
							arrays = append(arrays, fmt.Sprintf("d%d.%s",i+1, strings.Join(objects,".")))
						}
					} else {
						if i == 0 {
							// Its an nested Array
							a := fmt.Sprintf("(DISTINCT ARRAY d%d FOR d%d IN d%d.%s END, %s)", i, i,i+1, fieldInfo.Field, fieldInfo.Field)
							arrays = append(arrays, a)
						} else {
							arrays = append([]string{"DISTINCT ARRAY"}, arrays...)
							arrays = append(arrays, fmt.Sprintf("FOR d%d IN %s END, %s", i, fieldInfo.Field, fieldInfo.Field))
						}
						fields = strings.Join(arrays, " ")
					}
					i++
				}

				idx += fmt.Sprintf("CREATE INDEX `%s_%s` ON `%s`(%s) WHERE resourceType = `%s`", strings.ToLower(resource), fieldStack.Name, bucket, fields, resource)
				idx += fmt.Sprintf(" ")
			}
		}

		idx += " "
		index = append(index, idx)
	}

	return index
}
