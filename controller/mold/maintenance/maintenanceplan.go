/**
 * 模具保养计划
 */

package mold

import (
	"bytes"
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/common/cron"
	"crf-mold/dao"
	"crf-mold/model"
	"encoding/json"
	"errors"
	"fmt"
	cronExpress "github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//func TestSql() {
//	var taskIds []int64
//	dao.GetConn().Raw(`
//				SELECT
//					id
//				FROM mold_inout_bound_task_flow
//				WHERE job_id IN (
//					SELECT
//						job_id
//					FROM mold_inout_bound_task_flow
//					WHERE
//						station_unique_code = ? AND flow_order = ? AND flag = 0
//				) AND flow_order = ? AND flag = 1
//			`, "101", constant.IN_TASKFLOW_ASRS_PLC_WRITE_CMD, constant.IN_TASKFLOW_ASRS_READ_RFID).Scan(&taskIds)
//	fmt.Printf("task %v", taskIds)
//}

func CronInitAddTask() {
	// 添加有效的所有的CronJob信息
	result := GetAllCronJob()
	for _, vo := range result {
		if job, err := cron.CronJob.CronWithSeconds(vo.PlanCron).Do(GenerateTimingMaintenanceTask, vo.ID, vo.Molds); err == nil {
			job.Tag(strconv.FormatInt(vo.ID, 10))
		} else {
			logrus.Errorf("Plan ID：%d Error: %s", vo.ID, err.Error())
		}
	}
}

// 删除保养计划关联的模具
func deleteMaintenancePlanRel(tx *gorm.DB, planId int64, userId string) {
	tx.Table("mold_maintenance_plan_rel").Where("mold_maintenance_plan_id = ?", planId).Updates(&model.MoldMaintenancePlanRel{
		IsDeleted: "Y",
		UpdatedBy: userId,
	})
}

// SaveMaintenancePlanRel 保存计时保养计划关联的模具
func SaveMaintenancePlanRel(tx *gorm.DB, planId int64, userId string, moldIds []int64) {
	// 目前做法，先全部删除，再新增
	deleteMaintenancePlanRel(tx, planId, userId)

	// 新增
	saveBatch := make([]model.MoldMaintenancePlanRel, len(moldIds))
	for i := 0; i < len(moldIds); i++ {
		var item model.MoldMaintenancePlanRel
		item.MoldId = moldIds[i]
		item.MoldMaintenancePlanID = planId
		item.CreatedBy = userId
		item.UpdatedBy = userId
		saveBatch[i] = item
	}
	tx.Table("mold_maintenance_plan_rel").CreateInBatches(saveBatch, len(saveBatch))
}

// GetAllCronJob 获取所有的CronJob信息
func GetAllCronJob() []MaintenanceTimingPlanOutVO {
	var result []MaintenanceTimingPlanOutVO
	dao.GetConn().Table("mold_maintenance_plan").Select("id, name, plan_cron, timeout_hours").Where(
		"is_deleted = 'N' and plan_type = 'timing' and status = 'running'").Scan(&result)
	for i := 0; i < len(result); i++ {
		var molds []int64
		dao.GetConn().Table("mold_maintenance_plan_rel").Select("mold_id").Where(
			"is_deleted = 'N' and mold_maintenance_plan_id = ?", result[i].ID).Scan(&molds)
		if len(molds) > 0 {
			result[i].Molds = molds
		} else {
			result[i].Molds = nil
		}
	}
	return result
}

// removeMaintenancePlanCronTask 删除Cron任务
func removeMaintenancePlanCronTask(planId int64) {
	if err := cron.CronJob.RemoveByTag(strconv.FormatInt(planId, 10)); err != nil {
		logrus.Error(err.Error())
	}
}

// AddMaintenancePlanCronTask 新增计时保养Cron任务
func AddMaintenancePlanCronTask(planId int64, planCron string, moldIds []int64) error {
	// 从Cron里面先删除定时任务，再添加定时任务
	removeMaintenancePlanCronTask(planId)
	// 添加定时任务
	if job, err := cron.CronJob.CronWithSeconds(planCron).Do(GenerateTimingMaintenanceTask, planId, moldIds); err == nil {
		job.Tag(strconv.FormatInt(planId, 10))
		logrus.Infof("新增Cron任务成功，计划ID：%d，模具ID列表：%v", planId, moldIds)
		return nil
	} else {
		logrus.Errorf("Plan ID：%d Error: %s", planId, err.Error())
		return err
	}
}

// CheckMoldExistedPlan 检查模具是否已经设置其它保养计划
func CheckMoldExistedPlan(tx *gorm.DB, planId int64, checkMolds []int64) {
	var planRel model.MoldMaintenancePlanRel
	if err := tx.Table("mold_maintenance_plan_rel").Where(
		"is_deleted = 'N' and mold_maintenance_plan_id != ? and mold_id in ?", planId, checkMolds).First(&planRel).Error; err == nil {
		// 存在模具已经设置其它保养计划
		var moldInfo model.MoldInfo
		if err := dao.GetConn().Table("mold_info").Where("id = ? and is_deleted = 'N'", planRel.MoldId).First(&moldInfo).Error; err == nil {
			panic(base.ParamsError(fmt.Sprintf("模具 %s 已经设置其它保养计划", moldInfo.Code)))
		}
	}
}

// MaintenancePlanBodyReValidate 创建维保计划时数据校验不通过后 如必填字段 整数字段强制转换等
func MaintenancePlanBodyPreValidate(vo *CreateMoldMaintenancePlanVO, body []byte) error {
	// 对空字段进行处理
	if vo.Name == "" || vo.PlanType == "" {
		//c.JSON(http.StatusBadRequest, base.ParamsError("必填字段为空, 保养名称/类型"))
		return errors.New("必填字段为空, 保养名称/类型")
	}
	// 对原有的请求体进行校验 Unmarshal 转换为map
	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		return err
	}
	// 针对整数进行校验
	validateKeys := map[string]string{"taskStart": "TaskStart", "taskEnd": "TaskEnd", "taskStandard": "TaskStandard"}
	for k, kv := range validateKeys {
		// 该键是否存在
		if v, ok := requestBody[k]; ok {
			// 该值的数据类型是否为int
			if v, ok := v.(string); ok {
				// 为string进行强制转换 转换不了则直接返回异常 空串则直接转为0
				if v == "" {
					v = "0"
				}
				if v, err := strconv.ParseInt(v, 10, 64); err == nil {
					// 通过反射修改值
					reflectValue := reflect.ValueOf(vo).Elem()
					fieldValue := reflectValue.FieldByName(kv)
					if fieldValue.IsValid() && fieldValue.CanSet() {
						fieldValue.Set(reflect.ValueOf(v))
					} else {
						return errors.New("数据类型异常")
					}
				} else {
					return errors.New("数据类型异常")
				}
			}
		}
	}
	return nil
}

