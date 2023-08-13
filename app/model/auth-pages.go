package model

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	"strings"
	"sync"
)

type AuthPages struct {
	Id     	   int    				 `gorm:"type:int(32); comment:主键;" json:"id"`
	Name   	   string 				 `gorm:"comment:名称;" json:"name"`
	Path   	   string 				 `gorm:"comment:路径;" json:"path"`
	Icon   	   string 				 `gorm:"comment:图标;" json:"icon"`
	Svg    	   string 				 `gorm:"type:text; comment:SVG图标;" json:"svg"`
	Size   	   string 				 `gorm:"comment:图标大小; default:'16px';" json:"size"`
	Hash   	   string 				 `gorm:"comment:哈希值;" json:"hash"`
	Remark 	   string 				 `gorm:"comment:备注; default:Null;" json:"remark"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// AfterFind - 查询Hook
func (this *AuthPages) AfterFind(tx *gorm.DB) (err error) {

	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)

	return
}

// BeforeCreate - 创建前的Hook
func (this *AuthPages) BeforeCreate(tx *gorm.DB) (err error) {

	// 检查 hash 是否存在
	if exist := facade.DB.Model(&AuthRules{}).WithTrashed().Where("hash", this.Hash).Exist(); exist {
		return errors.New(fmt.Sprintf("hash: %s 已存在", this.Hash))
	}

	return
}

// InitAuthPages - 初始化AuthPages表
func InitAuthPages() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&AuthPages{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "AuthPages表迁移失败")
		return
	}

	// 页面列表
	pages := []AuthPages{
		{Name: "撰写文章", Icon: "article", Path: "/admin/article/write"},
		{Name: "文章列表", Icon: "article", Path: "/admin/article"},
		{Name: "文章分组", Icon: "group", Path: "/admin/article/group", Size: "14px"},
		{Name: "用户", Icon: "user", Path: "/admin/users"},
		{Name: "评论", Icon: "comment", Path: "/admin/comment"},
		{Name: "公告", Icon: "bell", Path: "/admin/placard"},
		{Name: "轮播", Icon: "banner", Path: "/admin/banner"},
		{Name: "标签", Icon: "tag", Path: "/admin/tags"},
		{Name: "友链", Icon: "link", Path: "/admin/links"},
		{Name: "系统配置", Icon: "system", Path: "/admin/system", Size: "15px"},
		{Name: "页面列表", Icon: "open", Path: "/admin/pages", Size: "17px"},
		{Name: "撰写页面", Icon: "article", Path: "/admin/pages/write"},
		{Name: "友链分组", Icon: "group", Path: "/admin/links/group", Size: "14px"},
		{Name: "权限规则", Icon: "rule", Path: "/admin/auth/rules", Size: "17px"},
		{Name: "权限分组", Icon: "group", Path: "/admin/auth/group", Size: "14px"},
		{Name: "页面权限", Icon: "open", Path: "/admin/auth/pages", Size: "17px"},
		{Name: "接口密钥", Icon: "key", Path: "/admin/api/keys", Size: "14px"},
		{Name: "IP黑名单", Icon: "qps", Path: "/admin/ip/black", Size: "14px"},
		{Name: "QPS预警", Icon: "black", Path: "/admin/qps/warn", Size: "14px"},
	}

	wg := sync.WaitGroup{}

	// 检查规则是否存在，不存在则添加
	for _, item := range pages {
		wg.Add(1)
		go func(item AuthPages, wg *sync.WaitGroup) {
			defer wg.Done()

			hash := utils.Hash.Sum32(item.Path)

			tx := facade.DB.Model(&item).Where("hash", hash)

			// 如果存在，就不要再添加了
			if exist := tx.Exist(); exist {
				return
			}

			res := tx.Save(&AuthPages{
				Hash: hash,
				Name: cast.ToString(item.Name),
				Path: cast.ToString(item.Path),
				Icon: cast.ToString(item.Icon),
				Size: cast.ToString(item.Size),
			})

			if res.Error != nil {
				if strings.Contains(res.Error.Error(), "已存在") {
					return
				}
				facade.Log.Error(map[string]any{"error": res.Error.Error()}, "自动添加页面失败")
			}
		}(item, &wg)
	}

	wg.Wait()
}