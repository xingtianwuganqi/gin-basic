package router

import (
	"gin-basic/handler"
	"gin-basic/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	userRouter := r.Group("/v1/user")
	{
		userRouter.POST("/event", middleware.OptionalJWTMiddleware(), handler.SubmitUserEvent)
		userRouter.POST("/register", handler.UserRegister)
		userRouter.POST("/login", handler.UserLogin)
		userRouter.POST("/device-token", middleware.JWTTokenMiddleware(), handler.UpdateUserDeviceToken)
		userRouter.GET("/info", handler.GetUserInfo)
		userRouter.GET("/credits", middleware.JWTTokenMiddleware(), handler.GetUserCredits)
		userRouter.GET("/credit-list", middleware.JWTTokenMiddleware(), handler.GetUserCreditList)
	}

	appRouter := r.Group("/v1/app")
	{
		appRouter.GET("/version", handler.GetAppVersionInfo)
	}
}
