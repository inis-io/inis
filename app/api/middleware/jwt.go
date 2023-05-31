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

		uid := jwt.Data["uid"]
		cacheName := fmt.Sprintf("user[%v]", uid)

		// 用户缓存不存在 - 从数据库中获取 - 并写入缓存
		if !facade.Cache.Has(cacheName) {

			item := facade.DB.Model(&model.Users{}).Find(uid)
			ctx.Set("user", item)
			if cast.ToBool(facade.CacheToml.Get("api")) {
				go func() {
					facade.Cache.Set(cacheName, item, time.Duration(jwt.Valid)*time.Second)
				}()
			}
		} else {

			ctx.Set("user", facade.Cache.Get(cacheName))
		}

		ctx.Next()
	}
}
