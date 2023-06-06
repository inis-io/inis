package model

import (
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
)

type Config struct {
	Id      int    `gorm:"type:int(32); comment:主键;" json:"id"`
	Key   	string `gorm:"size:32; comment:唯一键; default:Null;" json:"key"`
	Value 	string `gorm:"type:text; comment:值; default:Null;" json:"value"`
	Remark  string `gorm:"comment:备注; default:Null;" json:"remark"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitConfig - 初始化Config表
func InitConfig() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Config{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Config表迁移失败")
		return
	}

	// 初始化数据
	go func() {

		configs := []Config{
			{Key: "SYSTEM_API_KEY", Value: "0", Remark: "是否开启API KEY验证"},
			{Key: "SYSTEM_QPS", Value: "1", Json: utils.Json.Encode(map[string]any{
				"point": 10, "global": 50,
			}), Remark: "是否开启接口限流器（QPS）"},
			{Key: "SYSTEM_PAGE_LIMIT", Value: "1", Text: "50", Remark: "限制分页查询的最大条数量"},
		}

		for _, item := range configs {
			// 判断是否存在
			if facade.DB.Model(&Config{}).Where("key", item.Key).Exist() {
				continue
			}
			// 创建数据
			facade.DB.Model(&item).Create(&item)
		}
	}()
}

// AfterFind - 查询Hook
func (this *Config) AfterFind(tx *gorm.DB) (err error) {

	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)

	return
}