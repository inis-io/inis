package facade

import (
	"crypto/md5"
	"fmt"
	"github.com/fsnotify/fsnotify"
	JWT "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"time"
)

// CryptToml - Crypt配置文件
var CryptToml *utils.ViperResponse

// initCryptToml - 初始化Crypt配置文件
func initCryptToml() {

	key := fmt.Sprintf("%v-%v", uuid.New().String(), time.Now().Unix())
	secret := fmt.Sprintf("INIS-%x", md5.Sum([]byte(key)))

	item := utils.Viper(utils.ViperModel{
		Path: "config",
		Mode: "toml",
		Name: "crypt",
		Content: utils.Replace(TempCrypt, map[string]any{
			"${jwt.key}": 		secret,
			"${jwt.expire}":    "7 * 24 * 60 * 60",
			"${jwt.issuer}" :   "inis.cn",
			"${jwt.subject}":   "inis",
		}),
	}).Read()

	if item.Error != nil {
		Log.Error(map[string]any{
			"error":     item.Error,
			"func_name": utils.Caller().FuncName,
			"file_name": utils.Caller().FileName,
			"file_line": utils.Caller().Line,
		}, "Crypt配置初始化错误")
		return
	}

	CryptToml = &item
}

// 初始化加密配置
func initCrypt() {

}

func init() {
	// 初始化配置文件
	initCryptToml()
	// 初始化缓存
	initCrypt()

	// 监听配置文件变化
	CryptToml.Viper.WatchConfig()
	// 配置文件变化时，重新初始化配置文件
	CryptToml.Viper.OnConfigChange(func(event fsnotify.Event) {
		initCrypt()
	})
}

type JwtStruct struct {
	request  JwtRequest
	response JwtResponse
}

// JwtRequest - JWT请求
type JwtRequest struct {
	// 过期时间
	Expire  int64        `json:"expire"`
	// 颁发者签名
	Issuer  string       `json:"issuer"`
	// 主题
	Subject string       `json:"subject"`
	// 密钥
	Key     string       `json:"key"`
}

// JwtResponse - JWT响应
type JwtResponse struct {
	Text  string         `json:"text"`
	Data  map[string]any `json:"data"`
	Error error          `json:"error"`
	Valid int64          `json:"valid"`
}

// Jwt - 入口
func Jwt(request ...JwtRequest) *JwtStruct {

	if len(request) == 0 {
		request = append(request, JwtRequest{})
	}

	// 过期时间
	if request[0].Expire == 0 {
		request[0].Expire = cast.ToInt64(utils.Calc(CryptToml.Get("jwt.expire", "7200")))
	}

	// 颁发者签名
	if utils.Is.Empty(request[0].Issuer) {
		request[0].Issuer = cast.ToString(CryptToml.Get("jwt.issuer", "inis.cn"))
	}

	// 主题
	if utils.Is.Empty(request[0].Subject) {
		request[0].Subject = cast.ToString(CryptToml.Get("jwt.subject", "inis"))
	}

	// 密钥
	if utils.Is.Empty(request[0].Key) {
		request[0].Key = cast.ToString(CryptToml.Get("jwt.key", "inis"))
	}

	return &JwtStruct{
		request: request[0],
		response: JwtResponse{
			Data: make(map[string]any),
		},
	}
}

// Create - 创建JWT
func (this *JwtStruct) Create(data map[string]any) (result JwtResponse) {

	type JwtClaims struct {
		Data map[string]any `json:"data"`
		JWT.RegisteredClaims
	}

	IssuedAt  := JWT.NewNumericDate(time.Now())
	ExpiresAt := JWT.NewNumericDate(time.Now().Add(time.Second * time.Duration(this.request.Expire)))

	item, err := JWT.NewWithClaims(JWT.SigningMethodHS256, JwtClaims{
		Data: data,
		RegisteredClaims: JWT.RegisteredClaims{
			IssuedAt:  IssuedAt,				// 签发时间戳
			ExpiresAt: ExpiresAt,				// 过期时间戳
			Issuer:    this.request.Issuer,		// 颁发者签名
			Subject:   this.request.Subject,	// 签名主题
		},
	}).SignedString([]byte(this.request.Key))

	if err != nil {
		this.response.Error = err
		return this.response
	}

	this.response.Text = item

	return this.response
}

// Parse - 解析JWT
func (this *JwtStruct) Parse(token any) (result JwtResponse) {

	type JwtClaims struct {
		Data map[string]any `json:"data"`
		JWT.RegisteredClaims
	}

	item, err := JWT.ParseWithClaims(cast.ToString(token), &JwtClaims{}, func(token *JWT.Token) (any, error) {
		return []byte(this.request.Key), nil
	})

	if err != nil {
		Log.Error(map[string]any{
			"error":     err,
			"func_name": utils.Caller().FuncName,
			"file_name": utils.Caller().FileName,
			"file_line": utils.Caller().Line,
		}, "JWT解析错误")
		this.response.Error = err
		return this.response
	}

	if key, _ := item.Claims.(*JwtClaims); item.Valid {
		this.response.Data  = key.Data
		this.response.Valid = key.RegisteredClaims.ExpiresAt.Time.Unix() - time.Now().Unix()
	}

	return this.response
}