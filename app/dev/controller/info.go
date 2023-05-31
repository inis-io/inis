package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"net"
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

// system - 系统信息
func (this *Info) system(ctx *gin.Context) {

	// 获取本机mac地址
	var mac string
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				mac = addr.String()
				break
			}
		}
	}
	this.json(ctx, map[string]any{
		"mac"   : mac,
		"ip"    : this.get(ctx, "ip"),
		"domain": this.get(ctx, "domain"),
		"GOOS"  : runtime.GOOS,
		"GOARCH": runtime.GOARCH,
		"GOROOT": runtime.GOROOT(),
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