package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/internal/app"
	"github.com/teenyblog/models"
)

// Returns the current user and session (models) from the context.
func currentSession(c *gin.Context) (models.User, models.UserSession) {
	session := c.MustGet("user_session").(models.UserSession);
	
	return session.User, session
}


func (*ApiRoute) CurrentUser(c *gin.Context){
	user, _ := currentSession(c)
	
	c.JSON(200, gin.H{
		"status": 200,
		"message": "OK",
		"user": user,
	})
}

func (a *ApiRoute) LogoutCurrentUser(c *gin.Context) {
	_, session := currentSession(c)
	
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	app.DeleteSession(a.db, c.Request.Context(), session.ID);
	
	
	c.Redirect(http.StatusFound, "/")
}