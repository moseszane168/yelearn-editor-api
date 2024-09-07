/**
 * 模具文档
 */
package mold

import (
	"crf-mold/base"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 模具文档
// @Summary 模具文档分页
// @Accept json
// @Produce json
// @Param codeOrName query string false "模具文档名称"
// @Param moldCode query string true "模具code"
// @Param currentPage query int false "当前页"
// @Param size query int false "每页数量"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/doc/page [get]
func PageMoldDoc(c *gin.Context) {
	name := c.Query("codeOrName")
	code := c.Query("moldCode")

	if code == "" {
		panic(base.ParamsErrorN())
	}

	current, size := base.GetPageParams(c)
	var results []model.MoldDoc
	tx := dao.GetConn().Table("mold_doc")
	tx.Joins("LEFT JOIN user_info ON user_info.login_name = mold_doc.updated_by")
	tx.Where("mold_doc.is_deleted = 'N' and mold_doc.mold_code = ?", code).Order("mold_doc.gmt_created desc")
	if name != "" {
		tx.Where("(mold_doc.version like concat('%',?,'%') or mold_doc.content like concat('%',?,'%'))", name, name)
	}
	tx.Select("mold_doc.id, mold_doc.file_key, mold_doc.mold_code, mold_doc.name, mold_doc.version, mold_doc.content, " +
		"mold_doc.remark, mold_doc.is_deleted, mold_doc.created_by, COALESCE(user_info.name, mold_doc.updated_by) AS updated_by, mold_doc.gmt_created, mold_doc.gmt_updated")

	page := base.Page(tx, &results, current, size)

	if len(results) == 0 {
		page.List = []interface{}{}
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 模具文档
// @Summary 保存模具文档
// @Accept json
// @Produce json
// @Param Body body SaveMoldDocVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/doc [post]
func SaveMoldDoc(c *gin.Context) {
	var vo SaveMoldDocVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 名称+版本存在则修改，不存在则新增
	var moldDoc model.MoldDoc
	if err := dao.GetConn().Table("mold_doc").Where("name = ? and version = ? and is_deleted = 'N'", vo.Name, vo.Version).First(&moldDoc).Error; err != nil {
		// 不存在，新增
		dao.GetConn().Table("mold_doc").Create(&model.MoldDoc{
			FileKey:   vo.FileKey,
			MoldCode:  vo.MoldCode,
			Name:      vo.Name,
			Version:   vo.Version,
			Content:   vo.Content,
			Remark:    vo.Remark,
			CreatedBy: c.GetHeader(constant.USERID),
			UpdatedBy: c.GetHeader(constant.USERID),
		})
	} else {
		// 存在，修改
		// 是否重复
		var count int64
		dao.GetConn().Table("mold_doc").Where("id != ? and name = ? and version = ? and is_deleted = 'N'", moldDoc.ID, vo.Name, vo.Version).Count(&count)
		if count > 0 {
			panic(base.ResponseEnum[base.MOLD_DOC_EXIST])
		}

		// 修改
		dao.GetConn().Table("mold_doc").Where("id = ?", moldDoc.ID).Updates(&model.MoldDoc{
			FileKey:   vo.FileKey,
			MoldCode:  vo.MoldCode,
			Name:      vo.Name,
			Version:   vo.Version,
			Content:   vo.Content,
			Remark:    vo.Remark,
			UpdatedBy: c.GetHeader(constant.USERID),
		})
	}

	c.JSON(http.StatusOK, base.Success(true))
}
