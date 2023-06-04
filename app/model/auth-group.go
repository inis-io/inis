package model

import (
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
)

type AuthGroup struct {
	Id      int    `gorm:"type:int(32); comment:主键;" json:"id"`
	Name    string `gorm:"comment:权限名称;" json:"name"`
	Uids    string `gorm:"type:text; comment:用户ID;" json:"uids"`
	Root	int    `gorm:"type:int(32); comment:'是否拥有越权限操作数据的能力'; default:0;" json:"root"`
	Rules   string `gorm:"type:text; comment:权限规则;" json:"rules"`
	Default int    `gorm:"type:int(32); comment:默认权限; default:0;" json:"default"`
	Remark  string `gorm:"comment:备注; default:Null;" json:"remark"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitAuthGroup - 初始化AuthGroup表
func InitAuthGroup() {
	// 数据库
	DB := facade.NewDB(facade.DBModeMySql)
	// 迁移表
	err := DB.Drive().AutoMigrate(&AuthGroup{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "AuthGroup表迁移失败")
		return
	}
	// 初始化数据
	count := DB.Model(&AuthGroup{}).Count()
	if count != 0 {
		return
	}
	DB.Model(&AuthGroup{}).Create(&AuthGroup{
		Id: 	 1,
		Name:    "超级管理员",
		Rules:   "all",
		Uids:    "|1|",
		Root: 	 1,
		Default: 1,
		Remark:  "超级管理员，拥有所有权限！",
	})
}
