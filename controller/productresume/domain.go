package productresume

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/common/cron"
	"crf-mold/controller/email"
	"crf-mold/controller/message"
	"crf-mold/controller/user"
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type productionResume struct {
	topic string
	mda   MDA
}

func NewProductResume(topic string, mda MDA) *productionResume {
	return &productionResume{
		topic: topic,
		mda:   mda,
	}
}

func (pr *productionResume) getLineLevel() string {
	columns := strings.Split(pr.topic, "/")
	return columns[2]
}

func (pr *productionResume) getOrderCode() string {
	return fmt.Sprintf("%d", pr.mda.ShiftId)
}

func (pr *productionResume) getTotalCount() int64 {
	var totalCount int64

	for i := 0; i < len(pr.mda.PartList); i++ {
		totalCount += pr.mda.PartList[i].QtyOK
		totalCount += pr.mda.PartList[i].QtyNOK
	}

	return totalCount
}

func getChongKongMoldsByPartCode(tx *gorm.DB, partCode string) []model.MoldInfo {
	var molds []model.MoldInfo
	sql := `
		select 
			mi.ID,
			mi.code,
			mi.name,
			mi.flush_count
		from mold_info mi
		inner join mold_part_rel mpr on mpr.mold_code = mi.code and mpr.is_deleted = 'N'
		where mi.is_deleted = 'N' and mpr.part_code = ? and mi.type = 'chongkong'
		group by mi.ID,mi.flush_count
	`
	tx.Raw(sql, partCode).Find(&molds)
	return molds
}

func getChengXingMoldsByLineLevel(tx *gorm.DB, lineLevel string) []MoldMaintenanceInfo {
	var molds []MoldMaintenanceInfo
	tx.Table("mold_info").Where("is_deleted = 'N' and status = 'zhengchang' and type = 'chengxing' and line_level = ?", lineLevel).Find(&molds)
	return molds
}

func GenerateMoldTask(tx *gorm.DB, molds []model.MoldInfo, plans map[int64]model.MoldMaintenancePlan) (newTask []model.MoldMaintenanceTask, timeoutTask []model.MoldMaintenanceTask) {
	for i := 0; i < len(molds); i++ {
		moldInfo := molds[i]
		plan, ok := plans[moldInfo.ID]
		// 不存在正在运行的 不是正在运行 则跳过
		if !ok {
			continue
		}
		// 生成第i个区间的任务		TaskStandard=保养计划中的保养冲次
		//n := int64(math.Ceil(float64(moldInfo.CalcFlushCount) / float64(plan.TaskStandard*10000)))
		n := int64(math.Ceil(float64(moldInfo.CalcFlushCount) / float64(plan.TaskStandard)))
		for ; n > 0; n-- {
			// TaskStandard 保养冲次    TaskStart提取生成任务的冲次
			//interval := n * plan.TaskStandard * 10000 // 第n个区间的标准
			//min := interval - plan.TaskStart*10000    // 第n个区间的最小
			interval := n * plan.TaskStandard
			min := interval - plan.TaskStart
			if moldInfo.CalcFlushCount >= min { // 可以生成该区间的任务
				// 是否任务已经生成了
				var count int64
				//tx.Table("mold_maintenance_task").Where("mold_id = ? and task_gen_interval >= ? and task_gen_interval < ?", moldInfo.ID, interval, 99999*10000).Count(&count)
				// 只要存在未完成的任务 则表示有任务正在运行 则不需要生成
				tx.Table("mold_maintenance_task").Where("mold_id = ? and status != ? and is_deleted = 'N'", moldInfo.ID, constant.COMPLETE).Count(&count)
				if count > 0 { // 生成了就跳过了，不继续向下找了
					break
				}

				// 生成任务
				var task model.MoldMaintenanceTask
				task.Code = common.GenerateCode("RW")
				task.MoldID = moldInfo.ID
				task.MoldMaintenancePlainID = plan.ID
				task.Status = constant.WAIT
				task.PlanType = constant.PLAN_METERING_TYPE
				task.TaskGenInterval = interval
				//task.TimeoutCount = interval + plan.TaskEnd*10000
				task.TimeoutCount = interval + plan.TaskEnd
				task.UpdatedBy = "system"
				task.CreatedBy = "system"
				tx.Table("mold_maintenance_task").Create(&task)
				newTask = append(newTask, task)

				// 添加模具任务生成的生命周期记录
				var lifecycleRecord model.MoldMaintenanceTaskLifecycle
				lifecycleRecord.MoldMaintenanceTaskID = task.ID
				lifecycleRecord.Title = constant.TASKGENERATE // 任务生成
				lifecycleRecord.CreatedBy = "system"
				lifecycleRecord.UpdatedBy = "system"
				lifecycleRecord.Time = base.Now().String()
				tx.Table("mold_maintenance_task_lifecycle").Create(&lifecycleRecord)
			}
		}

		// 把前面区间进行中的任务都设置为超时
		var tasks []model.MoldMaintenanceTask

		tx.Table("mold_maintenance_task").Where("is_deleted = 'N' and mold_id = ? and timeout_count <= ? and status = ?",
			moldInfo.ID, moldInfo.CalcFlushCount, constant.WAIT).Find(&tasks)

		batchSave := make([]model.MoldMaintenanceTaskLifecycle, len(tasks))
		for i := 0; i < len(tasks); i++ {
			tx.Table("mold_maintenance_task").Where("id = ?", tasks[i].ID).Updates(&model.MoldMaintenanceTask{
				Status: constant.TIMEOUT,
			})
			timeoutTask = append(timeoutTask, tasks[i])

			// 添加模具任务超时的生命周期记录
			var lifecycleRecord model.MoldMaintenanceTaskLifecycle
			lifecycleRecord.MoldMaintenanceTaskID = tasks[i].ID
			lifecycleRecord.Title = constant.TASKTIMEOUT // 任务超时
			lifecycleRecord.CreatedBy = "system"
			lifecycleRecord.UpdatedBy = "system"
			lifecycleRecord.Time = base.Now().String()
			batchSave[i] = lifecycleRecord
		}

		tx.Table("mold_maintenance_task_lifecycle").CreateInBatches(batchSave, len(batchSave))
	}

	return
}

