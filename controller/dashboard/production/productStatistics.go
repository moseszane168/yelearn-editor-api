package dashboard

import (
	"crf-mold/base"
	"crf-mold/common/cron"
	"crf-mold/dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// @Tags 看板
// @Summary 维修停机率-周月报（柱图）
// @Accept json
// @Produce json
// @Param body body LineProductionRequestVO true "LineProductionRequestVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} OutageFactorVo
// @Router /dashboard/stopRate [post]
func ProductStatistics(c *gin.Context) {
	var vo LineProductionRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	var result []OutageFactorVo //old

	sql := `SELECT
				p1.line_level,
				p1.shift_no,
				p1.totle_time,
				SUM(p2.last_time) AS last_time,
				any_value(FORMAT(( p2.last_time / p1.totle_time )* 100, 2 )) AS rate_num
			FROM
				(
				SELECT
					t1.line_level,
					t1.shift_no,
					any_value (
					SUM( t1.shift_min )) AS totle_time 
				FROM
					line_product_resume t1 
				WHERE
					t1.gmt_created BETWEEN ? 
					AND ? 
					AND t1.is_deleted = 'N' 
				GROUP BY
					t1.shift_no,
					t1.line_level 
				ORDER BY
					t1.line_level,
					t1.shift_no 
				) p1
				LEFT JOIN (
				SELECT
					t2.line_level,
					t2.last_time,
				CASE
						WHEN TIME( t2.gmt_created ) BETWEEN '07:00:00' 
						AND '18:30:00' THEN
							'1' ELSE '2' 
							END AS shift_no 
					FROM
						mold_repair t2 
					WHERE
						t2.gmt_created BETWEEN ? 
						AND ? 
						AND t2.is_deleted = 'N'
					) p2 ON p1.line_level = p2.line_level 
				AND p1.shift_no = p2.shift_no GROUP BY p1.line_level,p1.shift_no ORDER BY p1.line_level+0,p1.shift_no+0`

	parse, _ := time.Parse("2006-01-02", vo.EndDate)
	dd, _ := time.ParseDuration("24h")
	dd1 := parse.Add(dd)
	endDate := dd1.Format("2006-01-02")

	dao.GetConn().Raw(sql, vo.StartDate, endDate, vo.StartDate, endDate).Scan(&result)

	// 创建一个新的结构体切片，使用更新后的结构体
	resultWithRate := make([]OutageFactorPercentVo, len(result))

	// 复制原始数据到新切片，并为新字段赋值
	for i, res := range result {
		resultWithRate[i] = OutageFactorPercentVo{
			OutageFactorVo: res,
			Rate:           strconv.FormatFloat(float64(res.RateNum), 'f', -1, 32) + "%",
		}
	}

	c.JSON(http.StatusOK, base.Success(resultWithRate))
}

// @Tags 看板
// @Summary 日报-质量列表
// @Accept json
// @Produce json
// @Param body body LineProductionRequestVO true "LineProductionRequestVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} RejectsProducts
// @Router /dashboard/rejects [post]
func DefectiveProducts(c *gin.Context) {
	var vo LineProductionRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	var result []RejectsProducts
	sql := `SELECT
			line_level,
			defect_desc,
			defect_count 
		FROM
			mold_quality 
		WHERE
			DATE_FORMAT( gmt_created, '%Y-%m-%d' )= ? AND is_deleted='N'`
	dao.GetConn().Raw(sql, vo.StartDate).Scan(&result)
	c.JSON(http.StatusOK, base.Success(result))
}

//备件消耗趋势
// @Tags 看板
// @Summary 备件消耗趋势-日周月报
// @Accept json
// @Produce json
// @Param body body LineProductionRequestVO true "LineProductionRequestVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} UserSparePartVo
// @Router /dashboard/userSparePart [post]
func UserSparePart(c *gin.Context) {
	var vo LineProductionRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	var result []UserSparePartVo
	parse, _ := time.Parse("2006-01-02", vo.EndDate)
	dd, _ := time.ParseDuration("24h")
	dd1 := parse.Add(dd)
	endDate := dd1.Format("2006-01-02")
	sql := `SELECT DATE_FORMAT(t1.gmt_created,'%Y-%m-%d') AS date_time, SUM(t1.count) AS totle_num FROM mold_replace_spare_rel t1 WHERE t1.gmt_created BETWEEN ? AND ? GROUP BY DATE_FORMAT(t1.gmt_created,'%Y-%m-%d')`
	dao.GetConn().Raw(sql, vo.StartDate, endDate).Scan(&result)
	dateList := GetBetweenDates(vo.StartDate, vo.EndDate)
	var uspl UserSparePartList
	uspl.UspvList = result
	uspl.DateList = dateList

	c.JSON(http.StatusOK, base.Success(uspl))
}

// @Tags 看板
// @Summary 故障top10(按时间按次数)
// @Accept json
// @Produce json
// @Param body body LineProductionReportRequestVO true "LineProductionReportRequestVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} ErrorTopTens
// @Router /dashboard/errorTopTen [post]
func ErrorTopTen(c *gin.Context) {
	var vo LineProductionReportRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	var result []ErrorTopTens
	parse, _ := time.Parse("2006-01-02", vo.EndDate)
	dd, _ := time.ParseDuration("24h")
	dd1 := parse.Add(dd)
	endDate := dd1.Format("2006-01-02")
	var sql = "SELECT t1.fault_desc,COUNT(t1.id) count_num, SUM(t1.last_time) last_time FROM mold_repair t1 WHERE t1.gmt_created BETWEEN ? AND ? AND t1.fault_desc IS NOT NULL AND t1.fault_desc <>'' GROUP BY t1.fault_desc "
	if vo.OrderField == "count" {
		sql = sql + " ORDER BY COUNT(t1.id) DESC LIMIT 10"
	} else {
		sql = sql + " ORDER BY SUM(t1.last_time) DESC LIMIT 10"
	}
	dao.GetConn().Raw(sql, vo.StartDate, endDate).Scan(&result)
	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 看板
// @Summary 日报-维修列表
// @Accept json
// @Produce json
// @Param body body LineProductionRequestVO true "LineProductionRequestVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} MaintainList
// @Router /dashboard/maintain [post]
func MaintainListPlan(c *gin.Context) {
	var vo LineProductionRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	var result []MaintainList
	sql := `SELECT
		t1.line_level,
		t1.repair_content,
		CASE WHEN TIME( t1.gmt_created ) BETWEEN '07:00:00' AND '18:30:00' THEN '1' ELSE '2' END AS shift_no,
		t1.last_time 
		FROM
			mold_repair t1 
	WHERE
		DATE_FORMAT( gmt_created, '%Y-%m-%d' )= ?`
	dao.GetConn().Raw(sql, vo.StartDate).Scan(&result)
	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 看板
// @Summary 趋势图-周/月趋势图
// @Accept json
// @Produce json
// @Param body body LineProductionReportRequestVO true "LineProductionReportRequestVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} TendencyChartVo
// @Router /dashboard/tendencyChart [post]
func TendencyChart(c *gin.Context) {
	var vo LineProductionReportRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	var tcvList []TendencyChartVo
	var resultList []TendencyChartVo
	startDate := fmt.Sprintf("%02s", vo.StartDate)
	endDate := fmt.Sprintf("%02s", vo.EndDate)
	var lineLevelList []string = strings.Split(vo.LineLevel, ",")
	var sql = "SELECT t1.line_level,t1.report_date,t1.report_value1,t1.report_value2,t1.value1_desc,t1.value2_desc FROM mold_product_report t1 WHERE t1.report_code=? AND t1.report_date BETWEEN ? AND ? "
	if !strings.HasPrefix(vo.ReportCode, "UserSparePart") && !strings.HasPrefix(vo.ReportCode, "ErrorTimes") && len(strings.TrimSpace(vo.LineLevel)) > 0 {
		sql = sql + " AND t1.line_level IN (?) "
		dao.GetConn().Raw(sql, vo.ReportCode, startDate, endDate, lineLevelList).Scan(&tcvList)
	} else {
		dao.GetConn().Raw(sql, vo.ReportCode, startDate, endDate).Scan(&tcvList)
	}
	start, _ := strconv.Atoi(vo.StartDate)
	end, _ := strconv.Atoi(vo.EndDate)
	lo, hi := start, end
	datetime := make([]int, hi-lo+1)
	for i := range datetime {
		datetime[i] = i + lo
	}

	var reqlineList []string
	var result []LineList
	if strings.HasPrefix(vo.ReportCode, "UserSparePart") || strings.HasPrefix(vo.ReportCode, "ErrorTimes") {
		reqlineList = append(reqlineList, vo.ReportCode)
	} else if vo.LineLevel == "" {
		sql2 := `SELECT t1.key AS line_level,t1.value_cn AS line_name FROM dict_property t1 WHERE t1.group_code='line' AND t1.is_deleted='N' ORDER BY t1.key+0`
		dao.GetConn().Raw(sql2).Scan(&result)
		for _, line1 := range result {
			reqlineList = append(reqlineList, line1.LineLevel)
		}
	} else {
		for _, line2 := range lineLevelList {
			reqlineList = append(reqlineList, line2)
		}
	}
	for _, lineKey := range reqlineList {
		fmt.Println(lineKey)
		for _, value := range datetime {
			flag := 0
			for _, item := range tcvList {
				reportDatei, _ := strconv.Atoi(item.ReportDate)
				if reportDatei == value && item.LineLevel == lineKey {
					resultList = append(resultList, item)
					flag++
					break
				}
			}
			if flag == 0 {
				addItem := TendencyChartVo{
					LineLevel:    lineKey,
					ReportDate:   fmt.Sprintf("%02s", strconv.Itoa(value)),
					ReportValue1: 0,
					ReportValue2: 0,
				}
				resultList = append(resultList, addItem)
			}
		}
	}
	c.JSON(http.StatusOK, base.Success(resultList))
}

// @Tags 看板
// @Summary 趋势图-产线列表
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {object} LineList
// @Router /dashboard/getLineList [post]
func GetLineList(c *gin.Context) {
	var result []LineList
	sql := `SELECT t1.key AS line_level,t1.value_cn AS line_name FROM dict_property t1 WHERE t1.group_code='line' AND t1.is_deleted='N' ORDER BY t1.key+0`
	dao.GetConn().Raw(sql).Scan(&result)
	c.JSON(http.StatusOK, base.Success(result))
}

func CronProductReport() {
	// 每周日23：50分生成周报月报趋势图   0 50 23 ? * SUN
	if _, err := cron.CronJob.CronWithSeconds("0 50 23 * * ?").Do(GenerateProductReportData); err != nil {
		logrus.Error(err.Error())
	}
}

// @Tags 看板
// @Summary 趋势图-数据生成手动触发接口
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {object} string
// @Router /dashboard/generateProductReport [post]
func GenerateProductReport(c *gin.Context) {
	GenerateProductReportData()
	c.JSON(http.StatusOK, base.SuccessN())
}

func GenerateProductReportData() {
	fmt.Println("定时任务GenerateProductReportData")
	ProductTotleByWeek()
	ProductTotleByMonth()

	StopLastTimeByWeek()
	StopLastTimeByMonth()

	ErrorTimesByWeek()
	ErrorTimesByMonth()

	OutageRateByWeek()
	OutageRateByMonth()

	UserSparePartByWeek()
	UserSparePartByMonth()
}

//按周统计生产数量
func ProductTotleByWeek() {
	year, week := GetWeek(time.Now())
	startDate := GetFirstDateOfWeek(time.Now()).Format("2006-01-02")
	endDate := GetNextFirstDateOfWeek(time.Now()).Format("2006-01-02")
	fmt.Println(year, week, startDate, endDate)
	var productTotleList []ProductTotle
	sql := `SELECT t1.line_level,any_value(SUM(t1.qty_ok)) AS totle_num FROM line_product_resume t1 WHERE t1.shift_date BETWEEN ? AND ? AND t1.is_deleted='N' GROUP BY t1.line_level ORDER BY t1.line_level+0`
	dao.GetConn().Raw(sql, startDate, endDate).Scan(&productTotleList)
	for _, item := range productTotleList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, item.LineLevel, "ProductTotleByWeek", week, item.TotleNum)
	}
}

//按月统计生产数量
func ProductTotleByMonth() {
	firstDay, lastDay, totalDays := GetMonthDay(time.Now())
	fmt.Println(firstDay, lastDay, totalDays)
	var productTotleList []ProductTotle
	sql := `SELECT t1.line_level,any_value(SUM(t1.qty_ok)) AS totle_num FROM line_product_resume t1 WHERE t1.shift_date BETWEEN ? AND ? AND t1.is_deleted='N' GROUP BY t1.line_level ORDER BY t1.line_level+0`
	dao.GetConn().Raw(sql, firstDay, lastDay).Scan(&productTotleList)
	for _, item := range productTotleList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, item.LineLevel, "ProductTotleByMonth", time.Now().Format("01"), item.TotleNum)
	}
}

