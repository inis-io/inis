package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"strings"
)

type Device struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Device) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"user": this.users,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Device) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"bind": this.bind,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPUT - PUT请求本体
func (this *Device) IPUT(ctx *gin.Context) {
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
func (this *Device) IDEL(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"bind": this.unbind,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// INDEX - GET请求本体
func (this *Device) INDEX(ctx *gin.Context) {

	this.json(ctx, nil, facade.Lang(ctx, "好的！"), 200)
}

// bind - 绑定设备
func (this *Device) bind(ctx *gin.Context) {

	// 获取参数
	params := this.params(ctx)

	user := this.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
			this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	if !this.meta.root(ctx) {
		this.json(ctx, nil, facade.Lang(ctx, "无权限！"), 403)
		return
	}

	if utils.Is.Empty(params["account"]) {
		this.json(ctx, nil, facade.Lang(ctx, "账号不能为空！"), 400)
		return
	}

	if utils.Is.Empty(params["password"]) {
		this.json(ctx, nil, facade.Lang(ctx, "密码不能为空！"), 400)
		return
	}

	item := utils.Curl(utils.CurlRequest{
		Url:    facade.Uri + "/api/comm/login",
		Method: "POST",
		Body: params,
	}).Send()

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "远程服务器错误：%v", item.Error.Error()), 500)
		return
	}

	if cast.ToInt(item.Json["code"]) != 200 {
		this.json(ctx, item.Json["data"], item.Json["msg"], item.Json["code"])
		return
	}

	// 登录社区返回的信息
	res := cast.ToStringMap(item.Json["data"])

	headers := facade.Comm.Signature(nil)
	headers["Authorization"] = res["token"]

	// 绑定设备
	item = utils.Curl(utils.CurlRequest{
		Url:    facade.Uri + "/sn/device/bind",
		Method: "POST",
		Headers: headers,
	}).Send()

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "远程服务器错误：%v", item.Error.Error()), 500)
		return
	}

	if cast.ToInt(item.Json["code"]) != 200 {
		this.json(ctx, item.Json["data"], item.Json["msg"], item.Json["code"])
		return
	}

	this.json(ctx, gin.H{ "user": res["user"] }, facade.Lang(ctx, "绑定成功！"), 200)
}

// unbind - 解绑设备
func (this *Device) unbind(ctx *gin.Context) {

	user := this.user(ctx)
	// 即便中间件已经校验过登录了，这里还进行二次校验是未了防止接口权限被改，而 uid 又是强制的，从而导致的意外情况
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	if !this.meta.root(ctx) {
		this.json(ctx, nil, facade.Lang(ctx, "无权限！"), 403)
		return
	}

	// 绑定设备
	item := utils.Curl(utils.CurlRequest{
		Url:    facade.Uri + "/sn/device/bind",
		Method: "DELETE",
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

	this.json(ctx, nil, facade.Lang(ctx, "解绑成功！"), 200)
}

// user - 获取用户信息
func (this *Device) users(ctx *gin.Context) {

	item := utils.Curl(utils.CurlRequest{
		Url:    facade.Uri + "/sn/users/info",
		Method: "GET",
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