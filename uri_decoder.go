package uriql

import "udhvabon.com/neuron/uriql/models"

type URIDecoder struct {
	Decoder *QueryDecoder
}

func GetURIDecoder() *URIDecoder {
	return &URIDecoder{
		Decoder: GetQueryDecoder(),
	}
}

func (b *URIDecoder) Decode(query string) models.QueryParam{

	return b.Decoder.DecodeQueryString(query)
}