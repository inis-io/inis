package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"inis/app/api/controller"
	middle "inis/app/api/middleware"
	global "inis/app/middleware"
)

func Route(Gin *gin.Engine) {

	// 全局中间件
	group := Gin.Group("/api/").Use(
		global.QpsPoint(),  // QPS限制 - 单接口限流
		global.QpsGlobal(), // QPS限制 - 全局限流
		global.Params(),    // 解析参数
		middle.Jwt(),       // 验证权限
		middle.Rule(),      // 验证规则
		middle.Security(),  // 安全校验
	)

	// 允许动态挂载的路由
	allow := map[string]controller.ApiInterface{
		"test":          &controller.Test{},
		"comm":          &controller.Comm{},
		"file":          &controller.File{},
		"tags":          &controller.Tags{},
		"pages":         &controller.Pages{},
		"users":         &controller.Users{},
		"oauth":         &controller.OAuth{},
		"links":         &controller.Links{},
		"proxy":         &controller.Proxy{},
		"banner":        &controller.Banner{},
		"config":        &controller.Config{},
		"article":       &controller.Article{},
		"comment":       &controller.Comment{},
		"placard":       &controller.Placard{},
		"api-keys":      &controller.ApiKeys{},
		"auth-group":    &controller.AuthGroup{},
		"auth-pages":    &controller.AuthPages{},
		"auth-rules":    &controller.AuthRules{},
		"links-group":   &controller.LinksGroup{},
		"article-group": &controller.ArticleGroup{},
	}

	// 动态配置路由
	for key, item := range allow {
		group.Any(key, item.INDEX)
		group.GET(fmt.Sprintf("%s/:method", key), item.IGET)
		group.PUT(fmt.Sprintf("%s/:method", key), item.IPUT)
		group.POST(fmt.Sprintf("%s/:method", key), item.IPOST)
		group.DELETE(fmt.Sprintf("%s/:method", key), item.IDEL)
	}
}
