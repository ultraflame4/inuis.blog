package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/teenyblog/routers/handlers"
)
func InitRouter() *gin.Engine {
	r := gin.Default()

	_api := r.Group("/api")
	{
		_api.POST("/user/login", handlers.LoginUser)
	}

	return r
}
