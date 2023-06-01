package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
	"time"
)

// Jwt - JWT 中间件
func Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tokenName := cast.ToString(facade.AppToml.Get("app.token_name", "INIS_LOGIN_TOKEN"))

		var token string
		if !utils.Is.Empty(ctx.Request.Header.Get("Authorization")) {
			token = ctx.Request.Header.Get("Authorization")
		} else {
			token, _ = ctx.Cookie(tokenName)
		}

		// 为空直接跳过
		if utils.Is.Empty(token) {
			ctx.Next()
			return
		}

		result := gin.H{"code": 401, "msg": facade.Lang(ctx, "禁止非法操作！"), "data": nil}

		jwt := facade.Jwt.Parse(token)
		if jwt.Error != nil {
			result["msg"] = utils.Ternary(jwt.Valid == 0, facade.Lang(ctx, "登录已过期，请重新登录！"), jwt.Error.Error())
			ctx.SetCookie(tokenName, "", -1, "/", "", false, false)
			ctx.JSON(200, result)
			ctx.Abort()
			return
		}

		var user map[string]any
		cacheName := fmt.Sprintf("user[%v]", jwt.Data["uid"])

		// 用户缓存不存在 - 从数据库中获取 - 并写入缓存
		if !facade.Cache.Has(cacheName) {

			user = facade.DB.Model(&model.Users{}).Find(jwt.Data["uid"])
			go func() {
				if cast.ToBool(facade.CacheToml.Get("open")) {
					facade.Cache.Set(cacheName, user, time.Duration(jwt.Valid)*time.Second)
				}
			}()
		} else {

			user = cast.ToStringMap(facade.Cache.Get(cacheName))
		}

		// 密码发生变化 - 强制退出
		if jwt.Data["hash"] != facade.Hash.Sum32(user["password"]) {
			result["msg"] = facade.Lang(ctx, "登录已过期，请重新登录！")
			ctx.SetCookie(tokenName, "", -1, "/", "", false, false)
			ctx.JSON(200, result)
			ctx.Abort()
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}
