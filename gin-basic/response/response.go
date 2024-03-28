package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Response(c *gin.Context, code uint, data interface{}, msg string) {
	if data == nil {
		data = gin.H{}
	}
	res := BaseResponse{}
	res.Code = http.StatusOK
	res.Msg = msg
	res.Data = data

	c.JSON(http.StatusOK, res)
}

// Success 成功
func Success(c *gin.Context, data interface{}) {
	Response(c, 200, data, "success")
}

// Fail 出错
func Fail(c *gin.Context, code uint, msg string) {
	Response(c, code, gin.H{}, msg)
}
