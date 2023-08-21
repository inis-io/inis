package timer

import (
	"github.com/jasonlvhit/gocron"
)

var Timer *gocron.Scheduler

func init() {
	Timer = gocron.NewScheduler()
}

func Run() {

	Log.Run()
	Device.Run()

	go func() {
		<- Timer.Start()
	}()
}