package router

import (
	"gin-basic/handler"
	"gin-basic/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSysDictRouter(r *gin.Engine) {
	adminRouter := r.Group("/v1/admin/dict", middleware.AdminOnly())
	{
		adminRouter.POST("/type", handler.CreateSysDictType)
		adminRouter.PUT("/type", handler.UpdateSysDictType)
		adminRouter.GET("/type", handler.ListSysDictType)
		adminRouter.GET("/type/:dict_type/:status", handler.UpdateSysDictTypeStatus)

		adminRouter.POST("/item", handler.CreateSysDictItem)
		adminRouter.PUT("/item", handler.UpdateSysDictItem)
		adminRouter.GET("/item", handler.ListSysDictItem)
		adminRouter.GET("/item/:item_value/:status", handler.UpdateSysDictItemStatus)
	}

	dictRouter := r.Group("/v1/dict")
	{
		dictRouter.GET("/item/:dict_type", handler.ListClientSysDictItem)
	}
}
