package tests

import (
	"testing"
	"udhvabon.com/neuron/uriql/models"
	search "udhvabon.com/neuron/uriql"
	"reflect"
)

func checkVal(t *testing.T, expect *models.QueryParam, got *models.QueryParam) {

	if reflect.DeepEqual(expect.Object, got.Object){

	}else{
		t.Fatalf("Object did not match : Exprected %+v got %+v", expect.Object, got.Object)
	}

	if reflect.DeepEqual(expect.Array, got.Array){

	} else {
		t.Fatalf("Could not find Array : Exprected %+v got %+v", expect.Array, got.Array)
	}

	if reflect.DeepEqual(got.Field, expect.Field) {

	} else {
		t.Fatalf("Field Value did not match : Exprected %+v got %+v ", expect.Field, got.Field)
	}

	if got.Condition != expect.Condition {
		t.Fatalf("Could not find Condition : Exprected %+v got %+v", expect.Condition, got.Condition)
	}
	if got.FHIRType != expect.FHIRType {
		t.Fatalf("Decode Type mismatched : Exprected %+v got %+v", expect.FHIRType, got.FHIRType)
	}

	if cap(got.Value) > 1 {
		for _, k := range got.Value {
			if k == expect.Value[0] || k == expect.Value[1] {

			} else {
				t.Fatalf("Could not find Array : Exprected %s or %s got %s ", expect.Value[0], expect.Value[1], k)
			}
		}
	} else {
		if got.Value[0] != expect.Value[0] {
			t.Fatalf("Could not Decode Value : Expected %s got %s", expect.Value[0], got.Value[0])
		}
	}
}

func TestDecode(t *testing.T) {

	dict := map[string]map[string]models.SearchParam{

		"Patient": map[string]models.SearchParam{

			"active": models.SearchParam{
				Type:      "token",
				FieldType: "boolean",
				Path:      []string{"active"},
			},
			"name": models.SearchParam{
				Type:      "string",
				FieldType: "string",
				Path:      []string{"[]name.[]family", "[]name.[]given"},
			},
		},
	}

	decode := search.GetQueryDecoder(dict)

	t.Log("Testing Universal Parameter : ")


	p := "Patient?name:contains=Mr."
	qp := decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Array: []string{"name"},
		Field: []models.FieldInfo{
			{
				"family",
				true,
			},
			{
				"given",
				true,
			},
		},
		Condition: "like",
		FHIRType: "string",
		Value: []string{"Mr."},
	}, &qp)

	p = "Patient?active=true"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Field: []models.FieldInfo{
			{
				"active",
				false,
			},
		},
		Condition: "=",
		FHIRType: "token",
		Value: []string{"true"},
	}, &qp)




}
