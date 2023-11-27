package model

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	"strings"
	"sync"
)

type AuthGroup struct {
	Id         int    				 `gorm:"type:int(32); comment:主键;" json:"id"`
	Name       string 				 `gorm:"comment:权限名称;" json:"name"`
	Key        string 				 `gorm:"size:256; comment:唯一键; default:Null;" json:"key"`
	Uids       string 				 `gorm:"type:text; comment:用户ID;" json:"uids"`
	Root	   int    				 `gorm:"type:int(32); comment:'是否拥有越权限操作数据的能力'; default:0;" json:"root"`
	Rules      string 				 `gorm:"type:text; comment:权限规则;" json:"rules"`
	Default    int    				 `gorm:"type:int(32); comment:默认权限; default:0;" json:"default"`
	Pages      string 				 `gorm:"type:text; comment:页面权限; default:Null;" json:"pages"`
	Remark     string 				 `gorm:"comment:备注; default:Null;" json:"remark"`
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
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&AuthGroup{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "AuthGroup表迁移失败")
		return
	}
	// 初始化数据
	count := facade.DB.Model(&AuthGroup{}).Count()
	if count != 0 {
		return
	}
	facade.DB.Model(&AuthGroup{}).Create(&AuthGroup{
		Id: 	 1,
		Name:    "超级管理员",
		Uids:    "|1|",
		Rules:   "all",
		Pages:   "all",
		Root: 	 1,
		Default: 1,
		Remark:  "超级管理员，拥有所有权限！",
	})
}

// AfterSave - 保存后的Hook（包括 create update）
func (this *AuthGroup) AfterSave(tx *gorm.DB) (err error) {

	// key 唯一处理
	if !utils.Is.Empty(this.Key) {
		exist := facade.DB.Model(&AuthGroup{}).WithTrashed().Where("id", "!=", this.Id).Where("key", this.Key).Exist()
		if exist {
			return errors.New("key 已存在！")
		}
	}
	return
}

// AfterFind - 查询Hook
func (this *AuthGroup) AfterFind(tx *gorm.DB) (err error) {

	this.Result = this.result()
	this.Text   = cast.ToString(this.Text)
	this.Json   = utils.Json.Decode(this.Json)
	return
}

// result - 返回结果
func (this *AuthGroup) result() (result map[string]any) {

	var users any
	wg := sync.WaitGroup{}
	wg.Add(1)

	go this.users(&wg, &users)

	wg.Wait()

	return map[string]any{
		"users"   : users,
	}
}

// tags - 标签
func (this *AuthGroup) users(wg *sync.WaitGroup, result *any) {

	defer wg.Done()

	// 标签信息
	tags  := utils.ArrayUnique(utils.ArrayEmpty(strings.Split(this.Uids, "|")))
	*result = facade.DB.Model(&[]Users{}).WhereIn("id", tags).Column("id", "nickname", "avatar", "account")
}

// Auth 应用权限
func (this *AuthGroup) Auth(uid any, group any, isRemove bool) {

	for _, id := range utils.Unity.Ids(group) {

		item := facade.DB.Model(&AuthGroup{}).WithTrashed().Where("id", id).Find()
		if utils.Is.Empty(item) {
			continue
		}

		uids := utils.Unity.Ids(item["uids"])

		// 移除权限
		if isRemove {
			// 判断 uid 是否在数组中
			if utils.InArray(cast.ToInt(uid), cast.ToIntSlice(uids)) {
				// 原生 Go 语言删除数组元素
				for key, val := range uids {
					if cast.ToInt(val) == cast.ToInt(uid) {
						uids = append(uids[:key], uids[key+1:]...)
					}
				}
			}
		} else {
			uids = append(uids, uid)
		}

		// 去重 - 去空
		uids = utils.Unity.Ids(uids)

		var result string
		if len(uids) > 0 {
			result = fmt.Sprintf("|%v|", strings.Join(cast.ToStringSlice(uids), "|"))
		}
		// 更新数据
		facade.DB.Model(&AuthGroup{}).WithTrashed().Where("id", id).Update(map[string]any{
			"uids": result,
		})
	}
}