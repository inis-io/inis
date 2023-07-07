package facade

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/unti-io/go-utils/utils"
	"runtime"
	"strings"
	"time"
)

type CommStruct struct {

}

var Comm *CommStruct

// sn - 获取机器码
func (this *CommStruct) sn() (result string) {

	result, err := machineid.ID()
	if err != nil {
		return utils.Get.Mac()
	}

	return result
}

// Device - 设备信息
func (this *CommStruct) Device() *utils.CurlResponse {

	// 1、把原始的 body 传输进行原样传递
	body := map[string]any{
		"sn":  this.sn(),
		"mac": utils.Get.Mac(),
		"port": map[string]any{
			"run":  Var.Get("port"),
			"real": AppToml.Get("app.port"),
		},
		"domain": Var.Get("domain"),
		"goos":   runtime.GOOS,
		"goarch": runtime.GOARCH,
		"cpu":    runtime.NumCPU(),
	}

	// 2、使用sn和mac进行 Token 16位 对称加密
	token := Token
	key   := utils.Hash.Token(body["sn"], 16, token)
	iv    := utils.Hash.Token(body["mac"], 16, token)

	// 3、接着再对整体参数进行64位的token大写加密
	aes   := utils.AES(key, iv)
	unix  := time.Now().Unix()

	return utils.Curl(utils.CurlRequest{
		Body   : body,
		Method : "POST",
		Headers: map[string]any{
			// X-Khronos(时间戳) - 当前的时间戳
			"X-Khronos": unix,
			// X-Argus(加密文本) - 真实有效的数据
			"X-Argus"  : aes.Encrypt(utils.Json.Encode(body)).Text,
			// X-Gorgon(加密文本)
			"X-Gorgon" : "8642" + utils.Hash.Token(body["sn"], 48, unix),
			// X-SS-STUB(MD5) - 用于检查 body 数据是否被篡改
			"X-SS-STUB": strings.ToUpper(utils.Hash.Token(utils.Map.ToURL(body), 32, unix)),
		},
		Url: Uri + "/dev/device/record",
	}).Send()
}