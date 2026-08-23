package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/teenyblog/models"
	"gorm.io/gorm"
)

type ApiHandler struct {
	db *gorm.DB
}

func NewApiHandler(db *gorm.DB) ApiHandler{
	return ApiHandler {
		db: db,
	}
}


