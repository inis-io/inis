package facade

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"strings"
)

func Lang(ctx *gin.Context, key string, args ...any) (result string) {

	var lang string
	// 获取语言
	lang, _ = ctx.Cookie("inis_lang")
	lang = ctx.DefaultQuery("inis_lang", lang)
	lang = strings.ToLower(lang)

	model := utils.LangModel{
		Directory: "config/i18n/",
	}

	// 设置语言
	if !utils.Is.Empty(lang) {
		model.Lang = lang
	}

	return cast.ToString(utils.Lang(model).Value(key, args...))
}
