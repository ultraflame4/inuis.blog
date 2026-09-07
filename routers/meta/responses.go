package meta

import "gorm.io/gorm"
type PageMeta struct {
	// Page number
	Page       int `json:"page"`
	// Items per page
  	PerPage    int `json:"per_page"`
   	// TotalItems number of items (includes other pages aswell)
  	TotalItems      int `json:"total_items"`
   	// Total number of pages
  	TotalPages int `json:"total_pages"`

}

func CreatePageMeta(q PageQuery, total_count int) PageMeta {
	total_pages := total_count / q.Limit
	// Handle edge case where total_count is not divisible by q.Limit
	// Increment total_pages to account for the remaining items
	if total_count%q.Limit != 0 {
		total_pages++
	}
	return PageMeta{
		Page:  q.Page,
		PerPage: q.Limit,
		TotalItems: total_count,
		TotalPages: max(1, total_pages), // Always a minimum of 1 page
	}
}



// Summarised version of a post.
type PostSummary struct {
	gorm.Model
	Title string
	Summary string
	AllowPublic bool
	HID string `json:"hid"`
}
