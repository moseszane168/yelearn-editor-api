package productresume

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/controller/excel"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags 生产履历
// @Summary 生产履历分页
// @Accept json
// @Produce json
// @Param query query PageProductResumeInputVO true "PageProductResumeInputVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} PageProductResumeOutputVO
// @Router /mold/productresume/page [get]
func PageProductResume(c *gin.Context) {
	var vo PageProductResumeInputVO
	var sumCount int64
	var moldInfo model.MoldInfo
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	var moldProductResume model.MoldProductResume
	base.CopyProperties(&moldProductResume, vo)

	if vo.RfidTid != "" {
		dao.GetConn().Table("mold_info").Select("code").Where("rfid = ? and is_deleted = 'N'", vo.RfidTid).First(&moldInfo)
	}
	if moldInfo.Code != "" && moldInfo.Code != vo.CodeOrName {
		vo.CodeOrName = moldInfo.Code
	}

	tx := dao.GetConn().Table("mold_product_resume")

	tx.Where("(mold_code like concat('%',?,'%') or order_code like concat('%',?,'%') or part_code like concat('%',?,'%'))", vo.CodeOrName, vo.CodeOrName, vo.CodeOrName)
	if vo.LineLevel != "" {
		tx.Where("lineLevel = ?", vo.LineLevel)
	}
	if vo.LineLevel != "" {
		tx.Where("lineLevel = ?", vo.LineLevel)
	}
	if vo.StartTime != "" {
		tx.Where("complete_time >= ?", vo.StartTime)
	}
	if vo.EndTime != "" {
		tx.Where("complete_time <= ?", vo.EndTime)
	}

	tx = tx.Where("is_deleted = 'N'").Order("gmt_created desc")
	if vo.CodeOrName != "" {
		tx.Select("sum(`count`)").Row().Scan(&sumCount)
	}
	tx.Select("*")

	var result []model.MoldProductResume
	page := base.Page(tx, &result, vo.GetCurrentPage(), vo.GetSize())

	// 转换为outvo
	if len(result) > 0 {
		vos := make([]PageProductResumeOutputVO, len(result))
		for i := 0; i < len(result); i++ {
			item := result[i]
			var vo PageProductResumeOutputVO
			base.CopyProperties(&vo, item)
			vo.SumCount = sumCount
			vos[i] = vo
		}
		page.List = vos
	} else {
		page.List = []interface{}{}
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 生产履历
// @Summary 导出
// @Accept json
// @Produce json
// @Param ids query array false "指定ID导出备件"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/productresume/excel [get]
func Export(c *gin.Context) {
	idsStr := c.QueryArray("ids")

	var res []model.MoldProductResume
	tx := dao.GetConn().Table("mold_product_resume")

	l := len(idsStr)
	if l == 0 {
		tx.Find(&res)
	} else {
		var ids []int
		for i := 0; i < l; i++ {
			id, err := strconv.Atoi(idsStr[i])
			if err != nil {
				panic(base.ParamsErrorN())
			}

			ids = append(ids, id)
		}

		tx.Where("id in ?", ids).Find(&res)
	}

	vos := make([]ProductResumeExportVO, len(res))
	for i := 0; i < len(res); i++ {
		var v ProductResumeExportVO
		base.CopyProperties(&v, res[i])
		vos[i] = v
	}

	fileName := excel.GenerateExcelFile(excel.PRODUCT_RESUME)
	info := common.ExportExcel(fileName, "", vos)

	c.JSON(http.StatusOK, base.Success(info))
}
