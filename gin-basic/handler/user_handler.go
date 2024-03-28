package handler

import (
	"errors"
	"pet-project/db"
	"pet-project/middleware"
	"pet-project/models"
	"pet-project/response"
	"pet-project/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginInfo struct {
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Password string `form:"password" json:"password"`
	Code     string `form:"code" json:"code"`
}

type LoginUserInfo struct {
	UserId uint   `json:"userId"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Token  string `json:"token"`
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
	var login LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	var findUser models.UserInfo
	findResult := db.DB.Where("phone = ?", login.Phone).First(&findUser)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		user := models.UserInfo{
			Phone:    login.Phone,
			Password: login.Password,
		}
		result := db.DB.Create(&user)
		if result.Error != nil {
			response.Fail(c, util.ApiCode.CreateErr, util.ApiMessage.CreateErr)
			return
		}
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
			return
		}
		data := LoginUserInfo{
			UserId: user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		response.Success(c, data)
	} else {
		response.Fail(c, util.ApiCode.UserExistsError, util.ApiMessage.UserExistsError)
	}

}

// UserPhoneLogin 用户登录
func UserPhoneLogin(c *gin.Context) {
	var login LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	var user models.UserInfo
	result := db.DB.Where("phone = ?", login.Phone).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, util.ApiCode.QueryError, "该手机号未注册")
		return
	}
	if user.Password == login.Password {
		// 密码正确, 生成token，登录完成
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
			return
		}
		data := LoginUserInfo{
			UserId: user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		response.Success(c, data)
	} else {
		response.Fail(c, 300, "密码错误")
	}
}
