package uriql

import "github.com/restra-social/uriql/models"

// URIDecoder :
type URIDecoder struct {
	Decoder *QueryDecoder
}

// GetURIDecoder : Get Query Decoder Object
func GetURIDecoder(dict *models.Dictionary) *URIDecoder {
	return &URIDecoder{
		Decoder: GetQueryDecoder(dict),
	}
}

// Decode : Decodes Query String from Request Info
func (b *URIDecoder) Decode(request models.RequestInfo) []models.QueryParam {
	return b.Decoder.DecodeQueryString(request)
}
