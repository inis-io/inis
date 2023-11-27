package controller

import (
	"context"
	"fmt"
	AliYunClient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	AliYunUtil "github.com/alibabacloud-go/openapi-util/service"
	AliYunUtilV2 "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	TencentCloud "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/unti-io/go-utils/utils"
	"gopkg.in/gomail.v2"
	"inis/app/facade"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
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
		"test-sms-email"  : this.testSMSEmail,
		"test-sms-aliyun" : this.testSMSAliyun,
		"test-sms-tencent": this.testSMSTencent,
		"test-redis"	  : this.testRedis,
		"test-oss"		  : this.testOSS,
		"test-cos"		  : this.testCOS,
		"test-kodo"	  	  : this.testKODO,
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
		"sms":     	       this.putSMS,
		"sms-email":       this.putSMSEmail,
		"sms-aliyun":      this.putSMSAliyun,
		"sms-tencent":     this.putSMSTencent,
		"sms-drive":       this.putSMSDrive,
		"crypt-jwt":       this.putCryptJWT,
		"cache-default":   this.putCacheDefault,
		"cache-redis":     this.putCacheRedis,
		"cache-file":      this.putCacheFile,
		"cache-ram":       this.putCacheRam,
		"storage-default": this.putStorageDefault,
		"storage-local":   this.putStorageLocal,
		"storage-oss":     this.putStorageOSS,
		"storage-cos":     this.putStorageCOS,
		"storage-kodo":    this.putStorageKODO,
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
	result["drive"] = item.Get("drive")

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
	case "drive":
		this.putSMSDrive(ctx)
	case "email":
		this.putSMSEmail(ctx)
	case "aliyun":
		this.putSMSAliyun(ctx)
	case "tencent":
		this.putSMSTencent(ctx)
	default:
		this.putSMSDrive(ctx)
	}
}

