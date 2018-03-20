package tests

import (
	decoder "github.com/restra-social/uriql"
	search "github.com/restra-social/uriql/builder"
	"github.com/restra-social/uriql/dictionary"
	"github.com/restra-social/uriql/models"
	"testing"
)

func printResult(t *testing.T, p string, qp interface{}, q string) {
	t.Logf("Decoding : %s", p)
	//t.Logf("Decoded to : %+v", qp)
	t.Logf("Query %s", q)
}

func TestN1QLBuild(t *testing.T) {

	dict := &models.Dictionary{Model: dictionary.N1QLDictionary(), Bucket: "test"}

	decode := decoder.GetQueryDecoder(dict)
	builder := search.GetN1QLQueryBuilder(dict.Bucket, "type")

	t.Log("Testing Restaurant Parameter : ")

	var filter models.Filter

	p := "order?store_id=1235"
	qp, err := decode.DecodeQueryString(models.RequestInfo{Type: "order", Query: p}, &filter)
	if err != nil {
		t.Errorf(err.Error())
	}
	q := builder.Build(qp)
	printResult(t, p, qp, q)

	/*	p := "profiles?name:contains=mr&hobbies=sports&_size=10&_page=2"
		qp, err := decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "profile", Query: p}, &filter)
		if err != nil {
			t.Errorf(err.Error())
		}
		q := builder.Build(qp)
		printResult(t, p, qp, q)*/

	/*p = "profiles?hobbies=sports&_size=10&_page=1"
	qp, err = decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "profile", Query: p}, &filter)
	if err != nil {
		t.Errorf(err.Error())
	}
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "profiles?name:contains=mr&_page=3"
	qp, err = decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "profile", Query: p}, &filter)
	if err != nil {
		t.Errorf(err.Error())
	}
	q = builder.Build(qp)
	printResult(t, p, qp, q)*/

	/*p = "Patient?name:contains=Mr."
	qp = decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "Patient", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?organization=Medical/10234"
	qp = decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "Patient", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?identifier=http://acme.org/patient|2345"
	qp = decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "Patient", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?general-practitioner=Practitioner/2345"
	qp = decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "Patient", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?active=true"
	qp = decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "Patient", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?address-use=Dhaka"
	qp = decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "Patient", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)*/

	///Observation?subject=Patient/23

	/*decode = search.GetQueryDecoder(dictionary.RestaurantItemsDictionary())
	builder = builder.GetN1QLQueryBuilder()

	t.Log("Testing Restaurant Items Parameter : ")


	p = "restaurant-items?foods=burger"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "restaurant-items?foods-price=45"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)*/

}
