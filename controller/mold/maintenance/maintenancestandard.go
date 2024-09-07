/**
 * 模具保养标准
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

// @Tags 模具保养标准
// @Summary 新增保养标准
// @Accept json
// @Produce json
// @Param Body body CreateMoldMaintenanceStandardVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} int64
// @Router /maintenance/standard [post]
func CreateMoldMaintenanceStandard(c *gin.Context) {
	var vo CreateMoldMaintenanceStandardVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	// 名称唯一
	var count int64
	tx.Table("mold_maintenance_standard").Where("name = ? and is_deleted = 'N'", vo.Name).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MAINTENANCE_NAME_EXIST])
	}

	// 通用标准已存在
	if vo.StandardType == constant.GENERAL {
		tx.Table("mold_maintenance_standard").Where("is_deleted = 'N' and standard_type = ? and level = ? and type = ?", constant.GENERAL, vo.Level, vo.Type).Count(&count)
		if count > 0 {
			panic(base.ResponseEnum[base.MAINTENANCE_STANDARD_GENERAL_EXIST])
		}
	}

	var moldMaintenanceStandard model.MoldMaintenanceStandard
	base.CopyProperties(&moldMaintenanceStandard, vo)

	userId := c.GetHeader(constant.USERID)
	moldMaintenanceStandard.CreatedBy = userId
	moldMaintenanceStandard.UpdatedBy = userId
	moldMaintenanceStandard.GmtCreated = vo.Time

	// 新增
	tx.Table("mold_maintenance_standard").Create(&moldMaintenanceStandard)

	// 标准内容
	if len(vo.Content) > 0 {
		SaveMaintenanceStandardContent(tx, moldMaintenanceStandard.ID, userId, vo.Content)
	} else {
		panic(base.ResponseEnum[base.MAINTENANCE_STANDARD_CONTENT_CAN_NOT_EMPTY])
	}

	// 关联模具，只有专用标准有
	if len(vo.MoldIds) > 0 {
		if vo.StandardType == constant.SPECIAL {
			SaveMaintenanceStandardRel(tx, moldMaintenanceStandard.ID, userId, vo.MoldIds)
		}
	} else {
		if vo.StandardType == constant.SPECIAL {
			panic(base.ResponseEnum[base.MAINTENANCE_STANDARD_REL_MOLD_CAN_NOT_EMPTY])
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, base.Success(moldMaintenanceStandard.ID))
}

// @Tags 模具保养标准
// @Summary 删除保养标准
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/standard [delete]
func DeleteMoldMaintenanceStandard(c *gin.Context) {
	idstr := c.Query("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	dao.GetConn().Table("mold_maintenance_standard").Where("id = ?", id).Update("is_deleted", "Y")

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具保养标准
// @Summary 更新保养标准
// @Accept json
// @Produce json
// @Param Body body UpdateMoldMaintenanceStandardVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} bool
// @Router /maintenance/standard [put]
func UpdateMoldMaintenanceStandard(c *gin.Context) {
	var vo UpdateMoldMaintenanceStandardVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	userId := c.GetHeader(constant.USERID)

	// 名称唯一
	var count int64
	tx.Table("mold_maintenance_standard").Where("name = ? and id != ? and is_deleted = 'N'", vo.Name, vo.ID).Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.MAINTENANCE_NAME_EXIST])
	}

	// 通用标准已存在
	if vo.StandardType == constant.GENERAL {
		tx.Table("mold_maintenance_standard").Where("is_deleted = 'N' and standard_type = ? and level = ? and type = ? and id != ?", constant.GENERAL, vo.Level, vo.Type, vo.ID).Count(&count)
		if count > 0 {
			panic(base.ResponseEnum[base.MAINTENANCE_STANDARD_GENERAL_EXIST])
		}
	}

	var moldMaintenanceStandard model.MoldMaintenanceStandard
	base.CopyProperties(&moldMaintenanceStandard, vo)
	moldMaintenanceStandard.UpdatedBy = userId

	tx.Table("mold_maintenance_standard").Updates(&moldMaintenanceStandard)

	// 标准内容
	if len(vo.Content) > 0 {
		SaveMaintenanceStandardContent(tx, moldMaintenanceStandard.ID, userId, vo.Content)
	} else {
		panic(base.ResponseEnum[base.MAINTENANCE_STANDARD_CONTENT_CAN_NOT_EMPTY])
	}

	// 关联模具，只有专用标准
	if len(vo.MoldIds) > 0 {
		if vo.StandardType == constant.SPECIAL {
			SaveMaintenanceStandardRel(tx, moldMaintenanceStandard.ID, userId, vo.MoldIds)
		}
	} else {
		if vo.StandardType == constant.SPECIAL {
			panic(base.ResponseEnum[base.MAINTENANCE_STANDARD_REL_MOLD_CAN_NOT_EMPTY])
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具保养标准
// @Summary 保养标准列表
// @Accept json
// @Produce json
// @Param name query string false "保养标准名称"
// @Param AuthToken header string false "Token"
// @Success 200 {object} model.MoldMaintenanceStandard
// @Router /maintenance/standard/list [get]
func ListMoldMaintenanceStandard(c *gin.Context) {
	name := c.Query("name")

	var results []model.MoldMaintenanceStandard
	tx := dao.GetConn().Table("mold_maintenance_standard").Where("is_deleted = 'N'")
	if name != "" {
		tx.Where("name like concat('%',?,'%')", name)
	}

	tx.Order("`gmt_created` desc").Find(&results)

	if len(results) > 0 {
		c.JSON(http.StatusOK, base.Success(results))
	} else {
		c.JSON(http.StatusOK, base.Success([]model.MoldMaintenanceStandard{}))
	}
}

// @Tags 模具保养标准
// @Summary 保养标准分页
// @Accept json
// @Produce json
// @Param query query PageMaintenanceStandardInputVO true "PageMaintenanceStandardInputVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} PageMaintenanceStandardOutputVO
// @Router /maintenance/standard/page [get]
func PageMoldMaintenanceStandard(c *gin.Context) {
	var vo PageMaintenanceStandardInputVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	current, size := base.GetPageParams(c)
	var results []model.MoldMaintenanceStandard
	tx := dao.GetConn().Table("mold_maintenance_standard").Where("is_deleted = 'N'").Order("`gmt_created` desc")
	if vo.CodeOrName != "" {
		tx.Where("(code like concat('%',?,'%') or name like concat('%',?,'%'))", vo.CodeOrName, vo.CodeOrName)
	} else {
		m := model.MoldMaintenanceStandard{}
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
		vos := make([]PageMaintenanceStandardOutputVO, len(results))
		for i := 0; i < len(results); i++ {
			var item PageMaintenanceStandardOutputVO
			base.CopyProperties(&item, results[i])

			createdBy := item.CreatedBy
			if createdBy != "" {
				var ui model.UserInfo
				if err := dao.GetConn().Table("user_info").Where("login_name = ?", createdBy).First(&ui).Error; err == nil {
					item.CreatedBy = ui.Name
				}
			}

			vos[i] = item
		}
		page.List = vos
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 模具保养标准
// @Summary 保养标准查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} MaintenanceStandardContentVO
// @Router /maintenance/standard/one [get]
func OneMoldMaintenanceStandard(c *gin.Context) {
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
			mi.line_level,
			mi.process,
			GROUP_CONCAT(DISTINCT mpr.part_code ORDER BY mpr.part_code SEPARATOR '/') as part_codes,
			mms.level
		from mold_info mi
		left join mold_part_rel mpr on mpr.mold_code = mi.code and mpr.is_deleted = 'N'
		left join mold_maintenance_standard_rel mmsr on mmsr.mold_id = mi.id and mmsr.is_deleted = 'N'
		left join mold_maintenance_standard mms on mms.id = mmsr.mold_maintenance_standard_id and mms.is_deleted = 'N'
		where mi.is_deleted = 'N' and mms.id = ?
		group by mi.id,mi.code,mi.line_level,mi.process,mms.level
	`

	var out MoldMaintenanceStandardOutVO
	var moldRel []MoldRelPageOutVO

	tx := dao.GetConn()

	var maintenanceStandard model.MoldMaintenanceStandard
	if err := tx.Table("mold_maintenance_standard").Where("id = ? and is_deleted = 'N'", id).First(&maintenanceStandard).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	base.CopyProperties(&out, maintenanceStandard)

	// 关联模具
	tx.Raw(sql, id).Scan(&moldRel)
	if len(moldRel) > 0 {
		out.Molds = moldRel
	} else {
		out.Molds = []MoldRelPageOutVO{}
	}

	// 标准详情
	var contents []MaintenanceStandardContentVO
	tx.Table("mold_maintenance_standard_content").Where("mold_maintenance_standard_id = ? and is_deleted = 'N'", id).Find(&contents)
	if len(contents) > 0 {
		out.Content = contents
	} else {
		out.Content = []MaintenanceStandardContentVO{}
	}

	c.JSON(http.StatusOK, base.Success(out))
}

// @Tags 模具保养标准
// @Summary 关联模具列表查询
// @Accept json
// @Produce json
// @Param codeOrName query string false "模具编号/零件号"
// @Param currentPage query int false "当前页"
// @Param size query int false "每页数量"
// @Param AuthToken header string false "Token"
// @Success 200 {object} MoldRelPageOutVO
// @Router /maintenance/standard/mold/page [get]
func MoldRelPage(c *gin.Context) {
	codeOrName := c.Query("codeOrName")
	current, size := base.GetPageParams(c)

	sql := `
		select
			mi.id,
			mi.name,
			mi.code,
			mi.line_level,
			mi.process,
			GROUP_CONCAT(DISTINCT mpr.part_code ORDER BY mpr.part_code SEPARATOR '/') as part_codes,
			mms.level
		from mold_info mi
		left join mold_part_rel mpr on mpr.mold_code = mi.code and mpr.is_deleted = 'N'
		left join mold_maintenance_standard_rel mmsr on mmsr.mold_id = mi.id and mmsr.is_deleted = 'N'
		left join mold_maintenance_standard mms on mms.id = mmsr.mold_maintenance_standard_id and mms.is_deleted = 'N'
		where mi.is_deleted = 'N' and (mi.code like concat('%',?,'%') or mpr.part_code like concat('%',?,'%'))
		group by mi.id,mi.code,mi.line_level,mi.process,mms.level
	`

	var result []MoldRelPageOutVO
	page := base.PageWithRawSQL(dao.GetConn(), &result, current, size, sql, codeOrName, codeOrName)

	c.JSON(http.StatusOK, base.Success(page))
}

func SaveMaintenanceStandardContent(tx *gorm.DB, id int64, userId string, vos []MaintenanceStandardContentVO) {
	// 删除
	tx.Table("mold_maintenance_standard_content").Where("mold_maintenance_standard_id = ?", id).Updates(&model.MoldMaintenanceStandardRel{
		IsDeleted: "Y",
		UpdatedBy: userId,
	})

	// 新增
	saveBatch := make([]model.MoldMaintenanceStandardContent, len(vos))
	for i := 0; i < len(vos); i++ {
		item := vos[i]
		var standardContent model.MoldMaintenanceStandardContent
		base.CopyProperties(&standardContent, item)
		standardContent.MoldMaintenanceStandardID = id
		standardContent.Order = i
		standardContent.CreatedBy = userId
		standardContent.UpdatedBy = userId
		saveBatch[i] = standardContent
	}

	tx.Table("mold_maintenance_standard_content").CreateInBatches(saveBatch, len(saveBatch))
}

func SaveMaintenanceStandardRel(tx *gorm.DB, id int64, userId string, vos []int64) {
	// 删除
	tx.Table("mold_maintenance_standard_rel").Where("mold_maintenance_standard_id = ?", id).Updates(&model.MoldMaintenanceStandardRel{
		IsDeleted: "Y",
		UpdatedBy: userId,
	})

	// 新增
	saveBatch := make([]model.MoldMaintenanceStandardRel, len(vos))
	for i := 0; i < len(vos); i++ {
		var item model.MoldMaintenanceStandardRel
		item.MoldId = vos[i]
		item.MoldMaintenanceStandardID = id
		item.CreatedBy = userId
		item.UpdatedBy = userId
		saveBatch[i] = item
	}

	tx.Table("mold_maintenance_standard_rel").CreateInBatches(saveBatch, len(saveBatch))
}

// @Tags 模具保养标准
// @Summary 获取模具保养任务可选择的标准列表
// @Accept json
// @Produce json
// @Param moldCode query string true "模具编码"
// @Param level query string false "保养级别"
// @Param AuthToken header string false "Token"
// @Success 200 {object} MoldStandardSelectVO
// @Router /maintenance/standard/select [get]
func MoldStandardSelectList(c *gin.Context) {
	moldCode := c.Query("moldCode")
	level := c.Query("level")

	var moldInfo model.MoldInfo
	if err := dao.GetConn().Table("mold_info").Where("code = ? and is_deleted = 'N'", moldCode).First(&moldInfo).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	var result []MoldStandardSelectVO

	// 通用标准
	var standards []model.MoldMaintenanceStandard
	dao.GetConn().Table("mold_maintenance_standard").Where("is_deleted = 'N' and type = ? and standard_type = 'general'", moldInfo.Type).
		Find(&standards)
	if len(standards) > 0 {
		// 加入结果集
		for _, standard := range standards {
			if level == "" || standard.Level == level {
				result = append(result, MoldStandardSelectVO{
					ID:               standard.ID,
					MaintenanceName:  standard.Name,
					MaintenanceLevel: standard.Level,
				})
			}
		}
	}

	// 专用标准
	var specialStandardRels []model.MoldMaintenanceStandardRel
	dao.GetConn().Table("mold_maintenance_standard_rel").Where("mold_id = ? and is_deleted = 'N'", moldInfo.ID).Find(&specialStandardRels)
	if len(specialStandardRels) > 0 {
		var standardIds []int64
		for _, v := range specialStandardRels {
			standardIds = append(standardIds, v.MoldMaintenanceStandardID)
		}

		var specialStandards []model.MoldMaintenanceStandard
		dao.GetConn().Table("mold_maintenance_standard").Where("id in ? and is_deleted = 'N' and type = ?", standardIds, moldInfo.Type).Find(&specialStandards)
		if len(specialStandards) > 0 {
			// 加入结果集
			for _, standard := range specialStandards {
				if level == "" || standard.Level == level {
					result = append(result, MoldStandardSelectVO{
						ID:               standard.ID,
						MaintenanceName:  standard.Name,
						MaintenanceLevel: standard.Level,
					})
				}
			}
		}
	}

	if len(result) > 0 {
		c.JSON(http.StatusOK, base.Success(result))
	} else {
		c.JSON(http.StatusOK, base.Success([]interface{}{}))
	}
}
