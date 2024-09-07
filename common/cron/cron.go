/**
 * gocron 工具类
 */

package cron

import (
	"github.com/go-co-op/gocron"
	"time"
)

var CronJob *gocron.Scheduler

func Init() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	CronJob = gocron.NewScheduler(timezone)
	CronJob.StartAsync()
}
