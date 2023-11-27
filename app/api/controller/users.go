package controller

import (
	"fmt"
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

type Users struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Users) IGET(ctx *gin.Context) {
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
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Users) IPOST(ctx *gin.Context) {

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
func (this *Users) IPUT(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"update":  this.update,
		"restore": this.restore,
		"email":   this.email,
		"phone":   this.phone,
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
func (this *Users) IDEL(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"remove":  this.remove,
		"delete":  this.delete,
		"clear":   this.clear,
		"destroy": this.destroy,
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
func (this *Users) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "没什么用！"), 202)
}

// 删除缓存
func (this *Users) delCache() {
	// 删除缓存
	facade.Cache.DelTags([]any{"[GET]", "users"})
}

// one 获取指定数据
func (this *Users) one(ctx *gin.Context) {

	code := 204
	msg  := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx)

	// 表数据结构体
	table := model.Users{}
	// 允许查询的字段
	allow := []any{"id", "email"}
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

		mold.WithoutField("password")

		user := this.user(ctx)
		// 越权 - 既没有管理权限，也不是自己的数据
		if !this.meta.root(ctx) && (table.Id != user.Id || user.Id == 0) {
			mold.WithoutField("account", "email", "phone")
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
func (this *Users) all(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 获取请求参数
	params := this.params(ctx, map[string]any{
		"page":  1,
		"order": "create_time desc",
	})

	// 表数据结构体
	table := model.Users{}
	// 允许查询的字段
	allow := []any{"source"}
	// 动态给结构体赋值
	for key, val := range params {
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			utils.Struct.Set(&table, key, val)
		}
	}

	page := cast.ToInt(params["page"])
	limit := this.meta.limit(ctx)
	var result []model.Users
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

		mold.WithoutField("password")

		// 越权 - 没有管理权限
		if !this.meta.root(ctx) {
			mold.WithoutField("account", "email", "phone")
		}

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
func (this *Users) rand(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 限制最大数量
	limit := this.meta.limit(ctx)

	// 排除的 id 列表
	except := utils.Unity.Ids(params["except"])

	onlyTrashed := cast.ToBool(params["onlyTrashed"])
	withTrashed := cast.ToBool(params["withTrashed"])

	item := facade.DB.Model(&model.Users{}).OnlyTrashed(onlyTrashed).WithTrashed(withTrashed)
	if !utils.Is.Empty(except) {
		item = item.Where("id", "NOT IN", except)
	}

	// 从全部的 id 中随机选取指定数量的 id
	ids := utils.Rand.Slice(utils.Unity.Ids(item.Column("id")), limit)

	// 查询条件
	mold := facade.DB.Model(&[]model.Users{}).Where("id", "IN", ids)
	mold.OnlyTrashed(onlyTrashed).WithTrashed(withTrashed).IWhere(params["where"]).IOr(params["or"])
	mold.ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])
	mold.WithoutField("password")

	// 越权 - 没有管理权限
	if !this.meta.root(ctx) {
		mold.WithoutField("account", "email", "phone")
	}

	// 查询并打乱顺序
	data := utils.Array.MapWithField(utils.Rand.MapSlice(mold.Select()), params["field"])

	if utils.Is.Empty(data) {
		this.json(ctx, nil, facade.Lang(ctx, "无数据！"), 204)
		return
	}

	this.json(ctx, data, facade.Lang(ctx, "数据请求成功！"), 200)
}

// save 保存数据 - 包含创建和更新
func (this *Users) save(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.create(ctx)
	} else {
		this.update(ctx)
	}
}

