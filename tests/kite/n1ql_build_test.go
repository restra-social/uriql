package kite

import (
	"testing"
	search "udhvabon.com/neuron/uriql"
	"udhvabon.com/neuron/uriql/dictionary"
)

func printResult(t *testing.T, p string , qp interface{}, q string) {
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	t.Logf("Query %s", q)
}

func TestN1QLBuild(t *testing.T) {

	decode := search.GetQueryDecoder(dictionary.RestaurantDictionary())
	builder := search.GetN1QLBuilder()

	t.Log("Testing Restaurant Parameter : ")

	p := "restaurant?title=Vuter"
	qp := decode.DecodeQueryString(p)
	q := builder.Build(qp)
	printResult(t, p, qp, q)

	p = "restaurant?address=dhaka"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "restaurant?test=ac"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	decode = search.GetQueryDecoder(dictionary.RestaurantItemsDictionary())
	builder = search.GetN1QLBuilder()

	t.Log("Testing Restaurant Items Parameter : ")


	p = "restaurant-items?foods=burger"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "restaurant-items?foods-price=45"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

}
