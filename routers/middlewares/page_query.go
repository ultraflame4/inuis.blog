package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/bytesblog/routers/meta"
)

// Extracts the pagination query parameters from the request and sets them on the context
//
// Aborts with 400 Bad Request if the query parameters are invalid
func PaginationQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		query, err := meta.BindPageQuery(c);
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Set("pagination_query", query)
		c.Next()
	}
}
