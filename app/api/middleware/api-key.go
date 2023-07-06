package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
)

// ApiKey - 安全校验中间件
func ApiKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var open bool
		// 缓存名字
		cacheName  := "config[SYSTEM_API_KEY]"
		cacheState := cast.ToBool(facade.CacheToml.Get("open"))

		// 检查缓存是否存在
		if cacheState && facade.Cache.Has(cacheName) {
			// 获取缓存
			open = cast.ToBool(facade.Cache.Get(cacheName))
		} else {

			// 获取配置
			item := facade.DB.Model(&model.Config{}).Where("key", "SYSTEM_API_KEY").Find()
			// 转换为布尔值
			open = cast.ToBool(item["value"])
			// 设置缓存
			if cacheState {
				go facade.Cache.Set(cacheName, open)
			}
		}

		// 如果关闭了安全校验，则直接跳过
		if !open {
			ctx.Next()
			return
		}

		// 获取请求头中的 key
		var key string
		if !utils.Is.Empty(ctx.Request.Header.Get("i-api-key")) {
			key = ctx.Request.Header.Get("i-api-key")
		} else {
			key, _ = ctx.GetQuery("i-api-key")
		}

		// 为空拦截请求
		if utils.Is.Empty(key) && !utils.In.Array(ctx.Request.URL.Path, []any{"/api/file/rand"}) {
			ctx.JSON(200, gin.H{"code": 403, "msg": facade.Lang(ctx, "禁止非法操作！"), "data": nil})
			ctx.Abort()
			return
		}

		var keys []string
		cacheColumn := "[GET]/api/api-keys/column[value]"

		// 检查缓存是否存在
		if cacheState && facade.Cache.Has(cacheColumn) {

			keys = cast.ToStringSlice(facade.Cache.Get(cacheColumn))

		} else {

			// 获取所有的 key
			keys = cast.ToStringSlice(facade.DB.Model(&model.ApiKeys{}).Column("value"))
			// 设置缓存
			if cacheState {
				go facade.Cache.Set(cacheColumn, keys)
			}
		}

		// 如果 key 不在 keys 中，并且不是随机图接口，则拦截请求
		if !utils.InArray(key, keys) && !utils.In.Array(ctx.Request.URL.Path, []any{"/api/file/rand"}) {
			ctx.JSON(200, gin.H{"code": 403, "msg": facade.Lang(ctx, "禁止非法操作！"), "data": nil})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
