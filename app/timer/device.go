package timer

import (
	"inis/app/facade"
)

type DeviceStruct struct {}

var Device *DeviceStruct

func (this *DeviceStruct) Run() {

	// 先执行一次
	go facade.Comm.Device()

	// 每30分钟执行一次
	_ = Timer.Every(30).Minutes().Do(facade.Comm.Device)
}