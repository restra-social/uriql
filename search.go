package uriql

type Uriql struct {
	Decoder *URIDecoder
	Builder *N1QLQueryBuilder
}

func GetUriql() *Uriql {
	return &Uriql{
		Decoder: GetURIDecoder(),
		Builder: GetN1QLBuilder(),
	}
}

func (f *Uriql) UrlToN1Ql(query string) string {
	decode := f.Decoder.Decode(query)
	return f.Builder.Build(decode)
}