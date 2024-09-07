/**
 * 模具维修管理
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 模具维修
// @Summary 模具维修分页
// @Accept json
// @Produce json
// @Param query query PageRepairVO true "PageRepairVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageRepairOutVO
// @Router /mold/repair/page [get]
func PageRepair(c *gin.Context) {
	var vo PageRepairVO
	var sumLastTime int64
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	var moldRepair model.MoldRepair
	base.CopyProperties(&moldRepair, vo)

	tx := dao.GetConn().Table("mold_repair")
	if vo.CodeOrName == "" {
		tx = dao.BuildWhereCondition(tx, moldRepair)
		// 额外处理时间字段
		if vo.ReportTimeBegin != "" {
			tx.Where("report_time >= ?", vo.ReportTimeBegin)
		}
		if vo.ReportTimeEnd != "" {
			tx.Where("report_time <= ?", vo.ReportTimeEnd)
		}
		if vo.ArriveTimeBegin != "" {
			tx.Where("arrive_time >= ?", vo.ArriveTimeBegin)
		}
		if vo.ArriveTimeEnd != "" {
			tx.Where("arrive_time <= ?", vo.ArriveTimeEnd)
		}
		if vo.FinishTimeBegin != "" {
			tx.Where("finish_time >= ?", vo.FinishTimeBegin)
		}
		if vo.FinishTimeEnd != "" {
			tx.Where("finish_time <= ?", vo.FinishTimeEnd)
		}
		if vo.RepairStations != "" {
			tx.Where("repair_station like concat('%',?,'%')", vo.RepairStations)
		}
		if vo.PartCodes != "" {
			var moldCodeList []string
			dao.GetConn().Table("mold_part_rel").Select("mold_code").Where(
				"is_deleted = 'N' and part_code like concat('%',?,'%')", vo.PartCodes).Scan(&moldCodeList)
			tx.Where("mold_code in (?)", moldCodeList)
		}
		// 编码模糊
		if vo.CodeLike != "" {
			tx.Where("code like concat('%',?,'%')", vo.CodeLike)
		}
	} else {
		var moldCodeList []string
		dao.GetConn().Table("mold_part_rel").Select("mold_code").Where(
			"is_deleted = 'N' and part_code like concat('%',?,'%')", vo.CodeOrName).Scan(&moldCodeList)
		tx.Where("mold_code in (?) or mold_code like concat('%',?,'%')", moldCodeList, vo.CodeOrName)
	}

	tx = tx.Where("is_deleted = 'N'").Order("gmt_created desc")

	// 统计总的维修时长
	tx.Select("sum(`last_time`)").Row().Scan(&sumLastTime)

	tx.Select("*")

	var result []model.MoldRepair
	page := base.Page(tx, &result, vo.GetCurrentPage(), vo.GetSize())

	// 转换为outvo
	vos := page.List.(*[]model.MoldRepair)
	pageVos, _ := base.CopyPropertiesList(reflect.TypeOf(PageRepairOutVO{}), vos).([]interface{})

	// 更换备件字段
	var vv []PageRepairOutVO
	for i := 0; i < len(pageVos); i++ {
		v := pageVos[i].(PageRepairOutVO)

		var moldRepairSpareRel []model.MoldReplaceSpareRel
		dao.GetConn().Table("mold_replace_spare_rel").Where("is_deleted = 'N' and repair_id = ? and type = 'repair'", v.ID).Find(&moldRepairSpareRel)

		// 维修备件
		var repairSpares []RepairSpare
		for i := 0; i < len(moldRepairSpareRel); i++ {
			repairSpares = append(repairSpares, RepairSpare{
				Code:  moldRepairSpareRel[i].SpareCode,
				Count: moldRepairSpareRel[i].Count,
			})
		}

		if len(repairSpares) == 0 {
			v.MoldRepairSpares = []RepairSpare{}
		} else {
			v.MoldRepairSpares = repairSpares
		}
		var moldPartRel []model.MoldPartRel
		dao.GetConn().Table("mold_part_rel").Where("is_deleted = 'N' and mold_code = ?", v.MoldCode).Find(&moldPartRel)
		if len(moldPartRel) > 0 {
			var partCodes []string
			for _, rel := range moldPartRel {
				partCodes = append(partCodes, rel.PartCode)
			}
			v.PartCodes = partCodes
			v.PartCodeStr = strings.Join(partCodes, "/")
		} else {
			v.PartCodes = []string{}
		}
		if v.RepairStation != "" {
			v.RepairStations = strings.Split(v.RepairStation, ",")
		} else {
			v.RepairStations = []string{}
		}

		v.SumLastTime = sumLastTime

		vv = append(vv, v)
	}

	if len(pageVos) == 0 {
		page.List = []interface{}{}
	} else {
		page.List = vv
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 模具维修
// @Summary 模具维修查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageRepairOutVO
// @Router /mold/repair/one [get]
func OneRepair(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var moldRepair []model.MoldRepair
	dao.GetConn().Table("mold_repair").Where("id = ? and is_deleted = 'N'", id).Find(&moldRepair)

	if len(moldRepair) > 0 {
		m := moldRepair[0]
		var vo PageRepairOutVO
		base.CopyProperties(&vo, m)

		// 更换备件字段
		var moldRepairSpareRel []model.MoldReplaceSpareRel
		dao.GetConn().Table("mold_replace_spare_rel").Where("is_deleted = 'N' and repair_id = ? and type = 'repair'", id).Find(&moldRepairSpareRel)

		var repairSpares []RepairSpare
		for i := 0; i < len(moldRepairSpareRel); i++ {
			repairSpares = append(repairSpares, RepairSpare{
				Code:  moldRepairSpareRel[i].SpareCode,
				Count: moldRepairSpareRel[i].Count,
			})
		}

		if len(repairSpares) == 0 {
			vo.MoldRepairSpares = []RepairSpare{}
		} else {
			vo.MoldRepairSpares = repairSpares
		}
		var moldPartRel []model.MoldPartRel
		dao.GetConn().Table("mold_part_rel").Where("is_deleted = 'N' and mold_code = ?", vo.MoldCode).Find(&moldPartRel)
		if len(moldPartRel) > 0 {
			var partCodes []string
			for _, rel := range moldPartRel {
				partCodes = append(partCodes, rel.PartCode)
			}
			vo.PartCodes = partCodes
			vo.PartCodeStr = strings.Join(partCodes, "/")
		} else {
			vo.PartCodes = []string{}
		}
		if vo.RepairStation != "" {
			vo.RepairStations = strings.Split(vo.RepairStation, ",")
		} else {
			vo.RepairStations = []string{}
		}

		c.JSON(http.StatusOK, base.Success(vo))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 模具维修
// @Summary 删除模具维修
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/repair [delete]
func DeleteRepair(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	ids := vo.IDS
	if len(ids) == 0 {
		panic(base.ParamsErrorN())
	}

	dao.GetConn().Table("mold_repair").Where("id in ?", ids).Updates(&model.MoldRepair{
		IsDeleted: "Y",
		UpdatedBy: c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具维修
// @Summary 添加模具维修
// @Accept json
// @Produce json
// @Param Body body CreateRepairVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/repair [post]
func CreateRepair(c *gin.Context) {
	var vo CreateRepairVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 参数校验,TODO

	var moldRepair model.MoldRepair
	base.CopyProperties(&moldRepair, vo)

	moldRepair.Code = common.GenerateCode("WX")

	userId := c.GetHeader(constant.USERID)
	moldRepair.CreatedBy = userId
	moldRepair.UpdatedBy = userId
	moldRepair.RepairStation = strings.Join(vo.RepairStations, ",")

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("mold_repair").Create(&moldRepair).Error; err != nil {
		panic(err)
	}

	// 模具ID
	var moldInfo model.MoldInfo
	tx.Table("mold_info").Where("is_deleted = 'N' and code = ?", moldRepair.MoldCode).Select("id").First(&moldInfo)

	// 维修更换备件
	if len(vo.RepairSpares) != 0 {
		if err := SaveRepairSpare(tx, vo.RepairSpares, moldRepair.ID, moldInfo.ID, userId); err != nil {
			panic(err)
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(moldRepair.ID))
}

// @Tags 模具维修
// @Summary 更新模具维修
// @Accept json
// @Produce json
// @Param Body body UpdateRepairVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/repair [put]
func UpdateRepair(c *gin.Context) {
	var vo UpdateRepairVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 参数校验,TODO

	if vo.ID == 0 {
		panic(base.ParamsErrorN())
	}

	// 存在
	var one model.MoldRepair
	if err := dao.GetConn().Table("mold_repair").Where("id = ?", vo.ID).First(&one).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 更新
	userId := c.GetHeader(constant.USERID)
	var moldRepair model.MoldRepair
	base.CopyProperties(&moldRepair, vo)
	moldRepair.UpdatedBy = userId
	moldRepair.RepairStation = strings.Join(vo.RepairStations, ",")

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("mold_repair").Where("id = ?", moldRepair.ID).Updates(&moldRepair).Error; err != nil {
		panic(err)
	}
	// 针对时间置空的更新
	if err := tx.Table("mold_repair").Where("id = ?", moldRepair.ID).Updates(map[string]interface{}{
		"report_time": moldRepair.ReportTime, "arrive_time": moldRepair.ArriveTime, "finish_time": moldRepair.FinishTime, "last_time": moldRepair.LastTime}).Error; err != nil {
		panic(err)
	}

	// 模具ID
	var moldInfo model.MoldInfo
	tx.Table("mold_info").Where("is_deleted = 'N' and code = ?", moldRepair.MoldCode).Select("id").First(&moldInfo)

	// 维修更换备件
	if len(vo.RepairSpares) != 0 {
		if err := SaveRepairSpare(tx, vo.RepairSpares, moldRepair.ID, moldInfo.ID, userId); err != nil {
			panic(err)
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}
