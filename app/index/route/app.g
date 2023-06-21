package route

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/unti-io/go-utils/utils"
	"inis/app/facade"
	"inis/app/index/controller"
	debugs "runtime/debug"
)

func Route(Gin *gin.Engine) {

	// 拦截异常
	defer func() {
		if err := recover(); err != nil {
			facade.Log.Error(map[string]any{
				"error":     err,
				"stack":     string(debugs.Stack()),
				"func_name": utils.Caller().FuncName,
				"file_name": utils.Caller().FileName,
				"file_line": utils.Caller().Line,
			}, "首页模板渲染发生错误")
		}
	}()

	go watch(Gin)

	Gin.LoadHTMLGlob("public/index.html")

	// 注册路由
	Gin.GET("/", controller.Index)
}

// watch - 监听 public/index.html 的文件变化
func watch(Gin *gin.Engine) {

	item, err := fsnotify.NewWatcher()
	if err != nil {
		facade.Log.Error(map[string]any{
			"error":     err,
			"stack":     string(debugs.Stack()),
			"func_name": utils.Caller().FuncName,
			"file_name": utils.Caller().FileName,
			"file_line": utils.Caller().Line,
		}, "监听 public/index.html 文件变化发生错误")
		return
	}

	defer func(item *fsnotify.Watcher) {
		err := item.Close()
		if err != nil {
			facade.Log.Error(map[string]any{
				"error":     err,
				"stack":     string(debugs.Stack()),
				"func_name": utils.Caller().FuncName,
				"file_name": utils.Caller().FileName,
				"file_line": utils.Caller().Line,
			}, "监听 public/index.html 文件变化发生错误")
			return
		}
	}(item)

	err = item.Add("public/index.html")
	if err != nil {
		facade.Log.Error(map[string]any{
			"error":     err,
			"stack":     string(debugs.Stack()),
			"func_name": utils.Caller().FuncName,
			"file_name": utils.Caller().FileName,
			"file_line": utils.Caller().Line,
		}, "监听 public/index.html 文件变化发生错误")
		return
	}

	for {
		select {
		case event, ok := <-item.Events:
			if !ok {
				return
			}
			// 重新加载模板
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("index.html 文件发生变化，重新加载模板")
				Gin.LoadHTMLGlob("public/index.html")
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println("index.html 文件被删除，重新加载模板")
			}
		case err, ok := <- item.Errors:
			if !ok {
				return
			}
			facade.Log.Error(map[string]any{
				"error":     err,
				"stack":     string(debugs.Stack()),
				"func_name": utils.Caller().FuncName,
				"file_name": utils.Caller().FileName,
				"file_line": utils.Caller().Line,
			}, "监听 public/index.html 文件变化发生错误")
		}
	}
}