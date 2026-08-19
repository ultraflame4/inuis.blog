package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	// Title of post
	Title   string
	// Generated summary of post. (First few lines of the post)
	Summary string
	// File location of the post content. (File is assumed to be markdown)
	FileLocation string

	// Author id of the post
	UserID uint
	// Author of the post
	User User
}
