package app

import (
	"context"
	"log/slog"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/ultraflame4/inuis/models"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

type PostMetadata struct {
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	Tags          []string `json:"tags"`
	ContentOffset uint   `json:"content_offset"`
}

var frontmatterFormat = frontmatter.NewFormat("---", "---", yaml.Unmarshal)

// Generates an assortment of metadata from a post's content (markdown / frontmatter)
func GeneratePostMetadata(content string) PostMetadata {
	md_content := strings.TrimSpace(content)
	// Read frontmatter data
	meta := PostMetadata{}
	reader := strings.NewReader(md_content)
	rest, err := frontmatter.Parse(reader, &meta, frontmatterFormat)
	if err != nil {
		slog.Error("Failed to process frontmatter data!:", "error", err)
	} else {
		meta.ContentOffset = uint(len(md_content) - len(rest))
		md_content = string(rest)
	}

	// Trim whitespaces again, there may be whitespaces between the frontmatter and markdown content
	md_content = strings.TrimSpace(md_content)

	// Content body without metadata and title
	body := ""

	// Note: The title can be defined in the frontmatter. But that only acts as the title stored in the db.
	//       Not as the main title shown in the content / markdown.
	//       The stored title acts as a search entry rather, and name of the post listing. Which is why we still need to
	//       extract / remove the content title from the md_content
	// Title extraction
	{
		// Extracted title. Will be used as fallback title entry
		fb_title := ""
		// Remove the title
		if title_len := strings.IndexByte(md_content, '\n'); title_len >= 0 {
			fb_title = md_content[:title_len]
			// Remove title line from body
			body = md_content[title_len:]
		} else {
			// If i == -1. Means no newlines. The entire content is the title
			fb_title = strings.Clone(md_content)
			// Also means that body has no content
			body = ""
		}
		// Use as entry if not found in frontmatter
		if len(meta.Title) == 0 {
			// Use [1:] to remove the '#' prefix
			meta.Title = strings.TrimPrefix(fb_title, "#")
		}
	}

	meta.Title = strings.TrimSpace(meta.Title)

	// Fallback summary generation, if summary not found in frontmatter
	if len(meta.Summary) == 0 {
		// Generate summary from body
		summary_len := 100
		summary_body := strings.ReplaceAll(body, "\n", "")
		meta.Summary = strings.TrimSpace(summary_body[0:min(len(summary_body), summary_len)])
	}

	return meta
}

func CreatePost(content string, userID uint, db *gorm.DB, ctx context.Context) (models.Post, error) {
	hid := GenerateHID()
	meta := GeneratePostMetadata(content)
	post := models.Post{
		Title:   meta.Title,
		Summary: meta.Summary,
		Content: content,
		ContentOffset: meta.ContentOffset,
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
	existingPost.Title = meta.Title
	existingPost.Summary = meta.Summary
	existingPost.ContentOffset = meta.ContentOffset
	existingPost.Content = content
	err := db.Save(existingPost).Error
	return existingPost, err
}
