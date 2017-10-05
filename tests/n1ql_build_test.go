package tests

import (
	decoder "github.com/kite-social/uriql"
	search "github.com/kite-social/uriql/builder"
	"github.com/kite-social/uriql/dictionary"
	"github.com/kite-social/uriql/models"
	"testing"
)

func printResult(t *testing.T, p string, qp interface{}, q string) {
	t.Logf("Decoding : %s", p)
	//t.Logf("Decoded to : %+v", qp)
	t.Logf("Query %s", q)
}

func TestN1QLBuild(t *testing.T) {

	dict := &models.Dictionary{Model: dictionary.N1QLDictionary(), Bucket: "kite"}

	decode := decoder.GetQueryDecoder(dict)
	builder := search.GetN1QLQueryBuilder(dict.Bucket)

	t.Log("Testing Restaurant Parameter : ")

	p := "Patient?language=http://acme.org/patient|BN"
	qp := decode.DecodeQueryString(models.RequestInfo{UserID: "1234567890", Type: "Patient", Query: p})
	q := builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?name:contains=Mr."
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
