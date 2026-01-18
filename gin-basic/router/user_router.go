package router

import (
	"gin-basic/handler"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	userRouter := r.Group("/v1/user")
	{
		userRouter.POST("/register", handler.UserRegister)
		userRouter.POST("/login", handler.UserLogin)
		userRouter.GET("/info", handler.GetUserInfo)
	}
}