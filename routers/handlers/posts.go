package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/models"
)

// Public unprotected api
//
// Get post for non-authenticated users
func (a *ApiHandler) PublicGetPost(c *gin.Context) {
	post_id := c.Param("id")
	db := a.db

	var post models.Post
	err := db.First(&post, post_id).Error
	if err == nil && post.Published {
		c.JSON(http.StatusOK, gin.H{
			"post": post,
		})
		return
	}

	c.AbortWithStatusJSON(404, gin.H{
		"status": "error",
		"message": "Cannot find post with id:" + post_id,
	})
}

func (a *ApiHandler) GetUserPost(c *gin.Context) {
	post_id := c.Param("id")
	db := a.db

	var post models.Post
	err := db.First(&post, post_id).Error
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"post": post,
		})
		return
	}

	c.AbortWithStatusJSON(404, gin.H{
		"status": "error",
		"message": "Cannot find post with id:" + post_id,
	})
	
}


func (a *ApiHandler) PublishUserPost(c *gin.Context) {
	user, _ := currentSession(c)
	post_id := c.Param("id")
	db := a.db

	var post models.Post
	err := db.First(&post, post_id).Error
	if err == nil && post.UserID == user.ID {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
		return
	}

	c.AbortWithStatusJSON(404, gin.H{
		"status": "error",
		"message": "Cannot find post with id:" + post_id,
	})
	
}

func (a *ApiHandler) SavePost(c *gin.Context) {
	db := a.db
	user, _ := currentSession(c)

	post_id := c.Query("id");
	var post models.Post

	
	title := ""
	content := c.PostForm("content")
	body := c.PostForm("content")
	
	// Remove the first line of the content. (Treated as title)
	if i := strings.IndexByte(body, '\n'); i >= 0 {
		title = content[:i]
        body =  body[i+1:]
    }

    summary_len := 100
    if len(body) < summary_len {
        summary_len = len(body)
    }
    
    // Generate summary from body
    summary :=  strings.ReplaceAll(body, "\n",  "")[0:summary_len]
	

	post = models.Post{
		UserID: user.ID,
		Summary: summary,
		Content: content,
		Title: title,
	}

	if post_id != "" {
		post_id, err := strconv.ParseUint(post_id, 10, 32);
		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{
				"status": "error",
				"message": "Post not found",
			})
			return
		}

		// Check if post exists and belongs to the user
		var existingPost models.Post
		err = db.First(&existingPost, post_id).Error
		if err != nil || existingPost.UserID != user.ID {
			c.AbortWithStatusJSON(404, gin.H{
				"status": "error",
				"message": "Post not found",
			})
			return
		}
		
		post.ID = uint(post_id)
		err = db.Save(&post).Error
		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{
				"status": "error",
				"message": "Failed to save post",
				"error": err.Error(),
			})
			return
		}
	} else {
		err := db.Create(&post).Error
		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{
				"status": "error",
				"message": "Failed to save post.",
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"post": post,
	})
}