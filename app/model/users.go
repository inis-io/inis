package model

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"inis/app/facade"
	regexp2 "regexp"
	"strings"
	"sync"
)

type Users struct {
	Id          int    				  `gorm:"type:int(32); comment:主键;" json:"id"`
	Account     string 				  `gorm:"size:32; comment:帐号; default:Null;" json:"account"`
	Password    string 				  `gorm:"comment:密码;" json:"password"`
	Nickname    string 				  `gorm:"size:32; comment:昵称;" json:"nickname"`
	Email       string 				  `gorm:"size:128; comment:邮箱;" json:"email"`
	Phone       string 				  `gorm:"size:32; comment:手机号;" json:"phone"`
	Avatar      string 				  `gorm:"comment:头像; default:Null;" json:"avatar"`
	Description string 				  `gorm:"comment:描述; default:Null;" json:"description"`
	Title       string 				  `gorm:"comment:头衔; default:Null;" json:"title"`
	Gender		string 				  `gorm:"comment:性别; default:Null;" json:"gender"`
	Exp  		int    				  `gorm:"type:int(32); comment:经验值; default:0;" json:"exp"`
	Source      string 				  `gorm:"size:32; default:'default'; comment:注册来源;" json:"source"`
	Remark      string 				  `gorm:"comment:备注; default:Null;" json:"remark"`
	// 以下为公共字段
	Json        any                   `gorm:"type:longtext; comment:用于存储JSON数据;" json:"json"`
	Text        any                   `gorm:"type:longtext; comment:用于存储文本数据;" json:"text"`
	Result      any                   `gorm:"type:varchar(256); comment:不存储数据，用于封装返回结果;" json:"result"`
	LoginTime   int64                 `gorm:"size:32; comment:登录时间; default:Null;" json:"login_time"`
	CreateTime  int64                 `gorm:"autoCreateTime; comment:创建时间;" json:"create_time"`
	UpdateTime  int64                 `gorm:"autoUpdateTime; comment:更新时间;" json:"update_time"`
	DeleteTime  soft_delete.DeletedAt `gorm:"comment:删除时间; default:0;" json:"delete_time"`
}

// InitUsers - 初始化Users表
func InitUsers() {
	// 迁移表
	err := facade.DB.Drive().AutoMigrate(&Users{})
	if err != nil {
		facade.Log.Error(map[string]any{"error": err}, "Users表迁移失败")
		return
	}
}

// AfterFind - 查询后的钩子
func (this *Users) AfterFind(tx *gorm.DB) (err error) {

	if utils.Is.Empty(this.Avatar) {

		// 正则匹配邮箱 [1-9]\d+@qq.com 是否匹配
		reg := regexp2.MustCompile(`[1-9]\d+@qq.com`).MatchString(this.Email)
		if reg {

			// 获取QQ号
			qq := regexp2.MustCompile(`[1-9]\d+`).FindString(this.Email)
			this.Avatar = "https://q1.qlogo.cn/g?b=qq&nk=" + qq + "&s=100"

		} else {
			avatars := utils.File(utils.FileRequest{
				Ext:    ".png, .jpg, .jpeg, .gif",
				Dir:    "public/assets/rand/avatar/",
				Domain: fmt.Sprintf("%v/", facade.Var.Get("domain")),
				Prefix: "public/",
			}).List()

			// 随机获取头像
			if len(avatars.Slice) > 0 {
				this.Avatar = cast.ToString(avatars.Slice[utils.Rand.Int(0, len(avatars.Slice)-1)])
			}
		}
	}

	// 替换 url 中的域名
	this.Avatar = utils.Replace(this.Avatar, DomainTemp1())

	this.Result = this.result()
	this.Text   = cast.ToString(this.Text)
	this.Json   = utils.Json.Decode(this.Json)
	return
}

// AfterSave - 保存后的Hook（包括 create update）
func (this *Users) AfterSave(tx *gorm.DB) (err error) {

	go func() {
		this.Avatar = utils.Replace(this.Avatar, DomainTemp2())
		tx.Model(this).UpdateColumn("avatar", this.Avatar)
	}()

	// 账号 唯一处理
	if !utils.Is.Empty(this.Account) {
		exist := facade.DB.Model(&Users{}).WithTrashed().Where("id", "!=", this.Id).Where("account", this.Account).Exist()
		if exist {
			return errors.New("账号已存在！")
		}
	}

	// 邮箱 唯一处理
	if !utils.Is.Empty(this.Email) {
		exist := facade.DB.Model(&Users{}).WithTrashed().Where("id", "!=", this.Id).Where("email", this.Email).Exist()
		if exist {
			return errors.New("邮箱已存在！")
		}
	}

	// 手机号 唯一处理
	if !utils.Is.Empty(this.Phone) {
		exist := facade.DB.Model(&Users{}).WithTrashed().Where("id", "!=", this.Id).Where("phone", this.Phone).Exist()
		if exist {
			return errors.New("手机号已存在！")
		}
	}

	return
}

