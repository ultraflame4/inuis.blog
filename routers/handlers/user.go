package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/teenyblog/models"
)

func (*ApiRoute) LoginUser(c *gin.Context){
	c.JSON(500, gin.H{
		"status": 500,
		"message": "Not implemented",
	})
}


// Returns the current user from the context.
func currentUser(c *gin.Context) models.User {
	return c.MustGet("user").(models.User)
}

func (*ApiRoute) CurrentUser(c *gin.Context){
	user := currentUser(c)
	
	c.JSON(200, gin.H{
		"status": 200,
		"message": "OK",
		"user": user,
	})
}