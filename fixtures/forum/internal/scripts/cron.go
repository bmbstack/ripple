package scripts

import (
	"fmt"
	. "github.com/bmbstack/ripple/helper"
	"github.com/go-co-op/gocron"
	"runtime"
	"time"
)

func RunCron() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	loc, _ := time.LoadLocation(TimeLocationName)
	cron := gocron.NewScheduler(loc)
	_, _ = cron.Every(1).Second().Do(runTask)
	cron.StartBlocking()
}

func runTask() {
	fmt.Println("===============run task, time: ", time.Now().Format(DateFullLayout))
}
