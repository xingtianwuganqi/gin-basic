package handler

import (
	"errors"
	"gin-basic/db"
	"gin-basic/logger"
	"gin-basic/models"
	"gin-basic/response"
	"gin-basic/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListUserEvents(c *gin.Context) {
	var query models.UserEventQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Logger.Warn("Invalid parameters for user event query", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	list, total, err := service.ListUserEvent(query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidQueryTime) {
			logger.Logger.Warn("Invalid time parameter for user event query", zap.Error(err))
			response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
			return
		}
		logger.Logger.Error("List user events failed", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageNum <= 0 {
		query.PageNum = 1
	}

	response.Success(c, models.PageResponse{
		List:     list,
		Total:    total,
		PageNum:  query.PageNum,
		PageSize: query.PageSize,
	})
}

func AdminSubscriptionList(c *gin.Context) {
	var page models.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	normalizePage(&page)

	var subscriptions []models.Subscription
	if err := db.DB.Model(&models.Subscription{}).
		Limit(page.PageSize).
		Offset((page.PageNum - 1) * page.PageSize).
		Order("created_at desc").
		Find(&subscriptions).Error; err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, subscriptions)
}

func AdminTransactionList(c *gin.Context) {
	var page models.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	normalizePage(&page)

	var transactions []models.Transaction
	if err := db.DB.Model(&models.Transaction{}).
		Limit(page.PageSize).
		Offset((page.PageNum - 1) * page.PageSize).
		Order("created_at desc").
		Find(&transactions).Error; err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, transactions)
}

func AdminSubUsageList(c *gin.Context) {
	var page models.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	normalizePage(&page)

	var usages []models.SubscriptionUsage
	if err := db.DB.Model(&models.SubscriptionUsage{}).
		Limit(page.PageSize).
		Offset((page.PageNum - 1) * page.PageSize).
		Order("created_at desc").
		Find(&usages).Error; err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, usages)
}

type adminCreditItem struct {
	models.UserCredit
	RemainingCredits int `json:"remaining_credits"`
}

func AdminCredits(c *gin.Context) {
	var page models.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	normalizePage(&page)

	credits, total, err := service.GetAllUserCredits(page.PageNum, page.PageSize)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	items := make([]adminCreditItem, 0, len(credits))
	for _, credit := range credits {
		items = append(items, adminCreditItem{
			UserCredit:       credit,
			RemainingCredits: credit.TotalCredits - credit.UsedCredits,
		})
	}

	response.Success(c, gin.H{
		"list":  items,
		"total": total,
	})
}

func AdminUpsertAppVersion(c *gin.Context) {
	var req models.AppVersionRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.UpsertAppVersion(req)
	if err != nil {
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	response.Success(c, result)
}

func AdminAppVersionList(c *gin.Context) {
	result, err := service.ListAppVersion()
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, result)
}

func normalizePage(page *models.Page) {
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 20
	}
}
