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

type Comment struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Comment) IGET(ctx *gin.Context) {

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
		"flat":   this.flat,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Comment) IPOST(ctx *gin.Context) {

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
func (this *Comment) IPUT(ctx *gin.Context) {
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
func (this *Comment) IDEL(ctx *gin.Context) {
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
func (this *Comment) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "没什么用！"), 202)
}

// 删除缓存
func (this *Comment) delCache() {
	// 删除缓存
	facade.Cache.DelTags([]any{"[GET]","comment"})
}

// one 获取指定数据
func (this *Comment) one(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx)

	// 表数据结构体
	table := model.Comment{}
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
func (this *Comment) all(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"page":        1,
		"order":       "create_time desc",
	})

	// 表数据结构体
	table := model.Comment{}
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
	var result []model.Comment
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
func (this *Comment) rand(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 限制最大数量
	limit  := this.meta.limit(ctx)

	// 排除的 id 列表
	except := utils.Unity.Ids(params["except"])

	onlyTrashed := cast.ToBool(params["onlyTrashed"])
	withTrashed := cast.ToBool(params["withTrashed"])

	item := facade.DB.Model(&model.Comment{}).OnlyTrashed(onlyTrashed).WithTrashed(withTrashed)
	if !utils.Is.Empty(except) {
		item = item.Where("id", "NOT IN", except)
	}

	// 从全部的 id 中随机选取指定数量的 id
	ids := utils.Rand.Slice(utils.Unity.Ids(item.Column("id")), limit)

	// 查询条件
	mold := facade.DB.Model(&[]model.Comment{}).Where("id", "IN", ids)
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
func (this *Comment) save(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.create(ctx)
	} else {
		this.update(ctx)
	}
}

// create 创建数据
func (this *Comment) create(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"bind_type": "article",
	})
	// 验证器
	err := validator.NewValid("comment", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	if utils.Is.Empty(params["bind_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bind_id"), 400)
		return
	}

	user := this.meta.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	// 评论配置
	var comment map[string]any

	// 从数据库里面找一下存不存在这个类型的数据
	switch params["bind_type"] {
	case "article":
		article := facade.DB.Model(&model.Article{}).Where("id", params["bind_id"]).Find()
		if utils.Is.Empty(article) {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的文章！"), 400)
			return
		}
		comment = cast.ToStringMap(cast.ToStringMap(article["json"])["comment"])
	case "page":
		page := facade.DB.Model(&model.Pages{}).Where("id", params["bind_id"]).Find()
		if utils.Is.Empty(page) {
			this.json(ctx, nil, facade.Lang(ctx, "不存在的页面！"), 400)
			return
		}
		comment = cast.ToStringMap(cast.ToStringMap(page["json"])["comment"])
	}

	// 允许评论选项继承了父级配置
	if cast.ToInt(comment["allow"]) == 0 {
		comment["allow"] = this.config("comment")["allow"]
	}

	// 评论开关
	if cast.ToInt(comment["allow"]) == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "评论功能已关闭！"), 400)
		return
	}

	// 表数据结构体
	table := model.Comment{
		Uid: user.Id,
		Agent: this.header(ctx, "User-Agent"),
		Ip: cast.ToString(this.get(ctx, "ip")),
		CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix(),
	}
	allow := []any{"pid", "content", "bind_id", "bind_type", "editor", "json", "text"}

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

	// 更新用户经验
	go func() {
		_ = (&model.EXP{}).Add(model.EXP{
			Uid:  user.Id,
			Type: "comment",
			BindId: table.BindId,
			BindType: table.BindType,
		})
	}()

	this.json(ctx, gin.H{ "id": table.Id }, facade.Lang(ctx, "创建成功！"), 200)
}

// update 更新数据
func (this *Comment) update(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "id"), 400)
		return
	}

	// 验证器
	err := validator.NewValid("comment", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.Comment{}
	allow := []any{"content", "editor", "json", "text"}
	async := utils.Async[map[string]any]()

	root := this.meta.root(ctx)

	// 越权 - 增加可选字段
	if root {
		allow = append(allow, "pid", "bind_id", "bind_type")
	}

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
	if !root && cast.ToInt(item.Find()["uid"]) != this.user(ctx).Id {
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
func (this *Comment) count(ctx *gin.Context) {

	// 表数据结构体
	table := model.Comment{}
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	this.json(ctx, item.Count(), facade.Lang(ctx, "查询成功！"), 200)
}

// sum 求和
func (this *Comment) sum(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.Comment
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
func (this *Comment) min(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.Comment
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
func (this *Comment) max(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.Comment
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
func (this *Comment) column(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table []model.Comment
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
func (this *Comment) remove(ctx *gin.Context) {

	// 表数据结构体
	table := model.Comment{}
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
func (this *Comment) delete(ctx *gin.Context) {

	// 表数据结构体
	table := model.Comment{}
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
func (this *Comment) clear(ctx *gin.Context) {

	// 表数据结构体
	table := model.Comment{}

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
func (this *Comment) restore(ctx *gin.Context) {

	// 表数据结构体
	table := model.Comment{}
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

// replies 递归获取子评论的 id 列表
func (this *Comment) replies(pid any, ctx *gin.Context) (ids []int) {

	// 获取请求参数
	params := this.params(ctx)

	var result []model.Comment
	mold := facade.DB.Model(&result).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
	mold.Where("pid", pid).Column("id")
	mold.Select()

	for _, val := range result {
		ids = append(ids, val.Id)
		ids = append(ids, this.replies(val.Id, ctx)...)
	}

	return ids
}

// flat 扁平化数据
func (this *Comment) flat(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"page":        1,
		"bind_type":   "article",
		"order":       "create_time desc",
	})

	if utils.Is.Empty(params["bind_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "bind_id 不能为空！"), 400)
		return
	}

	// 表数据结构体
	table := model.Comment{}
	// 允许查询的字段
	allow := []any{"bind_id", "bind_type"}
	// 动态给结构体赋值
	for key, val := range params {
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			utils.Struct.Set(&table, key, val)
		}
	}

	page := cast.ToInt(params["page"])
	limit := this.meta.limit(ctx)
	var result []model.Comment
	mold := facade.DB.Model(&result).Where("pid", 0).OnlyTrashed(cast.ToBool(params["onlyTrashed"])).WithTrashed(cast.ToBool(params["withTrashed"]))
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

		list := cast.ToSlice(data)

		// 根据 pid 递归获取子评论的 id 列表
		for key, val := range list {

			ids := this.replies(cast.ToStringMap(val)["id"], ctx)
			replies := facade.DB.Model(&[]model.Comment{}).WhereIn("id", ids).Order("create_time asc").Select()

			// 排除字段
			cast.ToStringMap(list[key])["replies"] = utils.ArrayMapWithField(replies, params["field"])
		}

		data = list

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

// config 配置
func (this *Comment) config(key ...any) (json map[string]any) {

	var config map[string]any

	// 缓存名称
	cacheName := "config[ARTICLE]"
	// 是否开启了缓存
	cacheState := cast.ToBool(facade.CacheToml.Get("open"))

	// 检查缓存是否存在
	if cacheState && facade.Cache.Has(cacheName) {

		config = cast.ToStringMap(facade.Cache.Get(cacheName))

	} else {

		config = facade.DB.Model(&model.Config{}).Where("key", "ARTICLE").Find()
		// 存储到缓存中
		if cacheState {
			go facade.Cache.Set(cacheName, config)
		}
	}

	if len(key) > 0 {
		return cast.ToStringMap(cast.ToStringMap(config["json"])[cast.ToString(key[0])])
	}

	return cast.ToStringMap(config["json"])
}