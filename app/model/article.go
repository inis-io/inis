package model

import (
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	"strings"
)

type Article struct {
	Id         int     				 `gorm:"type:int(32); comment:主键;" json:"id"`
	Uid		   int     				 `gorm:"type:int(32); comment:用户ID; default:0;" json:"uid"`
	Title      string  				 `gorm:"size:256; comment:标题; default:Null;" json:"title"`
	Abstract   string  				 `gorm:"size:512; comment:摘要; default:Null;" json:"abstract"`
	Content    string  				 `gorm:"type:longtext; comment:内容; default:Null;" json:"content"`
	Covers     string  				 `gorm:"type:text; comment:封面; default:Null;" json:"covers"`
	Top 	   int     				 `gorm:"type:int(1); comment:置顶; default:0;" json:"top"`
	Views	   int     				 `gorm:"type:int(32); comment:浏览量; default:0;" json:"views"`
	Tags 	   string  				 `gorm:"comment:标签; default:Null;" json:"tags"`
	Group	   string  				 `gorm:"comment:分组; default:Null;" json:"group"`
	Remark     string  				 `gorm:"comment:备注; default:Null;" json:"remark"`
	Editor     string  				 `gorm:"comment:编辑器; default:'tinymce';" json:"editor"`
	LastUpdate int64   				 `gorm:"comment:最后更新时间; default:0;" json:"last_update"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitArticle - 初始化Article表
func InitArticle() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Article{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Article表迁移失败")
		return
	}
}

// AfterFind - 查询Hook
func (this *Article) AfterFind(tx *gorm.DB) (err error) {

	// 作者信息
	author := make(map[string]any)
	allow  := []string{"id", "nickname", "avatar", "description", "result"}
	user   := facade.DB.Model(&Users{}).Find(this.Uid)

	if !utils.Is.Empty(user) {
		author = utils.Map.WithField(user, allow)
	}
	// 标签信息
	tags  := utils.ArrayUnique(utils.ArrayEmpty(strings.Split(this.Tags, "|")))
	// 分类信息
	group := utils.ArrayUnique(utils.ArrayEmpty(strings.Split(this.Group, "|")))

	// 当前的评论配置
	comment := cast.ToStringMap(cast.ToStringMap(utils.Json.Decode(this.Json))["comment"])
	config  := this.config("comment")

	// 允许评论选项继承了父级配置
	if cast.ToInt(comment["allow"]) == 0 {
		comment["allow"] = config["allow"]
	}
	// 显示评论选项继承了父级配置
	if cast.ToInt(comment["show"]) == 0 {
		comment["show"]  = config["show"]
	}

	this.Result = map[string]any{
		"author" : author,
		"comment": comment,
		"group"  : facade.DB.Model(&[]ArticleGroup{}).WhereIn("id", group).Column("id", "pid", "name", "avatar", "description"),
		"tags"   : facade.DB.Model(&[]Tags{}).WhereIn("id", tags).Column("id", "name", "avatar", "description"),
	}
	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)

	return
}

// config - 获取配置
func (this *Article) config(key ...any) (json map[string]any) {

	var config map[string]any

	// 缓存名称
	cacheName := "config[ARTICLE]"
	// 是否开启了缓存
	cacheState := cast.ToBool(facade.CacheToml.Get("open"))

	// 检查缓存是否存在
	if cacheState && facade.Cache.Has(cacheName) {

		config = cast.ToStringMap(facade.Cache.Get(cacheName))

	} else {

		config = facade.DB.Model(&Config{}).Where("key", "ARTICLE").Find()
		// 存储到缓存中
		if cacheState {
			go facade.Cache.Set(cacheName, config)
		}
	}

	if len(key) > 0 {
		return cast.ToStringMap(cast.ToStringMap(config["json"])[cast.ToString(key[0])])
	}

	return cast.ToStringMap(config["json"])
}