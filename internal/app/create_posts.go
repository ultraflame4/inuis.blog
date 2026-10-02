package app

import (
	"context"
	"strings"

	"github.com/ultraflame4/inuis/models"
	"gorm.io/gorm"
)

type PostMetadata struct {
	title   string
	summary string
}

// Generates an assortment of metadata from a post's content (markdown / frontmatter)
func GeneratePostMetadata(content string) PostMetadata {
	meta := PostMetadata{}
	body := ""

	// Remove the first line of the content. (Treated as title)
	if i := strings.IndexByte(content, '\n'); i >= 0 {
		meta.title = content[:i]
		body = content[i+1:]
	} else {
		// If i == -1. Means no newlines. The entire content is the first line.
		meta.title = strings.Clone(content)
		// Also means that body has no content
		body = ""
	}
	meta.title = strings.TrimPrefix(meta.title, "#")
	meta.title = strings.TrimSpace(meta.title)

	// Generate summary from body
	summary_len := 100
	summary_body := strings.ReplaceAll(body, "\n", "")
	meta.summary = strings.TrimSpace(summary_body[0:min(len(summary_body), summary_len)])
	return PostMetadata{}
}

func CreatePost(content string, userID uint, db *gorm.DB, ctx context.Context) (models.Post, error) {
	hid := GenerateHID()
	meta := GeneratePostMetadata(content)
	post := models.Post{
		Title:   meta.title,
		Summary: meta.summary,
		Content: content,
		UserID:  userID,
		HID:     hid,
	}
	err := db.Create(&post).Error
	return post, err
}

// Finds a post by post_hid
//
// This will return ok as long as the user is authorised to read the post.
//
// Params:
//   - `hid` - Hid of post
//   - `userId` - User finding the post
//
// Returns: (post_data, ok)
func FindPostByHID(hid string, userID uint, db *gorm.DB, ctx context.Context) (models.Post, bool) {
	post, err := gorm.G[models.Post](db).Preload("User", nil).Where("h_id = ?", hid).First(ctx)
	return post, err == nil && (post.UserID == userID || post.AllowPublic)
}

func FindPublicPostByHID(hid string, db *gorm.DB, ctx context.Context) (models.Post, bool) {
	post, err := gorm.G[models.Post](db).Preload("User", nil).Where("h_id = ?", hid).First(ctx)
	return post, err == nil && post.AllowPublic
}

func SavePost(existingPost *models.Post, content string, db *gorm.DB, ctx context.Context) (*models.Post, error) {
	if existingPost == nil {
		return nil, nil
	}
	meta := GeneratePostMetadata(content)
	existingPost.Title = meta.title
	existingPost.Summary = meta.summary
	existingPost.Content = content
	err := db.Save(existingPost).Error
	return existingPost, err
}
