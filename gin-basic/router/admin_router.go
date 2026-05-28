package router

import (
	"gin-basic/handler"
	"gin-basic/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(r *gin.Engine) {
	adminRouter := r.Group("/v1/admin", middleware.AdminOnly())
	{
		adminRouter.GET("/user/events", handler.ListUserEvents)
		adminRouter.GET("/subscription", handler.AdminSubscriptionList)
		adminRouter.GET("/transaction", handler.AdminTransactionList)
		adminRouter.GET("/sub-usage", handler.AdminSubUsageList)
		adminRouter.GET("/credits", handler.AdminCredits)
		adminRouter.POST("/app_version", handler.AdminUpsertAppVersion)
		adminRouter.GET("/app_version", handler.AdminAppVersionList)
	}
}
