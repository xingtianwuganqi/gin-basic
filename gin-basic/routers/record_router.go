package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/handler"
	"pet-project/middleware"
)

func RegisterRecordRouter(r *gin.Engine) {
	recordRouter := r.Group("/v1/record")
	{
		recordRouter.GET("/petAction/list", handler.GetPetActionList)
		recordRouter.POST("/pet/create", middleware.JWTTokenMiddleware(), handler.PetInfoCreate)
		recordRouter.POST("/pet/action", handler.CreatePetActionType)
		recordRouter.GET("/list", middleware.JWTTokenMiddleware(), handler.GetRecordList)

	}
}
