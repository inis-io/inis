package facade

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"runtime"
	"strings"
	"time"
)

type CommStruct struct {}

var Comm *CommStruct

// Sn - 获取机器序列号
func (this *CommStruct) Sn() (result string) {

	mac := utils.Get.Mac()
	result, err := machineid.ID()
	if err != nil {
		result = mac
	}

	return utils.Hash.Token(result, 32, mac)
}

// Device - 设备信息
func (this *CommStruct) Device() *utils.CurlResponse {

	// 内存信息
	var memory map[string]any

	vm, err := mem.VirtualMemory()
	if err == nil {
		memory = map[string]any{
			"free" : vm.Free,
			"used" : vm.Used,
			"total": vm.Total,
		}
	}

	// 1、把原始的 body 传输进行原样传递
	body := map[string]any{
		"sn":   this.Sn(),
		"mac":  utils.Get.Mac(),
		"port": map[string]any{
			"run":  Var.Get("port"),
			"real": AppToml.Get("app.port"),
		},
		"memory": utils.Json.String(memory),
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

// Signature - 签名算法
func (this *CommStruct) Signature(params map[string]any) (result map[string]any) {

	// 运行端口
	port  := AppToml.Get("app.port")
	// 当前时间戳
	unix  := time.Now().Unix()
	// AES加密密钥
	key   := utils.Hash.Token(port, 16, "AesKey")
	// AES加密向量
	iv    := utils.Hash.Token(unix, 16, "AesIv")

	return map[string]any{
		// X-Khronos(时间戳) - 当前的时间戳
		"X-Khronos": unix,
		// X-Argus(加密文本) - 真实有效的数据
		"X-Argus"  : utils.AES(key, iv).Encrypt(utils.Json.Encode(map[string]any{
			"sn"   : this.Sn(),
			"port" : port,
			"mac"  : utils.Get.Mac(),
		})).Text,
		// X-Gorgon(加密文本)
		"X-Gorgon" : cast.ToString(port) + utils.Hash.Token(this.Sn(), 48, unix),
		// X-SS-STUB(MD5) - 用于检查 body 数据是否被篡改
		"X-SS-STUB": strings.ToUpper(utils.Hash.Token(utils.Map.ToURL(params), 32, unix)),
	}
}

// WithField - 保留指定字段
func (this *CommStruct) WithField(data map[string]any, field any) (result map[string]any) {

	// 参数归一化
	keys := cast.ToStringSlice(utils.Unity.Keys(field))

	// 如果为空，返回原始数据
	if utils.Is.Empty(keys) {
		return data
	}

	return utils.Map.WithField(data, keys)
}