package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/internal/app"
	"github.com/teenyblog/routers/handlers"
	"github.com/teenyblog/routers/middlewares"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	app.RegisterOIDCProviders(context.Background())
	r := gin.Default()
	handlers := handlers.NewApiHandler(db)
	_api := r.Group("/api")
	{
		_api.GET("/oidc/providers", handlers.ListLoginProviders)
		_api.GET("/oidc/:oidc/login", handlers.StartOauthLoginFlow)
		_api.GET("/oidc/:oidc/callback", handlers.AcceptOAuthCallback)
	}
	{
		_me := _api.Group("/me")
		_me.Use(middlewares.CurrentUser(db))
		_me.GET("/", handlers.CurrentUser)
		_me.GET("/logout", handlers.LogoutCurrentUser)
		_me.POST("/posts/new", handlers.CreatePost)
	}

	return r
}
