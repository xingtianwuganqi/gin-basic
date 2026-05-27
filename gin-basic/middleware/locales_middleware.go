package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func normalizeLanguageTag(lang string) string {
	lang = strings.TrimSpace(strings.ToLower(lang))
	if lang == "" {
		return "en"
	}

	// Accept-Language may contain multiple values like: en-US,en;q=0.9
	if idx := strings.Index(lang, ","); idx >= 0 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, ";"); idx >= 0 {
		lang = lang[:idx]
	}

	switch {
	case strings.HasPrefix(lang, "zh-tw"), strings.HasPrefix(lang, "zh-hk"), strings.HasPrefix(lang, "zh-mo"), strings.HasPrefix(lang, "zh-hant"):
		return "zh-Hant"
	case strings.HasPrefix(lang, "zh"):
		return "zh"
	case strings.HasPrefix(lang, "ko"):
		return "ko"
	case strings.HasPrefix(lang, "ja"):
		return "ja"
	case strings.HasPrefix(lang, "th"):
		return "th"
	case strings.HasPrefix(lang, "fr"):
		return "fr"
	case strings.HasPrefix(lang, "de"):
		return "de"
	case strings.HasPrefix(lang, "es"):
		return "es"
	case strings.HasPrefix(lang, "en"):
		return "en"
	default:
		return "en"
	}
}

func LocaleMiddleware(bundle *i18n.Bundle) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		lang := normalizeLanguageTag(ctx.GetHeader("Accept-Language"))
		locale := i18n.NewLocalizer(bundle, lang)
		ctx.Set("lang", locale)
		ctx.Next()
	}

}
