package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"strings"
)

type SystemVersion struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *SystemVersion) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"next":     this.next,
		"download": this.download,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *SystemVersion) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPUT - PUT请求本体
func (this *SystemVersion) IPUT(ctx *gin.Context) {
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
func (this *SystemVersion) IDEL(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// INDEX - GET请求本体
func (this *SystemVersion) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "好的！"), 200)
}

// next - 获取下个版本
func (this *SystemVersion) next(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	item := utils.Curl(utils.CurlRequest{
		Url:    facade.Uri + "/sn/system-version/next",
		Method: "GET",
		Headers: facade.Comm.Signature(params),
		Query:   params,
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

// download - 获取下载地址
func (this *SystemVersion) download(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "id"), 400)
		return
	}

	user := this.meta.user(ctx)
	if user.Id == 0 {
		this.json(ctx, nil, facade.Lang(ctx, "请先登录！"), 401)
		return
	}

	// user 转 map
	json := cast.ToStringMap(utils.Json.Decode(utils.Json.Encode(user)))
	body := map[string]any{
		"ip"   : this.get(ctx, "ip"),
		"id"   : params["id"],
		"agent": this.header(ctx, "User-Agent"),
		"json" : utils.Map.WithField(json, []string{"id", "nickname", "email", "phone"}),
	}

	item := utils.Curl(utils.CurlRequest{
		Url:    facade.Uri + "/sn/system-version/download",
		Method: "GET",
		Headers: facade.Comm.Signature(body),
		Body:   body,
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