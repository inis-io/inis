package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middle "inis/app/api/middleware"
	"inis/app/inis/controller"
	global "inis/app/middleware"
)

func Route(Gin *gin.Engine) {

	// 全局中间件
	group := Gin.Group("/inis/").Use(
		middle.IpBlack(),   // IP黑名单
		global.QpsPoint(),  // QPS限制 - 单接口限流
		global.QpsGlobal(), // QPS限制 - 全局限流
		global.Params(),    // 解析参数
		middle.Jwt(),       // 解析JWT
		middle.Rule(),      // 验证规则
		middle.ApiKey(),    // 安全校验
	)

	// 允许动态挂载的路由
	allow := map[string]controller.ApiInterface{
		"test":           &controller.Test{},
		"order":          &controller.Order{},
		"users":          &controller.Users{},
		"device":         &controller.Device{},
		"theme-version":  &controller.ThemeVersion{},
		"system-version": &controller.SystemVersion{},
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