// 添加生产履历记录
func (pr *productionResume) SaveProductionResume(tx *gorm.DB) {
	// 删除工单号对应的履历
	tx.Table("mold_product_resume").Where("order_code = ?", pr.getOrderCode()).Updates(&model.MoldProductResume{
		IsDeleted: "Y",
	})

	// 所有该线别的模具
	molds := getChengXingMoldsByLineLevel(tx, pr.getLineLevel())
	var batchSave []model.MoldProductResume
	for i := 0; i < len(molds); i++ {
		mold := molds[i]
		batchSave = append(batchSave, model.MoldProductResume{
			MoldType:     mold.Type,
			MoldId:       mold.ID,
			MoldCode:     mold.Code,
			MoldName:     mold.Name,
			OrderCode:    pr.getOrderCode(),
			LineLevel:    pr.getLineLevel(),
			Count:        pr.getTotalCount(),
			IsChecked:    "N",
			CompleteTime: base.Now().String(),
		})
	}

	// 所有拥有零件号的模具
	partList := pr.mda.PartList
	for i := 0; i < len(partList); i++ {
		partCode := partList[i].PartNo
		count := partList[i].QtyOK + partList[i].QtyNOK
		molds := getChongKongMoldsByPartCode(tx, partCode)
		for i := 0; i < len(molds); i++ {
			mold := molds[i]
			batchSave = append(batchSave, model.MoldProductResume{
				MoldType:     mold.Type,
				MoldId:       mold.ID,
				MoldCode:     mold.Code,
				MoldName:     mold.Name,
				OrderCode:    pr.getOrderCode(),
				LineLevel:    pr.getLineLevel(),
				PartCode:     partCode,
				Count:        count,
				IsChecked:    "N",
				CompleteTime: base.Now().String(),
				CreatedBy:    "system",
				UpdatedBy:    "system",
			})
		}
	}

	tx.Table("mold_product_resume").CreateInBatches(batchSave, len(batchSave))
}