// getCache - 获取缓存服务配置
func (this *Toml) getCache(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// 允许的查询范围
	field := []any{"redis","file","ram"}

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

	result := cast.ToStringMap(item.Get(cast.ToString(params["name"])))
	result["open"]    = cast.ToBool(item.Get("open"))
	result["default"] = item.Get("default")

	// 获取指定
	this.json(ctx, result, facade.Lang(ctx, "数据请求成功！"), 200)
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

	result := cast.ToStringMap(item.Get(cast.ToString(params["name"])))
	result["default"] = item.Get("default")

	// 获取指定
	this.json(ctx, result, facade.Lang(ctx, "数据请求成功！"), 200)
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

// putSMSDrive - 修改SMS驱动配置
func (this *Toml) putSMSDrive(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"default": "email",
	})

	opts  := make(map[string]any)
	allow := []any{"email", "sms", "default"}

	for key, value := range params {
		if !utils.In.Array(key, allow) {
			continue
		}
		opts[fmt.Sprintf("${drive.%s}", key)] = value
	}

	temp := facade.TempSMS
	temp = utils.Replace(temp, opts)

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
func (this *Toml) testSMSEmail(ctx *gin.Context) {

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

// testSMSAliyun - 发送阿里云测试短信
func (this *Toml) testSMSAliyun(ctx *gin.Context) {

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

	if utils.Is.Empty(params["phone"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "phone"), 400)
		return
	}

	if !utils.Is.Phone(params["phone"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 格式不正确！", "phone"), 400)
		return
	}

	client, err := AliYunClient.NewClient(&AliYunClient.Config{
		// 访问的域名
		Endpoint: tea.String(cast.ToString(params["endpoint"])),
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(cast.ToString(params["access_key_id"])),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(cast.ToString(params["access_key_secret"])),
	})

	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试短信发送失败！"), 400)
		return
	}

	query := map[string]any{
		// 必填，接收短信的手机号码
		"PhoneNumbers": tea.String(cast.ToString(params["phone"])),
		// 必填，短信签名名称
		"SignName": tea.String(cast.ToString(params["sign_name"])),
		// 必填，短信模板ID
		"TemplateCode": tea.String(cast.ToString(params["verify_code"])),
	}

	query["TemplateParam"] = tea.String(utils.Json.Encode(map[string]any{
		"code": 6666,
	}))

	runtime := &AliYunUtilV2.RuntimeOptions{}
	request := &AliYunClient.OpenApiRequest{
		Query: AliYunUtil.Query(query),
	}

	// 返回值为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode
	result, err := client.CallApi((&facade.AliYunSMS{}).ApiInfo(), request, runtime)
	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试短信发送失败！"), 400)
		return
	}

	body := cast.ToStringMap(result["body"])
	if body["Code"] != "OK" {
		this.json(ctx, cast.ToString(body["Message"]), facade.Lang(ctx, "测试短信发送失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "测试短信发送成功！"), 200)
}

// testSMSTencent - 发送腾讯云测试短信
func (this *Toml) testSMSTencent(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	credential := common.NewCredential(
		cast.ToString(params["secret_id"]),
		cast.ToString(params["secret_key"]),
	)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = cast.ToString(params["endpoint"])
	client, err := TencentCloud.NewClient(
		credential,
		cast.ToString(params["region"]),
		clientProfile,
	)

	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试短信发送失败！"), 400)
		return
	}

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := TencentCloud.NewSendSmsRequest()

	request.PhoneNumberSet = common.StringPtrs([]string{cast.ToString(params["phone"])})
	request.SmsSdkAppId = common.StringPtr(cast.ToString(params["sms_sdk_app_id"]))
	request.SignName = common.StringPtr(cast.ToString(params["sign_name"]))
	request.TemplateId = common.StringPtr(cast.ToString(params["verify_code"]))
	request.TemplateParamSet = common.StringPtrs([]string{"6666"})

	item, err := client.SendSms(request)

	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试短信发送失败！"), 400)
		return
	}

	if item.Response == nil {
		this.json(ctx, "response is nil", facade.Lang(ctx, "测试短信发送失败！"), 400)
		return
	}

	if len(item.Response.SendStatusSet) == 0 {
		this.json(ctx, "response send status set is nil", facade.Lang(ctx, "测试短信发送失败！"), 400)
		return
	}

	if *item.Response.SendStatusSet[0].Code != "Ok" {
		this.json(ctx, item.Response.SendStatusSet[0].Message, facade.Lang(ctx, "测试短信发送失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "测试短信发送成功！"), 200)
}

// putCryptJWT - 修改JWT配置
func (this *Toml) putCryptJWT(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "key"), 400)
		return
	}

	if utils.Is.Empty(params["expire"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "expire"), 400)
		return
	}

	if utils.Is.Empty(params["issuer"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "issuer"), 400)
		return
	}

	if utils.Is.Empty(params["subject"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "subject"), 400)
		return
	}

	temp := facade.TempCrypt
	temp = utils.Replace(temp, map[string]any{
		"${jwt.key}":      params["key"],
		"${jwt.expire}":   params["expire"],
		"${jwt.issuer}":   params["issuer"],
		"${jwt.subject}":  params["subject"],
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.CryptToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/crypt.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// putCacheRedis - 修改Redis缓存配置
func (this *Toml) putCacheRedis(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"database": 0,
		"host":     "localhost",
		"port":     6379,
		"prefix":   "inis:",
		"expire":   "2 * 60 * 60",
	})

	if !utils.Is.Number(params["port"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 只能是数字！", "port"), 400)
		return
	}

	if !utils.Is.Number(params["database"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 只能是数字！", "database"), 400)
		return
	}

	temp := facade.TempCache
	temp = utils.Replace(temp, map[string]any{
		"${redis.host}":     params["host"],
		"${redis.port}":     params["port"],
		"${redis.database}": params["database"],
		"${redis.password}": params["password"],
		"${redis.prefix}":   params["prefix"],
		"${redis.expire}":   params["expire"],
		"${open}":           cast.ToBool(facade.CacheToml.Get("open")),
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.CacheToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/cache.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// putCacheFile - 修改File缓存配置
func (this *Toml) putCacheFile(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"path":     "runtime/cache",
		"prefix":   "inis_",
		"expire":   "2 * 60 * 60",
	})

	temp := facade.TempCache
	temp = utils.Replace(temp, map[string]any{
		"${file.path}":     params["path"],
		"${file.prefix}":   params["prefix"],
		"${file.expire}":   params["expire"],
		"${open}":          cast.ToBool(facade.CacheToml.Get("open")),
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.CacheToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/cache.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// putCacheRam - 修改Ram缓存配置
func (this *Toml) putCacheRam(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"expire":   "2 * 60 * 60",
	})

	temp := facade.TempCache
	temp = utils.Replace(temp, map[string]any{
		"${ram.expire}":   params["expire"],
		"${open}":         cast.ToBool(facade.CacheToml.Get("open")),
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.CacheToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/cache.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// testRedis - 测试Redis连接
func (this *Toml) testRedis(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"database": 0,
		"host":     "localhost",
		"port":     6379,
	})

	if !utils.Is.Number(params["port"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 只能是数字！", "port"), 400)
		return
	}

	if !utils.Is.Number(params["database"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 只能是数字！", "database"), 400)
		return
	}

	// 创建Redis连接客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", params["host"], cast.ToInt(params["port"])),
		DB:       cast.ToInt(params["database"]),
		Password: cast.ToString(params["password"]),
	})

	// Ping Redis
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试Redis连接失败！"), 400)
		return
	}

	this.json(ctx, pong, facade.Lang(ctx, "测试Redis连接成功！"), 200)
}

// putCacheDefault - 修改缓存默认服务类型
func (this *Toml) putCacheDefault(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"value": "file",
		"open":  "false",
	})

	allow := []any{"redis","file","ram"}

	if !utils.In.Array(params["value"], allow) {
		this.json(ctx, nil, facade.Lang(ctx, "value 只允许是 redis、file、ram ！"), 400)
		return
	}

	temp := facade.TempCache
	temp = utils.Replace(temp, map[string]any{
		"${open}":    utils.Ternary(cast.ToBool(params["open"]), "true", "false"),
		"${default}": params["value"],
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.CacheToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/cache.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// putStorageDefault - 修改存储默认服务类型
func (this *Toml) putStorageDefault(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["value"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "value"), 400)
		return
	}

	allow := []any{"local", "oss", "cos", "kodo"}

	if !utils.In.Array(params["value"], allow) {
		this.json(ctx, nil, facade.Lang(ctx, "value 只允许是 local、oss、cos、kodo 其中一个！"), 400)
		return
	}

	temp := facade.TempStorage
	temp = utils.Replace(temp, map[string]any{
		"${default}": params["value"],
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.StorageToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/storage.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// putStorageLocal - 修改本地存储配置
func (this *Toml) putStorageLocal(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"domain": this.get(ctx, "domain"),
	})

	temp := facade.TempStorage
	temp = utils.Replace(temp, map[string]any{
		"${local.domain}": params["domain"],
		"${local.path}"  : params["path"],
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.StorageToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/storage.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// testOSS - 测试OSS连接
func (this *Toml) testOSS(ctx *gin.Context) {

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

	id       := cast.ToString(params["access_key_id"])
	secret   := cast.ToString(params["access_key_secret"])
	endpoint := cast.ToString(params["endpoint"])

	client, err := oss.New(endpoint, id, secret)

	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试OSS连接失败！"), 400)
		return
	}

	exist, err := client.IsBucketExist(cast.ToString(params["bucket"]))
	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试OSS连接失败！"), 400)
		return
	}

	if !exist {
		this.json(ctx, nil, facade.Lang(ctx, "Bucket 不存在！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "测试OSS连接成功！"), 200)
}

// putStorageOSS - 修改OSS存储配置
func (this *Toml) putStorageOSS(ctx *gin.Context) {

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

	if utils.Is.Empty(params["bucket"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bucket"), 400)
		return
	}

	temp := facade.TempStorage
	temp = utils.Replace(temp, map[string]any{
		"${oss.access_key_id}": 	cast.ToString(params["access_key_id"]),
		"${oss.access_key_secret}": cast.ToString(params["access_key_secret"]),
		"${oss.endpoint}": 			cast.ToString(params["endpoint"]),
		"${oss.bucket}": 			cast.ToString(params["bucket"]),
		"${oss.domain}": 			cast.ToString(params["domain"]),
		"${oss.path}"  : 			cast.ToString(params["path"]),
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.StorageToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/storage.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// testCOS - 测试COS连接
func (this *Toml) testCOS(ctx *gin.Context) {

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

	if utils.Is.Empty(params["app_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "app_id"), 400)
		return
	}

	if utils.Is.Empty(params["bucket"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bucket"), 400)
		return
	}

	if utils.Is.Empty(params["region"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "region"), 400)
		return
	}

	appId     := cast.ToString(params["app_id"])
	secretId  := cast.ToString(params["secret_id"])
	secretKey := cast.ToString(params["secret_key"])
	bucket    := cast.ToString(params["bucket"])
	region    := cast.ToString(params["region"])

	BucketURL, err := url.Parse(fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", bucket, appId, region))
	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试COS连接失败！"), 400)
		return
	}

	client := cos.NewClient(&cos.BaseURL{
		BucketURL: BucketURL,
	}, &http.Client{
		// 设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})

	// 查询存储桶
	exist, err := client.Bucket.IsExist(context.Background())

	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试COS连接失败！"), 400)
		return
	}

	if !exist {
		this.json(ctx, nil, facade.Lang(ctx, "Bucket 不存在！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "测试COS连接成功！"), 200)
}

// putStorageCOS - 修改COS存储配置
func (this *Toml) putStorageCOS(ctx *gin.Context) {

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

	if utils.Is.Empty(params["app_id"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "app_id"), 400)
		return
	}

	if utils.Is.Empty(params["bucket"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bucket"), 400)
		return
	}

	if utils.Is.Empty(params["region"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "region"), 400)
		return
	}

	temp := facade.TempStorage
	temp = utils.Replace(temp, map[string]any{
		"${cos.secret_id}":  cast.ToString(params["secret_id"]),
		"${cos.secret_key}": cast.ToString(params["secret_key"]),
		"${cos.app_id}": 	 cast.ToString(params["app_id"]),
		"${cos.bucket}": 	 cast.ToString(params["bucket"]),
		"${cos.region}": 	 cast.ToString(params["region"]),
		"${cos.domain}": 	 cast.ToString(params["domain"]),
		"${cos.path}"  : 	 cast.ToString(params["path"]),
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.StorageToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/storage.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}

// testKODO - 测试KODO连接
func (this *Toml) testKODO(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["access_key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "access_key"), 400)
		return
	}

	if utils.Is.Empty(params["secret_key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "secret_key"), 400)
		return
	}

	if utils.Is.Empty(params["bucket"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bucket"), 400)
		return
	}

	if utils.Is.Empty(params["region"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "region"), 400)
		return
	}

	// KODO 对象存储
	client := qbox.NewMac(cast.ToString(params["access_key"]), cast.ToString(params["secret_key"]))

	bucket := storage.NewBucketManager(client, nil)
	_, err := bucket.GetBucketInfo(cast.ToString(params["bucket"]))

	if err != nil {
		this.json(ctx, err.Error(), facade.Lang(ctx, "测试KODO连接失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "测试KODO连接成功！"), 200)
}

// putStorageKODO - 修改KODO存储配置
func (this *Toml) putStorageKODO(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	if utils.Is.Empty(params["access_key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "access_key"), 400)
		return
	}

	if utils.Is.Empty(params["secret_key"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "secret_key"), 400)
		return
	}

	if utils.Is.Empty(params["bucket"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "bucket"), 400)
		return
	}

	if utils.Is.Empty(params["region"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "region"), 400)
		return
	}

	if utils.Is.Empty(params["domain"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "domain"), 400)
		return
	}

	temp := facade.TempStorage
	temp = utils.Replace(temp, map[string]any{
		"${kodo.access_key}":  cast.ToString(params["access_key"]),
		"${kodo.secret_key}":  cast.ToString(params["secret_key"]),
		"${kodo.bucket}": 	   cast.ToString(params["bucket"]),
		"${kodo.region}": 	   cast.ToString(params["region"]),
		"${kodo.domain}": 	   cast.ToString(params["domain"]),
	})

	// 正则匹配出所有的 ${?} 字符串
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindAllStringSubmatch(temp, -1)

	for _, match := range matches {
		temp = strings.Replace(temp, match[0], cast.ToString(facade.StorageToml.Get(match[1])), -1)
	}

	item := utils.File().Save(strings.NewReader(temp), "config/storage.toml")

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "修改失败！"), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "修改成功！"), 200)
}