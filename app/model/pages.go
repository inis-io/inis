package model

import (
	"errors"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	"strings"
)

type Pages struct {
	Id         int    				 `gorm:"type:int(32); comment:主键;" json:"id"`
	Uid        int    				 `gorm:"type:int(32); comment:用户ID; default:0;" json:"uid"`
	Key        string 				 `gorm:"size:256; comment:唯一键; default:Null;" json:"key"`
	Title      string 				 `gorm:"size:256; comment:标题; default:Null;" json:"title"`
	Content    string 				 `gorm:"type:longtext; comment:内容; default:Null;" json:"content"`
	Editor     string 				 `gorm:"comment:编辑器; default:'tinymce';" json:"editor"`
	Tags 	   string  				 `gorm:"comment:标签; default:Null;" json:"tags"`
	Remark     string 				 `gorm:"comment:备注; default:Null;" json:"remark"`
	Audit	   int    				 `gorm:"type:int(12); comment:审核; default:0;" json:"audit"`
	LastUpdate int64  				 `gorm:"comment:最后更新时间; default:0;" json:"last_update"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitPages - 初始化Pages表
func InitPages() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Pages{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Pages表迁移失败")
		return
	}
}

// AfterSave - 保存后的Hook（包括 create update）
func (this *Pages) AfterSave(tx *gorm.DB) (err error) {

	// key 唯一处理
	if !utils.Is.Empty(this.Key) {
		exist := facade.DB.Model(&Pages{}).WithTrashed().Where("id", "!=", this.Id).Where("key", this.Key).Exist()
		if exist {
			return errors.New("key 已存在！")
		}
	}
	return
}

// AfterFind - 查询Hook
func (this *Pages) AfterFind(tx *gorm.DB) (err error) {

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

	// 标签信息
	tags := utils.ArrayUnique(utils.ArrayEmpty(strings.Split(this.Tags, "|")))

	this.Result = map[string]any{
		"comment": comment,
		"tags"   : facade.DB.Model(&[]Tags{}).WhereIn("id", tags).Column("id", "name", "avatar", "description"),
	}
	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)

	return
}

// config - 获取配置
func (this *Pages) config(key ...any) (json map[string]any) {

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