package meta

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PageQuery struct {
	Page  int `form:"page,default=1,min=1"`
	Limit int `form:"limit,default=10,max=100"`
	Sort  string `form:"sort,default=none"`
}

func (q *PageQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

func ApplySqlPagination[T any](q PageQuery,sqlQuery gorm.ChainInterface[T}]) gorm.ChainInterface[T] {
	return  sql.
		Limit(q.Limit).
		Offset(q.Offset())
}

func BindPageQuery(c *gin.Context) (PageQuery, error) {
	var query PageQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		return query, err
	}

	return query, nil
}

// Retrieves the page query from the request. This requires the PaginationQuery middleware
func GetPagination(c *gin.Context) PageQuery {
	query := c.MustGet("pagination_query").(PageQuery)
	return query
}
