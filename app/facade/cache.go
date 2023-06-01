package facade

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"reflect"
	"strings"
	"time"
)

func init() {

	// 初始化配置文件
	initCacheToml()
	// 初始化缓存
	initCache()

	// 监听配置文件变化
	CacheToml.Viper.WatchConfig()
	// 配置文件变化时，重新初始化配置文件
	CacheToml.Viper.OnConfigChange(func(event fsnotify.Event) {
		initCache()
	})
}

const (
	// CacheModeRedis - Redis缓存
	CacheModeRedis = "redis"
)

// NewCache - 创建Cache实例
/**
 * @param mode 驱动模式
 * @return CacheInterface
 * @example：
 * 1. cache := facade.NewCache("redis")
 * 2. cache := facade.NewCache(facade.CacheModeRedis)
 */
func NewCache(mode any) CacheInterface {
	switch strings.ToLower(cast.ToString(mode)) {
	case CacheModeRedis:
		Cache = Redis
	}
	return Cache
}

// CacheToml - 缓存配置文件
var CacheToml *utils.ViperResponse

// initCacheToml - 初始化缓存配置文件
func initCacheToml() {
	item := utils.Viper(utils.ViperModel{
		Path: "config",
		Mode: "toml",
		Name: "cache",
		Content: utils.Replace(TempCache, map[string]any{
			"${open}":           "false",
			"${default}":        "redis",
			"${local.expire}":   300,
			"${redis.host}":     "localhost",
			"${redis.port}":     "6379",
			"${redis.password}": "",
			"${redis.expire}":   7200,
			"${redis.prefix}":   "inis:",
			"${redis.database}": 0,
		}),
	}).Read()

	if item.Error != nil {
		Log.Error(map[string]any{
			"error":     item.Error,
			"func_name": utils.Caller().FuncName,
			"file_name": utils.Caller().FileName,
			"file_line": utils.Caller().Line,
		}, "缓存配置初始化错误")
		return
	}

	CacheToml = &item
}

// 初始化缓存
func initCache() {

	host := cast.ToString(CacheToml.Get("redis.host"))
	port := cast.ToString(CacheToml.Get("redis.port"))

	redisPrefix := cast.ToString(CacheToml.Get("redis.prefix"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		DB:       cast.ToInt(CacheToml.Get("redis.database")),
		Password: cast.ToString(CacheToml.Get("redis.password")),
	})
	redisExpire := time.Duration(cast.ToInt(CacheToml.Get("redis.expire", 7200))) * time.Second

	// Redis 缓存
	Redis = &RedisCacheStruct{
		Client: redisClient,
		Prefix: redisPrefix,
		Expire: redisExpire,
	}

	switch cast.ToString(CacheToml.Get("default")) {
	case "redis":
		Cache = Redis
	}
}

// Cache - Cache实例
/**
 * @return CacheInterface
 * @example：
 * cache := facade.Cache.Set("test", "这是测试", 5 * time.Minute)
 */
var Cache CacheInterface
var Redis *RedisCacheStruct

type CacheInterface interface {
	// Has
	/**
	 * @name 判断缓存是否存在
	 * @param key 缓存的key
	 * @return bool
	 */
	Has(key any) (ok bool)
	// Get
	/**
	 * @name 获取缓存
	 * @param key 缓存的key
	 * @return any 缓存值
	 */
	Get(key any) (value any)
	// Set
	/**
	 * @name 设置缓存
	 * @param key 缓存的key
	 * @param value 缓存的值
	 * @param expire （可选）过期时间
	 * @return bool
	 */
	Set(key any, value any, expire ...any) (ok bool)
	// Del
	/**
	 * @name 删除缓存
	 * @param key 缓存的key
	 * @return bool
	 */
	Del(key any) (ok bool)
	// DelPrefix
	/**
	 * @name 删除前缀缓存
	 * @param prefix 缓存的前缀
	 * @return bool
	 */
	DelPrefix(prefix ...any) (ok bool)
	// DelTags
	/**
	 * @name 删除标签缓存
	 * @param tag 缓存的标签
	 * @return bool
	 */
	DelTags(tag ...any) (ok bool)
	// Clear
	/**
	 * @name 清空缓存
	 * @return bool
	 */
	Clear() (ok bool)
}