// @Tags 模具保养计划
// @Summary 查询保养计划Cron运行时间
// @Accept json
// @Produce json
// @Param Body body CronExpressionVo true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/query/cron [post]
func QueryCronRunTime(c *gin.Context) {
	var vo CronExpressionVo
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	parser := cronExpress.NewParser(cronExpress.Second | cronExpress.Minute | cronExpress.Hour | cronExpress.Dom | cronExpress.Month | cronExpress.Dow | cronExpress.Descriptor)
	schedule, err := parser.Parse(vo.CronExpression)
	if err != nil {
		panic(base.ParamsError("表达式不正确"))
	}

	now := time.Now()
	var result CronRuntimeOutVO
	var res []string
	for i := 0; i < 5; i++ {
		next := schedule.Next(now)
		res = append(res, next.Format("2006-01-02 15:04:05"))
		now = next
	}
	result.RunTimeList = res

	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 模具保养计划
// @Summary 新增保养计划
// @Accept json
// @Produce json
// @Param Body body CreateMoldMaintenancePlanVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/plan [post]
func CreateMoldMaintenancePlan(c *gin.Context) {
	var vo CreateMoldMaintenancePlanVO
	// 传参校验
	// 在检验前先保存请求体的内容 以便后续再次校验
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	// 重新设置请求体的内容 否则ShouldBindBodyWith无法读取到
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// 执行校验
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		if err := MaintenancePlanBodyPreValidate(&vo, bodyBytes); err != nil {
			panic(base.ParamsError(err.Error()))
		}
	}
	// 不允许在线边服务器上操作保养计划,避免Cron计划重复
	if viper.GetString("env") == "prod-gk" {
		panic(base.ResponseEnum[base.MAINTENANCE_PLAN_OPERATE_NOT_ALLOWED])
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	// 名称唯一
	var count int64
	tx.Table("mold_maintenance_plan").Where("name = ? and is_deleted = 'N'", vo.Name).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MAINTENANCE_PLAN_NAME_EXIST])
	}

	// 检查更新的关联模具是否已经设置其它保养计划
	CheckMoldExistedPlan(tx, 0, vo.MoldIds)

	var moldMaintenancePlan model.MoldMaintenancePlan
	base.CopyProperties(&moldMaintenancePlan, vo)
	moldMaintenancePlan.Code = common.GenerateCode(common.JH)
	moldMaintenancePlan.Status = constant.RUNNING
	moldMaintenancePlan.PlanType = vo.PlanType
	moldMaintenancePlan.TaskStandard = vo.TaskStandard
	moldMaintenancePlan.TaskStart = vo.TaskStart
	moldMaintenancePlan.TaskEnd = vo.TaskEnd
	moldMaintenancePlan.PlanCron = vo.PlanCron
	moldMaintenancePlan.TimeoutHours = vo.TimeoutHours

	userId := c.GetHeader(constant.USERID)
	moldMaintenancePlan.CreatedBy = userId
	moldMaintenancePlan.UpdatedBy = userId

	// 新增
	tx.Table("mold_maintenance_plan").Create(&moldMaintenancePlan)

	// 关联模具，只有计时类型保养计划才有
	if len(vo.MoldIds) > 0 {
		SaveMaintenancePlanRel(tx, moldMaintenancePlan.ID, userId, vo.MoldIds)
		// 计时的需要新增定时任务
		if vo.PlanType == constant.PLAN_TIMING_TYPE {
			if err := AddMaintenancePlanCronTask(moldMaintenancePlan.ID, vo.PlanCron, vo.MoldIds); err != nil {
				panic(base.ResponseEnum[base.MAINTENANCE_PLAN_CRON_INVALID])
			}
		} else if vo.PlanType == constant.PLAN_METERING_TYPE {
			// 计量的立刻进行生成保养任务 (计量条件若满足时)
			GenerateMeteringTimingTask(tx, vo.MoldIds, moldMaintenancePlan)
		}

	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(moldMaintenancePlan.ID))
}

