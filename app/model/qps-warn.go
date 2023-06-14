package model

import (
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
)

type QpsWarn struct {
	Id         int     				 `gorm:"type:int(32); comment:主键;" json:"id"`
	Ip         string  				 `gorm:"comment:IP; default:Null;" json:"ip"`
	Agent      string  				 `gorm:"type:varchar(512); comment:浏览器信息; default:Null;" json:"agent"`
	Path	   string  				 `gorm:"comment:请求路径; default:Null;" json:"path"`
	Method     string  				 `gorm:"type:varchar(32); comment:请求方法; default:Null;" json:"method"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitQpsWarn - 初始化QpsWarn表
func InitQpsWarn() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&QpsWarn{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "QpsWarn表迁移失败")
		return
	}
}