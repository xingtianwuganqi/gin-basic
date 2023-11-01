package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"pet-project/db"
	"pet-project/models"
	"pet-project/util"
	"reflect"
	"strconv"
)

// PetInfoCreate 提交宠物详情
func PetInfoCreate(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var petInfo models.PetInfo
	if err := c.ShouldBind(&petInfo); err != nil {
		log.Println(err)
		util.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}

	// 如果token的userId和参数的userId不一样，说明不是同一个人
	log.Println("+++++++", reflect.TypeOf(petInfo.UserId), petInfo.UserId, petInfo.PetType)
	log.Println("0000", c.PostForm("pet_type"))
	if petInfo.UserId != userId {
		util.Fail(c, util.ApiCode.QueryError, util.ApiMessage.QueryError)
		return
	}

	result := db.DB.Create(&petInfo)
	if result.Error != nil {
		log.Println(result.Error)
		util.Fail(c, util.ApiCode.CreateErr, util.ApiMessage.CreateErr)
		return
	}
	util.Success(c, nil)
}

// CreatePetActionType 添加宠物行为
func CreatePetActionType(c *gin.Context) {
	var actionModel models.PetActionType
	if err := c.ShouldBind(&actionModel); err != nil {
		util.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}

	result := db.DB.Create(&actionModel)
	if result.Error != nil {
		util.Fail(c, util.ApiCode.CreateErr, util.ApiMessage.CreateErr)
		return
	}
	util.Success(c, nil)
}

// GetPetActionList 获取宠物行为
func GetPetActionList(c *gin.Context) {
	var petActionList []models.PetActionType
	result := db.DB.Model(&models.PetActionType{}).First(&petActionList)
	if result != nil {
		util.Fail(c, util.ApiCode.QueryError, util.ApiMessage.QueryError)
		return
	}
	util.Success(c, petActionList)
}

// GetRecordList 查询记录列表
func GetRecordList(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var pageNum = c.PostForm("pageNum")
	var pageSize = c.Query("pageSize")
	num, err := strconv.Atoi(pageNum)
	if err != nil {
		util.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		util.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	offset := (num - 1) * size
	var recordList []models.RecordList
	result := db.DB.Where("userId=?", userId).Offset(offset).Limit(size).Find(&recordList)
	if result.Error != nil {
		util.Fail(c, util.ApiCode.QueryError, util.ApiMessage.QueryError)
		return
	}
	util.Success(c, recordList)
}
