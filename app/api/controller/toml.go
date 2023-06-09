package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gopkg.in/gomail.v2"
	"inis/app/facade"
	"regexp"
	"strings"
)

type Toml struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Toml) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"log":     this.getLog,
		"sms":     this.getSMS,
		"cache":   this.getCache,
		"crypt":   this.getCrypt,
		"storage": this.getStorage,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Toml) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"test-email": this.testEmail,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPUT - PUT请求本体
func (this *Toml) IPUT(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"sms":     	   this.putSMS,
		"sms-email":   this.putSMSEmail,
		"sms-aliyun":  this.putSMSAliyun,
		"sms-tencent": this.putSMSTencent,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IDEL - DELETE请求本体
func (this *Toml) IDEL(ctx *gin.Context) {
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
func (this *Toml) INDEX(ctx *gin.Context) {
	this.json(ctx, nil, facade.Lang(ctx, "没什么用！"), 202)
}

// getSMS - 获取SMS服务配置
func (this *Toml) getSMS(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 允许的查询范围
	field := []any{"email", "aliyun", "tencent"}

	item := facade.SMSToml
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "SMS配置文件异常！"), 400)
		return
	}

	// 获取全部
	if utils.Is.Empty(params["name"]) {
		this.json(ctx, item.Result, facade.Lang(ctx, "数据请求成功！"), 200)
		return
	}

	if !utils.In.Array(params["name"], field) {
		this.json(ctx, nil, facade.Lang(ctx, "不允许的查询范围！"), 400)
		return
	}

	result := cast.ToStringMap(item.Get(cast.ToString(params["name"])))
	result["default"] = item.Get("default")

	// 获取指定
	this.json(ctx, result, facade.Lang(ctx, "数据请求成功！"), 200)
}

// putSMS - 修改SMS服务配置
func (this *Toml) putSMS(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["name"]) {
		this.json(ctx, nil, facade.Lang(ctx, "配置名称不能为空！"), 400)
		return
	}

	// 允许的修改范围
	field := []any{"default", "email", "aliyun", "tencent"}

	if !utils.In.Array(params["name"], field) {
		this.json(ctx, nil, facade.Lang(ctx, "不允许的修改范围！"), 400)
		return
	}

	switch params["name"] {
	case "default":
		this.putSMSDefault(ctx)
	case "email":
		this.putSMSEmail(ctx)
	case "aliyun":
		this.putSMSAliyun(ctx)
	case "tencent":
		this.putSMSTencent(ctx)
	default:
		this.putSMSDefault(ctx)
	}
}

// getCache - 获取缓存服务配置
func (this *Toml) getCache(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 允许的查询范围
	field := []any{"redis"}

	item := facade.CacheToml
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "缓存配置文件异常！"), 400)
		return
	}

	// 获取全部
	if utils.Is.Empty(params["name"]) {
		this.json(ctx, item.Result, facade.Lang(ctx, "数据请求成功！"), 200)
		return
	}

	if !utils.In.Array(params["name"], field) {
		this.json(ctx, nil, facade.Lang(ctx, "不允许的查询范围！"), 400)
		return
	}

	// 获取指定
	this.json(ctx, item.Get(cast.ToString(params["name"])), facade.Lang(ctx, "数据请求成功！"), 200)
}

// getCrypt - 获取加密服务配置
func (this *Toml) getCrypt(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 允许的查询范围
	field := []any{"jwt"}

	item := facade.CryptToml
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "加密配置文件异常！"), 400)
		return
	}

	// 获取全部
	if utils.Is.Empty(params["name"]) {
		this.json(ctx, item.Result, facade.Lang(ctx, "数据请求成功！"), 200)
		return
	}

	if !utils.In.Array(params["name"], field) {
		this.json(ctx, nil, facade.Lang(ctx, "不允许的查询范围！"), 400)
		return
	}

	// 获取指定
	this.json(ctx, item.Get(cast.ToString(params["name"])), facade.Lang(ctx, "数据请求成功！"), 200)
}

// getStorage - 获取存储服务配置
func (this *Toml) getStorage(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 允许的查询范围
	field := []any{"local", "oss", "cos", "kodo"}

	item := facade.StorageToml
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "存储配置文件异常！"), 400)
		return
	}

	// 获取全部
	if utils.Is.Empty(params["name"]) {
		this.json(ctx, item.Result, facade.Lang(ctx, "数据请求成功！"), 200)
		return
	}

	if !utils.In.Array(params["name"], field) {
		this.json(ctx, nil, facade.Lang(ctx, "不允许的查询范围！"), 400)
		return
	}

	// 获取指定
	this.json(ctx, item.Get(cast.ToString(params["name"])), facade.Lang(ctx, "数据请求成功！"), 200)
}

// getStorage - 获取日志服务配置
func (this *Toml) getLog(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 允许的查询范围
	var field []any

	item := facade.LogToml
	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "日志配置文件异常！"), 400)
		return
	}

	// 获取全部
	if utils.Is.Empty(params["name"]) {
		this.json(ctx, item.Result, facade.Lang(ctx, "数据请求成功！"), 200)
		return
	}

	if !utils.In.Array(params["name"], field) {
		this.json(ctx, nil, facade.Lang(ctx, "不允许的查询范围！"), 400)
		return
	}

	// 获取指定
	this.json(ctx, item.Get(cast.ToString(params["name"])), facade.Lang(ctx, "数据请求成功！"), 200)
}

