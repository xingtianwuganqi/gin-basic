package handler

import (
	"gin-basic/logger"
	"gin-basic/models"
	"gin-basic/response"
	"gin-basic/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	if err := c.ShouldBind(&req); err != nil {
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

func UpdateUserDeviceToken(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var req models.DeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid parameters for device token update", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.UpdateDeviceToken(userId, req)
	if err != nil {
		logger.Logger.Error("Update device token failed", zap.Uint("userId", userId), zap.Error(err))
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}

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

func SubmitUserEvent(c *gin.Context) {
	var userId *uint
	if value, ok := c.Get("userId"); ok {
		if id, ok := value.(uint); ok && id > 0 {
			userId = &id
		}
	}

	var req models.UserEventRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Logger.Warn("Invalid parameters for user event", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.CreateUserEvent(userId, req, c.ClientIP())
	if err != nil {
		logger.Logger.Error("Create user event failed", zap.Any("userId", userId), zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}

	response.Success(c, result)
}

func GetUserCredits(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var page models.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 20
	}

	summary, err := service.GetUserCreditSummary(userId)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	logs, total, err := service.GetCreditLogs(userId, page.PageNum, page.PageSize)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, gin.H{
		"total_credits":     summary.TotalCredits,
		"used_credits":      summary.UsedCredits,
		"remaining_credits": summary.TotalCredits - summary.UsedCredits,
		"logs":              logs,
		"log_total":         total,
	})
}

func GetUserCreditList(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var page models.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 20
	}

	list, _, err := service.GetUserPurchaseLogs(userId, page.PageNum, page.PageSize)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, list)
}

func GetAppVersionInfo(c *gin.Context) {
	var query models.AppVersionQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.GetAppVersion(query.Platform)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
			return
		}
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

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
