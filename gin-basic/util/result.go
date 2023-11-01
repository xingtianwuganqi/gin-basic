package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 成功
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := BaseResponse{}
	res.Code = http.StatusOK
	res.Msg = "success"
	res.Data = data

	c.JSON(http.StatusOK, res)
}

// Fail 出错
func Fail(c *gin.Context, code uint, msg string) {
	res := BaseResponse{}
	res.Code = code
	res.Msg = msg
	res.Data = gin.H{}
	c.JSON(http.StatusOK, res)
}