// SaveLineProductionResume 添加产线生产履历记录
func (pr *productionResume) SaveLineProductionResume(tx *gorm.DB) {
	shiftId := pr.mda.ShiftId
	// 删除班次ID对应的履历
	tx.Table("line_product_resume").Where("shift_id = ?", shiftId).Updates(&model.LineProductResume{
		IsDeleted: "Y",
	})
	var batchSave []model.LineProductResume

	// 提取零件号
	partList := pr.mda.PartList
	for i := 0; i < len(partList); i++ {
		batchSave = append(batchSave, model.LineProductResume{
			LineLevel: pr.getLineLevel(),
			ShiftId:   int(shiftId),
			StartTime: pr.mda.StartTime,
			EndTime:   pr.mda.EndTime,
			ShiftNo:   int16(pr.mda.ShiftNo),
			ShiftMin:  int(pr.mda.ShiftMin),
			ShiftDate: pr.mda.ShiftDate,
			PartCode:  partList[i].PartNo,
			QtyOk:     int(partList[i].QtyOK),
			QtyNOk:    int(partList[i].QtyNOK),
		})
	}

	tx.Table("line_product_resume").CreateInBatches(batchSave, len(batchSave))
}

// AddMoldFlushCount 添加模具冲次
func AddMoldFlushCount(tx *gorm.DB, items []PendingCheckProductResumeMoldInfo) {
	// 以模具编号作为key，累计计算每个模具的最终次数
	addedCounts := make(map[int64]int64)
	for _, item := range items {
		addedCounts[item.ID] += item.ResumeCount
	}
	for _, item := range items {
		//tx.Table("mold_info").Where("id = ?", item.ID).Updates(&model.MoldInfo{
		//	FlushCount:     item.FlushCount,
		//	CalcFlushCount: item.CalcFlushCount,
		//})
		// 产量和计算产量均需要加该工单的数据
		tx.Table("mold_info").Where("id = ?", item.ID).Updates(&model.MoldInfo{ //原代码在这
			FlushCount:     item.FlushCount + item.ResumeCount,
			CalcFlushCount: item.CalcFlushCount + item.ResumeCount,
		})
		// 更新检查标识
		tx.Table("mold_product_resume").Where("id = ?", item.ResumeID).Updates(&model.MoldProductResume{
			IsChecked: "Y",
		})
		logrus.Infof("AddMoldFlushCount 模具编号:%s 单号:%s 产量:%d 原先总冲次:%d", item.Code, item.OrderCode, item.ResumeCount, item.CalcFlushCount)
	}
}

func getCheckMolds(tx *gorm.DB, items []PendingCheckProductResumeMoldInfo) ([]model.MoldInfo, []int64) {
	var molds []model.MoldInfo
	m := make(map[int64]bool)
	var moldIds []int64
	// 去重
	for _, item := range items {
		if !m[item.ID] {
			m[item.ID] = true
			moldIds = append(moldIds, item.ID)
		}
	}
	if len(moldIds) > 0 {
		tx.Table("mold_info").Where("id in ? and is_deleted = 'N' and status = 'zhengchang' and type = ? ", moldIds, "chongkong").Find(&molds)
	}

	return molds, moldIds
}

// 筛选正常的冲孔的模具
func getCheckMoldsV1(tx *gorm.DB) ([]model.MoldInfo, []int64) {
	var molds []model.MoldInfo
	var moldIds []int64
	// 不以外部的items来筛选数据 直接筛选全部
	tx.Table("mold_info").Where("is_deleted = 'N' and status = 'zhengchang' and type = ? ", "chongkong").Find(&molds)
	for _, mold := range molds {
		moldIds = append(moldIds, mold.ID)
	}
	return molds, moldIds
}

func getRunningMeteringPlans(tx *gorm.DB, moldIds []int64) map[int64]model.MoldMaintenancePlan {
	var planRel []MaintenancePlanRel
	sql := `
			select 
			    mpr.mold_id,
			    mp.id,
			    mp.code,
			    mp.name,
			    mp.plan_type,
			    mp.mold_type,
			    mp.task_start,
			    mp.task_end,
			    mp.task_standard,
			    mp.plan_cron,
			    mp.timeout_hours,
			    mp.status,
			    mp.is_deleted
			from mold_maintenance_plan_rel mpr
			inner join mold_maintenance_plan mp on mp.id = mpr.mold_maintenance_plan_id and mp.is_deleted = 'N' and mp.status = ?
			where mpr.is_deleted = 'N' and mp.plan_type = ? and mpr.mold_id in ?
		`
	tx.Raw(sql, constant.RUNNING, constant.PLAN_METERING_TYPE, moldIds).Find(&planRel)
	result := make(map[int64]model.MoldMaintenancePlan)

	for _, plan := range planRel {
		result[plan.MoldId] = plan.MoldMaintenancePlan
	}
	return result
}

