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


// Returns the current user and session (models) from the context.
func currentSession(c *gin.Context) (models.User, models.UserSession) {
	session := c.MustGet("user_session").(models.UserSession);
	
	return session.User, session
}

