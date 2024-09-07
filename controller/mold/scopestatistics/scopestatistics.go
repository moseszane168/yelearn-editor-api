/**
 * 首页统计
 */

package scopestatistics

import (
	"crf-mold/base"
	"crf-mold/dao"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 首页
// @Summary 模具状态统计
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {object} MoldStatusGroupOutVO
// @Router /mold/statistics/mold/status [get]
func MoldStatusGroup(c *gin.Context) {
	var result MoldStatusGroupOutVO
	var statusList []MoldStatus
	dao.GetConn().Raw(`
		select
			mi.status,
			count(1) as count
		from mold_info mi
		where mi.is_deleted = 'N'
		group by mi.status
	`).Scan(&statusList)
	mapData := make(map[string]int64)
	var allCount int64
	for _, v := range statusList {
		mapData[v.Status] = v.Count
		allCount += v.Count
	}
	result.All = allCount
	marshal, _ := json.Marshal(mapData)
	json.Unmarshal(marshal, &result)

	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 首页
// @Summary 保养任务统计
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {object} MoldMaintenanceTaskOutVO
// @Router /mold/statistics/maintenance/task [get]
func MaintenanceStatusGroup(c *gin.Context) {
	var result MoldMaintenanceTaskOutVO
	var taskList []MoldMaintenanceTask
	now := time.Now()
	beginToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dao.GetConn().Raw(`
		select
			mmt.status,
			count(1) as count
		from mold_maintenance_task mmt
		INNER JOIN  mold_info mi on mi.id = mmt.mold_id and mi.is_deleted = 'N'
		where mmt.is_deleted = 'N' and (mmt.status in ('timeout', 'wait', 'pause') or 
			(mmt.status='complete' and mmt.gmt_updated >= ?))
		group by mmt.status
	`, beginToday).Scan(&taskList)
	mapData := make(map[string]int64)
	var allCount int64
	for _, v := range taskList {
		mapData[v.Status] = v.Count
		allCount += v.Count
	}
	marshal, _ := json.Marshal(mapData)
	json.Unmarshal(marshal, &result)

	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 首页
// @Summary 模具异常汇总
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {object} MoldExceptionOutVO
// @Router /mold/statistics/repair [get]
func MoldException(c *gin.Context) {
	var result []MoldExceptionOutVO
	dao.GetConn().Raw(`
		select
			t.fault_desc,
			sum(t.cost) as cost,
			count(1) as cnt
		from 
		(
			select
				mr.fault_desc,
				TIMESTAMPDIFF(HOUR,mr.report_time,mr.finish_time) as cost
			from mold_repair mr
			where mr.is_deleted = 'N' and mr.gmt_created >= ?
		) t
		group by t.fault_desc
	`, time.Now().UnixMilli()-30*24*60*60*1000).Scan(&result)

	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 首页
// @Summary AGV任务统计
// @Accept json
// @Produce json
// @Param body body AgvTaskVO true "AgvTaskVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} AgvTaskScopeVO
// @Router /mold/statistics/agv/task [post]
func AgvTask(c *gin.Context) {
	var vo AgvTaskVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var result AgvTaskScopeVO
	result.X = make([]X, 7)
	result.Y = make([]int64, 7)

	now := time.Now().Local()

	if vo.Type == "day" {
		for i := 0; i < 7; i++ {
			item := now.AddDate(0, 0, -i)
			beginTime := time.Date(item.Year(), item.Month(), item.Day(), 0, 0, 0, 0, time.Local)
			endTime := time.Date(item.Year(), item.Month(), item.Day(), 23, 59, 59, 0, time.Local)
			var count int64
			dao.GetConn().Table("mold_inout_bound_job").Where("agv = 1 and gmt_created >= ? and gmt_created <= ?", beginTime, endTime).Count(&count)
			result.Y[i] = count
			result.X[i].Month = int64(item.Month())
			result.X[i].Year = int64(item.Year())
			result.X[i].Day = int64(item.Day())
		}
	} else if vo.Type == "week" {
		startTime := now.AddDate(0, 0, -49)
		endTime := now
		weekDate := base.GroupByWeekDate(startTime, endTime)

		for i := 0; i < 7; i++ {
			v := weekDate[i]
			var count int64
			dao.GetConn().Table("mold_inout_bound_job").Where("agv = 1 and gmt_created >= ? and gmt_created <= ?", v.StartTime, v.EndTime).Count(&count)
			result.Y[i] = count
			result.X[i].Year = int64(v.Year)
			result.X[i].Week = int64(v.Week)
		}

	} else if vo.Type == "month" {
		for i := 0; i < 7; i++ {
			item := now.AddDate(0, -i, 0)
			nextMonth := item.AddDate(0, 1, 0)
			beginTime := time.Date(item.Year(), item.Month(), 1, 0, 0, 0, 0, time.Local)
			endTime := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 0, 0, 0, 0, time.Local)
			var count int64
			dao.GetConn().Table("mold_inout_bound_job").Where("agv = 1 and gmt_created >= ? and gmt_created < ?", beginTime, endTime).Count(&count)
			result.Y[i] = count
			result.X[i].Month = int64(item.Month())
			result.X[i].Year = int64(item.Year())
		}
	} else {
		panic(base.ParamsErrorN())
	}

	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 首页
// @Summary  出入库任务统计
// @Accept json
// @Produce json
// @Param body body InOutBoundTaskVO true "InOutBoundTaskVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} InOutBoundTaskOutVO
// @Router /mold/statistics/inoutbound/task [post]
func InOutBoundTask(c *gin.Context) {
	var vo InOutBoundTaskVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var result []InOutBoundTaskOutVO

	sql := `
		SELECT 
		    DATE_FORMAT(m.gmt_created, '%Y-%m-%d') AS Date, 
		    COUNT(*) AS Count
		FROM mold_inout_bound_job m
		INNER JOIN mold_info mi ON m.mold_id = mi.id AND mi.is_deleted = 'N'
		WHERE DATE(m.gmt_created) BETWEEN ? AND ?
	`

	if vo.Type != "" {
		sql += fmt.Sprintf(" AND m.type = '%s' ", vo.Type)
	}

	if vo.Status != "" {
		sql += fmt.Sprintf(" AND m.status = '%s' ", vo.Status)
	}

	if vo.Platform != "" {
		sql += fmt.Sprintf(" AND mi.platform = '%s' ", vo.Platform)
	}
	sql += " GROUP BY DATE_FORMAT(m.gmt_created, '%Y-%m-%d'), m.gmt_created ORDER BY m.gmt_created"

	dao.GetConn().Raw(sql, vo.StartDate, vo.EndDate).Scan(&result)

	c.JSON(http.StatusOK, base.Success(result))
}
