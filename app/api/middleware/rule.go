package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
	"strings"
	"sync"
)

// Rule - 规则校验中间件
func Rule() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		async := sync.WaitGroup{}
		async.Add(2)

		var user model.Users
		// 获取用户信息
		go func(async *sync.WaitGroup) {
			defer async.Done()
			user = users(ctx)
		}(&async)

		var rule map[string]any
		// 获取规则
		go func(async *sync.WaitGroup) {
			defer async.Done()
			rule = rules(ctx)
			// 挂载到上下文中
			ctx.Set("route", rule)
		}(&async)

		async.Wait()

		// 如果是公共路由 - 直接放行
		if rule["type"] == "common" {
			ctx.Next()
			return
		}

		// 如果是需要登录的路由 - 判断是否登录
		if rule["type"] == "login" {
			if user.Id == 0 {
				ctx.JSON(200, gin.H{"code": 401, "msg": facade.Lang(ctx, "请先登录！"), "data": nil})
				ctx.Abort()
				return
			}
			ctx.Next()
			return
		}

		// ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ 以下为处理默认的权限规则 ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓

		// 获取用户权限
		rules := (&model.Users{}).Rules(user.Id)
		name := fmt.Sprintf("[%v][%v]", strings.ToUpper(ctx.Request.Method), ctx.Request.URL.Path)

		// 判断是否有权限
		if !utils.InArray[any](name, rules) {
			ctx.JSON(200, gin.H{"code": 403, "msg": facade.Lang(ctx, "无权限！"), "data": nil})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// 从login-token中解析用户信息
func users(ctx *gin.Context) (result model.Users) {

	// 表数据结构体
	table := model.Users{}
	keys := utils.Struct.Keys(&table)

	if user, ok := ctx.Get("user"); ok {
		for key, val := range cast.ToStringMap(user) {
			if utils.InArray[string](key, keys) && !utils.Is.Empty(val) {
				utils.Struct.Set(&table, key, val)
			}
		}
	}
	return table
}

// 获取规则
func rules(ctx *gin.Context) (result map[string]any) {

	// 缓存状态
	cacheState := cast.ToBool(facade.CacheToml.Get("open"))
	path       := strings.ReplaceAll(ctx.Request.URL.Path, "/", ".")
	// 缓存名字
	cacheName  := fmt.Sprintf("rule[%v][%v]", strings.ToUpper(ctx.Request.Method), path)

	// 如果开启了缓存 - 且缓存存在 - 直接返回缓存
	if cacheState && facade.Cache.Has(cacheName) {
		return cast.ToStringMap(facade.Cache.Get(cacheName))
	}

	// 规则缓存不存在 - 从数据库中获取 - 并写入缓存
	var table model.AuthRules

	// 规则列表
	result = facade.DB.Model(&table).Where([]any{
		[]any{"route", "=", ctx.Request.URL.Path},
		[]any{"method", "=", strings.ToUpper(ctx.Request.Method)},
	}).Find()

	// 规则列表写入缓存
	if !utils.Is.Empty(result) && cacheState {
		go facade.Cache.Set(cacheName, result, 0)
	}

	return result
}
