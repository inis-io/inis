package timer

import "inis/app/facade"

type LogStruct struct {}

var Log *LogStruct

func (this *LogStruct) Run() {

	// 1. 每天凌晨 00:01:00 执行
	_ = Timer.Every(1).Day().At("00:01:00").Do(facade.InitLog)
}