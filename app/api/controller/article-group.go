package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
	"inis/app/validator"
	"math"
	"strings"
	"time"
)

type ArticleGroup struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *ArticleGroup) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"one":    this.one,
		"all":    this.all,
		"rand":   this.rand,
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
func (this *ArticleGroup) IPOST(ctx *gin.Context) {

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
func (this *ArticleGroup) IPUT(ctx *gin.Context) {
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
func (this *ArticleGroup) IDEL(ctx *gin.Context) {
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
func (this *ArticleGroup) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "没什么用！"), 202)
}

// 删除缓存
func (this *ArticleGroup) delCache() {
	// 删除缓存
	facade.Cache.DelTags([]any{"[GET]","article-group"})
}

// one 获取指定数据
func (this *ArticleGroup) one(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx)

	// 表数据结构体
	table := model.ArticleGroup{}
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

		mold := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
		mold.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])
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
func (this *ArticleGroup) all(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"page":        1,
		"order":       "create_time desc",
	})

	// 表数据结构体
	table := model.ArticleGroup{}
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
	var result []model.ArticleGroup
	mold := facade.DB.Model(&result).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
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

// rand 随机获取
func (this *ArticleGroup) rand(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 限制最大数量
	limit  := this.meta.limit(ctx)

	// 排除的 id 列表
	except := utils.Unity.Ids(params["except"])

	onlyTrashed := cast.ToBool(params["onlyTrashed"])
	withTrashed := cast.ToBool(params["withTrashed"])

	item := facade.DB.Model(&model.ArticleGroup{}).OnlyTrashed(onlyTrashed).WithTrashed(withTrashed)
	if !utils.Is.Empty(except) {
		item = item.Where("id", "NOT IN", except)
	}

	// 从全部的 id 中随机选取指定数量的 id
	ids := utils.Rand.Slice(utils.Unity.Ids(item.Column("id")), limit)

	// 查询条件
	mold := facade.DB.Model(&[]model.ArticleGroup{}).Where("id", "IN", ids)
	mold.OnlyTrashed(onlyTrashed).WithTrashed(withTrashed).IWhere(params["where"]).IOr(params["or"])
	mold.ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	// 查询并打乱顺序
	data := utils.Array.MapWithField(utils.Rand.MapSlice(mold.Select()), params["field"])

	if utils.Is.Empty(data) {
		this.json(ctx, nil, facade.Lang(ctx, "无数据！"), 204)
		return
	}

	this.json(ctx, data, facade.Lang(ctx, "好的！"), 200)
}

// save 保存数据 - 包含创建和更新
func (this *ArticleGroup) save(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.create(ctx)
	} else {
		this.update(ctx)
	}
}

// create 创建数据
func (this *ArticleGroup) create(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)
	// 验证器
	err := validator.NewValid("article-group", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.ArticleGroup{CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	allow := []any{"pid", "name", "key", "description", "avatar", "json", "text"}

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

	this.json(ctx, gin.H{ "id": table.Id }, facade.Lang(ctx, "创建成功！"), 200)
}

// update 更新数据
func (this *ArticleGroup) update(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "id"), 400)
		return
	}

	// 验证器
	err := validator.NewValid("article-group", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.ArticleGroup{}
	allow := []any{"pid", "name", "key", "description", "avatar", "json", "text"}
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
	tx := facade.DB.Model(&table).WithTrashed().Where("id", params["id"]).Scan(&table).Update(async.Result())

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	this.json(ctx, gin.H{ "id": table.Id }, facade.Lang(ctx, "更新成功！"), 200)
}

// count 统计数据
func (this *ArticleGroup) count(ctx *gin.Context) {

	// 表数据结构体
	table := model.ArticleGroup{}
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	this.json(ctx, item.Count(), facade.Lang(ctx, "查询成功！"), 200)
}

// column 获取单列数据
func (this *ArticleGroup) column(ctx *gin.Context) {

	// 表数据结构体
	table := model.ArticleGroup{}
	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"field": "*",
	})

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"])).Order(params["order"])
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
func (this *ArticleGroup) remove(ctx *gin.Context) {

	// 表数据结构体
	table := model.ArticleGroup{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	item := facade.DB.Model(&table)

	// 得到允许操作的 id 数组
	ids = utils.Unity.Ids(item.WhereIn("id", ids).Column("id"))

	// 无可操作数据
	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "无可操作数据！"), 204)
		return
	}

	// 软删除
	tx := item.Delete(ids)

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "删除失败！"), 400)
		return
	}

	this.json(ctx, gin.H{ "ids": ids }, facade.Lang(ctx, "删除成功！"), 200)
}

// delete 真实删除
func (this *ArticleGroup) delete(ctx *gin.Context) {

	// 表数据结构体
	table := model.ArticleGroup{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	item := facade.DB.Model(&table).WithTrashed()

	// 得到允许操作的 id 数组
	ids = utils.Unity.Ids(item.WhereIn("id", ids).Column("id"))

	// 无可操作数据
	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "无可操作数据！"), 204)
		return
	}

	// 真实删除
	tx := item.Force().Delete(ids)

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "删除失败！"), 400)
		return
	}

	this.json(ctx, gin.H{ "ids": ids }, facade.Lang(ctx, "删除成功！"), 200)
}

// clear 清空回收站
func (this *ArticleGroup) clear(ctx *gin.Context) {

	// 表数据结构体
	table := model.ArticleGroup{}

	item  := facade.DB.Model(&table).OnlyTrashed()

	ids := utils.Unity.Ids(item.Column("id"))

	// 无可操作数据
	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "无可操作数据！"), 204)
		return
	}

	// 找到所有软删除的数据
	tx := item.Force().Delete()

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "清空失败！"), 400)
		return
	}

	this.json(ctx, gin.H{ "ids": ids }, facade.Lang(ctx, "清空成功！"), 200)
}

// restore 恢复数据
func (this *ArticleGroup) restore(ctx *gin.Context) {

	// 表数据结构体
	table := model.ArticleGroup{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	item := facade.DB.Model(&table).OnlyTrashed().WhereIn("id", ids)

	// 得到允许操作的 id 数组
	ids = utils.Unity.Ids(item.Column("id"))

	// 无可操作数据
	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "无可操作数据！"), 204)
		return
	}

	// 还原数据
	tx := facade.DB.Model(&table).OnlyTrashed().Restore(ids)

	if tx.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "恢复失败！"), 400)
		return
	}

	this.json(ctx, gin.H{ "ids": ids }, facade.Lang(ctx, "恢复成功！"), 200)
}