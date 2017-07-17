package kite

import (
	"testing"
	//"udhvabon.com/neuron/uriql/models"
)
//
//func checkVal(t *testing.T, expect *models.QueryParam, got *models.QueryParam) {
//
//	if reflect.DeepEqual(expect.Object, got.Object){
//
//	}else{
//		t.Fatalf("Object did not match : Exprected %+v got %+v", expect.Object, got.Object)
//	}
//
//	if reflect.DeepEqual(expect.Array, got.Array){
//
//	} else {
//		t.Fatalf("Could not find Array : Exprected %+v got %+v", expect.Array, got.Array)
//	}
//
//	if reflect.DeepEqual(got.Field, expect.Field) {
//
//	} else {
//		t.Fatalf("Field Value did not match : Exprected %+v got %+v ", expect.Field, got.Field)
//	}
//
//	if got.Condition != expect.Condition {
//		t.Fatalf("Could not find Condition : Exprected %+v got %+v", expect.Condition, got.Condition)
//	}
//	if got.FHIRType != expect.FHIRType {
//		t.Fatalf("Decode Type mismatched : Exprected %+v got %+v", expect.FHIRType, got.FHIRType)
//	}
//
//	if cap(got.Value) > 1 {
//		for _, k := range got.Value {
//			if k == expect.Value[0] || k == expect.Value[1] {
//
//			} else {
//				t.Fatalf("Could not find Array : Exprected %s or %s got %s ", expect.Value[0], expect.Value[1], k)
//			}
//		}
//	} else {
//		if got.Value[0] != expect.Value[0] {
//			t.Fatalf("Could not Decode Value : Expected %s got %s", expect.Value[0], got.Value[0])
//		}
//	}
//}

func TestDecode2(t *testing.T) {

	//decode := search.GetQueryDecoder(dictionary3.RestaurantDictionary())

	//t.Log("Testing Kite Parameter : ")

	/*p := "restaurant?address=dhaka"
	qp := decode.DecodeQueryString(p)
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	//checkVal(t, &models.QueryParam{*/
	//
	//	Condition: "=",
	//	FHIRType:  "universal",
	//	Value:     []string{"1234567890"},
	//}, &qp)
}
