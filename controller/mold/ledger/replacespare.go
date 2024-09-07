/**
 * 备件更换履历
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags 备件更换履历
// @Summary 备件更换履历分页
// @Accept json
// @Produce json
// @Param codeOrName query string false "单号/备件编号"
// @Param moldCode query string true "模具code"
// @Param currentPage query int false "当前页"
// @Param size query int false "每页数量"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageMoldMaintenanceOutVO
// @Router /spare/replace/page [get]
func PageReplaceSpare(c *gin.Context) {
	name := c.Query("codeOrName")
	code := c.Query("moldCode")

	if code == "" {
		panic(base.ParamsErrorN())
	}

	current, size := base.GetPageParams(c)

	sql := `
		select
			mrsr.id,
			mrsr.spare_code,
			si.name as spare_name,
			mrsr.gmt_created as time,
			mrsr.type as reason,
			mrsr.count,
			concat(ifnull(mr.code,''),ifnull(mmt.code,'')) as code
		from mold_replace_spare_rel mrsr
		left join spare_info si on si.code = mrsr.spare_code and si.is_deleted = 'N'
		left join mold_repair mr on mr.id = mrsr.repair_id and mr.is_deleted = 'N'
		left join mold_maintenance_task mmt on mmt.id = mrsr.maintenance_task_id and mmt.is_deleted = 'N'
		left join mold_info mi on mi.id = mrsr.mold_id and mi.is_deleted = 'N'
		where mrsr.is_deleted = 'N' and (
			si.code like concat('%',?,'%') or mr.code like concat('%',?,'%') or mmt.code like concat('%',?,'%')
		) and mi.code = ?
		order by mrsr.gmt_created desc
	`

	var result []PageReplaceSpareVO
	page := base.PageWithRawSQL(dao.GetConn(), &result, current, size, sql, name, name, name, code)

	if len(result) == 0 {
		page.List = []interface{}{}
	}

	c.JSON(http.StatusOK, base.Success(page))
}
