package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/models"
	"github.com/teenyblog/routers/meta"
	"gorm.io/gorm"
)


func (a *ApiHandler) ListUserPublicPosts(c *gin.Context) {
	user_id := c.Param("user_id")
	db := a.db

	posts, err := gorm.G[models.Post](db).Where("user_id = ? AND allow_public = true", user_id).Find(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	summaries := make([]meta.PostSummary, 0, len(posts))
	for _, post := range posts {
		summaries = append(summaries, meta.PostSummary{
			Model:   post.Model,
			Title:   post.Title,
			Summary: post.Summary,
			AllowPublic: post.AllowPublic,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": summaries,
	})
	
}
