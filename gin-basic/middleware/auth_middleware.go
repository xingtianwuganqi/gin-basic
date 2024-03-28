package middleware

import (
	"errors"
	"net/http"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"pet-project/util"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MyClaims struct {
	UserId uint `json:"userId"`
	jwt.StandardClaims
}

var mySecret = []byte("伍c七Alz1θVx2ψLHNpfωv九nξ捌τD六053λwGμrMνRuegsη八γ陆jOBX8ρ三E9πFS零bδοmkχ7K6PβϵϕoZ五iυU一Jq柒ydYt四QhW4玖κCIαζTaι二σ")

//创建token

func GenToken(userId uint) (string, error) {
	c := MyClaims{
		userId, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 10, 0).Unix(), // 过期时间
			Issuer:    "pet-project",                       // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if Claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return Claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTTokenMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		mc, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		// 查询这个user是不是空
		var user models.UserInfo
		userResult := db.DB.Where("ID = ?", mc.UserId).Find(&user)
		if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
			response.Fail(c, util.ApiCode.UserNotFont, util.ApiMessage.UserNotFound)
			return
		}

		// 将当前请求的userId信息保存到请求的上下文c上
		c.Set("userId", mc.UserId)
		c.Next()
	}
}