//按周统计停机工时
func StopLastTimeByWeek() {
	year, week := GetWeek(time.Now())
	startDate := GetFirstDateOfWeek(time.Now()).Format("2006-01-02")
	endDate := GetNextFirstDateOfWeek(time.Now()).Format("2006-01-02")
	fmt.Println(year, week, startDate, endDate)
	var productTotleList []ProductTotle
	sql := `SELECT t1.line_level,any_value(t1.last_time) AS totle_num FROM mold_repair t1 WHERE t1.gmt_created BETWEEN ? AND ? AND t1.is_deleted='N' GROUP BY t1.line_level ORDER BY t1.line_level+0`
	dao.GetConn().Raw(sql, startDate, endDate).Scan(&productTotleList)
	for _, item := range productTotleList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, item.LineLevel, "StopLastTimeByWeek", week, item.TotleNum)
	}
}

//按月统计停机工时
func StopLastTimeByMonth() {
	firstDay, lastDay, totalDays := GetMonthDay(time.Now())
	fmt.Println(firstDay, lastDay, totalDays)
	var productTotleList []ProductTotle
	sql := `SELECT t1.line_level,any_value(t1.last_time) AS totle_num FROM mold_repair t1 WHERE t1.gmt_created BETWEEN ? AND ? AND t1.is_deleted='N' GROUP BY t1.line_level ORDER BY t1.line_level+0`
	dao.GetConn().Raw(sql, firstDay, lastDay).Scan(&productTotleList)
	for _, item := range productTotleList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, item.LineLevel, "StopLastTimeByMonth", time.Now().Format("01"), item.TotleNum)
	}
}

