package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
	"inis/app/validator"
	"math"
	"strings"
	"time"
)

type Config struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Config) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"one":    this.one,
		"all":    this.all,
		"count":  this.count,
		"column": this.column,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Config) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"save":   this.save,
		"create": this.create,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}

	// 删除缓存
	go this.delCache()
}

// IPUT - PUT请求本体
func (this *Config) IPUT(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"update":  this.update,
		"restore": this.restore,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}

	// 删除缓存
	go this.delCache()
}

// IDEL - DELETE请求本体
func (this *Config) IDEL(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"remove": this.remove,
		"delete": this.delete,
		"clear":  this.clear,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}

	// 删除缓存
	go this.delCache()
}

// INDEX - GET请求本体
func (this *Config) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "没什么用！"), 202)
}

// 删除缓存
func (this *Config) delCache() {
	// 删除缓存
	facade.Cache.DelTags([]any{"[GET]","config"})
	facade.Cache.DelTags([]any{"[GET]","[?]"})
}

// one 获取指定数据
func (this *Config) one(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "key"), 400)
		return
	}

	// 表数据结构体
	table := model.Config{}
	// 允许查询的字段
	allow := []any{"key"}
	// 动态给结构体赋值
	for key, val := range params {
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			utils.Struct.Set(&table, key, val)
		}
	}

	cacheName := this.cache.name(ctx)
	// 开启了缓存 并且 缓存中有数据
	if this.cache.enable(ctx) && facade.Cache.Has(cacheName) {

		// 从缓存中获取数据
		msg[1] = "（来自缓存）"
		data = facade.Cache.Get(cacheName)

	} else {

		mold := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
		// 越权 - 无权限屏蔽系统配置
		if !this.meta.root(ctx) {
			mold.Not("key", "LIKE", "SYSTEM_%")
		}
		item := mold.Where(table).Find()

		// 排除字段
		data = facade.Comm.WithField(item, params["field"])

		// 缓存数据
		if this.cache.enable(ctx) {
			go facade.Cache.Set(cacheName, data)
		}
	}

	if !utils.Is.Empty(data) {
		code = 200
		msg[0] = "数据请求成功！"
	}

	this.json(ctx, data, facade.Lang(ctx, strings.Join(msg, "")), code)
}

// all 获取全部数据
func (this *Config) all(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"page":        1,
		"order":       "create_time desc",
	})

	// 表数据结构体
	table := model.Config{}
	// 允许查询的字段
	var allow []any
	// 动态给结构体赋值
	for key, val := range params {
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			utils.Struct.Set(&table, key, val)
		}
	}

	page := cast.ToInt(params["page"])
	limit := this.meta.limit(ctx)
	var result []model.Config
	mold := facade.DB.Model(&result).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
	mold.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	// 越权 - 无权限屏蔽系统配置
	if !this.meta.root(ctx) {
		mold.Not("key", "LIKE", "SYSTEM_%")
	}

	count := mold.Where(table).Count()

	cacheName := this.cache.name(ctx)
	// 开启了缓存 并且 缓存中有数据
	if this.cache.enable(ctx) && facade.Cache.Has(cacheName) {

		// 从缓存中获取数据
		msg[1] = "（来自缓存）"
		data = facade.Cache.Get(cacheName)

	} else {

		// 从数据库中获取数据
		item := mold.Where(table).Limit(limit).Page(page).Order(params["order"]).Select()

		// 排除字段
		data = utils.ArrayMapWithField(item, params["field"])

		// 缓存数据
		if this.cache.enable(ctx) {
			go facade.Cache.Set(cacheName, data)
		}
	}

	if !utils.Is.Empty(data) {
		code = 200
		msg[0] = "数据请求成功！"
	}

	this.json(ctx, gin.H{
		"data":  data,
		"count": count,
		"page":  math.Ceil(float64(count) / float64(limit)),
	}, facade.Lang(ctx, strings.Join(msg, "")), code)
}