// @Tags 模具保养计划
// @Summary 删除保养计划
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/plan [delete]
func DeleteMoldMaintenancePlan(c *gin.Context) {
	idstr := c.Query("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		panic(base.ParamsErrorN())
	}
	// 不允许在线边服务器上操作保养计划,避免Cron计划重复
	if viper.GetString("env") == "prod-gk" {
		panic(base.ResponseEnum[base.MAINTENANCE_PLAN_OPERATE_NOT_ALLOWED])
	}
	userId := c.GetHeader(constant.USERID)

	dao.GetConn().Table("mold_maintenance_plan").Where("id = ?", id).Update("is_deleted", "Y")
	// 删除计时保养计划关联的模具
	deleteMaintenancePlanRel(dao.GetConn(), int64(id), userId)

	var count int64
	dao.GetConn().Table("mold_maintenance_plan").Where("id = ? and is_deleted = 'N' and plan_type = ?", id, constant.PLAN_TIMING_TYPE).Count(&count)
	if count > 0 {
		// 计时保养计划需要删除对应Cron任务
		removeMaintenancePlanCronTask(int64(id))
	}

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具保养计划
// @Summary 更新保养计划
// @Accept json
// @Produce json
// @Param Body body UpdateMoldMaintenancePlanVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/plan [put]
func UpdateMoldMaintenancePlan(c *gin.Context) {
	var vo UpdateMoldMaintenancePlanVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 不允许在线边服务器上操作保养计划,避免Cron计划重复
	if viper.GetString("env") == "prod-gk" {
		panic(base.ResponseEnum[base.MAINTENANCE_PLAN_OPERATE_NOT_ALLOWED])
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	userId := c.GetHeader(constant.USERID)

	// 名称唯一
	var c1 int64
	tx.Table("mold_maintenance_plan").Where("name = ? and id != ? and is_deleted = 'N'", vo.Name, vo.ID).Count(&c1)
	if c1 > 0 {
		panic(base.ResponseEnum[base.MAINTENANCE_NAME_EXIST])
	}

	// 检查更新的关联模具是否已经设置其它保养计划
	CheckMoldExistedPlan(tx, vo.ID, vo.MoldIds)

	var moldMaintenancePlan model.MoldMaintenancePlan
	base.CopyProperties(&moldMaintenancePlan, vo)
	moldMaintenancePlan.UpdatedBy = userId

	tx.Table("mold_maintenance_plan").Updates(&moldMaintenancePlan)

	if len(vo.MoldIds) > 0 {
		SaveMaintenancePlanRel(tx, moldMaintenancePlan.ID, userId, vo.MoldIds)
		if vo.PlanType == constant.PLAN_TIMING_TYPE {
			var count int64
			tx.Table("mold_maintenance_plan").Where("status = 'running' and is_deleted = 'N' and id = ?", vo.ID).Count(&count)
			// 运行中的计时保养计划才会添加Cron任务
			if count > 0 {
				if err := AddMaintenancePlanCronTask(moldMaintenancePlan.ID, vo.PlanCron, vo.MoldIds); err != nil {
					panic(base.ResponseEnum[base.MAINTENANCE_PLAN_CRON_INVALID])
				}
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具保养计划
// @Summary 保养计划列表
// @Accept json
// @Produce json
// @Param name query string false "保养计划名称"
// @Param AuthToken header string false "Token"
// @Success 200 {object} model.MoldMaintenancePlan
// @Router /maintenance/plan/list [get]
func ListMoldMaintenancePlan(c *gin.Context) {
	name := c.Query("name")

	var results []model.MoldMaintenancePlan
	tx := dao.GetConn().Table("mold_maintenance_plan").Where("is_deleted = 'N'")
	if name != "" {
		tx.Where("name like concat('%',?,'%')", name)
	}

	tx.Order("`gmt_created` desc").Find(&results)

	if len(results) > 0 {
		c.JSON(http.StatusOK, base.Success(results))
	} else {
		c.JSON(http.StatusOK, base.Success([]model.MoldMaintenancePlan{}))
	}
}

// @Tags 模具保养计划
// @Summary 保养计划分页
// @Accept json
// @Produce json
// @Param query query PageMaintenancePlanInputVO true "PageMaintenancePlanInputVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} PageMaintenancePlanOutputVO
// @Router /maintenance/plan/page [get]
func PageMoldMaintenancePlan(c *gin.Context) {
	var vo PageMaintenancePlanInputVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	current, size := base.GetPageParams(c)
	var results []model.MoldMaintenancePlan
	tx := dao.GetConn().Table("mold_maintenance_plan").Where("is_deleted = 'N'").Order("`gmt_created` desc")
	if vo.CodeOrName != "" {
		tx.Where("(code like concat('%',?,'%') or name like concat('%',?,'%'))", vo.CodeOrName, vo.CodeOrName)
	} else {
		m := model.MoldMaintenancePlan{}
		base.CopyProperties(&m, vo)
		dao.BuildWhereCondition(tx, m)
		// 额外处理时间字段
		if vo.GmtCreatedBegin != "" {
			tx.Where("gmt_created >= ?", vo.GmtCreatedBegin)
		}
		if vo.GmtCreatedEnd != "" {
			tx.Where("gmt_created <= ?", vo.GmtCreatedEnd)
		}
	}

	page := base.Page(tx, &results, current, size)

	if len(results) == 0 {
		page.List = []interface{}{}
	} else {
		outVos := make([]PageMaintenancePlanOutputVO, len(results))
		list := page.List.(*[]model.MoldMaintenancePlan)
		for i := 0; i < len(*list); i++ {
			base.CopyProperties(&(outVos[i]), (*list)[i])

			// 设置一下部门
			createdBy := outVos[i].CreatedBy
			if createdBy != "" {
				var ui model.UserInfo
				if err := dao.GetConn().Table("user_info").Where("login_name = ?", outVos[i].CreatedBy).First(&ui).Error; err == nil {
					outVos[i].Department = ui.Department
					outVos[i].CreatedBy = ui.Name
				}
			}
		}

		page.List = outVos
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 模具保养计划
// @Summary 保养计划查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} model.MoldMaintenancePlan
// @Router /maintenance/plan/one [get]
func OneMoldMaintenancePlan(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	sql := `
		select
			mi.id,
			mi.name,
			mi.code,
			mi.project_name,
			mi.type,
			mi.line_level,
			mi.process,
			GROUP_CONCAT(DISTINCT mpr.part_code ORDER BY mpr.part_code SEPARATOR '/') as part_codes
		from mold_info mi
		left join mold_part_rel mpr on mpr.mold_code = mi.code and mpr.is_deleted = 'N'
		left join mold_maintenance_plan_rel mmpr on mmpr.mold_id = mi.id and mmpr.is_deleted = 'N'
		left join mold_maintenance_plan mmp on mmp.id = mmpr.mold_maintenance_plan_id and mmp.is_deleted = 'N'
		where mi.is_deleted = 'N' and mmp.id = ?
		group by mi.id,mi.code,mi.line_level,mi.process
	`
	var out MaintenancePlanOutVO
	var moldRel []MaintenancePlanMoldRelOutVO

	var plan model.MoldMaintenancePlan
	if err := dao.GetConn().Table("mold_maintenance_plan").Where("id = ? and is_deleted = 'N'", id).First(&plan).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	base.CopyProperties(&out, plan)

	// 关联模具
	dao.GetConn().Raw(sql, id).Scan(&moldRel)
	if len(moldRel) > 0 {
		out.Molds = moldRel
	} else {
		out.Molds = []MaintenancePlanMoldRelOutVO{}
	}
	c.JSON(http.StatusOK, base.Success(out))
}

// @Tags 模具保养计划
// @Summary 保养计划开启/暂停
// @Accept json
// @Produce json
// @Param Body body common.IDVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/plan/status [POST]
func ToggleMoldMaintenancePlan(c *gin.Context) {
	var vo common.IDVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 不允许在线边服务器上操作保养计划,避免Cron计划重复
	if viper.GetString("env") == "prod-gk" {
		panic(base.ResponseEnum[base.MAINTENANCE_PLAN_OPERATE_NOT_ALLOWED])
	}

	var entity model.MoldMaintenancePlan
	if err := dao.GetConn().Table("mold_maintenance_plan").Where("id = ?", vo.ID).First(&entity).Error; err == nil {
		status := entity.Status
		if status == constant.RUNNING {
			status = constant.PAUSE
			if entity.PlanType == constant.PLAN_TIMING_TYPE {
				// 计时保养计划需要删除对应Cron任务
				removeMaintenancePlanCronTask(vo.ID)
			}
		} else {
			status = constant.RUNNING
			if entity.PlanType == constant.PLAN_TIMING_TYPE {
				// 计时保养计划需要重新添加定时任务
				var molds []int64
				dao.GetConn().Table("mold_maintenance_plan_rel").Select("mold_id").Where(
					"is_deleted = 'N' and mold_maintenance_plan_id = ?", vo.ID).Scan(&molds)
				if len(molds) > 0 {
					if err := AddMaintenancePlanCronTask(entity.ID, entity.PlanCron, molds); err != nil {
						panic(base.ResponseEnum[base.MAINTENANCE_PLAN_CRON_INVALID])
					}
				}
			}
		}

		dao.GetConn().Table("mold_maintenance_plan").Where("id = ?", vo.ID).Updates(&model.MoldMaintenancePlan{
			Status:    status,
			UpdatedBy: c.GetHeader(constant.USERID),
		})
	}

	c.JSON(http.StatusOK, base.Success(true))
}
