package handlers

import "gorm.io/gorm"

type ApiRoute struct {
	db *gorm.DB
}

func NewApiHandler(db *gorm.DB) ApiRoute{
	return ApiRoute {
		db: db,
	}
}
