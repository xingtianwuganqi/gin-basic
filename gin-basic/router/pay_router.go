package router

import (
	"gin-basic/handler"
	"gin-basic/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPayRouter(r *gin.Engine) {
	iapRouter := r.Group("/v1/iap")
	{
		iapRouter.POST("/verify", middleware.JWTTokenMiddleware(), handler.VerisfyAppleIAP)
		iapRouter.POST("/notify", handler.NotifyAppleIAP)
	}
}
