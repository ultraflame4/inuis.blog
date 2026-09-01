package app

import (
	"context"
	"strings"

	"github.com/ultraflame4/bytesblog/models"
	"gorm.io/gorm"
)

// Generates a title and summary for the post based on the content.
// Returns (title, summary)
func GeneratePostTitleAndSummary(content string) (string, string) {
	title := ""
	body := ""
	// Remove the first line of the content. (Treated as title)
	if i := strings.IndexByte(content, '\n'); i >= 0 {
		title = content[:i]
		body = content[i+1:]
	} else {
		// If i == -1. Means no newlines. The entire content is the first line.
		title = strings.Clone(content)
		// Also means that body has no content
		body = ""
	}
	title = strings.TrimPrefix(title, "#")
	title = strings.TrimSpace(title)

	// Generate summary from body
	summary_len := 100
	summary_body := strings.ReplaceAll(body, "\n", "")
	summary := summary_body[0:min(len(summary_body), summary_len)]
	summary = strings.TrimSpace(summary)
	return title, summary
}

func CreatePost(content string, userID uint, db *gorm.DB, ctx context.Context) (models.Post, error) {
	hid := GenerateHID()
	title, summary := GeneratePostTitleAndSummary(content)
	post := models.Post{
		Title:   title,
		Summary: summary,
		Content: content,
		UserID:  userID,
		HID:     hid,
	}
	err := db.Create(&post).Error
	return post, err
}

func FindPostByHID(hid string, userID uint, db *gorm.DB, ctx context.Context) (models.Post, bool) {
	post, err := gorm.G[models.Post](db).Where("h_id = ?", hid).First(ctx)
	return post, err == nil && post.UserID == userID
}

func FindPublicPostByHID(hid string, db *gorm.DB, ctx context.Context) (models.Post, bool) {
	post, err := gorm.G[models.Post](db).Preload("User", nil).Where("h_id = ?", hid).First(ctx)
	return post, err == nil && post.AllowPublic
}

func SavePost(existingPost *models.Post, content string, db *gorm.DB, ctx context.Context) (*models.Post, error) {
	if existingPost == nil {
		return nil, nil
	}
	title, summary := GeneratePostTitleAndSummary(content)
	existingPost.Title = title
	existingPost.Summary = summary
	existingPost.Content = content
	err := db.Save(existingPost).Error
	return existingPost, err
}
