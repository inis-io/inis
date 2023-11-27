package model

import (
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
)

type Banner struct {
	Id         int    				 `gorm:"type:int(32); comment:主键;" json:"id"`
	Uid        int    				 `gorm:"type:int(32); comment:用户ID; default:0;" json:"uid"`
	Title      string 				 `gorm:"size:32; comment:标题; default:Null;" json:"title"`
	Content    string 				 `gorm:"comment:内容; default:Null;" json:"content"`
	Url        string 				 `gorm:"size:256; comment:链接; default:Null;" json:"url"`
	Image      string 				 `gorm:"size:256; comment:图片; default:Null;" json:"image"`
	Target     string 				 `gorm:"size:32; comment:打开方式; default:'_blank';" json:"target"`
	StartTime  int64 				 `gorm:"size:32; comment:开始时间; default:Null;" json:"start_time"`
	EndTime    int64 				 `gorm:"size:32; comment:结束时间; default:Null;" json:"end_time"`
	Remark     string 				 `gorm:"comment:备注; default:Null;" json:"remark"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitBanner - 初始化Banner表
func InitBanner() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Banner{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Banner表迁移失败")
		return
	}
}

// AfterFind - 查询Hook
func (this *Banner) AfterFind(*gorm.DB) (err error) {
	// 替换 url 中的域名
	this.Image = utils.Replace(this.Image, DomainTemp1())
	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)
	return
}

// AfterSave - 保存后的Hook（包括 create update）
func (this *Banner) AfterSave(tx *gorm.DB) (err error) {

	go func() {
		this.Image = utils.Replace(this.Image, DomainTemp2())
		tx.Model(this).UpdateColumn("image", this.Image)
	}()

	return
}
