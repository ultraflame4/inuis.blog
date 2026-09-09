package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/inuis/internal/app"
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

func (a *ApiHandler) ListProviders(c *gin.Context) {
	providers := app.GetProviders()
	c.JSON(http.StatusOK, providers)
}
