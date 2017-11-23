package tests

import (
	decoder "github.com/bhromor/uriql"
	search "github.com/bhromor/uriql/builder"
	"github.com/bhromor/uriql/dictionary"
	"github.com/bhromor/uriql/models"
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
