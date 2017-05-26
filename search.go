package uriql

import "udhvabon.com/neuron/uriql/models"

type Uriql struct {
	Decoder *URIDecoder
	Builder *N1QLQueryBuilder
}

func GetUriql(dict map[string]map[string]models.SearchParam) *Uriql {
	return &Uriql{
		Decoder: GetURIDecoder(dict),
		Builder: GetN1QLBuilder(),
	}
}

func (f *Uriql) UrlToN1Ql(query string) string {
	decode := f.Decoder.Decode(query)
	return f.Builder.Build(decode)
}