package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/internal/app"
	"github.com/teenyblog/models"
	"gorm.io/gorm"
)


func (*ApiHandler) CurrentUser(c *gin.Context){
	user, _ := currentSession(c)
	
	c.JSON(200, gin.H{
		"status": 200,
		"message": "OK",
		"user": user,
	})
}

func (a *ApiHandler) Logout(c *gin.Context) {
	_, session := currentSession(c)
	
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	app.DeleteSession(a.db, c.Request.Context(), session.ID);
	
	
	c.Redirect(http.StatusFound, "/")
}

func (a *ApiHandler) ListPosts(c *gin.Context) {
	user, _ := currentSession(c)
	db := a.db

	posts, err := gorm.G[models.Post](db).Where("user_id = ?", user.ID).Find(c.Request.Context())
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



// Get a private post by id. Note that this also works for public posts.
func (a *ApiHandler) GetPost(c *gin.Context) {
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
		"error": "Cannot find post with id:" + post_id,
	})
	
}


func (a *ApiHandler) PublishPost(c *gin.Context) {
	user, _ := currentSession(c)
	post_id := c.Param("id")
	make_public := c.Query("public") == "true"
	db := a.db

	var post models.Post
	err := db.First(&post, post_id).Error
	if err == nil && post.UserID == user.ID {
		post.AllowPublic = make_public
		db.Save(&post)
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
				"error": "Post not found",
			})
			return
		}

		// Check if post exists and belongs to the user
		var existingPost models.Post
		err = db.First(&existingPost, post_id).Error
		if err != nil || existingPost.UserID != user.ID {
			c.AbortWithStatusJSON(404, gin.H{
				"error": "Post not found",
			})
			return
		}
		
		post.ID = uint(post_id)
		err = db.Save(&post).Error
		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{

				"error": "Failed to save post :" + err.Error(),
			})
			return
		}
	} else {
		err := db.Create(&post).Error
		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{

				"error": "Failed to save post :" + err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"post": post,
	})
}