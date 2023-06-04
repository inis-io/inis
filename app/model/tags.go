package model

import (
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
)

type Tags struct {
	Id          int    `gorm:"type:int(32); comment:主键;" json:"id"`
	Name        string `gorm:"size:32; comment:昵称; default:Null;" json:"name"`
	Description string `gorm:"comment:描述; default:Null;" json:"description"`
	Avatar      string `gorm:"size:256; comment:头像; default:Null;" json:"avatar"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitTags - 初始化Tags表
func InitTags() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Tags{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Tags表迁移失败")
		return
	}
}

// AfterFind - 查询Hook
func (this *Tags) AfterFind(tx *gorm.DB) (err error) {

	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)

	return
}