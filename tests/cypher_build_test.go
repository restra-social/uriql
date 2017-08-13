package tests

import (
	"testing"
	search "github.com/kite-social/uriql/qbuilder"
	decoder "github.com/kite-social/uriql"
	"github.com/kite-social/uriql/dictionary"
	"github.com/kite-social/uriql/models"
)


func TestCypherBuild(t *testing.T) {

	dict := &models.Dictionary{ Model: dictionary.CypherDictionary()}

	decode := decoder.GetQueryDecoder(dict)
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