func getRunningMeteringPlansV1(tx *gorm.DB) map[int64]model.MoldMaintenancePlan {
	var planRel []MaintenancePlanRel
	sql := `
			select 
			    mpr.mold_id,
			    mp.id,
			    mp.code,
			    mp.name,
			    mp.plan_type,
			    mp.mold_type,
			    mp.task_start,
			    mp.task_end,
			    mp.task_standard,
			    mp.plan_cron,
			    mp.timeout_hours,
			    mp.status,
			    mp.is_deleted
			from mold_maintenance_plan_rel mpr
			inner join mold_maintenance_plan mp on mp.id = mpr.mold_maintenance_plan_id and mp.is_deleted = 'N' and mp.status = ?
			inner join mold_info mi on mi.id = mpr.mold_id 
			where mpr.is_deleted = 'N' and mp.plan_type = ? and mi.is_deleted = 'N' and mi.status = 'zhengchang' 
			and mi.type = ?
		`
	tx.Raw(sql, constant.RUNNING, constant.PLAN_METERING_TYPE, "chongkong").Find(&planRel)
	result := make(map[int64]model.MoldMaintenancePlan)

	for _, plan := range planRel {
		result[plan.MoldId] = plan.MoldMaintenancePlan
	}
	return result
}

// GenerateMeteringTaskAndRecordLifecycle 生成模具保养计量任务(目前只有冲孔模具是走计量)
func GenerateMeteringTaskAndRecordLifecycle(tx *gorm.DB, items []PendingCheckProductResumeMoldInfo) (newTask []model.MoldMaintenanceTask, timeoutTask []model.MoldMaintenanceTask) {
	// 累计过的模具是否有达到任务生成标准的
	molds, moldIds := getCheckMolds(tx, items) // 正常的 冲孔的模具 基于未确认的生产履历列表 学将代码
	//molds, _ := getCheckMoldsV1(tx) // 正常的 冲孔的模具 lzg

	runningPlans := getRunningMeteringPlans(tx, moldIds) // 上述模具保养计划为running的计划
	//runningPlans := getRunningMeteringPlansV1(tx) // 上述模具保养计划为running的计划 lzg

	tasks, timeout := GenerateMoldTask(tx, molds, runningPlans)

	newTask = append(newTask, tasks...)
	timeoutTask = append(timeoutTask, timeout...)

	return
}

