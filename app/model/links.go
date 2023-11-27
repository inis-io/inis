package model

import (
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
)

type Links struct {
	Id          int    				  `gorm:"type:int(32); comment:主键;" json:"id"`
	Uid         int    				  `gorm:"type:int(32); comment:用户ID; default:0;" json:"uid"`
	Nickname    string 				  `gorm:"size:32; comment:昵称; default:Null;" json:"nickname"`
	Description string 				  `gorm:"comment:描述; default:Null;" json:"description"`
	Url         string 				  `gorm:"size:256; comment:链接; default:Null;" json:"url"`
	Avatar      string 				  `gorm:"size:256; comment:头像; default:Null;" json:"avatar"`
	Target      string 				  `gorm:"size:32; comment:打开方式; default:'_blank';" json:"target"`
	Check 	 	string 				  `gorm:"size:32; comment:审核; default:'wait';" json:"check"`
	State       string 				  `gorm:"size:32; comment:状态; default:'wait';" json:"state"`
	Remark      string 				  `gorm:"comment:备注; default:Null;" json:"remark"`
	Group       int    				  `gorm:"size:32; comment:分组; default:0;" json:"group"`
	// 以下为公共字段
	Json        any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text        any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result      any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime  int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime  int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime  soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitLinks - 初始化Links表
func InitLinks() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Links{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Links表迁移失败")
		return
	}

	// 初始化数据
	go initLinksData()
}

// AfterFind - 查询Hook
func (this *Links) AfterFind(tx *gorm.DB) (err error) {

	// 替换 url 中的域名
	this.Avatar = utils.Replace(this.Avatar, DomainTemp1())

	group := map[string]any{
		"id":          0,
		"name":        "默认分组",
		"avatar":      "",
		"description": "默认分组",
	}

	if this.Group != 0 {

		item := facade.DB.Model(&LinksGroup{}).Find(this.Group)

		if !utils.Is.Empty(item) {

			group = utils.Map.WithField(item, []string{"id", "avatar", "name", "description"})

		} else {

			// 如果分组不存在，则将分组设置为默认分组
			tx.Model(this).UpdateColumn("group", 0)
		}
	}

	// 封装返回结果
	this.Result = map[string]any{
		"group": group,
	}
	this.Text = cast.ToString(this.Text)
	this.Json = utils.Json.Decode(this.Json)

	return
}

// AfterSave - 保存后的Hook（包括 create update）
func (this *Links) AfterSave(tx *gorm.DB) (err error) {

	go func() {
		this.Avatar = utils.Replace(this.Avatar, DomainTemp2())
		tx.Model(this).UpdateColumn("avatar", this.Avatar)
	}()

	return
}

// initLinksData - 初始化Links表数据
func initLinksData() {

	// 如果数据表中有数据，则不进行初始化
	if facade.DB.Model(&Links{}).Count() != 0 {
		return
	}

	array := []Links{
		{
			Nickname:    "兔子",
			Description: "许一人，以偏爱，尽此生，之慷慨！",
			Url:         "https://inis.cn",
			Avatar:      "https://q.qlogo.cn/g?b=qq&nk=97783391&s=640",
			Remark:      "如果可以，请不要删除我！开发不易，感谢支持！",
		},
	}

	// 创建数据
	for _, item := range array {
		facade.DB.Model(&item).Create(&item)
	}
}