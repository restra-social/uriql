package builder

import (
	"github.com/restra-social/uriql/models"
)


type FullTextSearchQueryBuilder struct {
	page                   int
	limit                  int
}

// GetFullTextSearchQueryBuilder : Get N1QL Builder Object
func GetFullTextSearchQueryBuilder() *FullTextSearchQueryBuilder {
	return &FullTextSearchQueryBuilder{}
}

func (builder *FullTextSearchQueryBuilder) Build(queryInfo *models.QueryInfo) string {

	return ""
}
