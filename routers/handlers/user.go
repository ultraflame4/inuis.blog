package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/inuis/internal/app"
	"github.com/ultraflame4/inuis/models"
	"github.com/ultraflame4/inuis/routers/meta"
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

	filterQuery := strings.ToLower(c.Query("filter"))
	searchQuery := strings.ToLower(c.Query("q"))
	pageQuery := meta.GetPagination(c)
	ctx := c.Request.Context()

	sqlQuery := gorm.G[models.Post](db).Where("user_id = ? AND allow_public = true", user.ID)

	sqlQuery = sqlQuery.Order("created_at desc")
	sqlQuery = meta.ApplySqlPagination(pageQuery, sqlQuery)
	if filterQuery == "published" {
		sqlQuery = sqlQuery.Where("allow_public = ?", true)
	} else if filterQuery == "private" {
		sqlQuery = sqlQuery.Where("allow_public = ?", false)
	}

	if searchQuery != "" {
		sqlQuery = sqlQuery.Where("title LIKE ?", "%"+searchQuery+"%")
	}

	
	total_count, err := sqlQuery.Count(ctx, "ID")
	posts, err := sqlQuery.Find(ctx)
	
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
		summaries = append(summaries, meta.CreatePostSummary(post))
	}

	c.JSON(http.StatusOK, Response{
		PageMeta: meta.CreatePageMeta(pageQuery, int(total_count)),
		Posts:    summaries,
	})

}
