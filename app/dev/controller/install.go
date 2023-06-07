package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"inis/app/facade"
	"inis/app/model"
	"inis/app/validator"
	"runtime"
	"strings"
)

type Install struct {
	// 继承
	base
}

// IGET - GET请求本体
func (this *Install) IGET(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"check": this.check,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPOST - POST请求本体
func (this *Install) IPOST(ctx *gin.Context) {

	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{
		"lock": this.lock,
		"init-db": this.initDB,
		"connect-db": this.connectDB,
		"create-admin": this.createAdmin,
	}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IPUT - PUT请求本体
func (this *Install) IPUT(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// IDEL - DELETE请求本体
func (this *Install) IDEL(ctx *gin.Context) {
	// 转小写
	method := strings.ToLower(ctx.Param("method"))

	allow := map[string]any{}
	err := this.call(allow, method, ctx)

	if err != nil {
		this.json(ctx, nil, facade.Lang(ctx, "方法调用错误：%v", err.Error()), 405)
		return
	}
}

// INDEX - GET请求本体
func (this *Install) INDEX(ctx *gin.Context) {

	// params := this.params(ctx)

	system := map[string]any{
		"GOOS":   runtime.GOOS,
		"GOARCH": runtime.GOARCH,
		"GOROOT": runtime.GOROOT(),
		"NumCPU": runtime.NumCPU(),
		"NumGoroutine": runtime.NumGoroutine(),
		"go": utils.Version.Go(),
		"inis": facade.Version,
		"compare": utils.Version.Compare("v1.0.0", "1 2 0"),
		"agent":  this.header(ctx, "User-Agent"),
	}

	this.json(ctx, map[string]any{
		"system": system,
	}, facade.Lang(ctx, "好的！"), 200)
}

// connectDB - 连接数据库
func (this *Install) connectDB(ctx *gin.Context) {

	// 请求参数
	params := this.params(ctx, map[string]any{
		"hostport": 3306,
		"charset" : "utf8mb4",
		"hostname": "localhost",
	})

	charset  := cast.ToString(params["charset"])
	hostname := cast.ToString(params["hostname"])
	hostport := cast.ToString(params["hostport"])

	if utils.Is.Empty(params["username"]) {
		this.json(ctx, nil, facade.Lang(ctx, "数据库用户名不能为空！"), 400)
		return
	}

	if utils.Is.Empty(params["database"]) {
		this.json(ctx, nil, facade.Lang(ctx, "数据库名不能为空！"), 400)
		return
	}

	if utils.Is.Empty(params["password"]) {
		this.json(ctx, nil, facade.Lang(ctx, "数据库密码不能为空！"), 400)
		return
	}

	username := cast.ToString(params["username"])
	database := cast.ToString(params["database"])
	password := cast.ToString(params["password"])

	// 数据库连接信息
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", username, password, hostname, hostport, database, charset)

	// 使用mysql驱动连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		this.json(ctx, nil, fmt.Sprintf("数据库连接失败：%v", err.Error()), 400)
		return
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		this.json(ctx, nil, fmt.Sprintf("数据库连接失败：%v", err.Error()), 400)
		return
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			this.json(ctx, nil, fmt.Sprintf("数据库连接失败：%v", err.Error()), 400)
			return
		}
	}(sqlDB)

	this.json(ctx, nil, facade.Lang(ctx, "好的！"), 200)

	// 创建配置文件
	utils.File().Save(strings.NewReader(utils.Replace(facade.TempDatabase, map[string]any{
		"${mysql.hostname}": hostname,
		"${mysql.hostport}": hostport,
		"${mysql.username}": username,
		"${mysql.database}": database,
		"${mysql.password}": password,
		"${mysql.charset}" : charset,
		"${mysql.migrate}" : "true",
	})), "config/database.toml")
}

// initDB - 初始化数据库
func (this *Install) initDB(ctx *gin.Context) {
	// 初始化数据库
	facade.WatchDB(false)
	// 初始化数据表
	model.InitTable()
	this.json(ctx, nil, facade.Lang(ctx, "好的！"), 200)
}

// createAdmin - 创建管理员
func (this *Install) createAdmin(ctx *gin.Context) {

	// 表数据结构体
	table  := model.Users{}
	// 请求参数
	params := this.params(ctx)

	// 验证器
	err := validator.NewValid("users", params)

	// 参数校验不通过
	if err != nil {
		this.json(ctx, nil, err.Error(), 400)
		return
	}

	// 帐号不能为空
	if utils.Is.Empty(params["account"]) {
		this.json(ctx, nil, facade.Lang(ctx, "帐号不能为空"), 400)
		return
	}

	// 判断帐号是否已经注册
	if exist := facade.DB.Model(&table).Where("account", params["account"]).Exist(); exist {
		this.json(ctx, nil, facade.Lang(ctx, "该帐号已经注册"), 400)
		return
	}

	// 邮箱不能为空
	if utils.Is.Empty(params["email"]) {
		this.json(ctx, nil, facade.Lang(ctx, "邮箱不能为空"), 400)
		return
	}

	// 判断邮箱是否已经注册
	if exist := facade.DB.Model(&table).Where("email", params["email"]).Exist(); exist {
		this.json(ctx, nil, facade.Lang(ctx, "该邮箱已经注册"), 400)
		return
	}

	if utils.Is.Empty(params["password"]) {
		this.json(ctx, nil, facade.Lang(ctx, "%s 不能为空！", "密码"), 400)
		return
	}

	// 允许存储的字段
	allow := []any{"account", "password", "email", "nickname", "avatar", "description"}
	// 动态给结构体赋值
	for key, val := range params {
		// 加密密码
		if key == "password" {
			val = utils.Password.Create(params["password"])
		}
		// 防止恶意传入字段
		if utils.In.Array(key, allow) {
			utils.Struct.Set(&table, key, val)
		}
	}
	utils.Struct.Set(&table, "pages", "all")

	// 创建用户
	facade.DB.Model(&table).Create(&table)

	jwt := facade.Jwt().Create(facade.H{
		"uid": table.Id,
	})

	// 删除密码
	table.Password = ""

	result := map[string]any{
		"user":  table,
		"token": jwt.Text,
	}

	// 往客户端写入cookie - 存储登录token
	setToken(ctx, jwt.Text)

	go func() {

		uids := []string{cast.ToString(table.Id)}
		group := facade.DB.Model(&model.AuthGroup{}).Find(1)
		if !utils.Is.Empty(group) {
			uids = append(uids, strings.Split(cast.ToString(group["uids"]), "|")...)
		}

		// uids 去重 去空
		uids = cast.ToStringSlice(utils.ArrayUnique(utils.ArrayEmpty(uids)))
		facade.DB.Model(&model.AuthGroup{}).Where("id", 1).Update(&model.AuthGroup{
			Uids: fmt.Sprintf("|%s|", strings.Join(uids, "|")),
		})
	}()

	this.json(ctx, result, facade.Lang(ctx, "注册成功！"), 200)
}

// lock - 上锁（安装锁）
func (this *Install) lock(ctx *gin.Context) {

	if ok := utils.File().Exist("config/database.toml"); !ok {
		this.json(ctx, nil, facade.Lang(ctx, "请先完成数据库配置！"), 400)
		return
	}

	// 删除安装锁
	item := utils.File().Path("install.lock").Remove()

	if item.Error != nil {
		this.json(ctx, nil, facade.Lang(ctx, "解除安装锁失败：%v", item.Error.Error()), 400)
		return
	}

	this.json(ctx, nil, facade.Lang(ctx, "好的！"), 200)
}

// check - 安装锁状态
func (this *Install) check(ctx *gin.Context) {
	this.json(ctx, !utils.File().Exist("install.lock"), facade.Lang(ctx, "好的！"), 200)
}

// 设置登录token到客户的cookie中
func setToken(ctx *gin.Context, token any) {

	host := ctx.Request.Host
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	expire := cast.ToInt(facade.AppToml.Get("jwt.expire", "7200"))
	tokenName := cast.ToString(facade.AppToml.Get("app.token_name", "INIS_LOGIN_TOKEN"))

	ctx.SetCookie(tokenName, cast.ToString(token), expire, "/", host, false, false)
}
