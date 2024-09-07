package dashboard

import (
	"crf-mold/base"
	"crf-mold/dao"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// @Tags 看板
// @Summary 产线生产统计
// @Accept json
// @Produce json
// @Param body body LineProductionRequestVO true "LineProductionRequestVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} LineProductionResponseVo
// @Router /dashboard/production [post]
func LineProduction(c *gin.Context) {
	var vo LineProductionRequestVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	//var result []LineProductionVo //old
	var result []LineProductionProjectVo //new
	//var response []LineProductionResponseVo //old
	var response []LineProductionProjectResponseVo //new

	var maintenanceDuration map[string]map[int16]int64
	var maintenanceDurationList []LineMaintenanceDurationVo

	maintenanceSql := `
		SELECT 
			line_level,
			CASE 
				WHEN TIME(gmt_created) BETWEEN '07:00:00' AND '18:30:00' THEN '1'
				ELSE '2'
			END as shift_no,
			SUM(last_time) AS maintenance_duration
		FROM 
			mold_repair 
		WHERE 
			DATE(gmt_created) BETWEEN ? AND ?
			AND is_deleted = 'N'
		GROUP BY 
			line_level,
			shift_no
		ORDER BY 
			line_level,
			shift_no
	`
	dao.GetConn().Raw(maintenanceSql, vo.StartDate, vo.EndDate).Scan(&maintenanceDurationList)
	for _, item := range maintenanceDurationList {
		if maintenanceDuration == nil {
			maintenanceDuration = make(map[string]map[int16]int64)
		}
		if _, ok := maintenanceDuration[item.LineLevel]; !ok {
			maintenanceDuration[item.LineLevel] = make(map[int16]int64)
		}
		maintenanceDuration[item.LineLevel][item.ShiftNo] = item.MaintenanceDuration
	}
	//dao.GetConn().Raw(maintenanceSql, vo.StartDate, vo.EndDate).Scan(&maintenanceDuration)

	sql := `
		SELECT 
			lpr.line_level,
			lpr.shift_no,
			lpr.shift_min,
			lpr.part_code AS part_no,
			SUM(lpr.qty_ok + lpr.qty_nok) AS total_qty,
			mi.project_name
		FROM 
			line_product_resume lpr 
		left join mold_part_rel mpr on 
			lpr.part_code = mpr.part_code and mpr.is_deleted = 'N'
		left join mold_info mi on 
			mpr.mold_code = mi.code and mi.is_deleted = 'N'
		WHERE 
			lpr.shift_date BETWEEN ? AND ?
			AND lpr.is_deleted = 'N'
		GROUP BY 
			lpr.line_level,
			lpr.shift_no, 
			lpr.shift_min,
			lpr.part_code
		ORDER BY 
			lpr.line_level,
			lpr.shift_no
	`
	dao.GetConn().Raw(sql, vo.StartDate, vo.EndDate).Scan(&result)
	resultLen := len(result)
	if resultLen > 0 {
		lineLevel := result[0].LineLevel
		responseVo := LineProductionProjectResponseVo{
			LineLevel: lineLevel,
			DayShift: ShiftProductionProjectVo{
				ShiftNo:             0,
				ShiftMin:            0,
				MaintenanceDuration: 0,
				TotalQty:            0,
				PartList:            make([]PartQtyProjectVo, 0),
			},
			NightShift: ShiftProductionProjectVo{
				ShiftNo:             0,
				ShiftMin:            0,
				MaintenanceDuration: 0,
				TotalQty:            0,
				PartList:            make([]PartQtyProjectVo, 0),
			},
		}
		for i := 0; i < resultLen; i++ {
			if result[i].LineLevel != lineLevel {
				response = append(response, responseVo)
				lineLevel = result[i].LineLevel
				responseVo = LineProductionProjectResponseVo{
					LineLevel: lineLevel,
					DayShift: ShiftProductionProjectVo{
						ShiftNo:             0,
						ShiftMin:            0,
						MaintenanceDuration: 0,
						TotalQty:            0,
						PartList:            make([]PartQtyProjectVo, 0),
					},
					NightShift: ShiftProductionProjectVo{
						ShiftNo:             0,
						ShiftMin:            0,
						MaintenanceDuration: 0,
						TotalQty:            0,
						PartList:            make([]PartQtyProjectVo, 0),
					},
				}
			}
			if result[i].ShiftNo == 1 {
				responseVo.DayShift.ShiftNo = 1
				responseVo.DayShift.ShiftMin = result[i].ShiftMin
				responseVo.DayShift.MaintenanceDuration = maintenanceDuration[lineLevel][int16(1)]
				responseVo.DayShift.TotalQty += result[i].TotalQty
				if result[i].TotalQty > 0 {
					responseVo.DayShift.PartList = append(responseVo.DayShift.PartList, PartQtyProjectVo{
						PartNo:      result[i].PartNo,
						Qty:         result[i].TotalQty,
						ProjectName: result[i].ProjectName,
					})
				}
			} else {
				responseVo.NightShift.ShiftNo = 2
				responseVo.NightShift.ShiftMin = result[i].ShiftMin
				responseVo.NightShift.MaintenanceDuration = maintenanceDuration[lineLevel][int16(2)]
				responseVo.NightShift.TotalQty += result[i].TotalQty
				if result[i].TotalQty > 0 {
					responseVo.NightShift.PartList = append(responseVo.NightShift.PartList, PartQtyProjectVo{
						PartNo:      result[i].PartNo,
						Qty:         result[i].TotalQty,
						ProjectName: result[i].ProjectName,
					})
				}
			}
			if i == resultLen-1 {
				response = append(response, responseVo)
				break
			}
		}
	} else {
		response = make([]LineProductionProjectResponseVo, 0)
	}

	c.JSON(http.StatusOK, base.Success(response))
}
