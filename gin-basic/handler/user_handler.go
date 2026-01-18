package handler

import (
	"gin-basic/models"
	"gin-basic/response"
	"gin-basic/service"
	"gin-basic/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	logger.Logger.Info("Received user registration request", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid parameters for user registration", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.RegisterUser(req)
	if err != nil {
		logger.Logger.Error("User registration failed", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}

	logger.Logger.Info("User registration successful", zap.String("email", req.Email))
	response.Success(c, gin.H{
		"message": "User registered successfully",
	})
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	logger.Logger.Info("Received user login request", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid parameters for user login", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.LoginUser(req)
	if err != nil {
		logger.Logger.Error("User login failed", zap.Error(err))
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}

	logger.Logger.Info("User login successful", zap.String("email", req.Email))
	response.Success(c, result)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		logger.Logger.Warn("Email parameter is required")
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	logger.Logger.Info("Received user info retrieval request", 
		zap.String("email", email),
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))

	result, err := service.GetUserInfo(email)
	if err != nil {
		logger.Logger.Error("Failed to retrieve user info", zap.Error(err))
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}

	logger.Logger.Info("User info retrieved successfully", zap.String("email", email))
	response.Success(c, result)
}

// HealthCheck 健康检查端点
func HealthCheck(c *gin.Context) {
	logger.Logger.Info("Health check endpoint called", zap.String("clientIP", c.ClientIP()))

	healthInfo := map[string]interface{}{
		"status":    "healthy",
		"timestamp": "2026-01-18T14:41:06+08:00",
		"service":   "gin-basic",
	}

	logger.Logger.Info("Health check completed", zap.String("status", healthInfo["status"].(string)))

	c.JSON(200, healthInfo)
}