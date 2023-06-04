package model

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	"net/url"
	"strings"
)

type AuthRules struct {
	Id     int    `gorm:"type:int(32); comment:主键;" json:"id"`
	Name   string `gorm:"comment:规则名称;" json:"name"`
	Method string `gorm:"comment:请求类型; default:'GET';" json:"method"`
	Route  string `gorm:"comment:路由;" json:"route"`
	Type   string `gorm:"default:'default'; comment:规则类型;" json:"type"`
	Hash   string `gorm:"comment:哈希值;" json:"hash"`
	Cost   int    `gorm:"type:int(32); comment:费用; default:1;" json:"cost"`
	Remark string `gorm:"comment:备注; default:Null;" json:"remark"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// AfterFind - 查询Hook
func (this *AuthRules) AfterFind(tx *gorm.DB) (err error) {

	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)

	return
}

// BeforeCreate - 创建前的Hook
func (this *AuthRules) BeforeCreate(tx *gorm.DB) (err error) {

	// 检查 hash 是否存在
	exist := facade.DB.Model(&AuthRules{}).Where("hash", this.Hash).Exist()
	if exist {
		return errors.New(fmt.Sprintf("hash: %s 已存在", this.Hash))
	}

	return
}

// InitAuthRules - 初始化AuthRules表
func InitAuthRules() {

	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&AuthRules{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "AuthRules表迁移失败")
		return
	}

	// 动态生成规则
	for _, item := range createAuthRules() {
		go saveAuthRules(item)
	}
}

// createAuthRules - 生成规则
func createAuthRules() (result []AuthRules) {

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
			"POST":   {
				"create",
				"path=save&remark=勾选后，拥有该权限的用户不仅可以修改所有人的用户信息，还可以通过GET请求直接获取到所有用户的全部信息（包括账号、邮箱和电话），请谨慎使用",
			},
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
			"PUT":    {"path=update&type=login","path=restore&type=login"},
			"POST":   {"path=save&type=login", "path=create&type=login"},
			"DELETE": {"path=remove&type=login", "path=delete&type=login", "path=clear&type=login"},
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

				result = append(result, AuthRules{
					Type  : param["type"],
					Method: strings.ToUpper(method),
					Route : "/api/" + key + utils.Ternary[string](utils.Is.Empty(param["path"]), "", "/" + param["path"]),
					Name  : names[key] + utils.Default(param["name"], methods[method][param["path"]]),
					Remark: param["remark"],
				})
			}
		}
	}

	return result
}

// saveAuthRules 保存权限规则
func saveAuthRules(item AuthRules) {

	method := strings.ToUpper(cast.ToString(item.Method))
	hash   := facade.Hash.Sum32(fmt.Sprintf("[%s]%s", method, item.Route))

	table := AuthRules{
		Hash:   hash,
		Type:   item.Type,
		Remark: item.Remark,
		Name:   cast.ToString(item.Name),
		Method: cast.ToString(item.Method),
		Route:  cast.ToString(item.Route),
	}

	// 查询条件
	query := facade.DB.Model(&item).Where("hash", hash)

	// 如果存在，就不要再添加了
	if exist := query.Exist(); exist {
		return
	}

	tx := query.Save(&table)
	if tx.Error != nil {
		if strings.Contains(tx.Error.Error(), "已存在") {
			return
		}
		facade.Log.Error(map[string]any{"error": tx.Error.Error()}, "自动添加规则失败")
	}
}