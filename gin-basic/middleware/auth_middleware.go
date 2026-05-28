package middleware

import (
	"errors"
	"gin-basic/db"
	"gin-basic/internal"
	"gin-basic/logger"
	"gin-basic/models"
	"gin-basic/response"
	"gin-basic/settings"
	"net/http"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type MyClaims struct {
	DeviceId string `json:"deviceId"`
	jwt.StandardClaims
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

var mySecret = []byte("伍c七Alz1θVx2ψLHNpfωv九nξ捌τD六053λwGμrMνRuegsη八γ陆jOBX8ρ三E9πFS零bδοmkχ7K6PβϵϕoZ五iυU一Jq柒ydYt四QhW4玖κCIαζTaι二σ")

// 创建token
func GenToken(DeviceId string) (string, error) {
	claims := jwt.MapClaims{}
	claims["deviceId"] = DeviceId
	claims["exp"] = time.Now().AddDate(30, 0, 0).Unix()
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		logger.Logger.Error("token error: " + err.Error())
		return nil, err
	}
	if Claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return Claims, nil
	}
	return nil, errors.New("invalid token")
}

func getBearerToken(c *gin.Context) string {
	authorization := strings.TrimSpace(c.GetHeader("Authorization"))
	if authorization == "" {
		return ""
	}

	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func JWTTokenMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := getBearerToken(c)
		if settings.Conf.App.Env != "production" {
			logger.Logger.Debug("token: " + token)
		}
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  internal.LocalizeMsg(c.MustGet("lang").(*i18n.Localizer), response.ApiMsg.AuthErr),
				"data": map[string]interface{}{},
			})
			c.Abort()
			return
		}
		mc, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  internal.LocalizeMsg(c.MustGet("lang").(*i18n.Localizer), response.ApiMsg.AuthErr),
				"data": map[string]interface{}{},
			})
			c.Abort()
			return
		}

		// 查询这个user是不是空
		var user models.User
		error := db.DB.Where("device_id = ?", mc.DeviceId).First(&user).Error
		logger.Logger.Debug("query user", zap.Any("user", user))
		if errors.Is(error, gorm.ErrRecordNotFound) {
			logger.Logger.Error("user not found")
			response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
			c.Abort()
			return
		}

		if error != nil {
			logger.Logger.Error("query user failed", zap.Error(error))
			response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
			c.Abort()
			return
		}

		// 将当前请求的userId信息保存到请求的上下文c上
		c.Set("userId", user.ID)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {

	return func(c *gin.Context) {
		token := getBearerToken(c)
		if settings.Conf.App.Env != "production" {
			logger.Logger.Debug("token: " + token)
		}
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  internal.LocalizeMsg(c.MustGet("lang").(*i18n.Localizer), response.ApiMsg.AuthErr),
				"data": map[string]interface{}{},
			})
			c.Abort()
			return
		}
		mc, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  internal.LocalizeMsg(c.MustGet("lang").(*i18n.Localizer), response.ApiMsg.AuthErr),
				"data": map[string]interface{}{},
			})
			c.Abort()
			return
		}

		// 查询这个user是不是空
		var user models.User
		userResult := db.DB.Where("device_id = ?", mc.DeviceId).First(&user)
		if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
			response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
			c.Abort()
			return
		}

		if userResult.Error != nil {
			response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
			c.Abort()
			return
		}

		if user.Role != RoleAdmin {
			c.AbortWithStatusJSON(403, gin.H{
				"msg": "permission denied",
			})
			return
		}

		// 将当前请求的userId信息保存到请求的上下文c上
		c.Set("userId", user.ID)
		c.Next()
	}
}

func OptionalJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 默认未登录
		c.Set("userId", uint(0))

		token := getBearerToken(c)
		if token == "" {
			c.Next()
			return
		}

		// 解析 token
		mc, err := ParseToken(token)
		if err != nil {
			// token 不合法，当未登录处理
			c.Next()
			return
		}

		// 查询用户是否存在
		var user models.User
		err = db.DB.
			Select("id").
			Where("device_id = ?", mc.DeviceId).
			First(&user).Error

		if err != nil {
			// 用户不存在 / 被删除
			c.Next()
			return
		}

		// 登录态有效
		c.Set("userId", user.ID)
		c.Next()
	}
}
