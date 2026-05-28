package main

import (
	"fmt"
	"gin-basic/db"
	"gin-basic/logger"
	"gin-basic/router"
	"gin-basic/settings"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	settings.InitConf()

	// 初始化日志
	logger.InitLogger()
	defer logger.Logger.Sync()

	// 初始化数据库
	db.LinkDataBase()

	// 设置路由
	r := router.SetupRouter()

	// 启动服务
	port := settings.Conf.App.Port
	env := settings.Conf.App.Env

	logger.Logger.Info("Starting server",
		zap.Int("port", port),
		zap.String("env", env))

	err := r.Run(":" + fmt.Sprint(port))
	if err != nil {
		logger.Logger.Error("Failed to start server",
			zap.Error(err))
	}
}
