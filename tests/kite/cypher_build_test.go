package kite

import (
	"testing"
	search "udhvabon.com/neuron/uriql/qbuilder"
	decoder "udhvabon.com/neuron/uriql"
	"udhvabon.com/kiteengine/knet/uriql/dictionary"
	"udhvabon.com/neuron/uriql/models"
)


func TestCypherBuild(t *testing.T) {

	decode := decoder.GetQueryDecoder(dictionary.UserFriendDictionary())
	builder := search.GetCypherBuilder()

	t.Log("Testing Friend Parameter : ")

	p := "friend?status=pending"
	qp := decode.DecodeQueryString(models.RequestInfo{UserId: "UFHFH35", Type: "User", Query: p})
	q := builder.Build(qp)
	printResult(t, p, qp, q)

	p = "friend?since=12-04-2017"
	qp = decode.DecodeQueryString(models.RequestInfo{UserId: "UFHFH35", Type: "User", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "friend?user=12345"
	qp = decode.DecodeQueryString(models.RequestInfo{UserId: "12345", Type: "User", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

}
