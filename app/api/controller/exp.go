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
	"sync"
	"time"
)

type EXP struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *EXP) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"one":    this.one,
		"all":    this.all,
		"sum":    this.sum,
		"min":    this.min,
		"max":    this.max,
		"rand":   this.rand,
		"count":  this.count,
		"column": this.column,
		"active": this.active,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *EXP) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"save":     this.save,
		"create":   this.create,
		"like":     this.like,
		"share" :   this.share,
		"collect" : this.collect,
		"check-in": this.checkIn,
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
func (this *EXP) IPUT(ctx *gin.Context) {
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
func (this *EXP) IDEL(ctx *gin.Context) {
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
func (this *EXP) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "没什么用！"), 202)
}

// 删除缓存
func (this *EXP) delCache() {
	// 删除缓存
	facade.Cache.DelTags([]any{"[GET]","exp"})
}

// one 获取指定数据
func (this *EXP) one(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx)

	// 表数据结构体
	table := model.EXP{}
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
func (this *EXP) all(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"page":        1,
		"order":       "create_time desc",
	})

	// 表数据结构体
	table := model.EXP{}
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
	var result []model.EXP
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
func (this *EXP) rand(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 限制最大数量
	limit  := this.meta.limit(ctx)

	// 排除的 id 列表
	except := utils.Unity.Ids(params["except"])

	onlyTrashed := cast.ToBool(params["onlyTrashed"])
	withTrashed := cast.ToBool(params["withTrashed"])

	item := facade.DB.Model(&model.EXP{}).OnlyTrashed(onlyTrashed).WithTrashed(withTrashed)
	if !utils.Is.Empty(except) {
		item = item.Where("id", "NOT IN", except)
	}

	// 从全部的 id 中随机选取指定数量的 id
	ids := utils.Rand.Slice(utils.Unity.Ids(item.Column("id")), limit)

	// 查询条件
	mold := facade.DB.Model(&[]model.EXP{}).Where("id", "IN", ids)
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
func (this *EXP) save(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.create(ctx)
	} else {
		this.update(ctx)
	}
}

// create 创建数据
func (this *EXP) create(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)
	// 验证器
	err := validator.NewValid("exp", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	uid := this.meta.user(ctx).Id
	if uid == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	// 表数据结构体
	table := model.EXP{Uid: uid, CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	allow := []any{"value", "type", "description", "json", "text"}

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
func (this *EXP) update(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "id"), 400)
		return
	}

	// 验证器
	err := validator.NewValid("exp", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.EXP{}
	allow := []any{"value", "type", "description", "json", "text"}
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

	item := facade.DB.Model(&table).WithTrashed().Where("id", params["id"])

	// 越权 - 既没有管理权限，也不是自己的数据
	if !this.meta.root(ctx) && cast.ToInt(item.Find()["uid"]) != this.user(ctx).Id {
		this.json(ctx, nil, facade.Lang(ctx, "无权限！"), 403)
		return
	}

	// 更新数据 - Scan() 解析结构体，防止 table 拿不到数据
	tx := item.Scan(&table).Update(async.Result())

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	this.json(ctx, gin.H{ "id": table.Id }, facade.Lang(ctx, "更新成功！"), 200)
}

// count 统计数据
func (this *EXP) count(ctx *gin.Context) {

	// 表数据结构体
	table := model.EXP{}
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table)
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	this.json(ctx, item.Count(), facade.Lang(ctx, "查询成功！"), 200)
}