type RedisCacheStruct struct {
	Client *redis.Client
	Prefix string
	Expire time.Duration
}

func (this *RedisCacheStruct) Has(key any) (ok bool) {

	ctx := context.Background()

	result, err := this.Client.Exists(ctx, this.Prefix+cast.ToString(key)).Result()
	return utils.Ternary[bool](err != nil, false, result == 1)
}

func (this *RedisCacheStruct) Get(key any) (value any) {

	ctx := context.Background()

	result, err := this.Client.Get(ctx, this.Prefix+cast.ToString(key)).Result()

	return utils.Ternary[any](err != nil, nil, utils.Json.Decode(result))
}

func (this *RedisCacheStruct) Set(key any, value any, expire ...any) (ok bool) {

	ctx := context.Background()
	// 设置过期时间
	if len(expire) == 0 {
		expire = append(expire, cast.ToInt(CacheToml.Get("redis.expire", this.Expire)))
	}

	// 如果 exp不为时间类型，则转码为时间类型
	if reflect.ValueOf(expire[0]).Kind() != reflect.Int64 && expire[0] != 0 {
		expire[0] = time.Duration(cast.ToInt(expire[0])) * time.Second
	}

	err := this.Client.Set(ctx, this.Prefix+cast.ToString(key), utils.Json.Encode(value), cast.ToDuration(expire[0])).Err()
	return utils.Ternary[bool](err != nil, false, true)
}

func (this *RedisCacheStruct) Del(key any) (ok bool) {

	ctx := context.Background()
	err := this.Client.Del(ctx, this.Prefix+cast.ToString(key)).Err()
	return utils.Ternary[bool](err != nil, false, true)
}

func (this *RedisCacheStruct) DelPrefix(prefix ...any) (ok bool) {

	var keys []string
	var prefixes []string
	ctx := context.Background()

	if len(prefix) == 0 {
		return false
	}

	for _, value := range prefix {

		// 判断是否为切片
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			for _, val := range cast.ToSlice(value) {
				prefixes = append(prefixes, this.Prefix+cast.ToString(val)+"*")
			}
		} else {
			prefixes = append(prefixes, this.Prefix+cast.ToString(value)+"*")
		}
	}

	// 获取 prefixes 所有的key
	for _, val := range prefixes {
		item, err := this.Client.Keys(ctx, val).Result()
		if err != nil {
			return false
		}
		keys = append(keys, item...)
	}

	// 去重 - 去空
	keys = cast.ToStringSlice(utils.Array.Empty(utils.ArrayUnique(keys)))

	if len(keys) > 0 {
		err := this.Client.Del(ctx, keys...).Err()
		if err != nil {
			return false
		}
	}

	return true
}

func (this *RedisCacheStruct) DelTags(tag ...any) (ok bool) {

	var keys []string
	var tags []string
	ctx := context.Background()

	if len(tag) == 0 {
		return false
	}

	for _, value := range tag {

		var item string

		// 判断是否为切片
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			var tmp []string
			for _, val := range cast.ToSlice(value) {
				tmp = append(tmp, cast.ToString(val))
			}
			// 数组分割字符串
			item = strings.Join(tmp, "*")
		} else {
			item = cast.ToString(value)
		}
		tags = append(tags, fmt.Sprintf("%s*%s*", this.Prefix, item))
	}

	// 获取 prefixes 所有的key
	for _, val := range tags {
		item, err := this.Client.Keys(ctx, val).Result()
		if err != nil {
			return false
		}
		keys = append(keys, item...)
	}

	// 去重 - 去空
	keys = cast.ToStringSlice(utils.Array.Empty(utils.ArrayUnique(keys)))

	if len(keys) > 0 {
		err := this.Client.Del(ctx, keys...).Err()
		if err != nil {
			return false
		}
	}

	return true
}

func (this *RedisCacheStruct) Clear() (ok bool) {

	ctx := context.Background()
	err := this.Client.FlushDB(ctx).Err()
	return utils.Ternary[bool](err != nil, false, true)
}