// create 创建数据
func (this *Users) create(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)
	// 验证器
	err := validator.NewValid("users", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.Users{CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	allow := []any{"account", "password", "nickname", "email", "phone", "avatar", "description", "source", "remark", "title", "gender", "json", "text"}

	if utils.Is.Empty(params["email"]) {
		this.json(ctx, nil, facade.Lang(ctx, "邮箱不能为空！"), 400)
		return
	}

	// 动态给结构体赋值
	for key, val := range params {
		// 加密密码
		if key == "password" {
			val = utils.Password.Create(params["password"])
		}
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

	// 创建用户
	tx := facade.DB.Model(&table).Create(&table)

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	this.json(ctx, gin.H{ "id": table.Id }, facade.Lang(ctx, "创建成功！"), 200)
}

// update 更新数据
func (this *Users) update(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "id"), 400)
		return
	}

	// 验证器
	err := validator.NewValid("users", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 表数据结构体
	table := model.Users{}
	allow := []any{"id", "account", "password", "nickname", "avatar", "description", "gender", "json", "text"}
	async := utils.Async[map[string]any]()

	root := this.meta.root(ctx)
	// 越权 - 增加可选字段
	if root {
		allow = append(allow, "source", "remark", "title", "email", "phone")
	}

	// 动态给结构体赋值
	for key, val := range params {
		// 加密密码
		if key == "password" {
			// 密码为空时不更新此项
			if utils.Is.Empty(val) {
				continue
			}
			val = utils.Password.Create(params["password"])
		}
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

	// 越权 - 既没有管理权限，也不是自己的数据
	if !root && cast.ToInt(params["id"]) != this.user(ctx).Id {
		this.json(ctx, nil, facade.Lang(ctx, "无权限！"), 403)
		return
	}

	// 更新用户
	tx := facade.DB.Model(&table).WithTrashed().Where("id", params["id"]).Scan(&table).Update(async.Result())

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	// 删除缓存
	facade.Cache.Del(fmt.Sprintf("user[%v]", params["id"]))

	this.json(ctx, gin.H{ "id": table.Id }, facade.Lang(ctx, "更新成功！"), 200)
}

// count 统计数据
func (this *Users) count(ctx *gin.Context) {

	// 表数据结构体
	table := model.Users{}
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table)
	item.IWhere(params["where"]).IOr(params["or"]).ILike(params["like"]).INot(params["not"]).INull(params["null"]).INotNull(params["notNull"])

	this.json(ctx, item.Count(), facade.Lang(ctx, "查询成功！"), 200)
}

