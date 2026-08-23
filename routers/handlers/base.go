package handlers

import "gorm.io/gorm"

type ApiHandler struct {
	db *gorm.DB
}

func NewApiHandler(db *gorm.DB) ApiHandler{
	return ApiHandler {
		db: db,
	}
}
