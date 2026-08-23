package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/models"
	"gorm.io/gorm"
)

type _SummarisedPost struct {
	gorm.Model
	Title string
	Summary string
}

func (a *ApiHandler) ListUserPublicPosts(c *gin.Context) {
	user_id := c.Param("user_id")
	db := a.db

	posts, err := gorm.G[models.Post](db).Where("user_id = ? AND allow_public = true", user_id).Find(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	summaries := make([]_SummarisedPost, 0, len(posts))
	for _, post := range posts {
		summaries = append(summaries, _SummarisedPost{
			Model:   post.Model,
			Title:   post.Title,
			Summary: post.Summary,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": summaries,
	})
	
}
