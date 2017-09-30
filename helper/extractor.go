package helper

import (
	"github.com/kite-social/uriql/models"
	"regexp"
	"strings"
)

type StackInfo struct {
	Fields     *Stack
	Length     int
	ArrayCount int
	Name       string
}

// Path could be []name.[]family , []address.state, active , managingOrganization.reference
func GetFieldInfoFromPath(str string) (s StackInfo) {
	s.ArrayCount = 0
	s.Length = 0
	s.Fields = NewStack()

	var fv models.FieldInfo
	//count = 0 // initial array count

	if strings.Contains(str, ".") {
		fields := strings.Split(str, ".")

		// loop through the end of the path except for the last
		for _, field := range fields {

			if strings.HasPrefix(field, "[]") {
				// lets say if []address.city
				fv.Array = true
				fv.Object = false
				fv.Field = field[2:len(field)] // address
				s.ArrayCount++
			} else {
				// lets say if managingOrganization.reference
				fv.Field = field // managingOrganization
				fv.Array = false
				fv.Object = true
			}
			s.Fields.Push(&Node{fv})
			s.Length++
		}
	} else {
		// if active, gender
		fv.Array = false
		fv.Object = false
		fv.Field = str

		s.Fields.Push(&Node{fv})
		s.Length++
	}

	r := regexp.MustCompile(`[^\\w[.\]]+`)
	indexName := r.FindAllString(str, -1)

	s.Name = strings.Join(indexName, "_")
	return s
}
