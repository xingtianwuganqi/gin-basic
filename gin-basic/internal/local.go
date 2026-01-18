package internal

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
)

// ReloadLocalBundle ReloadThird 加载本地国际化文件
func ReloadLocalBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载翻译文件
	//bundle.MustLoadMessageFile("locales/active.en.toml")
	//bundle.MustLoadMessageFile("locales/active.zh.toml")
	loadMessageFiles(bundle)
	return bundle
}

func loadMessageFiles(b *i18n.Bundle) {
	localesDir := "locales"
	files, err := os.ReadDir(localesDir)
	if err != nil {
		log.Fatalf("failed to read locales directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// 使用绝对路径
		b.MustLoadMessageFile(filepath.Join(localesDir, file.Name()))
	}
}

func LocalizeMsg(locale *i18n.Localizer, messageID string) string {
	return locale.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

// LocalizeMsgCount 根据指定的语言环境和数量模板本地化消息。
// 该函数使用 locale 参数来确定消息的语言环境，messageID 参数来指定消息的唯一标识符，
// 以及 count 参数来提供消息中需要本地化处理的数量信息。
// 返回值是根据提供的信息本地化后的消息字符串。
func LocalizeMsgCount(locale *i18n.Localizer, messageID string, count string) string {
	// 使用 MustLocalize 方法根据提供的配置本地化消息。
	// 这里使用 MustLocalize 而不是 Localize 是因为希望在本地化失败时程序能够 panic，
	// 表明这是一个不应该被静默处理的错误。
	return locale.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageID,
		TemplateData: map[string]interface{}{
			"Count": count,
		},
	})
}

func LocalizeMsgTemplateData(locale *i18n.Localizer, messageID string, templateData map[string]interface{}) string {
	return locale.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}