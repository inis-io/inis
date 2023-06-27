package controller

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	// JWT "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/google/uuid"
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

	// 请求参数
	// params := this.params(ctx)

	item := RSA.Generate(2048)

	if item.Error != nil {
		this.json(ctx, nil, item.Error.Error(), 400)
		return
	}

	PrivateKey := item.PrivateKey
	PublicKey  := item.PublicKey

	// 私钥加密
	encode := RSA.Encrypt(PublicKey, "123456")
	decode := RSA.Decrypt(PrivateKey, encode.Text)

	res := gin.H{
		// "root" : this.meta.root(ctx),
		// "user" : this.meta.user(ctx),
		// "route": this.meta.route(ctx),
		// "rules": this.meta.rules(ctx),
		// "json" : utils.Json.Encode(params["json"]),
		"encode": encode.Text,
		"decode": decode.Text,
		// "rsa": RSA.Generate(2048),
	}

	this.json(ctx, res, facade.Lang(ctx, "好的！"), 200)
}

var RSA *RSAStruct

type RSAStruct struct {}

type RSAResponse struct {
	// 私钥
	PrivateKey string
	// 公钥
	PublicKey  string
	// 错误信息
	Error error
	// 文本
	Text string
}

// Generate 生成 RSA 密钥对
func (this *RSAStruct) Generate(bits any) (result *RSAResponse) {

	result = &RSAResponse{}

	private, err := rsa.GenerateKey(rand.Reader, cast.ToInt(bits))
	if err != nil {
		result.Error = err
		return result
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(private)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyBytes,
	}

	// 生成私钥
	result.PrivateKey = string(pem.EncodeToMemory(&privateKeyBlock))

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		result.Error = err
		return result
	}
	publicKeyBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyBytes,
	}

	// 生成公钥
	result.PublicKey = string(pem.EncodeToMemory(&publicKeyBlock))

	return result
}

// Encrypt 加密
func (this *RSAStruct) Encrypt(publicKey, text string) (result *RSAResponse) {

	result = &RSAResponse{}

	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Errorf("%v", r)
		}
	}()

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		result.Error = errors.New("public key error")
		return result
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		result.Error = err
		return result
	}

	pub := pubInterface.(*rsa.PublicKey)
	encode, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(text))
	if err != nil {
		result.Error = err
		return result
	}

	result.Text = base64.StdEncoding.EncodeToString(encode)

	return result
}

// Decrypt 解密
func (this *RSAStruct) Decrypt(privateKey, text string) (result *RSAResponse) {

	result = &RSAResponse{}

	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Errorf("%v", r)
		}
	}()

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		result.Error = errors.New("private key error")
		return result
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		result.Error = err
		return result
	}

	decode, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		result.Error = err
		return result
	}

	encode, err := rsa.DecryptPKCS1v15(rand.Reader, priv, decode)
	if err != nil {
		result.Error = err
		return result
	}

	result.Text = string(encode)

	return result
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
