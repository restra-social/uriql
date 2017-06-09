package kite

import (
	"testing"
	search "udhvabon.com/neuron/uriql"
	dictionary2 "udhvabon.com/kiteengine/docMan/uriql/dictionary"
)

func printResult(t *testing.T, p string , qp interface{}, q string) {
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	t.Logf("Query %s", q)
}

func TestN1QLBuild(t *testing.T) {

	decode := search.GetQueryDecoder(dictionary2.KiteDictionary())

	builder := search.GetN1QLBuilder()

	t.Log("Testing Universal Parameter : ")

	p := "restaurant?title=dhaka"
	qp := decode.DecodeQueryString(p)
	q := builder.Build(qp)
	printResult(t, p, qp, q)

	p = "restaurant?city=dhaka"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	/*
	p = "restaurant?city2=dhaka"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)*/
	//
	//p = "restaurant?city3=dhaka"
	//qp = decode.DecodeQueryString(p)
	//q = builder.Build(qp)
	//printResult(t, p, qp, q)
	//
	//p = "restaurant?city4=dhaka"
	//qp = decode.DecodeQueryString(p)
	//q = builder.Build(qp)
	//printResult(t, p, qp, q)





}
