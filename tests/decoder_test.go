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

	decode := search.QueryDecoder{}

	t.Log("Testing Universal Parameter : ")

	p := "Patient?_id=1234567890"
	qp := decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Field: []models.FieldInfo{
			{
				"id",
				false,
			},
		},
		Condition: "=",
		FHIRType:  "universal",
		Value:     []string{"1234567890"},
	}, &qp)

	p = "Patient?_lastUpdated=gt2010-10-01"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Object: []string{"meta"},
		Field: []models.FieldInfo{
			{
				"lastUpdated",
				false,
			},
		},
		Condition: ">",
		FHIRType:  "universal",
		Value:     []string{"2010-10-01"},
	}, &qp)

	t.Log("Testing NUMBER Parameter : ")

	p = "Encounter?length=gt204"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Field: []models.FieldInfo{
			{
				"length",
				false,
			},
		},
		Condition: ">",
		FHIRType:  "number",
		Value:     []string{"204"},
	}, &qp)

	p = "Encounter?length=ge6000"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Field: []models.FieldInfo{
			{
				"length",
				false,
			},
		},
		Condition: ">=",
		FHIRType: "number",
		Value: []string{"6000"},
	}, &qp)

	p = "Encounter?length=le27.5"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Field: []models.FieldInfo{
			{
				"length",
				false,
			},
		},
		Condition: "=<",
		FHIRType: "number",
		Value: []string{"27.5"},
	}, &qp)

	p = "Encounter?length=1029"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Field: []models.FieldInfo{
			{
				"length",
				false,
			},
		},
		Condition: "=",
		FHIRType: "number",
		Value: []string{"1029"},
	}, &qp)

	t.Log("\nTesting STRING Parameter : ")

	p = "Patient?name:contains=Mr."
	qp = decode.DecodeQueryString(p)
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

	p = "Patient?name=Fahim"
	qp = decode.DecodeQueryString(p)
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
		Condition: "=",
		FHIRType: "string",
		Value: []string{"Fahim"},
	}, &qp)

	p = "Patient?name:exact=Shariar"
	qp = decode.DecodeQueryString(p)
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
		Condition: "=",
		FHIRType: "string",
		Value: []string{"Shariar"},
	}, &qp)


	t.Log("\nTesting TOKEN Parameter : ")

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

	p = "Patient?gender=male"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Field: []models.FieldInfo{
			{
				"gender",
				false,
			},
		},
		Condition: "=",
		FHIRType: "token",
		Value: []string{"male"},
	}, &qp)


	p = "Patient?language=https://code.repo.org.bn|BN"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Array: []string{"communication"},
		Field: []models.FieldInfo{
			{
				"language",
				false,
			},
		},
		Condition: "=",
		FHIRType: "token",
		Value: []string{"https://code.repo.org.bn", "BN"},
	}, &qp)

	p = "Patient?identifier=|1234567"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Array: []string{"identifier"},
		Field: []models.FieldInfo{
			{
				"value",
				false,
			},
		},
		Condition: "=",
		FHIRType: "token",
		Value: []string{"", "1234567"},
	}, &qp)


	t.Log("\nTesting REFERENCE Parameter : ")

	p = "Patient?general-practitioner=Practitioner/23"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Array: []string{"generalPractitioner"},
		Field: []models.FieldInfo{
			{
				"reference",
				false,
			},
		},
		Condition: "=",
		FHIRType: "reference",
		Value: []string{"Practitioner", "23"},
	}, &qp)

	p = "Patient?general-practitioner:Practitioner=23"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Array: []string{"generalPractitioner"},
		Field: []models.FieldInfo{
			{
				"reference",
				false,
			},
		},
		Condition: "=",
		FHIRType: "reference",
		Value: []string{"Practitioner", "23"},
	}, &qp)

	p = "Patient?organization=Organization/3456"
	qp = decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	checkVal(t, &models.QueryParam{
		Object: []string{"managingOrganization"},
		Field: []models.FieldInfo{
			{
				"reference",
				false,
			},
		},
		Condition: "=",
		FHIRType: "reference",
		Value: []string{"Organization", "3456"},
	}, &qp)



}
