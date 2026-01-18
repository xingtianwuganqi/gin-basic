package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func LocaleMiddleware(bundle *i18n.Bundle) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		lang := ctx.GetHeader("Accept-Language")
		if lang == "" {
			lang = "zh"
		}
		locale := i18n.NewLocalizer(bundle, lang)
		ctx.Set("lang", locale)
		ctx.Next()
	}

}