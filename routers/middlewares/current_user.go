package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/internal/app"
	"gorm.io/gorm"
)

// This middleware sets checks for the current user and sets it on the context
//
// If the user is not authenticated, it responds with a 401 Unauthorized status
func CurrentUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("session_token")
		if err != nil || sid == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		if session := app.GetSession(db, c.Request.Context(), sid); session != nil {
			s := *session
			c.Set("user_session", s)
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is invalid or has expired"})
		
	}
}