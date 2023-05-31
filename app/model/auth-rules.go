package model

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	"strings"
	"sync"
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

// AfterSave - 保存后的Hook（包括 create update）
func (this *AuthRules) AfterSave(tx *gorm.DB) (err error) {

	// 检查 hash 是否存在
	exist := facade.DB.Model(&AuthRules{}).Where("hash", this.Hash).Exist()
	if exist {
		return errors.New(fmt.Sprintf("hash: %s 已存在", this.Hash))
	}

	return
}

// InitAuthRules - 初始化AuthRules表
func InitAuthRules() {
	// 数据库
	DB := facade.NewDB(facade.DBModeMySql)
	// 迁移表
	err := DB.Drive().AutoMigrate(&AuthRules{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "AuthRules表迁移失败")
		return
	}

	// 规则列表
	rules := []AuthRules{
		{Method: "POST",   Route: "/api/comm/login", Name: "【公共 API】传统和加密登录"},
		{Method: "POST",   Route: "/api/comm/social-login", Name: "【公共 API】验证码登录"},
		{Method: "POST",   Route: "/api/comm/register", Name: "【公共 API】注册账户"},
		{Method: "POST",   Route: "/api/comm/check-token", Name: "【公共 API】校验登录"},
		{Method: "DELETE", Route: "/api/comm/logout", Name: "【公共 API】退出登录"},
		{Method: "GET",    Route: "/api/file/rand", Name: "【文件 API】随机图"},
		{Method: "GET",    Route: "/api/file/to-base64", Name: "【文件 API】网络图片转base64"},
		{Method: "POST",   Route: "/api/file/upload", Name: "【文件 API】简单上传"},
		{Method: "GET",    Route: "/api/test", Name: "【测试 API】兔子专用"},
		{Method: "GET",    Route: "/api/test/request", Name: "【测试 API】测试GET请求"},
		{Method: "PUT",    Route: "/api/test/request", Name: "【测试 API】测试PUT请求"},
		{Method: "POST",   Route: "/api/test/request", Name: "【测试 API】测试POST请求"},
		{Method: "DELETE", Route: "/api/test/request", Name: "【测试 API】测试DEL请求"},
		{Method: "PUT",    Route: "/api/auth-group/uids", Name: "【权限分组 API】更改用户权限"},
		{Method: "GET",    Route: "/api/proxy", Name: "【代理 API】代理发起 GET 请求"},
		{Method: "PUT",    Route: "/api/proxy", Name: "【代理 API】代理发起 PUT 请求"},
		{Method: "POST",   Route: "/api/proxy", Name: "【代理 API】代理发起 POST 请求"},
		{Method: "PATCH",  Route: "/api/proxy", Name: "【代理 API】代理发起 PATCH 请求"},
		{Method: "DELETE", Route: "/api/proxy", Name: "【代理 API】代理发起 DELETE 请求"},
	}

	// 基础接口
	basics := map[string]string{
		"tags":        "【标签 API】",
		"users":       "【用户 API】",
		"links":       "【友链 API】",
		"banner":      "【轮播 API】",
		"config":      "【配置 API】",
		"article":     "【文章 API】",
		"comment":     "【评论 API】",
		"placard":     "【公告 API】",
		"api-keys":    "【接口密钥 API】",
		"auth-group":  "【权限分组 API】",
		"auth-pages":  "【页面权限 API】",
		"auth-rules":  "【权限规则 API】",
		"links-group": "【友链分组 API】",
		"article-group": "【文章分组 API】",
	}

	// 基础方法
	methods := map[string][]map[string]string{
		"GET": {
			{"one": "获取指定"},
			{"all": "获取全部"},
			{"count": "查询数量"},
			{"column": "列查询"},
		},
		"POST": {
			{"save": "保存数据（推荐）"},
			{"create": "添加数据"},
		},
		"PUT": {
			{"update": "更新数据"},
			{"restore": "恢复数据"},
		},
		"DELETE": {
			{"remove": "软删除（回收站）"},
			{"delete": "彻底删除"},
			{"clear": "清空回收站"},
		},
	}

	// 遍历生成规则
	for keys, value := range basics {
		for key, method := range methods {
			for _, item := range method {
				for k, v := range item {
					rules = append(rules, AuthRules{
						Method: strings.ToUpper(key),
						Name:   value + v,
						Route:  fmt.Sprintf("/api/%s/%s", keys, k),
					})
				}
			}
		}
	}

	wg := sync.WaitGroup{}

	// 检查规则是否存在，不存在则添加
	for _, item := range rules {
		wg.Add(1)
		go func(item AuthRules, wg *sync.WaitGroup) {
			defer wg.Done()

			method := strings.ToUpper(cast.ToString(item.Method))

			tx := DB.Model(&item).Where("route", item.Route).Where("method", method).Save(&AuthRules{
				Name:   cast.ToString(item.Name),
				Method: cast.ToString(item.Method),
				Route:  cast.ToString(item.Route),
				Hash:   facade.Hash.Sum32(fmt.Sprintf("[%s]%s", method, item.Route)),
			})
			if tx.Error != nil {
				if strings.Contains(tx.Error.Error(), "已存在") {
					return
				}
				facade.Log.Error(map[string]any{"error": tx.Error.Error()}, "自动添加规则失败")
			}
		}(item, &wg)
	}

	wg.Wait()
	go setRuleType()
}