//按周统计故障次数
func ErrorTimesByWeek() {
	year, week := GetWeek(time.Now())
	startDate := GetFirstDateOfWeek(time.Now()).Format("2006-01-02")
	endDate := GetNextFirstDateOfWeek(time.Now()).Format("2006-01-02")
	fmt.Println(year, week, startDate, endDate)
	var productTotleList []ProductTotle
	sql := `SELECT COUNT(t1.id) AS totle_num FROM mold_repair t1 WHERE t1.gmt_created BETWEEN ? AND ? AND t1.is_deleted='N'`
	dao.GetConn().Raw(sql, startDate, endDate).Scan(&productTotleList)
	for _, item := range productTotleList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, "ErrorTimesByWeek", "ErrorTimesByWeek", week, item.TotleNum)
	}
}

//按月统计故障次数
func ErrorTimesByMonth() {
	firstDay, lastDay, totalDays := GetMonthDay(time.Now())
	fmt.Println(firstDay, lastDay, totalDays)
	var productTotleList []ProductTotle
	sql := `SELECT COUNT(t1.id) AS totle_num FROM mold_repair t1 WHERE t1.gmt_created BETWEEN ? AND ? AND t1.is_deleted='N'`
	dao.GetConn().Raw(sql, firstDay, lastDay).Scan(&productTotleList)
	for _, item := range productTotleList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, "ErrorTimesByMonth", "ErrorTimesByMonth", time.Now().Format("01"), item.TotleNum)
	}
}

