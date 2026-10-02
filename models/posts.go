package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	// Id used to retrieve the post to prevent ID iteration attacks
	HID string `json:"hid" gorm:"uniqueIndex"`
	// Title of post
	Title string
	// Generated summary of post. (First few lines of the post)
	Summary string
	// Contents of the post in markdown. This also includes the post title and raw metadata
	Content string
	// Where the metadata of the markdown (frontmatter) ends. This is needed as Content stores the entire markdown. including the frontmatter metadata.
	// This will also include the final newline.
	ContentOffset uint
	// Author id of the post
	UserID uint
	// Author of the post
	User User
	// Can the post be viewed by the public
	AllowPublic bool
}
