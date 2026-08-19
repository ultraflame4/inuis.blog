package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/routers/handlers"
	"gorm.io/gorm"
)
func InitRouter(db *gorm.DB) *gin.Engine {
	handlers.RegisterOIDCProviders(context.Background())
	r := gin.Default()

	_api := r.Group("/api")
	{
		apiroutes := handlers.NewApiHandler(db)
		_api.GET("/login/:oidc", apiroutes.LoginFlow)
		_api.GET("/oauth2/:oidc/callback", apiroutes.AcceptOAuthCallback)
		_api.POST("/user/login", apiroutes.LoginUser)
	}

	return r
}
