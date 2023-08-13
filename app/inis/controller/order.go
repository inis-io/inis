package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"strings"
)

type Order struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Order) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"theme": this.theme,
		"themes": this.themes,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Order) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPUT - PUT请求本体
func (this *Order) IPUT(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IDEL - DELETE请求本体
func (this *Order) IDEL(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// INDEX - GET请求本体
func (this *Order) INDEX(ctx *gin.Context) {

	this.json(ctx, nil, facade.Lang(ctx, "好的！"), 200)
}

// theme - 查询已购的指定主题
func (this *Order) theme(ctx *gin.Context) {

	// 获取请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "唯一识别码 key 不能为空！"), 400)
		return
	}

	// 绑定设备
	item := utils.Curl(utils.CurlRequest{
		Url:     facade.Uri + "/sn/order/theme",
		Method:  "GET",
		Query :  params,
		Headers: facade.Comm.Signature(params),
	}).Send()

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "远程服务器错误：%v", item.Error.Error()), 500)
		return
	}

	if cast.ToInt(item.Json["code"]) != 200 {
		this.json(ctx, item.Json["data"], item.Json["msg"], item.Json["code"])
		return
	}

	this.json(ctx, item.Json["data"], facade.Lang(ctx, cast.ToString(item.Json["msg"])), 200)
}

// themes - 查询已购的全部主题
func (this *Order) themes(ctx *gin.Context) {

	// 绑定设备
	item := utils.Curl(utils.CurlRequest{
		Url:     facade.Uri + "/sn/order/themes",
		Method:  "GET",
		Headers: facade.Comm.Signature(nil),
	}).Send()

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "远程服务器错误：%v", item.Error.Error()), 500)
		return
	}

	if cast.ToInt(item.Json["code"]) != 200 {
		this.json(ctx, item.Json["data"], item.Json["msg"], item.Json["code"])
		return
	}

	this.json(ctx, item.Json["data"], facade.Lang(ctx, "好的！"), 200)
}