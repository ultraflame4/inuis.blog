package handlers

import (
	"github.com/gin-gonic/gin"
)

func (a *ApiHandler) GetAppMeta(c *gin.Context) {
	// If there are no admin users, require an invite code
	adminless := a.config.IsAdminless(a.db, c.Request.Context())

	c.JSON(200, gin.H{
		"brand_name":          a.config.BrandName,
		"copyright_name":      a.config.CopyrightName,
		"require_invite_code": a.config.RequireInviteCode || adminless || !a.config.ALLOW_SIGNUP,
		"allow_signup":        a.config.ALLOW_SIGNUP,
	})
}
