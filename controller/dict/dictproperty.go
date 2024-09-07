/**
 * 字典属性
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

// 校验新增和更新的参数
func validateSaveParams(dictProperty *model.DictProperty) {
	if dictProperty.GroupCode == "" || dictProperty.Key == "" || dictProperty.Order == 0 {
		panic(base.ParamsErrorN())
	}

	if dictProperty.ValueCn == "" && dictProperty.ValueEn == "" {
		panic(base.ParamsErrorN())
	}

	var count int64
	dao.GetConn().Table("dict_group").Where("code = ? and is_deleted = 'N'", dictProperty.GroupCode).Count(&count)
	if count == 0 {
		panic(base.ParamsErrorN())
	}

	if dictProperty.ID == 0 { // 新增
		dao.GetConn().Table("dict_property").Where("`key` = ? and group_code = ? and is_deleted = 'N'", dictProperty.Key, dictProperty.GroupCode).Count(&count)
	} else { // 更新
		dao.GetConn().Table("dict_property").Where("`key` = ? and group_code = ? and is_deleted = 'N' and id != ?", dictProperty.Key, dictProperty.GroupCode, dictProperty.ID).Count(&count)
	}

	if count > 0 {
		panic(base.ParamsErrorN())
	}
}

// @Tags 字典
// @Summary 新增字典属性
// @Accept json
// @Produce json
// @Param Body body CreateDictPropertyVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict [post]
func CreateDictProperty(c *gin.Context) {
	var vo CreateDictPropertyVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var dictProperty model.DictProperty
	base.CopyProperties(&dictProperty, vo)

	validateSaveParams(&dictProperty)

	// 值不重复
	if vo.ValueCn != "" {
		var count int64
		dao.GetConn().Table("dict_property").Where("group_code = ? and value_cn = ? and is_deleted = 'N'", vo.GroupCode, vo.ValueCn).Count(&count)
		if count > 0 {
			panic(base.ResponseEnum[base.DICT_VALUE_EXIST])
		}
	}

	if vo.ValueEn != "" {
		var count int64
		dao.GetConn().Table("dict_property").Where("group_code = ? and value_en = ? and is_deleted = 'N'", vo.GroupCode, vo.ValueEn).Count(&count)
		if count > 0 {
			panic(base.ResponseEnum[base.DICT_VALUE_EXIST])
		}
	}

	userId := c.GetHeader(constant.USERID)
	dictProperty.CreatedBy = userId
	dictProperty.UpdatedBy = userId

	// 新增
	dao.GetConn().Table("dict_property").Create(&dictProperty)

	c.JSON(http.StatusOK, base.Success(dictProperty.ID))
}

// @Tags 字典
// @Summary 删除字典属性
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict [delete]
func DeleteDictProperty(c *gin.Context) {
	idstr := c.Query("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	dao.GetConn().Table("dict_property").Where("id = ?", id).Update("is_deleted", "Y")

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 字典
// @Summary 更新字典属性
// @Accept json
// @Produce json
// @Param Body body UpdateDictPropertyVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict [put]
func UpdateDictProperty(c *gin.Context) {
	var vo UpdateDictPropertyVO
	// 传参有误
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	var dictProperty model.DictProperty
	base.CopyProperties(&dictProperty, vo)
	validateSaveParams(&dictProperty)

	// 值不重复
	if vo.ValueCn != "" {
		var count int64
		dao.GetConn().Table("dict_property").Where("group_code = ? and value_cn = ? and id != ? and is_deleted = 'N'", vo.GroupCode, vo.ValueCn, vo.ID).Count(&count)
		if count > 0 {
			panic(base.ResponseEnum[base.DICT_VALUE_EXIST])
		}
	}

	if vo.ValueEn != "" {
		var count int64
		dao.GetConn().Table("dict_property").Where("group_code = ? and value_en = ? and id != ? and is_deleted = 'N'", vo.GroupCode, vo.ValueEn, vo.ID).Count(&count)
		if count > 0 {
			panic(base.ResponseEnum[base.DICT_VALUE_EXIST])
		}
	}

	dao.GetConn().Table("dict_property").Where("id = ?", vo.ID).Updates(&model.DictProperty{Key: dictProperty.Key, ValueCn: dictProperty.ValueCn, ValueEn: dictProperty.ValueEn, Order: dictProperty.Order})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 字典
// @Summary 字典属性列表
// @Accept json
// @Produce json
// @Param code query string true "字典组code"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict [get]
func ListDictProperty(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		panic(base.ParamsErrorN())
	}

	var results []*model.DictProperty
	dao.GetConn().Table("dict_property").Where("is_deleted = 'N' and group_code = ?", code).Order("`order` asc").Find(&results)

	c.JSON(http.StatusOK, base.Success(results))
}

// @Tags 字典
// @Summary 字典属性分页
// @Accept json
// @Produce json
// @Param code query string true "字典组code"
// @Param currentPage query int false "当前页"
// @Param size query int false "每页数量"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /dict/page [get]
func PageDictProperty(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		panic(base.ParamsErrorN())
	}

	current, size := base.GetPageParams(c)
	var results []model.DictProperty
	tx := dao.GetConn().Table("dict_property").Where("is_deleted = 'N' and group_code = ?", code).Order("`order` asc")
	page := base.Page(tx, &results, current, size)

	if len(results) == 0 {
		page.List = []interface{}{}
	}

	c.JSON(http.StatusOK, base.Success(page))
}

/**
 * 获取指定字典组下的指定value对应的key
 */
func GetKeyByValue(groupCode, value, language string) string {
	tx := dao.GetConn().Table("dict_property")
	tx.Where("group_code = ?", groupCode)
	if language == "en-US" {
		tx.Where("value_en = ?", value)
	} else {
		tx.Where("value_cn = ?", value)
	}

	var r model.DictProperty
	if err := tx.First(&r).Error; err != nil {
		return groupCode
	}

	return r.Key
}
