/**
 * 字典组
 */

package dict

import (
	"crf-mold/base"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 字典
// @Summary 新增字典分类
// @Schemes
// @Accept json
// @Produce json
// @Param Body body CreateDictGroupVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict/group [post]
func CreateDictGroup(c *gin.Context) {
	var vo CreateDictGroupVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// name唯一
	var dictGroup model.DictGroup
	dictGroup.Name = vo.Name
	dictGroup.Code = vo.Code

	var count int64
	dao.GetConn().Table("dict_group").Where("(name = ? or code = ?)", dictGroup.Name, dictGroup.Code).Where("is_deleted = 'N'").Count(&count)
	if count > 0 {
		panic(base.ResponseEnum[base.RESOURCE_EXIST])
	}

	// 插入
	userId := c.GetHeader(constant.USERID)
	dictGroup.CreatedBy = userId
	dictGroup.UpdatedBy = userId
	dao.GetConn().Table("dict_group").Create(&dictGroup)

	c.JSON(http.StatusOK, base.Success(dictGroup.ID))
}

// @Tags 字典
// @Summary 删除字典分类
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict/group [delete]
func DeleteDictGroup(c *gin.Context) {
	idstr := c.Query("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	var one model.DictGroup
	if err := tx.Table("dict_group").Where("id = ?", id).First(&one).Error; err != nil {
		panic(base.ResponseEnum["600"])
	}

	tx.Table("dict_group").Where("id = ?", id).Update("is_deleted", "Y")

	tx.Table("dict_property").Where("group_code = ?", one.Code).Update("is_deleted", "Y")

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 字典
// @Summary 更新字典分类
// @Accept json
// @Produce json
// @Param Body body UpdateDictGroupVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict/group [put]
func UpdateDictGroup(c *gin.Context) {

	var vo UpdateDictGroupVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 存在
	var one model.DictGroup
	if err := dao.GetConn().Table("dict_group").Where("id = ? and is_deleted = 'N'", vo.ID).First(&one).Error; err != nil {
		panic(base.ResponseEnum["600"])
	}

	// 禁止重复
	var count int64
	dao.GetConn().Table("dict_group").Where("id != ? and (name = ? or code = ?) and is_deleted = 'N'", vo.ID, vo.Name, vo.Code).Count(&count)
	if count > 0 {
		panic(base.ParamsErrorN())
	}

	// 更新
	userId := c.GetHeader(constant.USERID)

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	tx.Table("dict_group").Where("id = ?", vo.ID).Updates(&model.DictGroup{Name: vo.Name, Code: vo.Code, UpdatedBy: userId})
	tx.Table("dict_property").Where("group_code = ? and is_deleted = 'N'", one.Code).Updates(&model.DictProperty{GroupCode: vo.Code, UpdatedBy: userId})

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 字典
// @Summary 字典分类列表
// @Accept json
// @Produce json
// @Param name query string false "name"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict/group [get]
func ListDictGroup(c *gin.Context) {
	name := c.Query("name")

	var results []*model.DictGroup
	tx := dao.GetConn().Table("dict_group").Where("is_deleted = 'N'")
	if name != "" {
		tx = tx.Where("name like concat('%',?,'%')", name).Order("gmt_created desc")
	}
	tx.Find(&results)

	if len(results) > 0 {
		c.JSON(http.StatusOK, base.Success(results))
	} else {
		c.JSON(http.StatusOK, base.Success([]*model.DictGroup{}))
	}

}

// @Tags 字典
// @Summary 获取所有字典信息
// @Accept json
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict/all [get]
func ListAll(c *gin.Context) {
	var dictAllVO []DictAllVO

	dao.GetConn().Raw(`
		select
			dg.code as group_code,
			dp.id,
			dp.key,
			dp.value_cn,
			dp.value_en
		from dict_group dg
		left join dict_property dp on dp.group_code = dg.code and dp.is_deleted = 'N'
		where dg.is_deleted = 'N'
	`).Scan(&dictAllVO)

	result := make(map[string][]DictAllVO)

	for _, v := range dictAllVO {
		item, _ := result[v.GroupCode]
		item = append(item, v)
		result[v.GroupCode] = item
	}

	c.JSON(http.StatusOK, base.Success(result))
}