// CheckProductionResume 检查生产履历，统计超过12小时的履历，确认最终的产量
func CheckProductionResume() {
	var items []PendingCheckProductResumeMoldInfo
	// mr.is_checked = 'N' 表示未处理过的生产履历
	sql := `
			select 
				mr.mold_id as id,
				mr.mold_code as code,
				mr.order_code,
				mr.count as resume_count,
				mr.id as resume_id,
				mi.flush_count,
				mi.calc_flush_count
			from mold_product_resume mr
			inner join mold_info mi on mi.id = mr.mold_id and mi.is_deleted = 'N'
			where mr.is_deleted = 'N' and mr.is_checked = 'N' and mr.gmt_created <= ? and mi.status = 'zhengchang' and mi.type = 'chongkong'
		`
	//mi.type = 'chongkong'
	// 查询检查生成时间已经超过12小时的生产履历，以其产量为该单号下该模具最终的产量
	duration, _ := time.ParseDuration("-12h")
	//duration, _ := time.ParseDuration("-10s") // 测试 用10s即可
	checkTime := base.Now().Time().Add(duration).Format(base.TimeFormart)
	dao.GetConn().Raw(sql, checkTime).Find(&items)

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	// 更新模具冲次
	AddMoldFlushCount(tx, items)

	// 2、盘点冲次是否满足条件，满足则生成计量任务并添加生命周期记录
	newTask, timeoutTask := GenerateMeteringTaskAndRecordLifecycle(tx, items)

	// 新生成的保养任务
	for _, task := range newTask {
		// 给所有人发消息
		message.CreateMessage(tx, message.Message{
			Content:  email.FormatMessage(email.MAINTENANCE_CREATE, task.Code).Content,
			Operator: user.GetAllUserLoginName(),
			JobId:    task.ID,
			Type:     constant.MAINTENANCE,
		})
	}

	// 已经超时的任务
	for _, task := range timeoutTask {
		// 给所有人发消息
		message.CreateMessage(tx, message.Message{
			Content:  email.FormatMessage(email.MAINTENANCE_TIMEOUT, task.Code).Content,
			Operator: user.GetAllUserLoginName(),
			JobId:    task.ID,
			Type:     constant.MAINTENANCE,
		})
	}

	tx.Commit()

	var taskIds []int64
	for _, v := range timeoutTask {
		taskIds = append(taskIds, v.ID)
	}

	// 事务提交后发送超时邮件
	//result := maintenance.QueryPageMaintenanceTask(taskIds)
	result := QueryPageMaintenanceTask(taskIds)
	receivers, err := email.GetMoldTaskTimeoutEmailReceiver(email.MAINTENANCE_TIMEOUT)
	if err != nil {
		logrus.Error(err)
	} else {
		if len(result) > 0 {
			for _, receiver := range receivers {
				if err := email.SendEmail(receiver, email.MAINTENANCE_TIMEOUT, result); err != nil {
					logrus.WithField("address", receiver).Error(err)
				}
			}
		}
	}

	// 更新所有已连接的websocket用户未读消息数量
	for _, operator := range user.GetAllUserLoginName() {
		client, ok := message.ClientMap[operator]
		if ok {
			client.Write(fmt.Sprintf(`{"count": %d}`, message.GetUnreadMessageCount(operator)))
		}
	}

}

func CronInitProductionResumeCheck() {
	// 每30分钟检测一次检查生产履历
	if _, err := cron.CronJob.CronWithSeconds("0 0,30 * ? * *").Do(CheckProductionResume); err != nil {
		logrus.Error(err.Error())
	}
	// 测试 每min执行1次
	//if _, err := cron.CronJob.CronWithSeconds("0 0/1 * ? * *").Do(CheckProductionResume); err != nil {
	//	logrus.Error(err.Error())
	//}
}

// QueryPageMaintenanceTask 通过ids批量查询保养任务页面数据
func QueryPageMaintenanceTask(taskIds []int64) []email.MaintenanceTaskEmailVO {
	var result []email.MaintenanceTaskEmailVO
	sql := `
		select
		    mmt.id,
			mmt.code as task_code,
			mi.code as mold_code,
			mi.name as mold_name,
			mi.type,
			mi.project_name,
			mmt.task_type,
			mmt.standard_name,
			mmt.maintenance_level,
			mmt.gmt_created,
			mmt.gmt_updated,
			mmt.status,
			COALESCE(ui.name, '系统') AS created_by,
			COALESCE(uop.name, mmt.operator) AS operator,
			ifnull(GROUP_CONCAT(DISTINCT mpr.part_code SEPARATOR '/'),'') as part_codes
		from mold_maintenance_task mmt
		left join mold_info mi on mi.id = mmt.mold_id and mi.is_deleted = 'N'
		left join mold_part_rel mpr on mpr.mold_code = mi.code and mpr.is_deleted = 'N'
		left join user_info ui on ui.login_name = mmt.created_by
		left join user_info uop on uop.login_name = mmt.operator
		where mmt.is_deleted = 'N' and mmt.id in (?)
		group by mmt.id,mmt.code,mi.code,mi.type,mmt.standard_name,mmt.maintenance_level,mmt.gmt_created,mmt.status,ui.name, uop.name
		order by mmt.gmt_updated desc, mmt.gmt_created desc
	`
	tx := dao.GetConn()
	taskIdStrings := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(taskIds)), ", "), "[]")
	tx.Raw(sql, taskIdStrings).Scan(&result)

	return result
}
