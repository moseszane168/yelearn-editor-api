/**
 * 模具保养任务
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/common/cron"
	"crf-mold/controller/email"
	"crf-mold/controller/message"
	"crf-mold/controller/productresume"
	"crf-mold/controller/user"
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CronInitTimeoutCheck() {
	// 每30分钟检测一次计时计划超时的保养任务
	if _, err := cron.CronJob.CronWithSeconds("0 0,30 * ? * *").Do(TimeoutTimingMaintenanceTaskCheck); err != nil {
		logrus.Error(err.Error())
	}
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

// TimeoutTimingMaintenanceTaskCheck 检测超时的计时保养任务
func TimeoutTimingMaintenanceTaskCheck() {
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
		where mmt.is_deleted = 'N' and mmt.plan_type in (?,?) and mmt.status = ? and mmt.timeout_time < ?
		group by mmt.id,mmt.code,mi.code,mi.type,mmt.standard_name,mmt.maintenance_level,mmt.gmt_created,mmt.status,ui.name, uop.name
		order by mmt.gmt_updated desc, mmt.gmt_created desc
	`
	tx := dao.GetConn()
	tx.Raw(sql, constant.PLAN_TIMING_TYPE, constant.PLAN_MANUAL_TYPE, constant.WAIT, base.Now().String()).Scan(&result)
	for i := 0; i < len(result); i++ {
		tx.Table("mold_maintenance_task").Where("id = ?", result[i].ID).Updates(&model.MoldMaintenanceTask{
			Status:     constant.TIMEOUT,
			UpdatedBy:  "system",
			GmtUpdated: base.Now().Time(),
		})
	}
	// 已经超时的任务
	for _, task := range result {
		// 给所有人发消息
		message.CreateMessage(tx, message.Message{
			Content:  email.FormatMessage(email.MAINTENANCE_TIMEOUT, task.TaskCode).Content,
			Operator: user.GetAllUserLoginName(),
			JobId:    task.ID,
			Type:     constant.MAINTENANCE,
		})
	}

	// 事务提交后发送超时邮件
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
}

// GenerateTimingMaintenanceTask 新增模具计时保养任务
func GenerateTimingMaintenanceTask(planId int64, moldIds []int64) {
	logrus.Infof("[GenerateTimingMaintenanceTask]开始执行Cron任务，计划ID：%d，模具ID列表：%v", planId, moldIds)
	// 避免定时任务重复执行
	time.Sleep(time.Second * 2)
	var entity model.MoldMaintenancePlan
	if err := dao.GetConn().Table("mold_maintenance_plan").Where("id = ? and status = 'running'", planId).First(&entity).Error; err != nil {
		logrus.Error("对应的计时保养计划不存在或未运行")
	} else {
		tx := dao.GetConn().Begin()
		defer dao.TransactionRollback(tx)

		var newTask []model.MoldMaintenanceTask
		// 新增
		for i := 0; i < len(moldIds); i++ {
			var task model.MoldMaintenanceTask
			task.Code = common.GenerateCode("RW")
			task.MoldID = moldIds[i]
			task.MoldMaintenancePlainID = planId
			task.Status = constant.WAIT
			task.TaskType = "auto"
			task.PlanType = constant.PLAN_TIMING_TYPE
			task.CreatedBy = "system"
			task.UpdatedBy = "system"
			task.TimeoutTime = base.Now().Time().Add(time.Hour * time.Duration(entity.TimeoutHours))
			tx.Table("mold_maintenance_task").Create(&task)

			newTask = append(newTask, task)

			// 添加模具任务生成的生命周期记录
			tx.Table("mold_maintenance_task_lifecycle").Create(&model.MoldMaintenanceTaskLifecycle{
				MoldMaintenanceTaskID: task.ID,
				Title:                 constant.TASKGENERATE, // 任务生成
				UpdatedBy:             "system",
				CreatedBy:             "system",
				Time:                  base.Now().String(),
			})
		}

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
		tx.Commit()
	}
}

// GenerateMeteringTimingTask 生成计量保养任务, moldIds需要检查是否生成的列表 plan单个任务
func GenerateMeteringTimingTask(tx *gorm.DB, moldIds []int64, plan model.MoldMaintenancePlan) {
	// 拼凑出生成保养任务的参数出来 molds planRel
	var molds []model.MoldInfo
	// 定义map时采用make的方式顺便开辟空间
	planRel := make(map[int64]model.MoldMaintenancePlan)
	for _, moldID := range moldIds {
		planRel[moldID] = plan
	}
	tx.Table("mold_info").Where("id in ? and is_deleted = 'N' and status = 'zhengchang' and type = 'chongkong' ", moldIds).Find(&molds)
	// 调用学将写好的生成任务函数
	tasks, _ := productresume.GenerateMoldTask(tx, molds, planRel)
	// 新生成的保养任务
	for _, task := range tasks {
		// 给所有人发消息
		message.CreateMessage(tx, message.Message{
			Content:  email.FormatMessage(email.MAINTENANCE_CREATE, task.Code).Content,
			Operator: user.GetAllUserLoginName(),
			JobId:    task.ID,
			Type:     constant.MAINTENANCE,
		})
	}
}

// @Tags 模具保养任务
// @Summary 新增保养任务
// @Accept json
// @Produce json
// @Param Body body AddMoldMaintenanceTaskVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/task/new [post]
func AddMoldMaintenanceTask(c *gin.Context) {
	var vo AddMoldMaintenanceTaskVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	userId := c.GetHeader(constant.USERID)
	var plan model.MoldMaintenancePlan
	tx.Table("mold_maintenance_plan_rel").Where("is_deleted = 'N' and mold_id = ?", vo.MoldID).First(&plan)
	//if err := tx.Table("mold_maintenance_plan").Where("is_deleted = 'N' and mold_type = ?", vo.MoldType).First(&plan).Error; err != nil {
	//	panic(base.ResponseEnum[base.MAINTENANCE_PLAN_NOT_EXIST])
	//}
	// 生成任务
	var task model.MoldMaintenanceTask
	task.Code = common.GenerateCode("RW")
	task.MoldID = vo.MoldID
	task.MoldMaintenancePlainID = plan.ID

	//如果需审核
	if vo.IsApproval == 1 {
		task.Status = constant.APPROVAL
	} else {
		task.Status = constant.WAIT
	}

	task.TaskType = "manual"
	task.PlanType = constant.PLAN_MANUAL_TYPE
	task.TaskGenInterval = 99999 * 10000
	task.TimeoutCount = 99999 * 10000
	task.TimeoutTime = base.Now().Time().Add(time.Hour * time.Duration(vo.TimeoutHours)) // 设置手动任务的超时时间
	task.UpdatedBy = userId
	task.CreatedBy = userId
	task.ApprovalId = vo.ApprovalId
	task.IsApproval = vo.IsApproval
	tx.Table("mold_maintenance_task").Create(&task)

	// 添加模具任务生成的生命周期记录
	tx.Table("mold_maintenance_task_lifecycle").Create(&model.MoldMaintenanceTaskLifecycle{
		MoldMaintenanceTaskID: task.ID,
		Title:                 constant.TASKGENERATE, // 任务生成
		UpdatedBy:             userId,
		CreatedBy:             userId,
		Time:                  base.Now().String(),
		ApprovalId:            vo.ApprovalId,
		IsApproval:            vo.IsApproval,
	})

	tx.Commit()
	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具保养任务
// @Summary 提交保养任务
// @Accept json
// @Produce json
// @Param Body body SubmitMoldMaintenanceTaskVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/task [post]
func SubmitMoldMaintenanceTask(c *gin.Context) {
	var vo SubmitMoldMaintenanceTaskVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	userId := c.GetHeader(constant.USERID)

	var old_status string
	// 存在并查看旧的状态
	var entity model.MoldMaintenanceTask
	if err := tx.Table("mold_maintenance_task").Where("id = ? and is_deleted = 'N'", vo.ID).First(&entity).Error; err != nil {
		panic(base.ParamsErrorN())
	} else {
		old_status = entity.Status
		if old_status == "complete" {
			vo.Time = entity.Time
		}
	}

	// 保养任务
	tx.Table("mold_maintenance_task").Where("id = ?", vo.ID).Updates(&model.MoldMaintenanceTask{
		UpdatedBy:                 userId,
		Status:                    constant.COMPLETE,
		StandardName:              vo.StandardName,
		MoldMaintenanceStandardID: vo.MoldMaintenanceStandardId,
		Remark:                    vo.Remark,
		MaintenanceLevel:          vo.MaintenanceLevel,
		Time:                      vo.Time,
		Operator:                  vo.Operator,
	})

	// 保养任务内容
	l := len(vo.Content)
	if l > 0 {
		contents := make([]model.MoldMaintenanceTaskContent, l)
		for i := 0; i < l; i++ {
			base.CopyProperties(&contents[i], vo.Content[i])
			// 所有保养内容必须完成
			if vo.Content[i].Status != "Y" {
				panic(base.ResponseEnum[base.MAINTENANCE_TASK_CONTENT_UNCOMPLETE])
			}
			contents[i].MoldMaintenanceTaskID = vo.ID
			contents[i].Status = "Y"
			contents[i].Order = i
			contents[i].CreatedBy = userId
			contents[i].UpdatedBy = userId
		}
		// 先删除原有的
		tx.Table("mold_maintenance_task_content").Where("mold_maintenance_task_id = ?", vo.ID).Updates(&model.MoldMaintenanceTaskContent{
			IsDeleted: "Y",
			UpdatedBy: userId,
		})
		// 再新增新的
		tx.Table("mold_maintenance_task_content").CreateInBatches(contents, l)
	}

	// 更换备件
	l = len(vo.ReplaceSpare)
	if l > 0 {
		spares := make([]model.MoldReplaceSpareRel, l)
		for i := 0; i < l; i++ {
			base.CopyProperties(&spares[i], vo.ReplaceSpare[i])

			// 查询上一次冲次和当前模具冲次
			var lastFlushCount int64
			var flushCount int64
			var lastFlushCountRecord model.MoldReplaceSpareRel
			if err := tx.Table("mold_replace_spare_rel").Where("mold_id = ? and spare_code = ? and is_deleted = 'N'", entity.MoldID, spares[i].SpareCode).Order("gmt_created desc").First(&lastFlushCountRecord).Error; err == nil {
				lastFlushCount = lastFlushCountRecord.FlushCount
			}

			var moldInfo model.MoldInfo
			if err := tx.Table("mold_info").Where("id = ? and is_deleted = 'N'", entity.MoldID).First(&moldInfo).Error; err == nil {
				flushCount = moldInfo.FlushCount
			}

			spares[i].Type = constant.MAINTENANCE
			spares[i].MaintenanceTaskID = vo.ID
			spares[i].CreatedBy = userId
			spares[i].UpdatedBy = userId
			spares[i].MoldID = entity.MoldID
			spares[i].FlushCount = flushCount
			spares[i].LastFlushCount = lastFlushCount
		}
		// 删除原有数据
		tx.Table("mold_replace_spare_rel").Where("maintenance_task_id = ?", vo.ID).Updates(map[string]interface{}{
			"is_deleted": "Y",
		})
		tx.Table("mold_replace_spare_rel").CreateInBatches(spares, l)
	}

	// 生命周期
	if old_status != "complete" {
		var lifecycleRecord model.MoldMaintenanceTaskLifecycle
		lifecycleRecord.MoldMaintenanceTaskID = vo.ID
		lifecycleRecord.Title = constant.TASKCOMPLETE
		lifecycleRecord.CreatedBy = userId
		lifecycleRecord.UpdatedBy = userId
		lifecycleRecord.Time = base.Now().String()
		tx.Table("mold_maintenance_task_lifecycle").Create(&lifecycleRecord)
	}

	// 手动保养任务完成之后，对计算冲次进行处理（只有冲孔的模具需要），调整到最近的标准冲次
	//if entity.TaskType == "manual" && entity.MoldMaintenancePlainID != 0 {
	// 第一次任务完成后 即原status不为complete，对计算冲次进行清理
	//if old_status != "complete" && entity.MoldMaintenancePlainID != 0 {
	//	var plan model.MoldMaintenancePlan
	//	if err := tx.Table("mold_maintenance_plan").Where("is_deleted = 'N' and id = ? ", entity.MoldMaintenancePlainID).First(&plan).Error; err == nil {
	//		var moldInfo model.MoldInfo
	//		if err := tx.Table("mold_info").Where("id = ? and is_deleted = 'N'", entity.MoldID).First(&moldInfo).Error; err == nil {
	//			calcFlushCount := moldInfo.CalcFlushCount
	//			//n := int64(math.Ceil(float64(calcFlushCount) / float64(plan.TaskStandard*10000)))
	//			n := int64(math.Ceil(float64(calcFlushCount) / float64(plan.TaskStandard)))
	//			//interval := n * plan.TaskStandard * 10000
	//			interval := n * plan.TaskStandard
	//			tx.Table("mold_info").Where("id = ? and is_deleted = 'N'", entity.MoldID).Updates(map[string]interface{}{
	//				//"calc_flush_count": interval,
	//				// 修改为直接清零
	//				"calc_flush_count": 0,
	//			})
	//
	//			// 更新保养任务的TaskGenInterval
	//			tx.Table("mold_maintenance_task").Where("id = ?", vo.ID).Updates(&model.MoldMaintenanceTask{
	//				TaskGenInterval: interval,
	//			})
	//
	//		}
	//	}
	//}
	if old_status != "complete" {
		tx.Table("mold_info").Where("id = ? and is_deleted = 'N'", entity.MoldID).Updates(map[string]interface{}{
			// 修改为直接清零
			"calc_flush_count": 0,
		})
	}
	tx.Commit()
	c.JSON(http.StatusOK, base.Success(true))
}

//
//
//	for ; n > 0; n-- {
//		interval := n * plan.TaskStandard * 10000 // 第n个区间的标准
//		min := interval - plan.TaskStart*10000    // 第n个区间的最小
//if entity.TaskType == "manual" {
//	tx.Table("mold_info").Where("id = ? and is_deleted = 'N'", entity.MoldID).Updates(map[string]interface{}{
//		"calc_flush_count": 0,
//	})
//}
//

// @Tags 模具保养任务
// @Summary 保养任务分页
// @Accept json
// @Produce json
// @Param query query PageMaintenanceTaskInputVO true "PageMaintenanceTaskInputVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageMaintenanceTaskOutVO
// @Router /maintenance/task/page [get]
func PageMoldMaintenanceTask(c *gin.Context) {
	var vo PageMaintenanceTaskInputVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	condition := ""
	if vo.Status != "" {
		combineString := "'" + strings.Join(strings.Split(vo.Status, ","), "','") + "'"
		condition += "and mmt.status in (" + combineString + ")"
	}
	if vo.CodeOrName == "" {
		// 额外处理时间字段
		if vo.GmtCreatedBegin != "" {
			condition += "and mmt.gmt_created >= '" + vo.GmtCreatedBegin + "' "
		}
		if vo.GmtCreatedEnd != "" {
			condition += "and mmt.gmt_created <= '" + vo.GmtCreatedEnd + "' "
		}
		if vo.TaskCode != "" {
			condition += "and mmt.code = '" + vo.TaskCode + "' "
		}
		if vo.MoldCode != "" {
			condition += "and mi.code = '" + vo.MoldCode + "' "
		}
		if vo.Type != "" {
			condition += "and mi.type = '" + vo.Type + "' "
		}
		if vo.ProjectName != "" {
			condition += "and mi.project_name = '" + vo.ProjectName + "' "
		}
		if vo.TaskType != "" {
			condition += "and mmt.task_type = '" + vo.TaskType + "' "
		}
		if vo.Operator != "" {
			condition += fmt.Sprintf(`and (uop.name like concat('%%', '%s' '%%') or mmt.operator like concat('%%', '%s' '%%'))`, vo.Operator, vo.Operator)
		}
		if vo.StandardName != "" {
			condition += "and mmt.standard_name = '" + vo.StandardName + "' "
		}
		if vo.MaintenanceLevel != "" {
			condition += "and mmt.maintenance_level = '" + vo.MaintenanceLevel + "' "
		}
		if vo.PartCodes != "" {
			condition += "and mpr.part_code = '" + vo.PartCodes + "' "
		}
	}
	if vo.RfidTid != "" {
		condition += "and mi.rfid = '" + vo.RfidTid + "' "
	}

	//只看本人审核
	userId := c.GetHeader(constant.USERID)
	if vo.Status == "approval" {
		condition += "and mmt.approval_id = '" + userId + "' "
	}

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
		where mmt.is_deleted = 'N' and (mi.code like concat('%',?,'%') or mpr.part_code like concat('%',?,'%'))
		`

	sql = sql + condition +
		`group by mmt.id,mmt.code,mi.code,mi.type,mmt.standard_name,mmt.maintenance_level,mmt.gmt_created,mmt.status,ui.name, uop.name
		 order by mmt.gmt_updated desc, mmt.gmt_created desc`

	var result []PageMaintenanceTaskOutVO
	page := base.PageWithRawSQL(dao.GetConn(), &result, vo.GetCurrentPage(), vo.GetSize(), sql, vo.CodeOrName, vo.CodeOrName)
	if len(result) == 0 {
		page.List = []interface{}{}
	} else {
		// 模具保养结果时间与详情时间轴不一致，直接按生命周期表结束时间为准
		var result2 []PageMaintenanceTaskOutVO
		for _, pmt := range result {
			var updateTime base.Time
			pmtId := pmt.ID
			sql2 := `SELECT MAX(t1.gmt_updated) AS updateTime FROM mold_maintenance_task_lifecycle t1 WHERE t1.mold_maintenance_task_id=?`
			dao.GetConn().Raw(sql2, pmtId).Scan(&updateTime)
			pmt.GmtUpdated = updateTime
			result2 = append(result2, pmt)
		}
		page.List = result2
	}
	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 模具保养任务
// @Summary 保养任务查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} MaintenanceTaskOneOutVO
// @Router /maintenance/task/one [get]
func OneMoldMaintenanceTask(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var task model.MoldMaintenanceTask
	if err := dao.GetConn().Table("mold_maintenance_task").Where("id = ? and is_deleted = 'N'", id).First(&task).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	var outVo MaintenanceTaskOneOutVO
	base.CopyProperties(&outVo, task)

	if task.MoldMaintenanceStandardID == 0 {
		outVo.MoldMaintenanceStandardId = nil
	} else {
		outVo.MoldMaintenanceStandardId = &task.MoldMaintenanceStandardID
	}

	// 设置用户名称
	if outVo.Operator != "" {
		var userInfo model.UserInfo
		if err := dao.GetConn().Table("user_info").Where("is_deleted = 'N' and login_name = ?", outVo.Operator).First(&userInfo).Error; err == nil {
			outVo.Operator = userInfo.Name
		}
	}

	// 模具编号和类型
	var m model.MoldInfo
	if err := dao.GetConn().Table("mold_info").Where("id = ? and is_deleted = 'N'", task.MoldID).First(&m).Error; err == nil {
		outVo.MoldCode = m.Code
		outVo.MoldName = m.Name
		outVo.Type = m.Type
	}

	// 保养任务内容
	var contents []model.MoldMaintenanceTaskContent
	if err := dao.GetConn().Table("mold_maintenance_task_content").Where("mold_maintenance_task_id = ? and is_deleted = 'N'", id).Find(&contents).Error; err == nil {
		list := make([]MaintenanceTaskContentVO, len(contents))
		for i := 0; i < len(contents); i++ {
			item := MaintenanceTaskContentVO{}
			base.CopyProperties(&item, contents[i])
			list[i] = item
		}
		outVo.Content = list
	}

	// 更换备件
	var replaceSpares []model.MoldReplaceSpareRel
	if err := dao.GetConn().Table("mold_replace_spare_rel").Where("maintenance_task_id = ? and is_deleted = 'N'", id).Find(&replaceSpares).Error; err == nil {
		list := make([]ReplaceSpareOneVO, len(replaceSpares))
		for i := 0; i < len(replaceSpares); i++ {
			item := ReplaceSpareOneVO{}
			base.CopyProperties(&item, replaceSpares[i])
			// 备件名称和规格型号
			var spareInfo model.SpareInfo
			if err := dao.GetConn().Table("spare_info").Where("code = ?", item.SpareCode).First(&spareInfo).Error; err == nil {
				item.SpareName = spareInfo.Name
				item.Flavor = spareInfo.Flavor
			}

			list[i] = item
		}
		outVo.ReplaceSpare = list
	}

	c.JSON(http.StatusOK, base.Success(outVo))
}

// @Tags 模具保养任务
// @Summary 保养任务挂起/继续
// @Accept json
// @Produce json
// @Param Body body StatusVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/task/status [POST]
func ToggleMoldMaintenanceTask(c *gin.Context) {
	var vo StatusVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	var entity model.MoldMaintenanceTask
	if err := tx.Table("mold_maintenance_task").Where("id = ?", vo.ID).First(&entity).Error; err == nil {
		status := entity.Status

		// 生命周期
		var lifecycle model.MoldMaintenanceTaskLifecycle
		lifecycle.Operator = c.GetHeader(constant.USERID)
		lifecycle.Time = base.Now().String()
		lifecycle.MoldMaintenanceTaskID = vo.ID
		lifecycle.CreatedBy = "system"
		lifecycle.UpdatedBy = "system"

		if status == constant.WAIT { // 挂起
			if entity.HangUpCount == 3 { // 挂起次数不能超过三次
				panic(base.ResponseEnum[base.MAINTENANCE_TASK_HUANG_UP_CAN_NOT_THAN_THREE])
			}

			tx.Table("mold_maintenance_task").Where("id = ?", vo.ID).Updates(&model.MoldMaintenanceTask{
				Status:      constant.PAUSE,
				Reason:      vo.Reason,
				HangUpCount: entity.HangUpCount + 1,
				UpdatedBy:   c.GetHeader(constant.USERID),
			})
			lifecycle.Title = constant.TASKPAUSE

			lifecycle.Content = vo.Reason
		} else if status == constant.PAUSE { // 继续保养
			tx.Table("mold_maintenance_task").Where("id = ?", vo.ID).Updates(&model.MoldMaintenanceTask{
				Status:    constant.WAIT,
				UpdatedBy: c.GetHeader(constant.USERID),
			})
			lifecycle.Title = constant.TASKCONTINUE
		}

		tx.Table("mold_maintenance_task_lifecycle").Create(&lifecycle)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具保养任务
// @Summary 获取保养记录
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} model.MoldMaintenanceTaskLifecycle
// @Router /maintenance/task/lifecycle [get]
func LifeCycle(c *gin.Context) {
	var vo common.IDVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	id := vo.ID
	if id == 0 {
		panic(base.ParamsErrorN())
	}
	sql := `
		SELECT 
			ml.id,
			ml.mold_maintenance_task_id,
			ml.title,
			COALESCE(uop.name, ml.updated_by) AS operator,
			ml.time,
			ml.content
		FROM mold_maintenance_task_lifecycle ml
		LEFT JOIN user_info uop on uop.login_name = ml.updated_by
		WHERE ml.mold_maintenance_task_id = ?
	`
	var result []model.MoldMaintenanceTaskLifecycle
	dao.GetConn().Raw(sql, vo.ID).Find(&result)

	if len(result) > 0 {
		outVos := base.CopyPropertiesList(reflect.TypeOf(MoldMaintenanceTaskLifecycleVO{}), result)
		c.JSON(http.StatusOK, base.Success(outVos))
	} else {
		c.JSON(http.StatusOK, base.Success([]interface{}{}))
	}
}

// @Tags 模具保养任务
// @Summary 保养任务转派
// @Accept json
// @Produce json
// @Param Body body MaintenanceTaskChargeVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/task/charge [post]
func ChangeChargeMaintenanceTask(c *gin.Context) {
	var vo MaintenanceTaskChargeVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 存在
	var one model.MoldMaintenanceTask
	if err := dao.GetConn().Table("mold_maintenance_task").Where("id = ? and is_deleted = 'N'", vo.Id).First(&one).Error; err != nil {
		panic(base.ParamsErrorN())
	}
	if one.Operator != "" {
		userId := c.GetHeader(constant.USERID)
		if userId != "admin" {
			panic(base.ResponseEnum[base.MAINTENANCE_CHARGE_FAILED])
		}
	}

	one.Operator = vo.ChargeName
	one.GmtUpdated = base.Now().Time()
	dao.GetConn().Table("mold_maintenance_task").Updates(&one)

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具保养任务
// @Summary 保养任务审核
// @Accept json
// @Produce json
// @Param Body body MaintenanceTaskApprovalVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/task/approval [post]
func ApprovalMaintenanceTask(c *gin.Context) {
	var vo MaintenanceTaskApprovalVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	//var plan model.MoldMaintenancePlan
	//tx.Table("mold_maintenance_plan_rel").Where("is_deleted = 'N' and mold_id = ?", vo.MoldID).First(&plan)

	// 存在
	var one model.MoldMaintenanceTask
	if err := dao.GetConn().Table("mold_maintenance_task").Where("id = ? and is_deleted = 'N'", vo.Id).First(&one).Error; err != nil {
		panic(base.ParamsErrorN())
	}
	if one.Operator != "" {
		userId := c.GetHeader(constant.USERID)
		if userId != "admin" {
			panic(base.ResponseEnum[base.MAINTENANCE_CHARGE_FAILED])
		}
	}

	one.ApprovalStatus = vo.ApprovalStatus

	//审核：通过-状态改为已保养
	if vo.ApprovalStatus == 1 {
		one.Status = "complete"
	}

	//审核：驳回-状态改为待保养
	if vo.ApprovalStatus == 2 {
		one.Status = "wait"
	}

	one.ApprovalComment = vo.ApprovalComment
	one.GmtUpdated = base.Now().Time()
	dao.GetConn().Table("mold_maintenance_task").Updates(&one)

	// 添加模具任务生成的生命周期记录
	tx.Table("mold_maintenance_task_lifecycle").Create(&model.MoldMaintenanceTaskLifecycle{
		MoldMaintenanceTaskID: one.ID,
		Title:                 constant.TASKGENERATE, // 任务生成
		UpdatedBy:             one.UpdatedBy,
		CreatedBy:             one.CreatedBy,
		Time:                  base.Now().String(),
		ApprovalId:            one.ApprovalId,
		IsApproval:            one.IsApproval,
		ApprovalStatus:        one.ApprovalStatus,
		ApprovalComment:       one.ApprovalComment,
	})

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具保养任务
// @Summary 删除保养任务
// @Accept json
// @Produce json
// @Param Body body MaintenanceTaskDeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/task/delete [post]
func DeleteMaintenanceTask(c *gin.Context) {
	var vo MaintenanceTaskDeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}
	userId := c.GetHeader(constant.USERID)

	dao.GetConn().Table("mold_maintenance_task").Where("id = ? and status = ?", vo.Id, vo.Status).Updates(&model.MoldMaintenanceTask{
		IsDeleted: "Y",
		UpdatedBy: userId,
	})

	c.JSON(http.StatusOK, base.Success(true))
}
