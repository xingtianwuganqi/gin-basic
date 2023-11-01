package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/handler"
)

func RegisterUserRouter(r *gin.Engine) {
	userRouter := r.Group("/v1/user")
	{
		userRouter.POST("/register", handler.UserRegister)
		userRouter.POST("/login", handler.UserPhoneLogin)
	}

}