// save 保存数据 - 包含创建和更新
func (this *Config) save(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "key"), 400)
		return
	}

	if !facade.DB.Model(&model.Config{}).WithTrashed().Where("key", params["key"]).Exist() {
		this.create(ctx)
	} else {
		this.update(ctx)
	}
}

// create 创建数据
func (this *Config) create(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)
	// 验证器
	err := validator.NewValid("config", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.Config{CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	allow := []any{"key", "value", "remark", "json", "text"}

	// 动态给结构体赋值
	for key, val := range params {
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			switch utils.Get.Type(val) {
			case "map":
				val = utils.Json.Encode(val)
			case "2d slice":
				val = utils.Json.Encode(val)
			case "slice":
				val = strings.Join(cast.ToStringSlice(val), ",")
			}
			utils.Struct.Set(&table, key, val)
		}
	}

	// 添加数据
	tx := facade.DB.Model(&table).Create(&table)

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	this.json(ctx, gin.H{
		"id": table.Id,
		"key": table.Key,
	}, facade.Lang(ctx, "创建成功！"), 200)
}

// update 更新数据
func (this *Config) update(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	// 验证器
	err := validator.NewValid("config", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.Config{}
	allow := []any{"key", "value", "remark", "json", "text"}
	async := utils.Async[map[string]any]()

	// 动态给结构体赋值
	for key, val := range params {
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			switch utils.Get.Type(val) {
			case "map":
				val = utils.Json.Encode(val)
			case "2d slice":
				val = utils.Json.Encode(val)
			case "slice":
				val = strings.Join(cast.ToStringSlice(val), ",")
			}
			async.Set(key, val)
		}
	}

	// 更新数据 - Scan() 方法用于将数据扫描到结构体中，使用的位置很重要
	tx := facade.DB.Model(&table).WithTrashed().Where("key", params["key"]).Scan(&table).Update(async.Result())

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	// 监听器
	go this.watch()

	this.json(ctx, gin.H{
		"id": table.Id,
		"key": table.Key,
	}, facade.Lang(ctx, "更新成功！"), 200)
}

// count 统计数据
func (this *Config) count(ctx *gin.Context) {

	// 表数据结构体
	table := model.Config{}
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	// 越权 - 无权限屏蔽系统配置
	if !this.meta.root(ctx) {
		item.Not("key", "LIKE", "SYSTEM_%")
	}

	this.json(ctx, item.Count(), facade.Lang(ctx, "查询成功！"), 200)
}

// column 获取单列数据
func (this *Config) column(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table []model.Config
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"])).Order(params["order"])
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	// key 数组 - 参数归一化
	keys := utils.Unity.Keys(params["keys"])
	if !utils.Is.Empty(keys) {
		item.WhereIn("key", keys)
	}

	cacheName := this.cache.name(ctx)
	// 越权 - 无权限屏蔽系统配置
	if !this.meta.root(ctx) {
		item.Not("key", "LIKE", "SYSTEM_%")
		cacheName += "&root"
	}
	// 开启了缓存 并且 缓存中有数据
	if this.cache.enable(ctx) && facade.Cache.Has(cacheName) {

		// 从缓存中获取数据
		msg[1] = "（来自缓存）"
		data = facade.Cache.Get(cacheName)

	} else {

		// 从数据库中获取数据 - 排除字段
		data = utils.ArrayMapWithField(item.Select(), params["field"])

		// 缓存数据
		if this.cache.enable(ctx) {
			go facade.Cache.Set(cacheName, data)
		}
	}

	if !utils.Is.Empty(data) {
		code = 200
		msg[0] = "数据请求成功！"
	}

	this.json(ctx, data, facade.Lang(ctx, strings.Join(msg, "")), code)
}

