package uriql

import "udhvabon.com/neuron/uriql/models"

type URIDecoder struct {
	Decoder *QueryDecoder
}

func GetURIDecoder(dict map[string]map[string]models.SearchParam) *URIDecoder {
	return &URIDecoder{
		Decoder: GetQueryDecoder(dict),
	}
}

func (b *URIDecoder) Decode(query string) models.QueryParam{

	return b.Decoder.DecodeQueryString(query)
}