//按周统计维修停机率
func OutageRateByWeek() {
	year, week := GetWeek(time.Now())
	startDate := GetFirstDateOfWeek(time.Now()).Format("2006-01-02")
	endDate := GetNextFirstDateOfWeek(time.Now()).Format("2006-01-02")
	fmt.Println(year, week, startDate, endDate)
	var outageFactorList []OutageFactorVo
	sql := `SELECT
	p1.line_level,
	any_value(FORMAT(( p2.last_time / p1.totle_time )* 100, 2 )) AS rate_num
FROM
	(
	SELECT
		t1.line_level,
		any_value (
		SUM( t1.shift_min )) AS totle_time 
	FROM
		line_product_resume t1 
	WHERE
		t1.gmt_created BETWEEN ? 
			AND ?
		AND t1.is_deleted = 'N' 
	GROUP BY
		t1.line_level 
	ORDER BY
		t1.line_level+0
	) p1
	LEFT JOIN (
	SELECT
		t2.line_level,
		t2.last_time
		FROM
			mold_repair t2 
		WHERE
			t2.gmt_created BETWEEN ? 
			AND ? 
			AND t2.is_deleted = 'N'
		) p2 ON p1.line_level = p2.line_level 
	 GROUP BY p1.line_level ORDER BY p1.line_level+0`
	dao.GetConn().Raw(sql, startDate, endDate, startDate, endDate).Scan(&outageFactorList)
	for _, item := range outageFactorList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, item.LineLevel, "OutageRateByWeek", week, item.RateNum)
	}
}

