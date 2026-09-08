package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/inuis/internal/app"
)

// Public unprotected api
//
// Get post for non-authenticated users
func (a *ApiHandler) GetPublicPost(c *gin.Context) {
	post_hid := c.Param("hid")

	post, ok := app.FindPublicPostByHID(post_hid, a.db, c.Request.Context())
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"post": post,
		})
		return
	}

	c.AbortWithStatusJSON(404, gin.H{
		"error": "Cannot find post with hid: " + post_hid,
	})
}
