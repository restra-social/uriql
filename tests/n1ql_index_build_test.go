package tests

import (
	decoder "github.com/restra-social/uriql"
	"github.com/bhromor/uriql/dictionary"
	"github.com/bhromor/uriql/models"
	"testing"
)

func TestN1QLIndexBuild(t *testing.T) {

	dict := &models.Dictionary{Model: dictionary.N1QLDictionary(), Bucket: "kite"}

	d := decoder.GetQueryDecoder(dict)

	idxx := d.DecodeQueryIndex()

	for _, idx := range idxx {

		//t.Log(idx)
		for _, k := range idx.Indexes {
			t.Log(k)
		}
	}

}