// putSMSDefault - 修改SMS默认配置
func (this *Toml) putSMSDefault(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["value"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "value"), 400)
		return
	}

	allow := []any{"email", "aliyun", "tencent"}

	if !utils.In.Array(params["value"], allow) {
		this.json(ctx, nil, facade.Lang(ctx, "value 只允许是 email、aliyun、tencent 其中一个！"), 400)
		return
	}

	temp := facade.TempSMS
	temp = utils.Replace(temp, map[string]any{
		"${default}": params["value"],
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.SMSToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/sms.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// putSMSEmail - 修改SMS邮箱配置
func (this *Toml) putSMSEmail(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["host"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "host"), 400)
		return
	}

	if utils.Is.Empty(params["port"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "port"), 400)
		return
	}

	if !utils.Is.Number(params["port"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 只能是数字！", "port"), 400)
		return
	}

	if utils.Is.Empty(params["account"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "account"), 400)
		return
	}

	if utils.Is.Empty(params["password"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "password"), 400)
		return
	}

	if utils.Is.Empty(params["sign_name"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "sign_name"), 400)
		return
	}

	temp := facade.TempSMS
	temp = utils.Replace(temp, map[string]any{
		"${email.host}":      params["host"],
		"${email.port}":      cast.ToInt(params["port"]),
		"${email.account}":   params["account"],
		"${email.password}":  params["password"],
		"${email.nickname}":  params["nickname"],
		"${email.sign_name}": params["sign_name"],
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.SMSToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/sms.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// putSMSAliyun - 修改阿里云短信服务配置
func (this *Toml) putSMSAliyun(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["access_key_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "access_key_id"), 400)
		return
	}

	if utils.Is.Empty(params["access_key_secret"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "access_key_secret"), 400)
		return
	}

	if utils.Is.Empty(params["endpoint"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "endpoint"), 400)
		return
	}

	if utils.Is.Empty(params["sign_name"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "sign_name"), 400)
		return
	}

	if utils.Is.Empty(params["verify_code"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "verify_code"), 400)
		return
	}

	temp := facade.TempSMS
	temp = utils.Replace(temp, map[string]any{
		"${aliyun.access_key_id}":     params["access_key_id"],
		"${aliyun.access_key_secret}": params["access_key_secret"],
		"${aliyun.endpoint}":  		   params["endpoint"],
		"${aliyun.sign_name}":  	   params["sign_name"],
		"${aliyun.verify_code}": 	   params["verify_code"],
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.SMSToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/sms.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// putSMSTencent - 修改腾讯云短信服务配置
func (this *Toml) putSMSTencent(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["secret_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "secret_id"), 400)
		return
	}

	if utils.Is.Empty(params["secret_key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "secret_key"), 400)
		return
	}

	if utils.Is.Empty(params["endpoint"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "endpoint"), 400)
		return
	}

	if utils.Is.Empty(params["sms_sdk_app_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "sms_sdk_app_id"), 400)
		return
	}

	if utils.Is.Empty(params["sign_name"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "sign_name"), 400)
		return
	}

	if utils.Is.Empty(params["verify_code"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "verify_code"), 400)
		return
	}

	if utils.Is.Empty(params["region"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "region"), 400)
		return
	}

	temp := facade.TempSMS
	temp = utils.Replace(temp, map[string]any{
		"${tencent.secret_id}":      params["secret_id"],
		"${tencent.secret_key}": 	 params["secret_key"],
		"${tencent.endpoint}":  	 params["endpoint"],
		"${tencent.sms_sdk_app_id}": params["sms_sdk_app_id"],
		"${tencent.sign_name}":		 params["sign_name"],
		"${tencent.verify_code}":	 params["verify_code"],
		"${tencent.region}":		 params["region"],
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.SMSToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/sms.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// testEmail - 测试邮件服务
func (this *Toml) testEmail(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["host"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "host"), 400)
		return
	}

	if utils.Is.Empty(params["port"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "port"), 400)
		return
	}

	if !utils.Is.Number(params["port"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 只能是数字！", "port"), 400)
		return
	}

	if utils.Is.Empty(params["account"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "account"), 400)
		return
	}

	if utils.Is.Empty(params["password"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "password"), 400)
		return
	}

	if utils.Is.Empty(params["sign_name"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "sign_name"), 400)
		return
	}

	if utils.Is.Empty(params["email"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "email"), 400)
		return
	}

	if !utils.Is.Email(params["email"]) {
		this.json(ctx, nil, facade.Lang(ctx, "邮箱格式不正确！"), 400)
		return
	}

	client := gomail.NewDialer(
		cast.ToString(params["host"]),
		cast.ToInt(params["port"]),
		cast.ToString(params["account"]),
		cast.ToString(params["password"]),
	)

	item := gomail.NewMessage()
	nickname := cast.ToString(params["nickname"])
	account  := cast.ToString(params["account"])
	item.SetHeader("From", nickname+"<"+account+">")
	// 发送给多个用户
	item.SetHeader("To", cast.ToString(params["email"]))
	// 设置邮件主题
	item.SetHeader("Subject", cast.ToString(params["sign_name"]))
	// 设置邮件正文
	item.SetBody("text/html", "当您收到此封邮件时，说明您的邮件服务配置正确！")

	// 发送邮件
	err := client.DialAndSend(item)

	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试邮件发送失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "测试邮件发送成功！"), 200)
}