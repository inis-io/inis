package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"golang.org/x/time/rate"
	"inis/app/facade"
	"inis/app/model"
	"strings"
	"sync"
	"time"
)

var (
	// QoSPoint - 单接口限流器
	QoSPoint  = make(map[string]*rate.Limiter)
	// QoSGlobal - 全局接口限流器
	QoSGlobal = make(map[string]*rate.Limiter)
	// mutex - 互斥锁
	mutex = &sync.Mutex{}
)

// QpsPoint - 单接口限流器
func QpsPoint() gin.HandlerFunc {

	go qpsDelete()
	go qpsReset()

	return func(ctx *gin.Context) {

		var config map[string]any

		// 缓存名称
		cacheName  := "config[SYSTEM_QPS]"
		// 是否开启了缓存
		cacheState := cast.ToBool(facade.CacheToml.Get("open"))

		// 检查缓存是否存在
		if cacheState && facade.Cache.Has(cacheName) {

			config = cast.ToStringMap(facade.Cache.Get(cacheName))

		} else {

			config = facade.DB.Model(&model.Config{}).Where("key", "SYSTEM_QPS").Find()
			// 存储到缓存中
			if cacheState {
				go facade.Cache.Set(cacheName, config)
			}
		}

		// 如果未开启接口限流器 - 直接跳过
		if !cast.ToBool(config["value"]) {
			ctx.Next()
			return
		}

		speed := cast.ToInt(cast.ToStringMap(config["json"])["point"])
		speed = utils.Ternary[int](utils.Is.Empty(speed), 10, speed)

		// 获取IP
		ip := ctx.ClientIP()
		// 获取URL路径
		path := ctx.Request.URL.Path
		// 获取请求方法
		method := ctx.Request.Method
		// 生成 IP+Path+Method Key
		key := fmt.Sprintf("ip=%s&path=%s&method=%s", ip, path, method)
		mutex.Lock()
		// 从Map中获取对应的访问频率限制器
		limit := QoSPoint[key]
		// 如果不存在则创建一个新的访问频率限制器
		if limit == nil {
			limit = rate.NewLimiter(rate.Every(10 * time.Millisecond), speed)
			QoSPoint[key] = limit
		}
		mutex.Unlock()

		// 尝试获取令牌
		if !limit.Allow() {

			// 记录 QPS 警告
			go QpsWarn(ctx)

			ctx.AbortWithStatusJSON(200, gin.H{"code": 429, "msg": facade.Lang(ctx, "请求过于频繁！"), "data": nil})
			return
		}

		ctx.Next()
	}
}

// QpsGlobal - 全局接口限流器
func QpsGlobal() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var config map[string]any

		cacheName  := "config[SYSTEM_QPS]"
		cacheState := cast.ToBool(facade.CacheToml.Get("open"))

		// 检查缓存是否存在
		if cacheState && facade.Cache.Has(cacheName) {

			config = cast.ToStringMap(facade.Cache.Get(cacheName))

		} else {

			config = facade.DB.Model(&model.Config{}).Where("key", "SYSTEM_QPS").Find()
			// 存储到缓存中
			if cacheState {
				go facade.Cache.Set(cacheName, config)
			}
		}

		// 如果未开启接口限流器 - 直接跳过
		if !cast.ToBool(config["value"]) {
			ctx.Next()
			return
		}

		speed := cast.ToInt(cast.ToStringMap(config["json"])["global"])
		speed = utils.Ternary[int](utils.Is.Empty(speed), 50, speed)

		// 获取IP
		ip := ctx.ClientIP()
		mutex.Lock()
		// 从Map中获取对应的访问频率限制器
		limit := QoSGlobal[ip]
		// 如果不存在则创建一个新的访问频率限制器
		if limit == nil {
			limit = rate.NewLimiter(rate.Every(10 * time.Millisecond), speed)
			QoSGlobal[ip] = limit
		}
		mutex.Unlock()
		// 尝试获取令牌
		if !limit.Allow() {

			// 记录 QPS 警告
			go QpsWarn(ctx)

			ctx.AbortWithStatusJSON(200, gin.H{"code": 429, "msg": facade.Lang(ctx, "请求过于频繁！"), "data": nil})
			return
		}

		ctx.Next()
	}
}

// qpsDelete - 监控QPSPoint和QoSGlobal的协程
func qpsDelete() {
	for {
		time.Sleep(time.Second)
		mutex.Lock()
		for key, item := range QoSPoint {
			if item.Allow() {
				delete(QoSPoint, key)
			}
		}
		for key, item := range QoSGlobal {
			if item.Allow() {
				delete(QoSGlobal, key)
			}
		}
		mutex.Unlock()
	}
}

// qpsReset - 重置QPSPoint和QoSGlobal的协程
func qpsReset() {
	// 每分钟检查一次
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		mutex.Lock()
		if len(QoSPoint) == 0 {
			for key := range QoSPoint {
				delete(QoSPoint, key)
			}
		}
		if len(QoSGlobal) == 0 {
			for key := range QoSGlobal {
				delete(QoSGlobal, key)
			}
		}
		mutex.Unlock()
	}
}

// QpsWarn - QPS警告
func QpsWarn(ctx *gin.Context) {

	var QpsBlock map[string]any

	cacheName  := "config[SYSTEM_QPS_BLOCK]"
	cacheState := cast.ToBool(facade.CacheToml.Get("open"))

	// 检查缓存是否存在
	if cacheState && facade.Cache.Has(cacheName) {

		QpsBlock = cast.ToStringMap(facade.Cache.Get(cacheName))

	} else {

		QpsBlock = facade.DB.Model(&model.Config{}).Where("key", "SYSTEM_QPS_BLOCK").Find()
		// 存储到缓存中
		if cacheState {
			go facade.Cache.Set(cacheName, QpsBlock)
		}
	}

	// 未初始化 - 直接跳过
	if utils.Is.Empty(QpsBlock) {
		return
	}

	// 如果未开启QPS警告 - 直接跳过
	if !cast.ToBool(QpsBlock["value"]) {
		return
	}

	// 自动拉黑配置
	config := cast.ToStringMap(QpsBlock["json"])
	// 获取区间时间戳
	unix   := time.Now().Add(- cast.ToDuration(utils.Calc(config["second"])) * time.Second).Unix()
	// 检查 QPS 上限
	count  := facade.DB.Model(&model.QpsWarn{}).Where("ip", ctx.ClientIP()).Where("create_time", ">", unix).Count()

	// 如果超过上限 - 自动拉黑
	if count >= cast.ToInt64(config["count"]) {

		ip := ctx.ClientIP()

		// 拉黑IP
		facade.DB.Model(&model.IpBlack{}).Where("ip", ip).Save(&model.IpBlack{
			Ip:     ip,
			Agent:  ctx.GetHeader("User-Agent"),
			Cause:  "触发QPS警告上限，自动拉黑！",
		})

		return
	}

	// 往数据库中写入警告信息
	tx := facade.DB.Model(&model.QpsWarn{}).Create(&model.QpsWarn{
		Ip:     ctx.ClientIP(),
		Agent:  ctx.GetHeader("User-Agent"),
		Path:   ctx.Request.URL.Path,
		Method: strings.ToUpper(ctx.Request.Method),
	})

	if tx.Error != nil {
		facade.Log.Error(map[string]any{
			"error":     tx.Error.Error(),
			"func_name": utils.Caller().FuncName,
			"file_name": utils.Caller().FileName,
			"file_line": utils.Caller().Line,
		}, "QPS警告写入失败！")
		return
	}
}
