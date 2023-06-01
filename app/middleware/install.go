package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unti-io/go-utils/utils"
	"strings"
)

// Install 安装引导中间件
func Install() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		path   := ctx.Request.URL.Path
		method := strings.ToUpper(ctx.Request.Method)

		// 检查运行目录是否存在 install.lock 文件
		if utils.File().Exist("install.lock") {

			if path == "/" && method == "GET" {

				if ok, _ := ctx.Cookie("install"); !utils.Is.Empty(ok) {
					ctx.Next()
					return
				}

				ctx.SetCookie("install", "1", 3, "/", "", false, true)

				// 重定向到后台首页 http.StatusMovedPermanently
				ctx.Redirect(301, "/#/install")
				ctx.Abort()

				return
			}

			// 安装引导没完成之前，禁止访问API应用
			if strings.HasPrefix(path, "/api/") {
				ctx.JSON(200, map[string]any{ "code": 412, "msg": "安装引导未完成，禁止访问！", "data": nil })
				ctx.Abort()
				return
			}
		} else {

			// 安装引导完成之后，禁止访问安装引导程序
			if strings.HasPrefix(path, "/dev/install") {
				ctx.JSON(200, map[string]any{ "code": 412, "msg": "程序已完成安装，禁止访问！", "data": nil })
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}