package meta

import (
	"github.com/gin-gonic/gin"
	"github.com/teenyblog/models"
)

// Returns the current user and session (models) from the context.
func GetSession(c *gin.Context) (models.User, models.UserSession) {
	session := c.MustGet("user_session").(models.UserSession);
	
	return session.User, session
}

