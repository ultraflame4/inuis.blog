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
	_api.Use(middlewares.PaginationQuery())
	{
		_api.GET("/oidc/providers", handlers.ListLoginProviders)
		_api.GET("/oidc/:oidc/login", handlers.StartOauthLoginFlow)
		_api.GET("/oidc/:oidc/callback", handlers.AcceptOAuthCallback)
	}
	{
		_me := _api.Group("/me")
		_me.Use(middlewares.CurrentUser(db))
		_me.GET("/", handlers.CurrentUser)
		_me.GET("/logout", handlers.Logout)
		_me.GET("/oidc-providers", handlers.ListConnectedOidcProviders)
		_me.GET("/oidc/:oidc/disconnect", handlers.UserDisconnectOidcProvider)
		_me.POST("/posts/save", handlers.SavePost)
		_me.GET("/posts/:id", handlers.GetPost)
		_me.DELETE("/posts/:id", handlers.DeletePost)
		_me.GET("/posts", handlers.ListPosts)
		_me.POST("/posts/:id/publish", handlers.PublishPost)
	}
	{
		_posts := _api.Group("/posts")
		_posts.GET("/:id", handlers.GetPublicPost)
	}
	{
		_user := _api.Group("/user")
		_user.GET("/:user_id/posts", handlers.ListUserPublicPosts)
	}

	return r
}
