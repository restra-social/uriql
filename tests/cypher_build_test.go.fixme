package tests

import (
	decoder "github.com/restra-social/uriql"
	search "github.com/restra-social/uriql/builder"
	"github.com/restra-social/uriql/dictionary"
	"github.com/restra-social/uriql/models"
	"testing"
)

func TestCypherBuild(t *testing.T) {

	dict := &models.Dictionary{Model: dictionary.CypherDictionary()}

	decode := decoder.GetQueryDecoder(dict)
	builder := search.GetCypherBuilder()

	t.Log("Testing Friend Parameter : ")

	p := "friend?status=pending"
	qp := decode.DecodeQueryString(models.RequestInfo{UserID: "UFHFH35", Type: "User", Query: p})
	q := builder.Build(qp)
	printResult(t, p, qp, q)

	p = "friend?since=12-04-2017"
	qp = decode.DecodeQueryString(models.RequestInfo{UserID: "UFHFH35", Type: "User", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

	p = "friend?user=12345"
	qp = decode.DecodeQueryString(models.RequestInfo{UserID: "12345", Type: "User", Query: p})
	q = builder.Build(qp)
	printResult(t, p, qp, q)

}
