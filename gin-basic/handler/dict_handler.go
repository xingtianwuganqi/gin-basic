package handler

import (
	"errors"
	"strconv"

	"gin-basic/models"
	"gin-basic/response"
	"gin-basic/service"

	"github.com/gin-gonic/gin"
)

func CreateSysDictType(c *gin.Context) {
	var req models.CreateSysDictTypeRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.CreateSysDictType(req)
	if err != nil {
		handleSysDictError(c, err, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, result)
}

func UpdateSysDictType(c *gin.Context) {
	var req models.UpdateSysDictTypeRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.UpdateSysDictType(req)
	if err != nil {
		handleSysDictError(c, err, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	response.Success(c, result)
}

func ListSysDictType(c *gin.Context) {
	page := models.Page{}
	if err := c.ShouldBind(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.ListSysDictTypes(page)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, result)
}

func UpdateSysDictTypeStatus(c *gin.Context) {
	dictType := c.Param("dict_type")
	statusValue, err := strconv.Atoi(c.Param("status"))
	if err != nil || (statusValue != 0 && statusValue != 1) {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.UpdateSysDictTypeStatus(dictType, int8(statusValue))
	if err != nil {
		handleSysDictError(c, err, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	response.Success(c, result)
}

func CreateSysDictItem(c *gin.Context) {
	var req models.CreateSysDictItemRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if req.Status != nil && *req.Status != 0 && *req.Status != 1 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if req.IsDefault != nil && *req.IsDefault != 0 && *req.IsDefault != 1 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.CreateSysDictItem(req)
	if err != nil {
		handleSysDictError(c, err, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, result)
}

func UpdateSysDictItem(c *gin.Context) {
	var req models.UpdateSysDictItemRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if req.Status != nil && *req.Status != 0 && *req.Status != 1 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if req.IsDefault != nil && *req.IsDefault != 0 && *req.IsDefault != 1 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.UpdateSysDictItem(req)
	if err != nil {
		handleSysDictError(c, err, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	response.Success(c, result)
}

func ListSysDictItem(c *gin.Context) {
	page := models.Page{}
	if err := c.ShouldBind(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.ListSysDictItems(page, c.Query("dict_type"))
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, result)
}

func UpdateSysDictItemStatus(c *gin.Context) {
	itemValue := c.Param("item_value")
	statusValue, err := strconv.Atoi(c.Param("status"))
	if err != nil || (statusValue != 0 && statusValue != 1) {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.UpdateSysDictItemStatus(itemValue, int8(statusValue), c.Query("dict_type"))
	if err != nil {
		handleSysDictError(c, err, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	response.Success(c, result)
}

func ListClientSysDictItem(c *gin.Context) {
	dictType := c.Param("dict_type")
	if dictType == "" {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.ListClientSysDictItems(dictType)
	if err != nil {
		handleSysDictError(c, err, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, result)
}

func handleSysDictError(c *gin.Context, err error, defaultCode uint, defaultMsg string) {
	switch {
	case errors.Is(err, service.ErrSysDictTypeExists), errors.Is(err, service.ErrSysDictItemExists), errors.Is(err, service.ErrSysDictItemAmbiguous):
		response.Fail(c, response.ApiCode.ConflictErr, response.ApiMsg.ConflictErr)
	case errors.Is(err, service.ErrSysDictTypeNotFound), errors.Is(err, service.ErrSysDictItemNotFound):
		response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
	case errors.Is(err, service.ErrSysDictTypeDisabled):
		response.Fail(c, response.ApiCode.RejectErr, response.ApiMsg.RejectErr)
	default:
		response.Fail(c, defaultCode, defaultMsg)
	}
}
