package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/internal/app"
	"github.com/teenyblog/routers/handlers"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	app.RegisterOIDCProviders(context.Background())
	r := gin.Default()
	apiroutes := handlers.NewApiHandler(db)
	_api := r.Group("/api")
	{
		_api.GET("/oidc/providers", apiroutes.ListLoginProviders)
		_api.GET("/oidc/:oidc/login", apiroutes.StartOauthLoginFlow)
		_api.GET("/oidc/:oidc/callback", apiroutes.AcceptOAuthCallback)
	}
	{
		_api.POST("/user/login", apiroutes.LoginUser)
	}

	return r
}
