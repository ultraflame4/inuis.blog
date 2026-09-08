package routers

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/inuis/internal/app"
	"github.com/ultraflame4/inuis/routers/handlers"
	"github.com/ultraflame4/inuis/routers/middlewares"
	"gorm.io/gorm"
)

func subFS(parent fs.FS, dir string) fs.FS {
	subFS, err := fs.Sub(parent, dir)
	if err != nil {
		panic(err)
	}
	return subFS
}

func InitRouter(db *gorm.DB, frontendFS embed.FS) *gin.Engine {
	app.RegisterOIDCProviders(context.Background())
	r := gin.Default()

	// Serve SPA routes
	{

		// 1. Get the sub-filesystem for the "_app" folder
		publicFS := subFS(frontendFS, "dist")

		// Serve main _app resources dir
		r.StaticFS("_app", http.FS(subFS(frontendFS, "dist/_app")))
		r.StaticFileFS("/robots.txt", "robots.txt", http.FS(publicFS))

		// Serve

		// Fallback route for SPA routing
		r.NoRoute(func(c *gin.Context) {
			c.FileFromFS("spa.html", http.FS(publicFS))
		})

	}

	handlers := handlers.NewApiHandler(db)
	_api := r.Group("/api")
	_api.Use(middlewares.PaginationQuery())

	// These endpoints return information about the application
	{
		_api.GET("/meta", handlers.GetAppMeta)
	}
	{
		_api.GET("/oidc/providers", handlers.ListLoginProviders)
		_api.GET("/oidc/:oidc/login", handlers.StartOauthLoginFlow)
		_api.GET("/oidc/:oidc/callback", handlers.AcceptOAuthCallback)
	}
	{
		_me := _api.Group("/me")
		_me.Use(middlewares.CurrentUser(db))
		_me.GET("/", handlers.CurrentUser)
		_me.PUT("/username", handlers.UpdateUsername)
		_me.GET("/logout", handlers.Logout)
		_me.GET("/logout/all", handlers.LogoutEverywhere)
		_me.GET("/oidc-providers", handlers.ListConnectedOidcProviders)
		_me.GET("/oidc/:oidc/disconnect", handlers.UserDisconnectOidcProvider)
		_me.POST("/posts/save", handlers.SavePost)
		_me.GET("/posts/:hid", handlers.GetPost)
		_me.DELETE("/posts/:hid", handlers.DeletePost)
		_me.GET("/posts", handlers.ListPosts)
		_me.POST("/posts/:hid/publish", handlers.PublishPost)
	}
	{
		_posts := _api.Group("/posts")
		_posts.GET("/:hid", handlers.GetPublicPost)
	}
	{
		_user := _api.Group("/user")
		_user.GET("/:user_hid", handlers.GetUserPublicInfo)
		_user.GET("/:user_hid/posts", handlers.ListUserPublicPosts)
	}

	return r
}