// setRuleType - 设置接口类型
func setRuleType() {

	rules := []AuthRules{
		{Method: "POST", 	 Route: "/api/comm/login", 		  Type: "common"},
		{Method: "POST", 	 Route: "/api/comm/social-login", Type: "common"},
		{Method: "POST", 	 Route: "/api/comm/register", 	  Type: "common"},
		{Method: "POST", 	 Route: "/api/comm/check-token",  Type: "common"},
		{Method: "DELETE",   Route: "/api/comm/logout",       Type: "common"},
		{Method: "GET", 	 Route: "/api/test", 		      Type: "common"},
		{Method: "GET", 	 Route: "/api/test/request",      Type: "common"},
		{Method: "PUT", 	 Route: "/api/test/request",      Type: "common"},
		{Method: "POST",     Route: "/api/test/request",      Type: "common"},
		{Method: "DELETE",   Route: "/api/test/request",      Type: "common"},
		{Method: "GET", 	 Route: "/api/file/rand", 		  Type: "common"},
		{Method: "GET", 	 Route: "/api/file/to-base64", 	  Type: "common"},
		{Method: "GET", Route: "/api/users/one", 		      Type: "common"},
		{Method: "GET", Route: "/api/users/all", 		      Type: "common"},
		{Method: "GET", Route: "/api/users/count",  	      Type: "common"},
		{Method: "GET", Route: "/api/users/column", 	      Type: "common"},
		{Method: "GET", Route: "/api/banner/one", 		      Type: "common"},
		{Method: "GET", Route: "/api/banner/all", 		      Type: "common"},
		{Method: "GET", Route: "/api/banner/count", 	      Type: "common"},
		{Method: "GET", Route: "/api/banner/column", 	      Type: "common"},
		{Method: "GET", Route: "/api/tags/one", 		      Type: "common"},
		{Method: "GET", Route: "/api/tags/all", 		      Type: "common"},
		{Method: "GET", Route: "/api/tags/count", 		      Type: "common"},
		{Method: "GET", Route: "/api/tags/column", 		      Type: "common"},
		{Method: "GET", Route: "/api/placard/one", 		      Type: "common"},
		{Method: "GET", Route: "/api/placard/all", 		      Type: "common"},
		{Method: "GET", Route: "/api/placard/count", 	      Type: "common"},
		{Method: "GET", Route: "/api/placard/column", 	      Type: "common"},
		{Method: "GET", Route: "/api/links/one", 		      Type: "common"},
		{Method: "GET", Route: "/api/links/all", 		      Type: "common"},
		{Method: "GET", Route: "/api/links/count", 		      Type: "common"},
		{Method: "GET", Route: "/api/links/column", 	      Type: "common"},
		{Method: "GET", Route: "/api/links-group/one", 	      Type: "common"},
		{Method: "GET", Route: "/api/links-group/all",        Type: "common"},
		{Method: "GET", Route: "/api/links-group/count",      Type: "common"},
		{Method: "GET", Route: "/api/links-group/column",     Type: "common"},
		{Method: "GET", Route: "/api/config/one", 		      Type: "common"},
		{Method: "GET", Route: "/api/config/all", 		      Type: "common"},
		{Method: "GET", Route: "/api/config/count", 	      Type: "common"},
		{Method: "GET", Route: "/api/config/column", 	      Type: "common"},
		{Method: "GET", Route: "/api/article/one", 		      Type: "common"},
		{Method: "GET", Route: "/api/article/all", 		      Type: "common"},
		{Method: "GET", Route: "/api/article/count", 	      Type: "common"},
		{Method: "GET", Route: "/api/article/column", 	      Type: "common"},
		{Method: "GET", Route: "/api/article-group/one",      Type: "common"},
		{Method: "GET", Route: "/api/article-group/all",      Type: "common"},
		{Method: "GET", Route: "/api/article-group/count",    Type: "common"},
		{Method: "GET", Route: "/api/article-group/column",   Type: "common"},
		{Method: "GET", Route: "/api/comment/one", 			  Type: "common"},
		{Method: "GET", Route: "/api/comment/all", 			  Type: "common"},
		{Method: "GET", Route: "/api/comment/count", 		  Type: "common"},
		{Method: "GET", Route: "/api/comment/column", 		  Type: "common"},
		{Method: "GET", Route: "/api/auth-pages/one", 		  Type: "login"},
		{Method: "GET", Route: "/api/auth-pages/all", 		  Type: "login"},
		{Method: "GET", Route: "/api/auth-pages/count", 	  Type: "login"},
		{Method: "GET", Route: "/api/auth-pages/column", 	  Type: "login"},
		{Method: "POST",Route: "/api/users/save", Remark:"勾选后，拥有该权限的用户不仅可以修改所有人的用户信息，还可以通过GET请求直接获取到所有用户的全部信息（包括账号、邮箱和电话），请谨慎使用"},
	}

	wg := sync.WaitGroup{}

	for _, item := range rules {
		wg.Add(1)
		go func(item AuthRules, wg *sync.WaitGroup) {
			defer wg.Done()
			tx := facade.DB.Model(&AuthRules{}).Where("method", item.Method).Where("route", item.Route).Update(&item)
			if tx.Error != nil {
				facade.Log.Error(map[string]any{"error": tx.Error.Error()}, "自动设置公共接口规则失败")
			}
		}(item, &wg)
	}

	wg.Wait()
}
