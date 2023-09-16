package model

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	"sync"
	"time"
)

type EXP struct {
	Id          int    `gorm:"type:int(32); comment:主键;" json:"id"`
	Uid 	    int    `gorm:"type:int(32); comment:用户ID;" json:"uid"`
	Value  		int    `gorm:"type:int(32); comment:经验值; default:0;" json:"value"`
	Type 		string `gorm:"comment:类型; default:'default';" json:"type"`
	BindType 	string `gorm:"comment:绑定类型; default:'default';" json:"bind_type"`
	BindId 		int    `gorm:"type:int(32); comment:绑定ID; default:0;" json:"bind_id"`
	State 		int    `gorm:"type:int(32); comment:状态; default:1;" json:"state"`
	Description string `gorm:"comment:描述; default:Null;" json:"description"`
	// 以下为公共字段
	Json       any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text       any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result     any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	CreateTime int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitEXP - 初始化EXP表
func InitEXP() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&EXP{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "EXP表迁移失败")
		return
	}
}

// AfterFind - 查询Hook
func (this *EXP) AfterFind(tx *gorm.DB) (err error) {

	this.Result = this.result()
	this.Text   = cast.ToString(this.Text)
	this.Json   = utils.Json.Decode(this.Json)
	return
}

// Add - 增加经验值
func (this *EXP) Add(table EXP) (err error) {

	// 拦截异常
	defer func() {
		if bug := recover(); bug != nil {
			facade.Log.Error(map[string]any{"error": bug}, "增加经验值失败")

		}
	}()

	if table.Uid == 0 {
		return errors.New("请先登录！")
	}

	limit := map[string][]any{
		"like":     {"点赞", 1, 10},	// 点赞 - 每天10次，一次1经验值
		"collect":  {"收藏", 1, 10},	// 收藏 - 每天10次，一次1经验值
		"visit":    {"访问", 1, 10},	// 访问 - 每天10次，一次1经验值
		"share":    {"分享", 1, 10},	// 分享 - 每天10次，一次1经验值
		"login":    {"登录", 5, 1},	// 登录 - 每天1次，一次5经验值
		"comment":  {"评论", 1, 10},	// 评论 - 每天10次，一次1经验值
		"check-in": {"签到", 30, 1},	// 签到 - 每天1次，一次30经验值
	}

	// 今天开始的时间戳
	now   := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 从数据库里面找到今天的记录
	count := facade.DB.Model(&EXP{}).Where([]any{
		[]any{"uid", "=", table.Uid},
		[]any{"type", "=", table.Type},
		[]any{"create_time", ">=", today.Unix()},
	}).Count()

	// 检查 limit[table.Type] 是否存在
	if _, ok := limit[table.Type]; !ok {
		return errors.New("未知的经验值类型！")
	}

	// 如果超过了限制，不增加经验值
	if count >= cast.ToInt64(limit[table.Type][2]) {
		switch table.Type {
		case "check-in":
			return errors.New("今天已经签到过了！")
		default:
			return errors.New(fmt.Sprintf("本日%s奖励经验值次数已经用完了！", limit[table.Type][0]))
		}
	}

	// 每次增加的经验值
	table.Value = cast.ToInt(limit[table.Type][1])
	if utils.Is.Empty(table.Description) {
		table.Description = fmt.Sprintf("%s奖励", limit[table.Type][0])
	}

	// 添加经验日志
	tx := facade.DB.Model(&EXP{}).Create(&table)

	if tx.Error != nil {
		return tx.Error
	}

	// 增加经验值
	facade.DB.Model(&Users{}).Where("id", table.Uid).Inc("exp", table.Value)

	return err
}

// result - 返回结果
func (this *EXP) result() (result map[string]any) {

	var author any
	wg := sync.WaitGroup{}
	wg.Add(1)

	go this.author(&wg, &author)

	wg.Wait()

	return map[string]any{
		"author"   : author,
	}
}

// tags - 标签
func (this *EXP) author(wg *sync.WaitGroup, result *any) {

	defer wg.Done()

	// 作者信息
	author := make(map[string]any)
	allow  := []string{"id", "nickname", "avatar", "description", "result", "title"}
	user   := facade.DB.Model(&Users{}).Find(this.Uid)

	if !utils.Is.Empty(user) {
		author = utils.Map.WithField(user, allow)
	}

	*result = author
}