// sum 求和
func (this *EXP) sum(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.EXP
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"])).Order(params["order"])
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	// id 数组 - 参数归一化
	ids := utils.Unity.Keys(params["ids"])
	if !utils.Is.Empty(ids) {
		item.WhereIn("id", ids)
	}

	// field 数组 - 参数归一化
	fields := utils.Unity.Keys(params["field"])

	if utils.Is.Empty(fields) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "field"), 400)
		return
	}

	cacheName := this.cache.name(ctx)
	// 开启了缓存 并且 缓存中有数据
	if this.cache.enable(ctx) && facade.Cache.Has(cacheName) {

		// 从缓存中获取数据
		msg[1] = "（来自缓存）"
		data = facade.Cache.Get(cacheName)

	} else {

		result := make(map[string]any)

		for _, val := range fields {
			result[cast.ToString(val)] = item.Sum(val)
		}

		// 从数据库中获取数据 - 排除字段
		data = result

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

// min 求最小值
func (this *EXP) min(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.EXP
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"])).Order(params["order"])
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	// id 数组 - 参数归一化
	ids := utils.Unity.Keys(params["ids"])
	if !utils.Is.Empty(ids) {
		item.WhereIn("id", ids)
	}

	// field 数组 - 参数归一化
	fields := utils.Unity.Keys(params["field"])

	if utils.Is.Empty(fields) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "field"), 400)
		return
	}

	cacheName := this.cache.name(ctx)
	// 开启了缓存 并且 缓存中有数据
	if this.cache.enable(ctx) && facade.Cache.Has(cacheName) {

		// 从缓存中获取数据
		msg[1] = "（来自缓存）"
		data = facade.Cache.Get(cacheName)

	} else {

		result := make(map[string]any)

		for _, val := range fields {
			result[cast.ToString(val)] = item.Min(val)
		}

		// 从数据库中获取数据 - 排除字段
		data = result

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

// max 求最大值
func (this *EXP) max(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.EXP
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"])).Order(params["order"])
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	// id 数组 - 参数归一化
	ids := utils.Unity.Keys(params["ids"])
	if !utils.Is.Empty(ids) {
		item.WhereIn("id", ids)
	}

	// field 数组 - 参数归一化
	fields := utils.Unity.Keys(params["field"])

	if utils.Is.Empty(fields) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "field"), 400)
		return
	}

	cacheName := this.cache.name(ctx)
	// 开启了缓存 并且 缓存中有数据
	if this.cache.enable(ctx) && facade.Cache.Has(cacheName) {

		// 从缓存中获取数据
		msg[1] = "（来自缓存）"
		data = facade.Cache.Get(cacheName)

	} else {

		result := make(map[string]any)

		for _, val := range fields {
			result[cast.ToString(val)] = item.Max(val)
		}

		// 从数据库中获取数据 - 排除字段
		data = result

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

// column 获取单列数据
func (this *EXP) column(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table []model.EXP
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"])).Order(params["order"])
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	// id 数组 - 参数归一化
	ids := utils.Unity.Keys(params["ids"])
	if !utils.Is.Empty(ids) {
		item.WhereIn("id", ids)
	}

	cacheName := this.cache.name(ctx)
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
func (this *EXP) remove(ctx *gin.Context) {

	// 表数据结构体
	table := model.EXP{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	item := facade.DB.Model(&table)

	// 越权 - 既没有管理权限，只能删除自己的数据
	if !this.meta.root(ctx) {
		item.Where("uid", this.user(ctx).Id)
	}

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
func (this *EXP) delete(ctx *gin.Context) {

	// 表数据结构体
	table := model.EXP{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	item := facade.DB.Model(&table).WithTrashed()

	// 越权 - 既没有管理权限，只能删除自己的数据
	if !this.meta.root(ctx) {
		item.Where("uid", this.user(ctx).Id)
	}

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
func (this *EXP) clear(ctx *gin.Context) {

	// 表数据结构体
	table := model.EXP{}

	item  := facade.DB.Model(&table).OnlyTrashed()

	// 越权 - 既没有管理权限，只能删除自己的数据
	if !this.meta.root(ctx) {
		item.Where("uid", this.user(ctx).Id)
	}

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
func (this *EXP) restore(ctx *gin.Context) {

	// 表数据结构体
	table := model.EXP{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	item := facade.DB.Model(&table).OnlyTrashed().WhereIn("id", ids)

	// 越权 - 既没有管理权限，只能删除自己的数据
	if !this.meta.root(ctx) {
		item.Where("uid", this.user(ctx).Id)
	}

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

// checkIn 每日签到
func (this *EXP) checkIn(ctx *gin.Context) {

	user := this.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	err := (&model.EXP{}).Add(model.EXP{
		Uid:	user.Id,
		Type:	"check-in",
	})

	if err != nil {
		this.json(ctx, gin.H{ "value": 0 }, err.Error(), 202)
		return
	}

	this.json(ctx, gin.H{ "value": 30 }, facade.Lang(ctx, "签到成功！"), 200)
}

// share 分享
func (this *EXP) share(ctx *gin.Context)  {

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"bind_type": "article",
	})

	allow := []any{"article", "page"}

	if !utils.In.Array(params["bind_type"], allow) {
		this.json(ctx, nil, facade.Lang(ctx, "不存在的分享类型！"), 400)
		return
	}

	if utils.Is.Empty(params["bind_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bind_id"), 400)
		return
	}

	user := this.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	// 从数据库里面找一下存不存在这个类型的数据
	switch params["bind_type"] {
	case "article":
		if exist := facade.DB.Model(&model.Article{}).Where("id", params["bind_id"]).Exist(); !exist {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的文章！"), 400)
			return
		}
	case "page":
		if exist := facade.DB.Model(&model.Pages{}).Where("id", params["bind_id"]).Exist(); !exist {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的页面！"), 400)
			return
		}
	}

	err := (&model.EXP{}).Add(model.EXP{
		Type:	  	 "share",
		Uid:	  	 user.Id,
		BindId:   	 cast.ToInt(params["bind_id"]),
		BindType: 	 cast.ToString(params["bind_type"]),
		Description: cast.ToString(params["description"]),
	})

	if err != nil {
		this.json(ctx, gin.H{ "value": 0 }, err.Error(), 202)
		return
	}

	this.json(ctx, gin.H{ "value": 1 }, facade.Lang(ctx, "分享成功！"), 200)
}

// collect 收藏
func (this *EXP) collect(ctx *gin.Context)  {

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"state":	 1,
		"bind_type": "article",
	})

	if !utils.InArray(cast.ToInt(params["state"]), []int{0, 1}) {
		this.json(ctx, nil, facade.Lang(ctx, "state 只能是 0 或 1"), 400)
		return
	}

	allow := []any{"article", "page"}

	if !utils.In.Array(params["bind_type"], allow) {
		this.json(ctx, nil, facade.Lang(ctx, "不存在的收藏类型！"), 400)
		return
	}

	if utils.Is.Empty(params["bind_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bind_id"), 400)
		return
	}

	user := this.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	// 从数据库里面找一下存不存在这个类型的数据
	switch params["bind_type"] {
	case "article":
		if exist := facade.DB.Model(&model.Article{}).Where("id", params["bind_id"]).Exist(); !exist {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的文章！"), 400)
			return
		}
	case "page":
		if exist := facade.DB.Model(&model.Pages{}).Where("id", params["bind_id"]).Exist(); !exist {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的页面！"), 400)
			return
		}
	}

	// 检查是否已经收藏过了
	item := facade.DB.Model(&model.EXP{}).Where([]any{
		[]any{"uid", "=", user.Id},
		[]any{"type", "=", "collect"},
		[]any{"bind_id", "=", params["bind_id"]},
		[]any{"bind_type", "=", params["bind_type"]},
	}).Find()

	// 存在记录，不允许刷经验
	if !utils.Is.Empty(item) {

		// 取消收藏
		if cast.ToInt(params["state"]) == 0 {
			tx := facade.DB.Model(&model.EXP{}).Where(item["id"]).Update(map[string]any{
				"state": 0,
			})
			if tx.Error != nil {
				this.json(ctx, nil, tx.Error.Error(), 400)
				return
			}
			this.json(ctx, gin.H{ "value": 0 }, facade.Lang(ctx, "取消收藏成功！"), 200)
			return
		}

		// 重复收藏
		if cast.ToInt(params["state"]) == 1 && cast.ToInt(item["state"]) == 1 {
			this.json(ctx, gin.H{ "value": 0 }, facade.Lang(ctx, "已经收藏过了！"), 400)
			return
		}

		// 重新收藏
		tx := facade.DB.Model(&model.EXP{}).Where(item["id"]).Update(map[string]any{
			"state": 1,
		})
		if tx.Error != nil {
			this.json(ctx, gin.H{ "value": 0 }, tx.Error.Error(), 400)
			return
		}

		this.json(ctx, gin.H{ "value": 0 }, facade.Lang(ctx, "收藏成功！"), 200)
		return
	}

	// ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ 以下为没有收藏过的情况 ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓

	if cast.ToInt(params["state"]) == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "您未收藏该内容！"), 400)
		return
	}

	err := (&model.EXP{}).Add(model.EXP{
		Type:	  	 "collect",
		Uid:	  	 user.Id,
		BindId:   	 cast.ToInt(params["bind_id"]),
		BindType: 	 cast.ToString(params["bind_type"]),
		Description: cast.ToString(params["description"]),
	})

	if err != nil {
		this.json(ctx, gin.H{ "value": 0 }, err.Error(), 202)
		return
	}

	this.json(ctx, gin.H{ "value": 1 }, facade.Lang(ctx, "收藏成功！"), 200)
}

