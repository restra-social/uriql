package uriql

import "github.com/kite-social/uriql/models"

type URIDecoder struct {
	Decoder *QueryDecoder
}

func GetURIDecoder(dict map[string]map[string]models.SearchParam) *URIDecoder {
	return &URIDecoder{
		Decoder: GetQueryDecoder(dict),
	}
}

func (b *URIDecoder) Decode(request models.RequestInfo) []models.QueryParam{
	return b.Decoder.DecodeQueryString(request)
}