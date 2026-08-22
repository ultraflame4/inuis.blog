package handlers

import "github.com/gin-gonic/gin"

func (*ApiRoute) LoginUser(c *gin.Context){
	c.JSON(500, gin.H{
		"status": 500,
		"message": "Not implemented",
	})
}

