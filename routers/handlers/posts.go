package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/bytesblog/models"
)

// Public unprotected api
//
// Get post for non-authenticated users
func (a *ApiHandler) GetPublicPost(c *gin.Context) {
	post_id := c.Param("id")
	db := a.db

	var post models.Post
	err := db.First(&post, post_id).Error
	if err == nil && post.AllowPublic {
		c.JSON(http.StatusOK, gin.H{
			"post": post,
		})
		return
	}

	c.AbortWithStatusJSON(404, gin.H{
		"error": "Cannot find post with id: " + post_id,
	})
}


