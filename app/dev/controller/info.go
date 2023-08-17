package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"os"
	"os/exec"
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
		"time"   : this.time,
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

// system - 系统信息
func (this *Info) system(ctx *gin.Context) {

	// 内存信息
	var memory map[string]any

	vm, err := mem.VirtualMemory()
	if err == nil {
		memory = map[string]any{
			"free" : vm.Free,
			"used" : vm.Used,
			"total": vm.Total,
		}
	}

	info := map[string]any{
		"path"  : utils.Get.Pwd(),
		"pid"   : utils.Get.Pid(),
		"port"  : map[string]any{
			"run" : this.get(ctx, "port"),
			"real": facade.AppToml.Get("app.port"),
		},
		"memory": memory,
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
		"text": "这是一个奇迹！",
	}, facade.Lang(ctx, "好的！"), 200)
}

// device - 设备信息
func (this *Info) device(ctx *gin.Context) {

	item := facade.Comm.Device()

	if item.Error != nil {
		this.json(ctx, nil, item.Error.Error(), 400)
		return
	}

	this.json(ctx, item.Json["data"], facade.Lang(ctx, "好的！"), 200)
}

// time - 时间信息
func (this *Info) time(ctx *gin.Context) {
	this.json(ctx, map[string]any{
		"unix": time.Now().Unix(),
		"date": time.Now().Format("2006-01-02 15:04:05"),
	}, facade.Lang(ctx, "好的！"), 200)
}

// renew - 更新
func (this *Info) renew(ctx *gin.Context) {

	path, err := os.Executable()
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	args := os.Args[1:]

	_, err = os.StartProcess(path, append([]string{path}, args...), &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "好的！"), 200)
}

// kill 杀死进程
func (this *Info) kill(ctx *gin.Context) {

	time.Sleep(5 * time.Second)

	// 根据操作系统选择不同的命令
	var cmd *exec.Cmd
	// kill -SIGHUP PID
	// kill -HUP pid
	cmd = exec.Command("kill", "-SIGHUP", cast.ToString(utils.Get.Pid()))
	// cmd = exec.Command("kill", "-SIGHUP", cast.ToString(utils.Get.Pid()))
	// cmd = exec.Command("taskkill", "/F", "/PID", cast.ToString(utils.Get.Pid()))
	// 守护进程
	// nohup /www/wwwroot/inis.cn/inis 1>/dev/null 2>&1 &

	// 执行命令
	err := cmd.Run()
	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "关闭进程失败：%v", err.Error()), 400)
		os.Exit(1)
		return
	}
}