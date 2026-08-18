package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/routers/handlers"
)
func InitRouter() *gin.Engine {
	handlers.RegisterOIDCProviders(context.Background())
	r := gin.Default()

	_api := r.Group("/api")
	{
		handlers.RegisterLoginHandlers(_api)
		_api.POST("/user/login", handlers.LoginUser)
	}

	return r
}
