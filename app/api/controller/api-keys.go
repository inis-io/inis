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

type ApiKeys struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *ApiKeys) IGET(ctx *gin.Context) {
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
func (this *ApiKeys) IPOST(ctx *gin.Context) {

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
func (this *ApiKeys) IPUT(ctx *gin.Context) {
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
func (this *ApiKeys) IDEL(ctx *gin.Context) {
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
func (this *ApiKeys) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "没什么用！"), 202)
}

// 删除缓存
func (this *ApiKeys) delCache() {
	// 删除缓存
	facade.Cache.DelTags([]any{"<GET>","api-keys"})
}

// one 获取指定数据
func (this *ApiKeys) one(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx)

	// 表数据结构体
	table := model.ApiKeys{}
	// 允许查询的字段
	allow := []any{"id"}
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

		mold := facade.DB.Model(&table).OnlyTrashed(params["onlyTrashed"]).WithTrashed(params["withTrashed"])
		mold.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])
		item := mold.Where(table).Find()

		// 缓存数据
		if this.cache.enable(ctx) {
			go facade.Cache.Set(cacheName, item)
		}

		data = item
	}

	if !utils.Is.Empty(data) {
		code = 200
		msg[0] = "数据请求成功！"
	}

	this.json(ctx, data, facade.Lang(ctx, strings.Join(msg, "")), code)
}

// all 获取全部数据
func (this *ApiKeys) all(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"page":        1,
		"limit":       5,
		"order":       "create_time desc",
	})

	// 表数据结构体
	table := model.ApiKeys{}
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
	limit := cast.ToInt(params["limit"])
	var result []model.ApiKeys
	mold := facade.DB.Model(&result).OnlyTrashed(params["onlyTrashed"]).WithTrashed(params["withTrashed"])
	mold.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])
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

		// 缓存数据
		if this.cache.enable(ctx) {
			go facade.Cache.Set(cacheName, item)
		}

		data = item
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
func (this *ApiKeys) save(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.create(ctx)
	} else {
		this.update(ctx)
	}
}

// create 创建数据
func (this *ApiKeys) create(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)
	// 验证器
	err := validator.NewValid("api-keys", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.ApiKeys{CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	allow := []any{"value", "remark", "json", "text"}

	// 动态给结构体赋值
	for key, val := range params {
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			if key == "value" {
				val = strings.ToUpper(cast.ToString(val))
			}
			utils.Struct.Set(&table, key, val)
		}
	}

	// 判断 value 是否是空的 - 自动创建
	if utils.Is.Empty(table.Value) {
		// 生成一个随机的UUID
		UUID := uuid.New().String()
		// 去除UUID中的横杠
		UUID = strings.Replace(UUID, "-", "", -1)
		// 转换成大写
		utils.Struct.Set(&table, "value", strings.ToUpper(UUID))
	}

	// 检查 value 唯一性
	if facade.DB.Model(&table).Where("value", table.Value).Exist() {
		this.json(ctx, nil, facade.Lang(ctx, "%s 已经存在！", table.Value), 400)
		return
	}

	// 添加数据
	tx := facade.DB.Model(&table).Create(&table)

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	this.json(ctx, map[string]any{
		"id": table.Id,
	}, facade.Lang(ctx, "创建成功！"), 200)
}