// like 点赞
func (this *EXP) like(ctx *gin.Context)  {

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"state":	 1,
		"bind_type": "article",
	})

	if !utils.InArray(cast.ToInt(params["state"]), []int{0, 1}) {
		this.json(ctx, nil, facade.Lang(ctx, "state 只能是 0 或 1"), 400)
		return
	}

	allow := []any{"article", "page", "comment"}

	if !utils.In.Array(params["bind_type"], allow) {
		this.json(ctx, nil, facade.Lang(ctx, "不存在的点赞类型！"), 400)
		return
	}

	if utils.Is.Empty(params["bind_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bind_id"), 400)
		return
	}

	user := this.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	// 从数据库里面找一下存不存在这个类型的数据
	switch params["bind_type"] {
	case "article":
		if exist := facade.DB.Model(&model.Article{}).Where("id", params["bind_id"]).Exist(); !exist {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的文章！"), 400)
			return
		}
	case "page":
		if exist := facade.DB.Model(&model.Pages{}).Where("id", params["bind_id"]).Exist(); !exist {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的页面！"), 400)
			return
		}
	case "comment":
		if exist := facade.DB.Model(&model.Comment{}).Where("id", params["bind_id"]).Exist(); !exist {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的评论！"), 400)
			return
		}
	}

	// 检查是否已经收藏过了
	item := facade.DB.Model(&model.EXP{}).Where([]any{
		[]any{"uid", "=", user.Id},
		[]any{"type", "=", "like"},
		[]any{"bind_id", "=", params["bind_id"]},
		[]any{"bind_type", "=", params["bind_type"]},
	}).Find()

	// 存在记录，不允许刷经验
	if !utils.Is.Empty(item) {

		// 取消点赞
		if cast.ToInt(params["state"]) == 0 {
			tx := facade.DB.Model(&model.EXP{}).Where(item["id"]).Update(map[string]any{
				"state": 0,
			})
			if tx.Error != nil {
				this.json(ctx, nil, tx.Error.Error(), 400)
				return
			}
			this.json(ctx, gin.H{ "value": 0 }, facade.Lang(ctx, "点踩成功！"), 200)
			return
		}

		// 重复点赞
		if cast.ToInt(params["state"]) == 1 && cast.ToInt(item["state"]) == 1 {
			this.json(ctx, gin.H{ "value": 0 }, facade.Lang(ctx, "已经点过赞啦！"), 400)
			return
		}

		// 重新点赞
		tx := facade.DB.Model(&model.EXP{}).Where(item["id"]).Update(map[string]any{
			"state": 1,
		})
		if tx.Error != nil {
			this.json(ctx, gin.H{ "value": 0 }, tx.Error.Error(), 400)
			return
		}

		this.json(ctx, gin.H{ "value": 0 }, facade.Lang(ctx, "点赞成功！"), 200)
		return
	}

	// ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ 以下为没有点赞过的情况 ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓

	msg := "点赞"
	if cast.ToInt(params["state"]) == 0 {
		msg = "点踩"
	}

	err := (&model.EXP{}).Add(model.EXP{
		Type:	  	 "like",
		Uid:	  	 user.Id,
		State: 		 cast.ToInt(params["state"]),
		BindId:   	 cast.ToInt(params["bind_id"]),
		BindType: 	 cast.ToString(params["bind_type"]),
		Description: utils.Default(cast.ToString(params["description"]), msg + "奖励"),
	})

	if err != nil {
		this.json(ctx, gin.H{ "value": 0 }, err.Error(), 202)
		return
	}

	this.json(ctx, gin.H{ "value": 1 }, facade.Lang(ctx, msg + "成功！"), 200)
}

