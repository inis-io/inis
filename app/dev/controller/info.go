package controller

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"runtime"
	"strings"
	"time"
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
		"device" : this.device,
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

	info := map[string]any{
		"port"  : map[string]any{
			"run" : this.get(ctx, "port"),
			"real": facade.AppToml.Get("app.port"),
		},
		"domain": this.get(ctx, "domain"),
		"GOOS"  : runtime.GOOS,
		"GOARCH": runtime.GOARCH,
		"NumCPU": runtime.NumCPU(),
		"NumGoroutine": runtime.NumGoroutine(),
	}

	this.json(ctx, info, facade.Lang(ctx, "好的！"), 200)
}

// version - 版本信息
func (this *Info) version(ctx *gin.Context) {
	this.json(ctx, map[string]any{
		"go": utils.Version.Go(),
		"inis": facade.Version,
	}, facade.Lang(ctx, "好的！"), 200)
}

func (this *Info) device(ctx *gin.Context) {

	body := map[string]any{
		"sn"    : sn(),
		"mac"   : utils.Get.Mac(),
		"port"  : map[string]any{
			"run" : this.get(ctx, "port"),
			"real": facade.AppToml.Get("app.port"),
		},
		"domain": this.get(ctx, "domain"),
		"goos"  : runtime.GOOS,
		"goarch": runtime.GOARCH,
		"cpu"   : runtime.NumCPU(),
	}

	// X-SS-STUB(MD5) X-Argus(加密文本) X-Khronos(时间戳) X-Gorgon
	encode := facade.Cipher(facade.Hash.Token(body["sn"]), facade.Hash.Token(body["mac"]))

	unix := time.Now().Unix()

	item := utils.Curl(utils.CurlRequest{
		Method: "POST",
		Url   : "http://localhost:8642/api/test/request",
		Body  : body,
		Headers: map[string]any{
			"X-Khronos" : unix,
			"X-Argus"   : encode.Encrypt(utils.Json.Encode(body)).Text,
			"X-Gorgon"  : "8642" + facade.Hash.Token(uuid.New().String(), 48),
			"X-SS-STUB" : strings.ToUpper(facade.Hash.Token(utils.Map.ToURL(body), 32, unix)),
		},
	}).Send()

	if item.Error != nil {
		this.json(ctx, nil, item.Error.Error(), 500)
		return
	}

	this.json(ctx, item.Json["data"], facade.Lang(ctx, "好的！"), 200)
}