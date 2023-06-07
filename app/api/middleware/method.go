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

// Method - 请求类型校验中间件
func Method() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tokenName := cast.ToString(facade.AppToml.Get("app.token_name", "INIS_LOGIN_TOKEN"))

		var token string
		if !utils.Is.Empty(ctx.Request.Header.Get("Authorization")) {
			token = ctx.Request.Header.Get("Authorization")
		} else {
			token, _ = ctx.Cookie(tokenName)
		}

		method := []any{"POST", "PUT", "DELETE", "PATCH"}

		if utils.In.Array(ctx.Request.Method, method) {

			result := gin.H{"code": 401, "msg": facade.Lang(ctx, "禁止非法操作！"), "data": nil}

			if utils.Is.Empty(token) {
				ctx.JSON(200, result)
				ctx.Abort()
				return
			}

			// 解析token
			jwt := facade.Jwt().Parse(token)
			if jwt.Error != nil {
				result["msg"] = utils.Ternary(jwt.Valid == 0, facade.Lang(ctx, "登录已过期，请重新登录！"), jwt.Error.Error())
				ctx.SetCookie(tokenName, "", -1, "/", "", false, false)
				ctx.JSON(200, result)
				ctx.Abort()
				return
			}

			uid := jwt.Data["uid"]
			cacheName  := fmt.Sprintf("user[%v]", uid)
			cacheState := cast.ToBool(facade.CacheToml.Get("open"))

			// 如果开启了缓存 - 且缓存存在 - 直接从缓存中获取
			if cacheState && facade.Cache.Has(cacheName) {

				ctx.Set("user", facade.Cache.Get(cacheName))

			} else {

				item := facade.DB.Model(&model.Users{}).Find(uid)
				ctx.Set("user", item)

				if cacheState {
					go facade.Cache.Set(cacheName, item, time.Duration(jwt.Valid)*time.Second)
				}
			}
		}

		ctx.Next()
	}
}
