package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	// Title of post
	Title   string
	// Generated summary of post. (First few lines of the post)
	Summary string
	// Contents of the post in markdown. This also includes the post title.
	Content string
	// Author id of the post
	UserID uint
	// Author of the post
	User User
}
