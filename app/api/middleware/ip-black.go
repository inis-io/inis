package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
)

// IpBlack - IP黑名单
func IpBlack() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var column []string

		// 缓存名字
		cacheName  := "[GET][ip-black][column]"
		cacheState := cast.ToBool(facade.CacheToml.Get("open"))

		// 检查缓存是否存在
		if cacheState && facade.Cache.Has(cacheName) {

			// 获取缓存
			column = cast.ToStringSlice(facade.Cache.Get(cacheName))

		} else {

			// 获取黑名单列表
			list := facade.DB.Model(&model.IpBlack{}).Column("ip")
			// 去空 - 去重
			column = cast.ToStringSlice(utils.ArrayEmpty(utils.ArrayUnique(cast.ToStringSlice(list))))
		}

		// 检查客户端IP是否在黑名单中
		if utils.InArray[string](ctx.ClientIP(), column) {
			ctx.JSON(200, gin.H{"code": 406, "msg": "IP已被列入黑名单", "data": nil})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}