package controller

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/gin-gonic/gin"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"runtime"
	"strings"
)

type Info struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Info) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"system" : this.system,
		"version": this.version,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Info) IPOST(ctx *gin.Context) {

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
func (this *Info) IPUT(ctx *gin.Context) {
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
func (this *Info) IDEL(ctx *gin.Context) {
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
func (this *Info) INDEX(ctx *gin.Context) {

	// params := this.params(ctx)

	system := map[string]any{
		"GOOS":   runtime.GOOS,
		"GOARCH": runtime.GOARCH,
		"GOROOT": runtime.GOROOT(),
		"NumCPU": runtime.NumCPU(),
		"NumGoroutine": runtime.NumGoroutine(),
		"go": utils.Version.Go(),
		"inis": facade.Version,
		"compare": utils.Version.Compare("v1.0.0", "1 2 0"),
		"agent":  this.header(ctx, "User-Agent"),
	}

	this.json(ctx, map[string]any{
		"system": system,
	}, facade.Lang(ctx, "好的！"), 200)
}

// sn - 获取机器码
func sn() (result string) {

	result, err := machineid.ID()
	if err != nil {
		return utils.Get.Mac()
	}

	return result
}

// system - 系统信息
func (this *Info) system(ctx *gin.Context) {

	this.json(ctx, map[string]any{
		"sn"    : sn(),
		"mac"   : utils.Get.Mac(),
		"port"  : facade.H{
			"run" : this.get(ctx, "port"),
			"real": facade.AppToml.Get("app.port"),
		},
		"domain": this.get(ctx, "domain"),
		"GOOS"  : runtime.GOOS,
		"GOARCH": runtime.GOARCH,
		"NumCPU": runtime.NumCPU(),
		"NumGoroutine": runtime.NumGoroutine(),
	}, facade.Lang(ctx, "好的！"), 200)
}

// version - 版本信息
func (this *Info) version(ctx *gin.Context) {
	this.json(ctx, map[string]any{
		"go": utils.Version.Go(),
		"inis": facade.Version,
	}, facade.Lang(ctx, "好的！"), 200)
}