// update 更新数据
func (this *ApiKeys) update(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "id"), 400)
		return
	}

	// 验证器
	err := validator.NewValid("api-keys", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.ApiKeys{}
	allow := []any{"value", "remark", "json", "text"}
	async := utils.Async[map[string]any]()

	// 动态给结构体赋值
	for key, val := range params {
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			if key == "value" {
				val = strings.ToUpper(cast.ToString(val))
			}
			async.Set(key, val)
		}
	}

	// 判断 value 是否是空的 - 自动创建
	if utils.Is.Empty(async.Get("value")) {
		// 生成一个随机的UUID
		UUID := uuid.New().String()
		// 去除UUID中的横杠
		UUID = strings.Replace(UUID, "-", "", -1)
		// 转换成大写
		async.Set("value", strings.ToUpper(UUID))
	}

	key  := cast.ToString(async.Get("value"))
	item := facade.DB.Model(&table).Where("value", key).Find()

	// 检查 value 唯一性
	if !utils.IsEmpty(item) && item["value"] != key {
		this.json(ctx, nil, facade.Lang(ctx, "%s 已经存在！", key), 400)
		return
	}

	// 更新数据 - Scan() 方法用于将数据扫描到结构体中，使用的位置很重要
	tx := facade.DB.Model(&table).WithTrashed().Where("id", params["id"]).Scan(&table).Update(async.Result())

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	this.json(ctx, map[string]any{
		"id": table.Id,
	}, facade.Lang(ctx, "更新成功！"), 200)
}

// count 统计数据
func (this *ApiKeys) count(ctx *gin.Context) {

	// 表数据结构体
	table := model.ApiKeys{}
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(params["onlyTrashed"]).WithTrashed(params["withTrashed"])
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	this.json(ctx, item.Count(), facade.Lang(ctx, "查询成功！"), 200)
}

// column 获取单列数据
func (this *ApiKeys) column(ctx *gin.Context) {

	// 表数据结构体
	table := model.ApiKeys{}
	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"field": "*",
	})

	item := facade.DB.Model(&table).OnlyTrashed(params["onlyTrashed"]).WithTrashed(params["withTrashed"]).Order(params["order"])
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	if !strings.Contains(cast.ToString(params["field"]), "*") {
		item.Field(params["field"])
	}

	// id 数组 - 参数归一化
	ids := utils.Unity.Keys(params["ids"])
	if !utils.Is.Empty(ids) {
		item.WhereIn("id", ids)
	}

	code := 200
	data := item.Column()
	msg  := facade.Lang(ctx, "查询成功！")

	if utils.Is.Empty(data) {
		code = 204
		msg  = facade.Lang(ctx, "无数据！")
	}

	this.json(ctx, data, msg, code)
}

// remove 软删除
func (this *ApiKeys) remove(ctx *gin.Context) {

	// 表数据结构体
	table := model.ApiKeys{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	// 软删除
	tx := facade.DB.Model(&table).Delete(ids)

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "删除失败！"), 400)
		return
	}

	// 如果是全部删除，检查是否开启了 API_KEY
	if facade.DB.Model(&model.ApiKeys{}).Count() == 0 {

		item := facade.DB.Model(&model.Config{}).Where("key", "SYSTEM_API_KEY")

		if cast.ToInt(item.Find()["value"]) == 1 {
			res := item.Update(map[string]any{
				"value": 0,
			})
			if res.Error == nil {
				go facade.Cache.DelTags("SYSTEM_API_KEY")
				this.json(ctx, nil, facade.Lang(ctx, "删除成功！<br>同时检测到您开启了API_KEY，但无密钥可用。<br>兔子已为您自动关闭API_KEY功能！"), 200)
				return
			}
		}
	}

	this.json(ctx, nil, facade.Lang(ctx, "删除成功！"), 200)
}

// delete 真实删除
func (this *ApiKeys) delete(ctx *gin.Context) {

	// 表数据结构体
	table := model.ApiKeys{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	// 真实删除
	tx := facade.DB.Model(&table).WithTrashed().Force().Delete(ids)

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "删除失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "删除成功！"), 200)
}

// clear 清空回收站
func (this *ApiKeys) clear(ctx *gin.Context) {

	// 表数据结构体
	table := model.ApiKeys{}

	// 找到所有软删除的数据
	tx := facade.DB.Model(&table).OnlyTrashed().Force().Delete()

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "清空失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "清空成功！"), 200)
}

// restore 恢复数据
func (this *ApiKeys) restore(ctx *gin.Context) {

	// 表数据结构体
	table := model.ApiKeys{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, params, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	// 还原数据
	tx := facade.DB.Model(&table).Restore(ids)

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "恢复失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "恢复成功！"), 200)
}