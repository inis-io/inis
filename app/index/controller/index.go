package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
)

func Index(ctx *gin.Context) {

	ctx.Header("Content-Type", "text/html; charset=utf-8")

	// 检查运行目录是否存在 install.lock 文件
	if utils.File().Exist("install.lock") {
		ctx.HTML(200, "index.html", gin.H{
			"title": "安装引导",
		})
		return
	}

	// 获取 API_KEY
	key := facade.DB.Model(&model.ApiKeys{}).Find()
	ctx.HTML(200, "index.html", gin.H{
		"title": "欢迎使用",
		"INIS" : utils.Json.Encode(map[string]any{
			"api": map[string]any{
				"key": key["value"],
			},
		}),
	})
}