package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/bytesblog/internal/app"
	"github.com/ultraflame4/bytesblog/models"
	"github.com/ultraflame4/bytesblog/routers/meta"
	"gorm.io/gorm"
)

func (*ApiHandler) CurrentUser(c *gin.Context) {
	user, s := meta.GetSession(c)

	c.JSON(200, gin.H{
		"status":           "success",
		"message":          "OK",
		"user":             user,
		"session_elevated": s.IsElevated(),
	})
}
func (a *ApiHandler) UpdateUsername(c *gin.Context) {
	user, _ := meta.GetSession(c)
	new_username, ok := c.GetPostForm("username")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "field: username is required",
		})
		return
	}
	user.Username = new_username
	if err := a.db.WithContext(c.Request.Context()).Save(user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to save user",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
	})
}

func (a *ApiHandler) Logout(c *gin.Context) {
	_, session := meta.GetSession(c)

	c.SetCookie("session_token", "", -1, "/", "", false, true)
	app.DeleteSession(a.db, c.Request.Context(), session.ID)

	c.Redirect(http.StatusFound, "/")
}

func (a *ApiHandler) LogoutEverywhere(c *gin.Context) {
	user, _ := meta.GetSession(c)

	app.RevokeAllSessions(a.db, c.Request.Context(), user.ID)
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func (a *ApiHandler) ListPosts(c *gin.Context) {

	ctx := c.Request.Context()
	filterQuery := strings.ToLower(c.Query("filter"))
	searchQuery := strings.ToLower(c.Query("q"))
	pageQuery := meta.GetPagination(c)
	user, _ := meta.GetSession(c)
	db := a.db

	sql := gorm.G[models.Post](db).
		Where("user_id = ?", user.ID)

	total_count, err := sql.Count(ctx, "ID")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	sqlQuery := meta.ApplySqlPagination(pageQuery, sql);

	if filterQuery == "published" {
		sqlQuery = sqlQuery.Where("allow_public = ?", true)
	} else if filterQuery == "private" {
		sqlQuery = sqlQuery.Where("allow_public = ?", false)
	}

	if searchQuery != "" {
		sqlQuery = sqlQuery.Where("title LIKE ?", "%"+searchQuery+"%")
	}

	posts, err := sqlQuery.Find(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type Response struct {
		meta.PageMeta
		Items []meta.PostSummary `json:"items"`
	}

	summaries := make([]meta.PostSummary, 0, len(posts))
	for _, post := range posts {
		summaries = append(summaries, meta.PostSummary{
			Model:       post.Model,
			Title:       post.Title,
			Summary:     post.Summary,
			AllowPublic: post.AllowPublic,
			HID:         post.HID,
		})
	}

	c.JSON(http.StatusOK, Response{
		PageMeta: meta.CreatePageMeta(pageQuery, int(total_count)),
		Items:    summaries,
	})

}

// Get a private post by id. Note that this also works for public posts.
func (a *ApiHandler) GetPost(c *gin.Context) {
	user, _ := meta.GetSession(c)
	post_hid := c.Param("hid")
	db := a.db

	var post models.Post
	post, ok := app.FindPostByHID(post_hid, user.ID, db, c.Request.Context())
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"post": post,
		})
		return
	}

	c.AbortWithStatusJSON(404, gin.H{
		"error": "Cannot find post with hid:" + post_hid,
	})

}

// Also used to make post private / unpublish post
func (a *ApiHandler) PublishPost(c *gin.Context) {
	user, _ := meta.GetSession(c)
	post_hid := c.Param("hid")
	make_public := c.Query("public") == "true"

	if post, ok := app.FindPostByHID(post_hid, user.ID, a.db, c.Request.Context()); ok{
		post.AllowPublic = make_public
		a.db.Save(&post)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
		return
	}

	c.AbortWithStatusJSON(404, gin.H{
		"status":  "error",
		"message": "Cannot find post with id:" + post_hid,
	})

}

func (a *ApiHandler) DeletePost(c *gin.Context) {
	user, _ := meta.GetSession(c)
	db := a.db
	post_hid := c.Param("hid")

	err := db.Transaction(func(tx *gorm.DB) error {
		post, ok := app.FindPostByHID(post_hid, user.ID, tx, c.Request.Context())
		if !ok {
			return errors.New("post not found")
		}

		return tx.Delete(&post).Error
	})

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
		return
	}

	c.AbortWithStatusJSON(404, gin.H{
		"status":  "error",
		"message": "Cannot find post with id:" + post_hid,
	})
}

func (a *ApiHandler) SavePost(c *gin.Context) {
	db := a.db
	user, _ := meta.GetSession(c)

	post_hid := c.Query("hid")
	content := c.PostForm("content")

	var post models.Post

	if post_hid != "" {
		// Check if post exists and belongs to the user
		existingPost, ok := app.FindPostByHID(post_hid, user.ID, db, c.Request.Context());
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		}

		if savedPost, err := app.SavePost(&existingPost, content, db, c.Request.Context()); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save post: " + err.Error(),
			})
			return
		} else {
			post = *savedPost
		}
	} else {
		var err error
		post, err = app.CreatePost(content, user.ID, db, c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create post: " + err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"post":   post,
	})
}

func (a *ApiHandler) ListConnectedOidcProviders(c *gin.Context) {
	user, _ := meta.GetSession(c)
	existing, err := gorm.G[models.UserIdentity](a.db).Where("user_id = ?", user.ID).Find(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Failed to retrieve connected OIDC providers: " + err.Error(),
		})
		return
	}

	providers := []string{}
	for _, value := range existing {
		provider_name := strings.ToLower(app.FindProviderName(value.Issuer))
		providers = append(providers, provider_name)
	}


	c.JSON(200, gin.H{
		"status":    "success",
		"providers": providers,
	})
}

func (a *ApiHandler) UserDisconnectOidcProvider(c *gin.Context) {
	user, s := meta.GetSession(c)
	provider := c.Param("oidc")
	// Check for elevated session
	if !s.IsElevated() {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Elevated session required",
		})
		return
	}

	conf, ok := app.GetOidcProvider(provider)
	if !ok {
		c.AbortWithStatusJSON(404, gin.H{
			"error": "OIDC provider " + provider + " not found",
		})
		return
	}

	err := conf.RemoveFromUser(a.db, c.Request.Context(), user.ID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Failed to remove OIDC provider: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
	})
}
