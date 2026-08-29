package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/bytesblog/models"
	"github.com/ultraflame4/bytesblog/routers/meta"
	"gorm.io/gorm"
)


func (a *ApiHandler) ListUserPublicPosts(c *gin.Context) {
	db := a.db
	user_id := c.Param("user_id")
	pageQuery := meta.GetPagination(c)
	ctx := c.Request.Context()

	sql := gorm.G[models.Post](db).Where("user_id = ? AND allow_public = true", user_id)
	total_count, err := sql.Count(ctx, "ID")
	posts, err := sql.Find(ctx)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	summaries := make([]meta.PostSummary, 0, len(posts))

	type Response struct {
		meta.PageMeta
		Posts []meta.PostSummary `json:"posts"`
	}
	
	for _, post := range posts {
		summaries = append(summaries, meta.PostSummary{
			Model:   post.Model,
			Title:   post.Title,
			Summary: post.Summary,
			AllowPublic: post.AllowPublic,
		})
	}

	c.JSON(http.StatusOK, Response{
		PageMeta: meta.CreatePageMeta(pageQuery, int(total_count)),
		Posts: summaries,
	})
	
}
