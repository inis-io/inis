package model

import (
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	"sync"
)

type Comment struct {
	Id         int     				 `gorm:"type:int(32); comment:主键;" json:"id"`
	Pid		   int     				 `gorm:"type:int(32); comment:父级ID; default:0;" json:"pid"`
	Uid		   int     				 `gorm:"type:int(32); comment:用户ID; default:0;" json:"uid"`
	Content    string  				 `gorm:"type:varchar(1024); comment:内容; default:Null;" json:"content"`
	Ip         string  				 `gorm:"comment:IP; default:Null;" json:"ip"`
	Agent      string  				 `gorm:"type:varchar(512); comment:浏览器信息; default:Null;" json:"agent"`
	BindId 	   int     				 `gorm:"type:int(32); comment:绑定ID; default:0;" json:"bind_id"`
	BindType   string  				 `gorm:"comment:绑定类型; default:'article';" json:"bind_type"`
	Editor     string 				 `gorm:"comment:编辑器; default:'text';" json:"editor"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitComment - 初始化Comment表
func InitComment() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Comment{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Comment表迁移失败")
		return
	}
}

// AfterFind - 查询Hook
func (this *Comment) AfterFind(tx *gorm.DB) (err error) {

	this.Result = this.result()
	this.Text   = cast.ToString(this.Text)
	this.Json   = utils.Json.Decode(this.Json)
	return
}

// result - 返回结果
func (this *Comment) result() (result map[string]any) {

	var page, author, article any

	wg := sync.WaitGroup{}
	wg.Add(3)

	go this.page(&wg, &page)
	go this.author(&wg, &author)
	go this.article(&wg, &article)

	wg.Wait()

	return map[string]any{
		"page":   page,
		"author": author,
		"article": article,
	}
}

// author - 解析作者信息
func (this *Comment) author(wg *sync.WaitGroup, result *any) {

	defer wg.Done()

	// 作者信息
	user := facade.DB.Model(&Users{}).Find(this.Uid)
	*result = utils.Map.WithField(user, []string{"id", "nickname", "avatar", "description", "result"})
}

// article - 解析文章信息
func (this *Comment) article(wg *sync.WaitGroup, result *any) {

	defer wg.Done()

	if this.BindType != "article" {
		return
	}

	*result = utils.Map.WithField(facade.DB.Model(&Article{}).Find(this.BindId), []string{"id", "title"})
}

// page - 解析页面信息
func (this *Comment) page(wg *sync.WaitGroup, result *any) {

	defer wg.Done()

	if this.BindType != "page" {
		return
	}

	*result = utils.Map.WithField(facade.DB.Model(&Pages{}).Find(this.BindId), []string{"id", "key", "title"})
}