// Rules - 生成用户权限列表
func (this *Users) Rules(uid any) (slice []any) {

	// 生成规则列表
	item := func(uid any) (slice []any) {

		var table []AuthGroup

		// 从规则分组里面查找
		group := facade.DB.Model(&table).Like("uids", "%|"+cast.ToString(uid)+"|%").Select()

		var hashes []any

		for _, item := range group {
			// 判断字符串中是否包含all
			if strings.Contains(cast.ToString(item["rules"]), "all") {
				hashes = append(hashes, "all")
				continue
			}
			// 逗号分隔数组
			list := strings.Split(cast.ToString(item["rules"]), ",")
			for _, val := range list {
				hashes = append(hashes, val)
			}
		}

		var list []any
		var rules []map[string]any
		var AuthRules []AuthRules

		if utils.Is.Empty(hashes) {
			return list
		}

		// 判断是否拥有全部权限
		if utils.In.Array("all", hashes) {

			rules = facade.DB.Model(&AuthRules).Select()

		} else {

			// hashes 去重 去空
			hashes = utils.Array.Empty(utils.Array.Unique(hashes))
			rules = facade.DB.Model(&AuthRules).WhereIn("hash", hashes).Select()
		}

		// 扁平化
		for _, item := range rules {
			list = append(list, fmt.Sprintf("[%v][%v]", item["method"], item["route"]))
		}

		// 去重 去空
		return utils.Array.Empty(utils.Array.Unique(list))
	}

	// 用户组缓存
	cacheName  := fmt.Sprintf("user[%v][rule-group]", uid)
	// 缓存状态
	cacheState := cast.ToBool(facade.CacheToml.Get("open"))

	var rules []any
	// 已经登录的用户 - 检查缓存中是否存在该用户的权限 - 存在则直接返回
	if cacheState && facade.Cache.Has(cacheName) {
		return cast.ToSlice(facade.Cache.Get(cacheName))
	}

	rules = item(uid)

	if cacheState {
		go facade.Cache.Set(cacheName, rules, 0)
	}

	return rules
}

// result - 返回结果
func (this *Users) result() (result map[string]any) {

	var auth, level any

	wg := sync.WaitGroup{}
	wg.Add(2)

	go this.auth(&wg, &auth)
	go this.level(&wg, &level)

	wg.Wait()

	return map[string]any{
		"auth": auth,
		"level": level,
	}
}

// auth - 解析用户权限
func (this *Users) auth(wg *sync.WaitGroup, result *any) {

	defer wg.Done()

	// 查询自己拥有的权限
	group := facade.DB.Model(&AuthGroup{}).Like("uids", "%|"+cast.ToString(this.Id)+"|%").Column("id", "rules", "name", "root", "pages", "key")

	var ids []int
	var rules []string
	var pages []string

	for _, val := range cast.ToSlice(group) {
		item := cast.ToStringMap(val)
		ids   = append(ids, cast.ToInt(item["id"]))
		// 逗号分隔的权限
		rules = append(rules, strings.Split(cast.ToString(item["rules"]), ",")...)
		// 逗号分隔的页面
		pages = append(pages, strings.Split(cast.ToString(item["pages"]), ",")...)
	}

	// 去重 去空
	rules = utils.Array.Filter(cast.ToStringSlice(utils.ArrayUnique[string](rules)))
	pages = utils.Array.Filter(cast.ToStringSlice(utils.ArrayUnique[string](pages)))

	*result = map[string]any{
		"all"  : utils.InArray("all", rules),
		"group": map[string]any{
			"ids": ids,
			"list": group,
		},
		"pages": map[string]any{
			"hash": pages,
		},
		"rules": map[string]any{
			"hash": rules,
		},
	}
}

// level - 解析用户等级
func (this *Users) level(wg *sync.WaitGroup, result *any) {

	defer wg.Done()

	// 查询字段
	field := []string{"name", "value", "description", "exp", "text", "json"}

	// 查询当前等级
	item1 := facade.DB.Model(&Level{}).Field(field).Limit(1).Where("exp", "<=", this.Exp).Order("exp desc")

	// 查询下一等级
	item2 := facade.DB.Model(&Level{}).Field(field).Limit(1).Where("exp", ">", this.Exp).Order("exp asc")

	currents := cast.ToSlice(item1.Column())
	var current any
	if len(currents) > 0 {
		current = currents[0]
	}

	nexts    := cast.ToSlice(item2.Column())
	var next any
	if len(nexts) > 0 {
		next = nexts[0]
	}

	*result = map[string]any{
		"current": current,
		"next"   : next,
	}
}

// Destroy - 注销后，清空用户数据
func (this *Users) Destroy(uid any) {

	// 清空权限
	if ids := facade.DB.Model(&[]AuthGroup{}).WithTrashed().Like("uids", "|" + cast.ToString(uid) + "|").Column("id"); !utils.Is.Empty(ids) {
		go (&AuthGroup{}).Auth(uid, ids, true)
	}

	// 表名
	tables := []any{
		Article{},	// 文章
		Comment{},	// 评论
		EXP{},		// 经验值
		Links{},	// 友链
		Pages{},	// 页面
		Banner{},	// 轮播
	}

	for _, table := range tables {
		go facade.DB.Model(&table).WithTrashed().Where("uid", uid).Delete()
	}
}