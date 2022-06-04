package route

import (
	"goframework/cmd/http/route/api"
	"goframework/cmd/http/route/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {

	// 1. 注册全局中间件
	router.Use(
		middleware.LimitBodySize,
		middleware.Log,
		gin.Recovery(),
		middleware.Cors,
	)

	// 2. 注册 swagger 文档
	api.RegisterSwagger(router)

	// 3. 注册路由
	api.RegisterV1(router)
}
