package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/model"
	"mime/multipart"
	"net/url"
	"strings"
)

type Test struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Test) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"request": this.request,
		"alipay":  this.alipay,
		"system":  this.system,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Test) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"return-url": this.returnUrl,
		"notify-url": this.notifyUrl,
		"request":    this.request,
		"upload":     this.upload,
		"qq":         this.qq,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPUT - PUT请求本体
func (this *Test) IPUT(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"request": this.request,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IDEL - DELETE请求本体
func (this *Test) IDEL(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"request": this.request,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// INDEX - GET请求本体
func (this *Test) INDEX(ctx *gin.Context) {

	var rules []model.AuthRules

	batch := map[string]map[string][]string{
		"test": {
			"GET":    {
				"path=&name=兔子专用&type=common",
				"path=request&name=测试GET请求&type=common",
			},
			"PUT":    {"path=request&name=测试GET请求&type=common"},
			"POST":   {"path=request&name=测试GET请求&type=common"},
			"DELETE": {"path=request&name=测试GET请求&type=common"},
		},
		"proxy": {
			"GET":    {"path=&name=代理 GET 请求&type=login"},
			"PUT":    {"path=&name=代理 PUT 请求&type=login"},
			"POST":   {"path=&name=代理 POST 请求&type=login"},
			"PATCH":  {"path=&name=代理 PATCH 请求&type=login"},
			"DELETE": {"path=&name=代理 DELETE 请求&type=login"},
		},
		"file": {
			"GET": {
				"path=rand&name=随机图&type=common",
				"path=to-base64&name=网络图片转base64&type=common",
			},
			"POST": {"path=upload&name=简单上传&type=login"},
		},
		"comm": {
			"POST": {
				"path=login&name=传统和加密登录&type=common",
				"path=social-login&name=验证码登录&type=common",
				"path=register&name=注册账户&type=common",
				"path=check-token&name=校验登录&type=common",
			},
			"DELETE": {"path=logout&name=退出登录&type=common"},
		},
		"tags": {
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"users":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"links":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"path=update&type=login","path=restore&type=login"},
			"POST":   {
				"path=save&type=login",
				"path=create&type=login",
			},
			"DELETE": {
				"path=remove&type=login",
				"path=delete&type=login",
				"path=clear&type=login",
			},
		},
		"pages":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"banner":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"config":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"article":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"placard":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"comment":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"api-keys": {
			"GET":    {"one", "all", "count", "column"},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"auth-group":{
			"GET":    {"one", "all", "count", "column"},
			"PUT":    {"update", "restore", "path=uids&name=更改用户权限"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"auth-rules":{
			"GET":    {"one", "all", "count", "column"},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"auth-pages":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"links-group":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
		"article-group":{
			"GET":    {
				"path=one&type=common",
				"path=all&type=common",
				"path=count&type=common",
				"path=column&type=common",
			},
			"PUT":    {"update","restore"},
			"POST":   {"save", "create"},
			"DELETE": {"remove", "delete", "clear"},
		},
	}

	// 接口名称
	names := map[string]string{
		"test":          "【测试 API】",
		"proxy":         "【代理 API】",
		"file":          "【文件 API】",
		"comm":          "【公共 API】",
		"tags":          "【标签 API】",
		"pages":         "【页面 API】",
		"users":         "【用户 API】",
		"links":         "【友链 API】",
		"banner":        "【轮播 API】",
		"config":        "【配置 API】",
		"article":       "【文章 API】",
		"comment":       "【评论 API】",
		"placard":       "【公告 API】",
		"api-keys":      "【接口密钥 API】",
		"auth-group":    "【权限分组 API】",
		"auth-pages":    "【页面权限 API】",
		"auth-rules":    "【权限规则 API】",
		"links-group":   "【友链分组 API】",
		"article-group": "【文章分组 API】",
	}

	// 基础方法
	methods := map[string]map[string]string{
		"GET": {
			"one": "获取指定",
			"all": "获取全部",
			"count": "查询数量",
			"column": "列查询",
		},
		"POST": {
			"save": "保存数据（推荐）",
			"create": "添加数据",
		},
		"PUT": {
			"update": "更新数据",
			"restore": "恢复数据",
		},
		"DELETE": {
			"remove": "软删除（回收站）",
			"delete": "彻底删除",
			"clear": "清空回收站",
		},
	}

	// 批量生成公共接口
	for key, value := range batch {
		for method, items := range value {
			for _, item := range items {

				param := map[string]string{
					"type": "default",
				}

				// 检查 item 是否包含 = 号
				if !strings.Contains(item, "=") {

					param["path"] = item

				} else {

					// 解析 "name=代理 GET 请求&path=&type=common"
					values, _ := url.ParseQuery(item)

					for name, text := range values {
						if len(text) == 1 {
							param[name] = text[0]
						} else {
							param[name] = cast.ToString(text)
						}
					}
				}

				rules = append(rules, model.AuthRules{
					Type  : param["type"],
					Method: strings.ToUpper(method),
					Route : "/api/" + key + utils.Ternary[string](utils.Is.Empty(param["path"]), "", "/" + param["path"]),
					Name  : names[key] + utils.Default(param["name"], methods[method][param["path"]]),
				})
			}
		}
	}

	// res := map[string]any{
	// 	"user" : this.meta.user(ctx),
	// 	"route": this.meta.route(ctx),
	// 	"rules": this.meta.rules(ctx),
	// }

	this.json(ctx, rules, facade.Lang(ctx, "好的！"), 200)
}

func (this *Test) qq(ctx *gin.Context) {

	params := this.params(ctx)

	if params["message_type"] == "private" {
		fmt.Println(utils.Json.Encode(params))

		item := utils.Curl(utils.CurlRequest{
			Method: "GET",
			Url:    "http://localhost:5700/send_private_msg",
			Query: map[string]any{
				"user_id": cast.ToString(params["user_id"]),
				"message": cast.ToString(params["message"]),
			},
		}).Send()

		if item.Error != nil {
			fmt.Println("发送失败", item.Error.Error())
			return
		}

		fmt.Println("发送成功", item.Json)
	}

	this.json(ctx, params, facade.Lang(ctx, "好的！"), 200)
}

func (this *Test) system(ctx *gin.Context) {

	params := this.params(ctx)

	this.json(ctx, params, facade.Lang(ctx, "好的！"), 200)
}

// INDEX - GET请求本体
func (this *Test) upload(ctx *gin.Context) {

	params := this.params(ctx)

	// 上传文件
	file, err := ctx.FormFile("file")
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 文件数据
	bytes, err := file.Open()
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}
	defer func(bytes multipart.File) {
		err := bytes.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(bytes)

	// 文件后缀
	suffix := file.Filename[strings.LastIndex(file.Filename, "."):]
	params["suffix"] = suffix

	item := facade.Storage.Upload(facade.Storage.Path()+suffix, bytes)
	if item.Error != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	params["item"] = item

	fmt.Println("url: ", item.Domain+item.Path)

	this.json(ctx, params, facade.Lang(ctx, "好的！"), 200)
}

func (this *Test) alipay(ctx *gin.Context) {

	// 初始化 BodyMap
	body := make(gopay.BodyMap)
	body.Set("subject", "统一收单下单并支付页面接口")
	body.Set("out_trade_no", uuid.New().String())
	body.Set("total_amount", "0.01")
	body.Set("product_code", "FAST_INSTANT_TRADE_PAY")

	payUrl, err := facade.Alipay().TradePagePay(context.Background(), body)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			fmt.Println(bizErr)
			return
		}
		fmt.Println(err)
		return
	}

	fmt.Println(payUrl)

	this.json(ctx, payUrl, "数据请求成功！", 200)
}

func (this *Test) returnUrl(ctx *gin.Context) {

	params := this.params(ctx)

	fmt.Println("==================== returnUrl：", params)
}

func (this *Test) notifyUrl(ctx *gin.Context) {

	params := this.params(ctx)

	fmt.Println("==================== notifyUrl：", params)
}

// 测试网络请求
func (this *Test) request(ctx *gin.Context) {

	params := this.params(ctx)

	this.json(ctx, map[string]any{
		"method":  ctx.Request.Method,
		"params":  params,
		"headers": this.headers(ctx),
	}, facade.Lang(ctx, "数据请求成功！"), 200)
}
