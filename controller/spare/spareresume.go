package spare

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/controller/excel"
	"crf-mold/dao"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 备件履历
// @Summary 分页
// @Accept json
// @Produce json
// @Param body body SpareResumePageQueryVO true "SpareResumePageQueryVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} SpareResumePageVO
// @Router /spare/resume/page [post]
func SpareResumePage(c *gin.Context) {
	var vo SpareResumePageQueryVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	sql, v := BuildSql(vo.SpareCode, vo.MoldCode, vo.SpareName, vo.Type, vo.Time)

	var result []SpareResumePageVO
	page := base.PageWithRawSQL(dao.GetConn(), &result, vo.GetCurrentPage(), vo.GetSize(), sql, v...)

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 备件履历
// @Summary 导出
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /spare/resume/excel [get]
func SpareResumeExport(c *gin.Context) {

	sql, v := BuildSql("", "", "", "", []string{})

	var result []SpareDetailsExportVO
	dao.GetConn().Raw(sql, v...).Scan(&result)

	fileName := excel.GenerateExcelFile(excel.SPARE_RESUME)
	info := common.ExportExcel(fileName, "", result)

	c.JSON(http.StatusOK, base.Success(info))
}

func BuildSql(spareCode, moldCode, spareName, typ string, t []string) (string, []interface{}) {
	sql := `
		select
			rel.spare_code,
			mi.code as mold_code,
			si.name as spare_name,
			IFNULL(mmt.code,mr.code) as rel_code,
			rel.created_by,
			rel.gmt_created,
			rel.flush_count as count,
			rel.flush_count - rel.last_flush_count as flush_count
		from mold_replace_spare_rel rel
		left join mold_info mi on rel.mold_id = mi.id and mi.is_deleted = 'N'
		left join spare_info si on si.code = rel.spare_code and si.is_deleted = 'N'
		left join mold_repair mr on mr.id = rel.repair_id and rel.repair_id is not null and mr.is_deleted = 'N'
		left join mold_maintenance_task mmt on mmt.id = rel.maintenance_task_id and mmt.is_deleted = 'N'
		where rel.is_deleted = 'N'
		and si.code like concat('%',?,'%')
		and mi.code like concat('%',?,'%')
		and si.name like concat('%',?,'%')
	`

	v := []interface{}{spareCode, moldCode, spareName}

	if len(t) == 2 {
		sql += "and rel.gmt_created >= ? "
		sql += "and rel.gmt_created <= ? "
		v = append(v, t[0])
		v = append(v, t[1])
	}
	if typ != "" {
		sql += "and rel.type = ? "
		v = append(v, typ)
	}
	sql += "order by rel.gmt_created desc "

	return sql, v
}
