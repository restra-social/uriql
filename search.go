package uriql

import (
	"github.com/kite-social/uriql/models"
	builder "github.com/kite-social/uriql/qbuilder"
)

type Uriql struct {
	Decoder *URIDecoder
	N1QLBuilder *builder.N1QLQueryBuilder
	CypherBuilder *builder.CypherQueryBuilder
}

func GetUriql(dict map[string]map[string]models.SearchParam) *Uriql {
	return &Uriql{
		Decoder: GetURIDecoder(dict),
		N1QLBuilder: builder.GetN1QLQueryBuilder(),
		CypherBuilder: builder.GetCypherBuilder(),
	}
}

func (f *Uriql) UrlToN1Ql(request models.RequestInfo) string {
	decode := f.Decoder.Decode(request)
	return f.N1QLBuilder.Build(decode)
}

func (f *Uriql) UrlToCypher(request models.RequestInfo) string {
	decode := f.Decoder.Decode(request)
	return f.CypherBuilder.Build(decode)
}