// remove 软删除
func (this *Config) remove(ctx *gin.Context) {

	// 表数据结构体
	var table []model.Config
	// 获取请求参数
	params := this.params(ctx)

	// key 数组 - 参数归一化
	keys := utils.Unity.Keys(params["keys"])

	if utils.Is.Empty(keys) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "keys"), 400)
		return
	}

	item := facade.DB.Model(&table)

	// 得到允许操作的 key 数组
	keys = utils.Unity.Keys(item.WhereIn("key", keys).Column("key"))

	// 无可操作数据
	if utils.Is.Empty(keys) {
		this.json(ctx, nil, facade.Lang(ctx, "无可操作数据！"), 204)
		return
	}

	// 软删除
	tx := item.WhereIn("key", keys).Delete()

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "删除失败！"), 400)
		return
	}

	this.json(ctx, gin.H{ "keys": keys }, facade.Lang(ctx, "删除成功！"), 200)
}

// delete 真实删除
func (this *Config) delete(ctx *gin.Context) {

	// 表数据结构体
	var table []model.Config
	// 获取请求参数
	params := this.params(ctx)

	// key 数组 - 参数归一化
	keys := utils.Unity.Keys(params["keys"])

	if utils.Is.Empty(keys) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "keys"), 400)
		return
	}

	item := facade.DB.Model(&table).WithTrashed()

	// 得到允许操作的 key 数组
	keys = utils.Unity.Keys(item.WhereIn("key", keys).Column("key"))

	// 无可操作数据
	if utils.Is.Empty(keys) {
		this.json(ctx, nil, facade.Lang(ctx, "无可操作数据！"), 204)
		return
	}

	// 软删除
	tx := item.WhereIn("key", keys).Force().Delete()

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "删除失败！"), 400)
		return
	}

	this.json(ctx, gin.H{ "keys": keys }, facade.Lang(ctx, "删除成功！"), 200)
}

// clear 清空回收站
func (this *Config) clear(ctx *gin.Context) {

	// 表数据结构体
	table := model.Config{}

	item  := facade.DB.Model(&table).OnlyTrashed()

	keys  := utils.Unity.Keys(item.Column("key"))

	// 无可操作数据
	if utils.Is.Empty(keys) {
		this.json(ctx, nil, facade.Lang(ctx, "无可操作数据！"), 204)
		return
	}

	// 找到所有软删除的数据
	tx := item.Force().Delete()

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "清空失败！"), 400)
		return
	}

	this.json(ctx, gin.H{ "keys": keys }, facade.Lang(ctx, "清空成功！"), 200)
}

// restore 恢复数据
func (this *Config) restore(ctx *gin.Context) {

	// 表数据结构体
	var table []model.Config
	// 获取请求参数
	params := this.params(ctx)

	// key 数组 - 参数归一化
	keys := utils.Unity.Keys(params["keys"])

	if utils.Is.Empty(keys) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "keys"), 400)
		return
	}

	item := facade.DB.Model(&table).OnlyTrashed().WhereIn("key", keys)

	// 得到允许操作的 key 数组
	keys = utils.Unity.Keys(item.Column("key"))

	// 无可操作数据
	if utils.Is.Empty(keys) {
		this.json(ctx, nil, facade.Lang(ctx, "无可操作数据！"), 204)
		return
	}

	// 软删除
	tx := facade.DB.Model(&table).OnlyTrashed().WhereIn("key", keys).Restore()

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "恢复失败！"), 400)
		return
	}

	this.json(ctx, gin.H{ "keys": keys }, facade.Lang(ctx, "恢复成功！"), 200)
}

// watch 监听数据
func (this *Config) watch() {

	item := facade.DB.Model(&model.Config{}).Where("key", "SYSTEM_API_KEY").Find()
	if cast.ToInt(item["value"]) == 1 {

		ApiKeys := facade.DB.Model(&model.ApiKeys{})
		// 检查 密钥 是否为空
		if ApiKeys.Count() == 0 {
			// 生成一个随机的UUID
			UUID := uuid.New().String()
			// 去除UUID中的横杠
			UUID = strings.Replace(UUID, "-", "", -1)
			ApiKeys.Create(&model.ApiKeys{
				Value: strings.ToUpper(UUID),
				Remark: "检测到您开启了API_KEY功能，但无可用密钥，兔子贴心的为您创建了一个密钥！不用谢，已为您自动生成五星好评！",
			})
		}
	}
}