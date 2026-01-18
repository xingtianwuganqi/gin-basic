package router

import (
	"gin-basic/handler"
	"gin-basic/internal"
	"gin-basic/middleware"
	"gin-basic/settings"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if settings.Conf.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	bundle := internal.ReloadLocalBundle()
	r := gin.Default()
	
	// 添加健康检查路由
	r.GET("/health", handler.HealthCheck)
	
	// 添加国际化中间件
	r.Use(middleware.LocaleMiddleware(bundle))
	
	// 注册用户相关路由
	RegisterUserRouter(r)
	
	return r
}