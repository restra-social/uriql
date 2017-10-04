package helper

import (
	"fmt"
	"github.com/kite-social/uriql/models"
	"regexp"
	"strings"
)

type StackInfo struct {
	//Fields *Queue
	Length     int
	ArrayPath  []string
	ObjectPath string
	Name       string
}

// Path could be []name.[]family , []address.state, active , managingOrganization.reference
func GetFieldInfoFromPath(str string) (s StackInfo) {
	var arr string
	s.Length = 0
	var objPath []string
	//s.Fields = NewQueue(6) // #todo dynamic queue length !! very bad practice

	var fv models.FieldInfo
	//count = 0 // initial array count

	if strings.Contains(str, ".") {
		fields := strings.Split(str, ".")

		var tempObj []string

		// loop through the end of the path except for the last
		for _, field := range fields {

			if strings.HasPrefix(field, "[]") {
				// lets say if []address.city
				fv.Array = true
				fv.Object = false
				fv.Field = field[2:len(field)] // address

				if len(tempObj) > 0 {
					arr = fmt.Sprintf("%s.%s", strings.Join(tempObj, "."), fv.Field)
				} else {
					arr = fmt.Sprintf("%s", fv.Field)
				}
				s.ArrayPath = append(s.ArrayPath, arr)
				// empty the temp object
				tempObj = []string{}
				objPath = []string{}
			} else {
				// lets say if managingOrganization.reference
				fv.Field = field // managingOrganization
				fv.Array = false
				fv.Object = true
				objPath = append(objPath, fv.Field)

				tempObj = append(tempObj, fv.Field)

			}
			//s.Fields.Push(&Node{fv})
			s.Length++
		}
	} else if strings.HasPrefix(str, "[]") {
		// lets say if []address
		fv.Array = true
		fv.Object = false
		fv.Field = str[2:len(str)] // address

		s.ArrayPath = append(s.ArrayPath, fv.Field)

	} else {
		// if active, gender
		fv.Array = false
		fv.Object = false
		fv.Field = str

		//s.Fields.Push(&Node{fv})
		objPath = append(objPath, fv.Field)
		s.Length++
	}

	r := regexp.MustCompile(`[^\\w[.\]]+`)
	indexName := r.FindAllString(str, -1)

	s.Name = strings.Join(indexName, "_")
	s.ObjectPath = strings.Join(objPath, ".")

	return s
}
