package model

import (
	"errors"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
)

type Pages struct {
	Id         int    `gorm:"type:int(32); comment:主键;" json:"id"`
	Uid        int    `gorm:"type:int(32); comment:用户ID; default:0;" json:"uid"`
	Key        string `gorm:"size:256; comment:唯一键; default:Null;" json:"key"`
	Title      string `gorm:"size:256; comment:标题; default:Null;" json:"title"`
	Content    string `gorm:"type:longtext; comment:内容; default:Null;" json:"content"`
	Remark     string `gorm:"comment:备注; default:Null;" json:"remark"`
	LastUpdate int64  `gorm:"comment:最后更新时间; default:0;" json:"last_update"`
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
	// 数据库
	DB := facade.NewDB(facade.DBModeMySql)
	// 迁移表
	err := DB.Drive().AutoMigrate(&Pages{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Pages表迁移失败")
		return
	}
}

// AfterSave - 保存后的Hook（包括 create update）
func (this *Pages) AfterSave(tx *gorm.DB) (err error) {

	// key 唯一处理
	if !utils.Is.Empty(this.Key) {
		exist := facade.DB.Model(&Pages{}).Where("id", "!=", this.Id).Where("key", this.Key).Exist()
		if exist {
			return errors.New("key 已存在！")
		}
	}

	return
}
