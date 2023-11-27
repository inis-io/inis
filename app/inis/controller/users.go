package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"strings"
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
		"access-token": this.accessToken,
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

	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPUT - PUT请求本体
func (this *Users) IPUT(ctx *gin.Context) {
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
func (this *Users) IDEL(ctx *gin.Context) {
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
func (this *Users) INDEX(ctx *gin.Context) {

	this.json(ctx, nil, facade.Lang(ctx, "好的！"), 200)
}

// accessToken - 获取访问令牌
func (this *Users) accessToken(ctx *gin.Context) {

	item := utils.Curl(utils.CurlRequest{
		Url:    facade.Uri + "/sn/users/access-token",
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