// sum 求和
func (this *Users) sum(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.Users
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
func (this *Users) min(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.Users
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
func (this *Users) max(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table model.Users
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
func (this *Users) column(ctx *gin.Context) {

	code := 204
	msg := []string{"无数据！", ""}
	var data any

	// 表数据结构体
	var table []model.Users
	// 获取请求参数
	params := this.params(ctx)

	item := facade.DB.Model(&table).Order(params["order"])
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

		item.WithoutField("password")

		// 越权 - 没有管理权限
		if !this.meta.root(ctx) {
			item.WithoutField("account", "email", "phone")
		}

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
func (this *Users) remove(ctx *gin.Context) {

	// 表数据结构体
	table := model.Users{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	if utils.In.Array(this.meta.user(ctx).Id, ids) {
		this.json(ctx, nil, facade.Lang(ctx, "不能删除自己！"), 400)
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
func (this *Users) delete(ctx *gin.Context) {

	// 表数据结构体
	table := model.Users{}
	// 获取请求参数
	params := this.params(ctx)

	// id 数组 - 参数归一化
	ids := utils.Unity.Ids(params["ids"])

	if utils.Is.Empty(ids) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "ids"), 400)
		return
	}

	if utils.In.Array(this.meta.user(ctx).Id, ids) {
		this.json(ctx, nil, facade.Lang(ctx, "不能删除自己！"), 400)
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
func (this *Users) clear(ctx *gin.Context) {

	// 表数据结构体
	table := model.Users{}

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
func (this *Users) restore(ctx *gin.Context) {

	// 表数据结构体
	table := model.Users{}
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

// email 修改邮箱
func (this *Users) email(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["email"]) {
		this.json(ctx, nil, facade.Lang(ctx, "邮箱不能为空！"), 400)
		return
	}

	user := this.meta.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	// 驱动
	drive := cast.ToString(facade.SMSToml.Get("drive.email"))

	if utils.Is.Empty(drive) {
		this.json(ctx, nil, facade.Lang(ctx, "管理员未开启邮箱服务，无法发送验证码！"), 400)
		return
	}

	// 从数据库里面找一下这个邮箱是否已经存在
	exist := facade.DB.Model(&model.Users{}).Where("email", params["email"]).Where("id", "!=", user.Id).Exist()
	if exist {
		this.json(ctx, nil, facade.Lang(ctx, "该邮箱已绑定其它账号！"), 400)
		return
	}

	// 缓存名称
	cacheName := fmt.Sprintf("%v-%v", drive, params["email"])

	// 验证码为空，发送验证码
	if utils.Is.Empty(params["code"]) {

		sms := facade.NewSMS(drive).VerifyCode(params["email"])
		if sms.Error != nil {
			this.json(ctx, nil, sms.Error.Error(), 400)
			return
		}

		// 缓存验证码 - 5分钟
		go facade.Cache.Set(cacheName, sms.VerifyCode, 5 * time.Minute)

		msg := fmt.Sprintf("验证码发送至您的邮箱：%s，请注意查收！", params["email"])
		this.json(ctx, nil, facade.Lang(ctx, msg), 201)
		return
	}

	// 获取缓存里面的验证码
	cacheCode := facade.Cache.Get(cacheName)

	if cast.ToString(params["code"]) != cast.ToString(cacheCode) {
		this.json(ctx, nil, facade.Lang(ctx, "验证码错误！"), 400)
		return
	}

	// 更新邮箱
	tx := facade.DB.Model(&model.Users{}).Where("id", user.Id).UpdateColumn("email", params["email"])
	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	// 删除验证码
	go facade.Cache.Del(cacheName)

	this.json(ctx, gin.H{ "id": user.Id }, facade.Lang(ctx, "修改成功！"), 200)
}

// phone 修改手机号
func (this *Users) phone(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["phone"]) {
		this.json(ctx, nil, facade.Lang(ctx, "手机号不能为空！"), 400)
		return
	}

	user := this.meta.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	// 驱动
	drive := cast.ToString(facade.SMSToml.Get("drive.sms"))

	if utils.Is.Empty(drive) {
		this.json(ctx, nil, facade.Lang(ctx, "管理员未开启短信服务，无法发送验证码！"), 400)
		return
	}

	// 从数据库里面找一下这个手机号是否已经存在
	exist := facade.DB.Model(&model.Users{}).Where("phone", params["phone"]).Where("id", "!=", user.Id).Exist()
	if exist {
		this.json(ctx, nil, facade.Lang(ctx, "该邮箱已绑定其它账号！"), 400)
		return
	}

	// 缓存名称
	cacheName := fmt.Sprintf("%v-%v", drive, params["phone"])

	// 验证码为空，发送验证码
	if utils.Is.Empty(params["code"]) {

		sms := facade.NewSMS(drive).VerifyCode(params["phone"])
		if sms.Error != nil {
			this.json(ctx, nil, sms.Error.Error(), 400)
			return
		}

		// 缓存验证码 - 5分钟
		go facade.Cache.Set(cacheName, sms.VerifyCode, 5 * time.Minute)

		msg := fmt.Sprintf("验证码发送至您的手机：%s，请注意查收！", params["phone"])
		this.json(ctx, nil, facade.Lang(ctx, msg), 201)
		return
	}

	// 获取缓存里面的验证码
	cacheCode := facade.Cache.Get(cacheName)

	if cast.ToString(params["code"]) != cast.ToString(cacheCode) {
		this.json(ctx, nil, facade.Lang(ctx, "验证码错误！"), 400)
		return
	}

	// 更新手机号
	tx := facade.DB.Model(&model.Users{}).Where("id", user.Id).UpdateColumn("phone", params["phone"])
	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	// 删除验证码
	go facade.Cache.Del(cacheName)

	this.json(ctx, gin.H{ "id": user.Id }, facade.Lang(ctx, "修改成功！"), 200)
}

// 注销 - 邮箱、手机号
func (this *Users) destroy(ctx *gin.Context) {

	table := model.Users{}
	params := this.params(ctx, map[string]any{
		"source": "default",
	})

	user := this.meta.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	var social string
	social = utils.Ternary(utils.Is.Email(user.Email), "email", social)
	social = utils.Ternary(utils.Is.Phone(user.Phone), "phone", social)

	if utils.Is.Empty(social) {
		this.json(ctx, nil, facade.Lang(ctx, "您未绑定手机或邮箱，无法验证注销安全性！"), 400)
		return
	}

	var contact string
	if social == "email" {
		contact = user.Email
	} else {
		contact = user.Phone
	}

	cacheName := fmt.Sprintf("[login][%v=%v]", social, contact)

	// 验证码为空 - 发送验证码
	if utils.Is.Empty(params["code"]) {

		drive := utils.Ternary(social == "email", "email", "sms")
		sms := facade.NewSMS(drive).VerifyCode(contact)
		if sms.Error != nil {
			this.json(ctx, nil, sms.Error.Error(), 400)
			return
		}
		// 缓存验证码 - 5分钟
		facade.Cache.Set(cacheName, sms.VerifyCode, 5 * time.Minute)
		this.json(ctx, nil, facade.Lang(ctx, "验证码发送成功！"), 201)
		return
	}

	// 获取缓存里面的验证码
	cacheCode := facade.Cache.Get(cacheName)

	if cast.ToString(params["code"]) != cacheCode {
		this.json(ctx, nil, facade.Lang(ctx, "验证码错误！"), 400)
		return
	}

	// 删除验证码
	go facade.Cache.Del(cacheName)

	// 清空数据
	(&model.Users{}).Destroy(user.Id)

	// 删除用户
	tx := facade.DB.Model(&table).Force().Delete(user.Id)

	if tx.Error != nil {
		this.json(ctx, nil, tx.Error.Error(), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "注销成功！"), 200)
}