/**
 * 保养履历
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags 保养履历
// @Summary 保养履历分页
// @Accept json
// @Produce json
// @Param codeOrName query string false "单号/名称"
// @Param moldCode query string true "模具code"
// @Param currentPage query int false "当前页"
// @Param size query int false "每页数量"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageMoldMaintenanceOutVO
// @Router /mold/maintenance/page [get]
func PageMoldMaintenance(c *gin.Context) {
	name := c.Query("codeOrName")
	code := c.Query("moldCode")

	if code == "" {
		panic(base.ParamsErrorN())
	}

	current, size := base.GetPageParams(c)

	sql := `
		select
			mmt.id,
			mmt.code,
			mmt.maintenance_level,
			mmt.standard_name,
			mmt.time,
			COALESCE(uop.name, mmt.operator) AS updated_by
		from mold_info mi
		inner join mold_maintenance_task mmt on mmt.mold_id = mi.id
		left join user_info uop on uop.login_name = mmt.operator
		where mi.code = ? and mi.is_deleted = 'N'
		and mmt.is_deleted = 'N' 
		and mmt.status = 'complete'
		and (mmt.code like concat('%',?,'%') or mmt.standard_name like concat('%',?,'%'))
		order by mmt.time desc
	`

	var result []PageMoldMaintenanceOutVO
	page := base.PageWithRawSQL(dao.GetConn(), &result, current, size, sql, code, name, name)

	if len(result) == 0 {
		page.List = []interface{}{}
	}
	c.JSON(http.StatusOK, base.Success(page))
}
