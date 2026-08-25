package handlers

import (
	"github.com/teenyblog/internal/app"
	"gorm.io/gorm"
)

type ApiHandler struct {
	db *gorm.DB
	config *app.Config
}

func NewApiHandler(db *gorm.DB) ApiHandler{
	return ApiHandler {
		db: db,
		config: app.GetConfig(),
	}
}
