package meta

import "gorm.io/gorm"
type PageMeta struct {
	// Page number
	Page       int `json:"page,omitempty"`
	// Items per page
  	PerPage    int `json:"per_page,omitempty"`
   	// TotalItems number of items (includes other pages aswell)
  	TotalItems      int `json:"total_items,omitempty"`
   	// Total number of pages
  	TotalPages int `json:"total_pages,omitempty"`	
  	
}

func CreatePageMeta(q PageQuery, total_count int) PageMeta {
	return PageMeta{
		Page:  q.Page,
		PerPage: q.Limit,
		TotalItems: total_count,
		TotalPages: total_count / q.Limit + 1,
	}
}



// Summarised version of a post.
type PostSummary struct {
	gorm.Model
	Title string
	Summary string
	AllowPublic bool
}