// active 活跃度排行
func (this *EXP) active(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx)

	// 原生Go获取本月开始时间戳
	now := time.Now()
	year, month, _ := now.Date()
	start := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	// 原生Go获取本月结束时间戳
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// 默认从本月开始到本月结束
	if utils.Is.Empty(params["start"]) {
		params["start"] = start.Unix()
	}
	if utils.Is.Empty(params["end"]) {
		params["end"] = end.Unix()
	}

	// 表数据结构体
	var table []model.EXP

	sql   := "SELECT uid, SUM(value) AS total, COUNT(id) as number FROM inis_exp WHERE create_time >= ? AND create_time <= ? GROUP BY uid ORDER BY SUM(value) DESC LIMIT ?"
	total := facade.DB.Model(&table).Query(sql, params["start"], params["end"], this.meta.limit(ctx)).Column("uid", "total", "number")
	list  := cast.ToSlice(total)

	cacheName := this.cache.name(ctx)
	// 开启了缓存 并且 缓存中有数据
	if this.cache.enable(ctx) && facade.Cache.Has(cacheName) {

		// 从缓存中获取数据
		msg[1] = "（来自缓存）"
		data = facade.Cache.Get(cacheName)

	} else {

		result := make([]any, len(list))

		wg := sync.WaitGroup{}

		for key, val  := range list {
			wg.Add(1)
			go func(key int, val any) {
				defer wg.Done()
				value      := cast.ToStringMap(val)
				field      := []string{"id", "nickname", "avatar", "description", "login_time", "title", "gender", "result"}
				author     := facade.DB.Model(&model.Users{}).Where("id", value["uid"]).Find()
				item       :=  facade.Comm.WithField(author, field)
				item["exp"] = cast.ToInt(value["total"])
				item["count"] = value["number"]
				result[key] = item
			}(key, val)
		}

		wg.Wait()

		data = result

		// 缓存数据
		if this.cache.enable(ctx) {
			go facade.Cache.Set(cacheName, result)
		}
	}

	if !utils.Is.Empty(data) {
		code = 200
		msg[0] = "数据请求成功！"
	}

	this.json(ctx, data, facade.Lang(ctx, strings.Join(msg, "")), code)
}