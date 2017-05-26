package tests

import (
	"testing"
	search "udhvabon.com/neuron/uriql"
	"udhvabon.com/neuron/soma/uriql/dictionary"
)

func printResult(t *testing.T, p string , qp interface{}, q string) {
	t.Logf("Decoding : %s", p)
	t.Logf("Decoded to : %+v", qp)
	t.Logf("Query %s", q)
}

func TestN1QLBuild(t *testing.T) {

	decode := search.GetQueryDecoder(dictionary.FHIRDictionary())

	builder := search.GetN1QLBuilder()

	t.Log("Testing Universal Parameter : ")

	p := "Patient/1234567890"
	qp := decode.DecodeQueryString(p)
	q := builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?_id=1234567890"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)


	p = "Patient?_lastUpdated=le2010-10-01"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)


	t.Log("Testing NUMBER Parameter : ")

	p = "Encounter?length=gt204"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Encounter?length=ge6000"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Encounter?length=le27.5"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Encounter?length=1029"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)


	t.Log("\nTesting STRING Parameter : ")

	p = "Patient?name:contains=Mr."
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?name=Fahim"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?name:exact=Shariar"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)


	t.Log("\nTesting TOKEN Parameter : ")

	p = "Patient?active=true"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?gender=male"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?address-use=home"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)


	p = "Patient?language=https://code.repo.org.bn|BN"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?identifier=|1234567"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?identifier=https://code.neuron.health/identifier|1234567"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)


	t.Log("\nTesting REFERENCE Parameter : ")

	p = "Patient?general-practitioner=Practitioner/23"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?general-practitioner:Practitioner=23"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "Patient?organization=Organization/3456"
	qp = decode.DecodeQueryString(p)
	q = builder.Build(qp)
	printResult(t, p, qp, q)



}
