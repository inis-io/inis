package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
)

func Index(ctx *gin.Context) {

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	// 获取 API_KEY
	key := cast.ToStringMap(facade.DB.Model(&model.ApiKeys{}).Find())

	ctx.HTML(200, "index.html", gin.H{
		"title": "欢迎使用",
		"INIS":  utils.Json.Encode(map[string]any{
			"api": map[string]any{
				"key": key["value"],
			},
		}),
	})
}