/**
 * 模具质量管理
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"strings"
)

// @Tags 模具质量
// @Summary 模具质量分页
// @Accept json
// @Produce json
// @Param query query PageQualityVO true "PageQualityVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageQualityOutVO
// @Router /mold/quality/page [get]
func PageQuality(c *gin.Context) {
	var vo PageQualityVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	condition := ""
	if vo.LikeQuery == "" {
		if vo.ProductionCountMin > 0 {
			condition += fmt.Sprintf(" and mq.production_count >= %d", vo.ProductionCountMin)
		}
		if vo.ProductionCountMax > 0 {
			condition += fmt.Sprintf(" and mq.production_count <= %d", vo.ProductionCountMax)
		}
		if vo.DefectCountMin > 0 {
			condition += fmt.Sprintf(" and mq.defect_count >= %d", vo.DefectCountMin)
		}
		if vo.DefectCountMax > 0 {
			condition += fmt.Sprintf(" and mq.defect_count <= %d", vo.DefectCountMax)
		}
		if vo.LineLevel != "" {
			condition += fmt.Sprintf(" and mq.line_level = '%s'", vo.LineLevel)
		}
		if vo.MoldCode != "" {
			condition += fmt.Sprintf(" and mi.code = '%s'", vo.MoldCode)
		}
		if vo.OrderCode != "" {
			condition += fmt.Sprintf(" and mq.order_code = '%s'", vo.OrderCode)
		}
		if vo.PartCodes != "" {
			condition += fmt.Sprintf(" and mpr.part_code = '%s'", vo.PartCodes)
		}
		if vo.DefectDesc != "" {
			condition += fmt.Sprintf(` and mq.defect_desc like concat('%%', '%s' '%%')`, vo.DefectDesc)
		}
		if vo.QualityContent != "" {
			condition += fmt.Sprintf(` and mq.quality_content like concat('%%', '%s' '%%')`, vo.QualityContent)
		}
		if vo.CreatedBy != "" {
			condition += fmt.Sprintf(` and (ui.name like concat('%%', '%s' '%%') or mq.created_by like concat('%%', '%s' '%%'))`, vo.CreatedBy, vo.CreatedBy)
		}
	} else {
		condition += fmt.Sprintf(` and (mi.code like concat('%%', '%s' '%%') or mpr.part_code like concat('%%', '%s' '%%') or mq.line_level like concat('%%', '%s' '%%'))`, vo.LikeQuery, vo.LikeQuery, vo.LikeQuery)
	}

	sql := `
		SELECT 
			mq.id,
			mq.line_level,
			mq.mold_id,
			mi.code AS mold_code,
			mq.order_code,
			mq.production_count,
			mq.defect_count,
			mq.defect_desc,
			mq.quality_content,
			ui.name AS created_by,
			uop.name AS updated_by,
			mq.gmt_created,
			mq.gmt_updated,
			ifnull(GROUP_CONCAT(DISTINCT mpr.part_code SEPARATOR '/'),'') as part_codes
		FROM mold_quality mq
		LEFT JOIN mold_info mi on mi.id = mq.mold_id and mi.is_deleted = 'N'
		LEFT JOIN mold_part_rel mpr on mpr.mold_code = mi.code and mpr.is_deleted = 'N'
		LEFT JOIN user_info ui on ui.login_name = mq.created_by
		LEFT JOIN user_info uop on uop.login_name = mq.updated_by
		WHERE mq.is_deleted = 'N' 
	`
	sql = sql + condition + `
		GROUP BY mq.id, mq.line_level, mq.mold_id, mi.code, mq.order_code, ui.name, uop.name, mq.gmt_created, mq.gmt_updated
		ORDER BY mq.gmt_updated desc, mq.gmt_created desc
	`
	var result []PageQualityOutVO
	page := base.PageWithRawSQL(dao.GetConn(), &result, vo.GetCurrentPage(), vo.GetSize(), sql)
	if len(result) == 0 {
		page.List = []interface{}{}
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 模具质量
// @Summary 模具质量查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageQualityOutVO
// @Router /mold/quality/one [get]
func OneQuality(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var moldQuality []model.MoldQuality
	dao.GetConn().Table("mold_quality").Where("id = ? and is_deleted = 'N'", id).Find(&moldQuality)

	if len(moldQuality) > 0 {
		m := moldQuality[0]
		var vo PageQualityOutVO
		base.CopyProperties(&vo, m)

		// 设置用户名称
		var userInfo model.UserInfo
		if err := dao.GetConn().Table("user_info").Where("is_deleted = 'N' and login_name = ?", vo.CreatedBy).First(&userInfo).Error; err == nil {
			vo.CreatedBy = userInfo.Name
		}
		// 模具编号

		var moldCode string
		row := dao.GetConn().Table("mold_info").Select("code").Where("id = ? and is_deleted = 'N'", vo.MoldId).Row()
		if row != nil {
			if err := row.Scan(&moldCode); err == nil {
				vo.MoldCode = moldCode
			}
		}

		var moldPartRel []model.MoldPartRel
		dao.GetConn().Table("mold_part_rel").Where("is_deleted = 'N' and mold_code = ?", vo.MoldCode).Find(&moldPartRel)
		if len(moldPartRel) > 0 {
			var partCodes []string
			for _, rel := range moldPartRel {
				partCodes = append(partCodes, rel.PartCode)
			}
			vo.PartCodes = strings.Join(partCodes, "/")
		}
		c.JSON(http.StatusOK, base.Success(vo))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 模具质量
// @Summary 删除模具质量
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/quality [delete]
func DeleteQuality(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	ids := vo.IDS
	if len(ids) == 0 {
		panic(base.ParamsErrorN())
	}

	dao.GetConn().Table("mold_quality").Where("id in ?", ids).Updates(&model.MoldQuality{
		IsDeleted: "Y",
		UpdatedBy: c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具质量
// @Summary 添加模具质量
// @Accept json
// @Produce json
// @Param Body body CreateQualityVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/quality [post]
func CreateQuality(c *gin.Context) {
	var vo CreateQualityVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 参数校验,TODO

	var moldQuality model.MoldQuality
	base.CopyProperties(&moldQuality, vo)

	moldQuality.Code = common.GenerateCode("Q")

	userId := c.GetHeader(constant.USERID)
	moldQuality.CreatedBy = userId
	moldQuality.UpdatedBy = userId

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("mold_quality").Create(&moldQuality).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(moldQuality.ID))
}

// @Tags 模具质量
// @Summary 更新模具质量
// @Accept json
// @Produce json
// @Param Body body UpdateQualityVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/quality [put]
func UpdateQuality(c *gin.Context) {
	var vo UpdateQualityVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 参数校验,TODO

	if vo.ID == 0 {
		panic(base.ParamsErrorN())
	}

	// 存在
	var one model.MoldQuality
	if err := dao.GetConn().Table("mold_quality").Where("id = ?", vo.ID).First(&one).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 更新
	userId := c.GetHeader(constant.USERID)
	var moldQuality model.MoldQuality
	base.CopyProperties(&moldQuality, vo)
	moldQuality.UpdatedBy = userId

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("mold_quality").Where("id = ?", moldQuality.ID).Updates(&moldQuality).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}
