package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/inuis/internal/app"
)

func (a *ApiHandler) GetAppMeta(c *gin.Context){
	// If there are no admin users, require an invite code
	adminless := app.CountUsersWithRole("admin", a.db, c.Request.Context()) == 0
	
	c.JSON(200, gin.H{
		"brand_name": a.config.BrandName,
		"copyright_name": a.config.CopyrightName,
		"requires_invite_code": a.config.RequireInviteCode || adminless,
	})
}