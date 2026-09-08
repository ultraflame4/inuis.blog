package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/bytesblog/internal/app"
	"github.com/ultraflame4/bytesblog/models"
	"github.com/ultraflame4/bytesblog/routers/meta"
	"gorm.io/gorm"
)

func (a *ApiHandler) GetUserPublicInfo(c *gin.Context) {
	db := a.db
	user_hid := c.Param("user_hid")
	user, found := app.FindUserByHid(user_hid, db, c.Request.Context())
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
	return
}

func (a *ApiHandler) ListUserPublicPosts(c *gin.Context) {
	db := a.db
	user_hid := c.Param("user_hid")
	user, found := app.FindUserByHid(user_hid, db, c.Request.Context())
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	pageQuery := meta.GetPagination(c)
	ctx := c.Request.Context()

	sql := gorm.G[models.Post](db).Where("user_id = ? AND allow_public = true", user.ID)
	sql = sql.Order("created_at desc")
	sql = meta.ApplySqlPagination(pageQuery, sql)
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
			Model:       post.Model,
			Title:       post.Title,
			Summary:     post.Summary,
			AllowPublic: post.AllowPublic,
		})
	}

	c.JSON(http.StatusOK, Response{
		PageMeta: meta.CreatePageMeta(pageQuery, int(total_count)),
		Posts:    summaries,
	})

}
