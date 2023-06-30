package model

import (
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
)

type Level struct {
	Id          int    `gorm:"type:int(32); comment:主键;" json:"id"`
	Name        string `gorm:"size:32; comment:名称; default:'LV0';" json:"name"`
	Value 	    int    `gorm:"type:int(32); comment:等级值; default:0;" json:"value"`
	Description string `gorm:"comment:描述; default:Null;" json:"description"`
	Exp  		int    `gorm:"type:int(32); comment:经验值; default:0;" json:"exp"`
	Remark      string `gorm:"comment:备注; default:Null;" json:"remark"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitLevel - 初始化Level表
func InitLevel() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Level{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Level表迁移失败")
		return
	}

	// 初始化数据
	go initLevelData()
}

// AfterFind - 查询Hook
func (this *Level) AfterFind(tx *gorm.DB) (err error) {

	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)

	return
}

// initLevelData - 初始化Level表数据
func initLevelData() {

	// 如果数据表中有数据，则不进行初始化
	if facade.DB.Model(&Level{}).Count() != 0 {
		return
	}

	array := []Level{
		{
			Value: 0,
			Name: "新手",
			Description: "刚开始接触学习的人",
			Exp: 0,
		},
		{
			Value: 1,
			Name: "入门",
			Description: "已经有一定学习经验的人",
			Exp: 2000,
		},
		{
			Value: 2,
			Name: "爱好者",
			Description: "喜欢学习并持续探索的人",
			Exp: 4000,
		},
		{
			Value: 3,
			Name: "专家",
			Description: "在某个领域拥有一定学习经验和专业知识的人",
			Exp: 6000,
		},
		{
			Value: 4,
			Name: "领袖",
			Description: "对某个领域拥有广泛学识和深度见解的人",
			Exp: 10000,
		},
		{
			Value: 5,
			Name: "导师",
			Description: "在某个领域掌握精深，并能分享经验和引领他人的人",
			Exp: 15000,
		},
		{
			Value: 6,
			Name: "大师",
			Description: "对某个领域拥有深厚学识、长期经验积累，并能以高超的能力指导他人的人",
			Exp: 20000,
		},
	}

	// 创建数据
	for _, item := range array {
		facade.DB.Model(&item).Create(&item)
	}
}