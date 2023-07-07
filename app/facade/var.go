package facade

import "github.com/unti-io/go-utils/utils"

// Var - 全局变量
var Var *utils.AsyncStruct[map[string]any]

func init() {
	Var = utils.Async[map[string]any]()
}