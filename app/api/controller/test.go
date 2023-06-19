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
	"mime/multipart"
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
		"cache":   this.getCache,
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
		"cache":      this.setCache,
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
		"cache":   this.delCache,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// INDEX - GET请求本体
func (this *Test) INDEX(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	// res := gin.H{
	//
	// 	// "root" : this.meta.root(ctx),
	// 	// "user" : this.meta.user(ctx),
	// 	// "route": this.meta.route(ctx),
	// 	// "rules": this.meta.rules(ctx),
	// 	// "json" : utils.Json.Encode(params["json"]),
	// }

	// text, ok := facade.CacheFile.Get(cast.ToString(params["key"]))
	// if !ok {
	// 	this.json(ctx, nil, "缓存不存在或已过期！", 400)
	// 	return
	// }

	this.json(ctx, facade.Cache.Get(params["key"]), facade.Lang(ctx, "好的！"), 200)
}

func (this *Test) getCache(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx)

	this.json(ctx, map[string]any{
		"value": facade.Cache.Get(params["key"]),
		// "has":  facade.CacheFile.Has(params["key"]),
	}, facade.Lang(ctx, "好的！"), 200)

	// this.json(ctx, map[string]any{
	// 	"value": string(facade.CacheFile.Get(params["key"])),
	// 	"has":  facade.CacheFile.Has(params["key"]),
	// }, facade.Lang(ctx, "好的！"), 200)
}

func (this *Test) setCache(ctx *gin.Context) {

	params := this.params(ctx)

	key   := cast.ToString(params["key"])
	value := cast.ToString(params["value"])

	// ok := facade.CacheFile.Set(key, []byte(value), params["exp"])
	ok := facade.Cache.Set(key, value, params["exp"])

	this.json(ctx, ok, "ok", 200)
}


func (this *Test) delCache(ctx *gin.Context) {

	ok := facade.Cache.DelTags([]any{"inis","test"}, "unti", "admin")

	if !ok {
		this.json(ctx, nil, "no", 400)
		return
	}

	this.json(ctx, nil, "ok", 200)
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