//按月统计维修停机率
func OutageRateByMonth() {
	firstDay, lastDay, totalDays := GetMonthDay(time.Now())
	fmt.Println(firstDay, lastDay, totalDays)
	var outageFactorList []OutageFactorVo
	sql := `SELECT
	p1.line_level,
	any_value(FORMAT(( p2.last_time / p1.totle_time )* 100, 2 )) AS rate_num
FROM
	(
	SELECT
		t1.line_level,
		any_value (
		SUM( t1.shift_min )) AS totle_time 
	FROM
		line_product_resume t1 
	WHERE
		t1.gmt_created BETWEEN ? 
			AND ?
		AND t1.is_deleted = 'N' 
	GROUP BY
		t1.line_level 
	ORDER BY
		t1.line_level+0
	) p1
	LEFT JOIN (
	SELECT
		t2.line_level,
		t2.last_time
		FROM
			mold_repair t2 
		WHERE
			t2.gmt_created BETWEEN ? 
			AND ? 
			AND t2.is_deleted = 'N'
		) p2 ON p1.line_level = p2.line_level 
	 GROUP BY p1.line_level ORDER BY p1.line_level+0`
	dao.GetConn().Raw(sql, firstDay, lastDay, firstDay, lastDay).Scan(&outageFactorList)
	for _, item := range outageFactorList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, item.LineLevel, "OutageRateByMonth", time.Now().Format("01"), item.RateNum)
	}
}

//按周统计备件消耗趋势
func UserSparePartByWeek() {
	year, week := GetWeek(time.Now())
	startDate := GetFirstDateOfWeek(time.Now()).Format("2006-01-02")
	endDate := GetNextFirstDateOfWeek(time.Now()).Format("2006-01-02")
	fmt.Println(year, week, startDate, endDate)
	var uspList []UserSparePartVo
	sql := `SELECT DATE_FORMAT(t1.gmt_created,'%Y-%m-%d') AS date_time, SUM(t1.count) totle_num FROM mold_replace_spare_rel t1 WHERE t1.gmt_created BETWEEN ? AND ? GROUP BY DATE_FORMAT(t1.gmt_created,'%Y-%m-%d')`
	dao.GetConn().Raw(sql, startDate, endDate).Scan(&uspList)
	for _, item := range uspList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?, ?)`
		dao.GetConn().Exec(sql, "UserSparePartByWeek", "UserSparePartByWeek", week, item.TotleNum)
	}
}

//按月统计备件消耗趋势
func UserSparePartByMonth() {
	firstDay, lastDay, totalDays := GetMonthDay(time.Now())
	fmt.Println(firstDay, lastDay, totalDays)
	var uspList []UserSparePartVo
	sql := `SELECT DATE_FORMAT(t1.gmt_created,'%Y-%m-%d') AS date_time, SUM(t1.count) totle_num FROM mold_replace_spare_rel t1 WHERE t1.gmt_created BETWEEN ? AND ? GROUP BY DATE_FORMAT(t1.gmt_created,'%Y-%m-%d')`
	dao.GetConn().Raw(sql, firstDay, lastDay).Scan(&uspList)
	for _, item := range uspList {
		sql := `REPLACE INTO mold_product_report(line_level, report_code, report_date, report_value1) VALUES (?, ?, ?)`
		dao.GetConn().Exec(sql, "UserSparePartByMonth", "UserSparePartByMonth", time.Now().Format("01"), item.TotleNum)
	}
}

//获取当前日期是该年第几周
func GetWeek(t time.Time) (y, w int) {
	return t.ISOWeek()
}

// GetFirstDateOfWeek 获取本周周一的日期
func GetFirstDateOfWeek(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).
		AddDate(0, 0, offset)
}

// GetNextFirstDateOfWeek 获取下周周一
func GetNextFirstDateOfWeek(t time.Time) time.Time {
	return GetFirstDateOfWeek(t).AddDate(0, 0, 7)
}

//获取当月第一天、最后一天及总天数
func GetMonthDay(monthAt time.Time) (firstDay, lastDay string, totalDays int) {
	firstDay = monthAt.AddDate(0, 0, -monthAt.Day()+1).Format("2006-01-02")
	lastDay = monthAt.AddDate(0, 1, -monthAt.Day()).Add(24 * time.Hour).Format("2006-01-02")
	totalDays = monthAt.AddDate(0, 1, -monthAt.Day()).Day()
	return
}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
func GetBetweenDates(sdate, edate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02"
	if len(timeFormatTpl) != len(sdate) {
		timeFormatTpl = timeFormatTpl[0:len(sdate)]
	}
	date, err := time.Parse(timeFormatTpl, sdate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, edate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}
