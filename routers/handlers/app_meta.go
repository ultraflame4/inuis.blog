package handlers

import "github.com/gin-gonic/gin"

func (a *ApiHandler) GetAppMeta(c *gin.Context){
	c.JSON(200, gin.H{
		"brand_name": a.config.BrandName,
		
	})
}