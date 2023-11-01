package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/settings"
)

func RegisterRouter() *gin.Engine {
	if settings.Conf.App.Debug == true {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	RegisterUserRouter(r)
	return r
}
