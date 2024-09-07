package common

import (
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"strconv"
	"time"
)

const (
	WX = "WX" // 维修
	GZ = "GZ" // 改造
	RW = "RW" // 模具保养任务
	JH = "JH" // 模具保养计划
)

var codeCountMap = make(map[string]int64)

func getCurrentDate() string {
	currentTime := time.Now()
	strTime := currentTime.Format("20060102")
	return strTime
}

func CodeGenInit() {
	date := getCurrentDate()

	// 维修
	var moldRepair model.MoldRepair
	if err := dao.GetConn().Table("mold_repair").Where("`code` like concat('WX',?,'%')", date).Order("id desc").Limit(1).First(&moldRepair).Error; err != nil {
		codeCountMap[WX] = 0
	} else {
		code := moldRepair.Code
		countStr := code[10:]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			panic(err)
		}
		codeCountMap[WX] = int64(count)
	}

	// 改造
	var moldRemodel model.MoldRemodel
	if err := dao.GetConn().Table("mold_remodel").Where("`code` like concat('GZ',?,'%')", date).Order("id desc").Limit(1).First(&moldRemodel).Error; err != nil {
		codeCountMap[GZ] = 0
	} else {
		code := moldRemodel.Code
		countStr := code[10:]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			panic(err)
		}
		codeCountMap[GZ] = int64(count)
	}

	// 计划
	var moldMaintenancePlan model.MoldMaintenancePlan
	if err := dao.GetConn().Table("mold_maintenance_plan").Where("`code` like concat('JH',?,'%')", date).Order("id desc").Limit(1).First(&moldMaintenancePlan).Error; err != nil {
		codeCountMap[JH] = 0
	} else {
		code := moldMaintenancePlan.Code
		countStr := code[10:]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			panic(err)
		}
		codeCountMap[JH] = int64(count)
	}

	// 任务
	var moldMaintenanceTask model.MoldMaintenanceTask
	if err := dao.GetConn().Table("mold_maintenance_task").Where("`code` like concat('RW',?,'%')", date).Order("id desc").Limit(1).First(&moldMaintenanceTask).Error; err != nil {
		codeCountMap[RW] = 0
	} else {
		code := moldMaintenanceTask.Code
		countStr := code[10:]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			panic(err)
		}
		codeCountMap[RW] = int64(count)
	}
}

func GenerateCode(pre string) string {
	date := getCurrentDate()
	count := codeCountMap[pre]
	count++
	codeCountMap[pre] = count
	return fmt.Sprintf("%s%s%04d", pre, date, count)
}
