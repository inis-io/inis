package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Upgrade struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Upgrade) IGET(ctx *gin.Context) {
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

// IPOST - POST请求本体
func (this *Upgrade) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"theme" : this.theme,
		"system": this.system,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPUT - PUT请求本体
func (this *Upgrade) IPUT(ctx *gin.Context) {
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

// IDEL - DELETE请求本体
func (this *Upgrade) IDEL(ctx *gin.Context) {
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
func (this *Upgrade) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "没什么用！"), 202)
}

// theme 主题升级
func (this *Upgrade) theme(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"files": []string{"static"},
	})

	if utils.Is.Empty(params["id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "主题版本ID"), 400)
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

	// 获取主题的下载地址
	curl := utils.Curl(utils.CurlRequest{
		Body:    body,
		Method:  "GET",
		Headers: facade.Comm.Signature(body),
		Url:     facade.Uri + "/sn/theme-version/download",
	}).Send()

	if curl.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "远程服务器错误：%v", curl.Error.Error()), 500)
		return
	}

	if cast.ToInt(curl.Json["code"]) != 200 {
		this.json(ctx, curl.Json["data"], curl.Json["msg"], curl.Json["code"])
		return
	}

	download := cast.ToStringMap(curl.Json["data"])

	// 检查 url 是否合法
	_, err := url.ParseRequestURI(cast.ToString(download["url"]))
	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不合法！", "下载地址"), 400)
		return
	}

	item := utils.File().Download(cast.ToString(download["url"]), "public/theme.zip")
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "下载失败：%v", item.Error.Error()), 400)
		return
	}

	// 批量删除文件和目录
	for _, path := range utils.Unity.Keys(params["files"]) {
		utils.File().Remove("public/" + cast.ToString(path))
	}

	// 解压文件
	item = utils.File().Dir("public").Name("theme.zip").UnZip()
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "解压失败：%v", item.Error.Error()), 400)
		return
	}

	// 删除压缩包
	utils.File().Remove("public/theme.zip")

	this.json(ctx, nil, facade.Lang(ctx, "升级成功！"), 200)
}

// system 系统升级
func (this *Upgrade) system(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)
	// 压缩包名称
	name := "system.zip"

	if utils.Is.Empty(params["id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "系统版本ID"), 400)
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

	// 获取系统版本的下载地址
	curl := utils.Curl(utils.CurlRequest{
		Url:    facade.Uri + "/sn/system-version/download",
		Method: "GET",
		Headers: facade.Comm.Signature(body),
		Body:   body,
	}).Send()

	if curl.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "远程服务器错误：%v", curl.Error.Error()), 500)
		return
	}

	if cast.ToInt(curl.Json["code"]) != 200 {
		this.json(ctx, curl.Json["data"], curl.Json["msg"], curl.Json["code"])
		return
	}

	download := cast.ToStringMap(curl.Json["data"])

	// 检查 url 是否合法
	_, err := url.ParseRequestURI(cast.ToString(download["url"]))
	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不合法！", "下载地址"), 400)
		return
	}

	item := utils.File().Download(cast.ToString(download["url"]), name)
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "下载失败：%v", item.Error.Error()), 400)
		return
	}

	// 解压文件
	item = utils.File().Dir("./").Name(name).UnZip()
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "解压失败：%v", item.Error.Error()), 400)
		return
	}

	// 删除压缩包
	utils.File().Remove(name)

	this.json(ctx, nil, facade.Lang(ctx, "升级成功！"), 200)

	// 杀死进程
	go this.kill(ctx)
}

// kill 杀死进程
func (this *Upgrade) kill(ctx *gin.Context) {

	time.Sleep(5 * time.Second)

	// 根据操作系统选择不同的命令
	var cmd *exec.Cmd
	cmd = exec.Command("taskkill", "/F", "/PID", cast.ToString(utils.Get.Pid()))
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