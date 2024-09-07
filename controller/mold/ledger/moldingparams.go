/**
 * 模具成型参数
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 模具成型参数
// @Summary 保存成型参数
// @Accept json
// @Produce json
// @Param Body body SaveMoldingParamsVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/molding [put]
func SaveMoldingParams(c *gin.Context) {
	var vo SaveMoldingParamsVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 存在
	var one model.MoldParams
	if err := dao.GetConn().Table("mold_params").Where("mold_code = ? and is_deleted = 'N'", vo.MoldCode).First(&one).Error; err != nil {
		// 新增
		var entity model.MoldParams
		base.CopyProperties(&entity, vo)
		entity.CreatedBy = c.GetHeader(constant.USERID)
		entity.UpdatedBy = c.GetHeader(constant.USERID)

		dao.GetConn().Table("mold_params").Create(&entity)
	} else { // 不存在
		// 更新
		var entity model.MoldParams
		base.CopyProperties(&entity, vo)
		entity.UpdatedBy = c.GetHeader(constant.USERID)

		dao.GetConn().Table("mold_params").Where("mold_code = ? and is_deleted = 'N'", one.MoldCode).Updates(&entity)
	}

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具成型参数
// @Summary 用户成型参数查看
// @Accept json
// @Produce json
// @Param code query string true "模具编码"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageMoldingParamsOutVO
// @Router /mold/molding/one [get]
func OneMoldingParams(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		panic(base.ParamsErrorN())
	}

	var moldParams model.MoldParams
	if err := dao.GetConn().Table("mold_params").Where("is_deleted = 'N' and mold_code = ?", code).First(&moldParams).Error; err != nil {
		c.JSON(http.StatusOK, base.SuccessN())
		return
	}

	var moldParamsCustoms []model.MoldParamsCustom
	dao.GetConn().Table("mold_params_custom").Where("is_deleted = 'N' and mold_code = ?", code).Find(&moldParamsCustoms)

	result := PageMoldingParamsOutVO{}
	base.CopyProperties(&result, moldParams)

	customs := []MoldingCustomVO{}
	for _, v := range moldParamsCustoms {
		var custom MoldingCustomVO
		base.CopyProperties(&custom, v)
		customs = append(customs, custom)
	}

	if len(customs) == 0 {
		result.Customs = nil
	} else {
		result.Customs = customs
	}

	c.JSON(http.StatusOK, base.Success(result))
}

// @Tags 模具成型参数
// @Summary 保存成型参数自定义字段
// @Accept json
// @Produce json
// @Param Body body CreateMoldingParamsCustomsVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/molding/customs [put]
func SaveMoldingParamsCustoms(c *gin.Context) {
	var vo CreateMoldingParamsCustomsVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	customs := vo.Customs
	moldCode := vo.MoldCode
	// key不能重复
	keyMap := make(map[string]bool)
	for _, v := range customs {
		_, ok := keyMap[v.Key]
		if ok {
			panic(base.ResponseEnum[base.MOLD_MOLDING_CUSTOM_KEY_REPEAT])
		}
		keyMap[v.Key] = true
	}

	userId := c.GetHeader(constant.USERID)

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	// 先删除所有
	tx.Table("mold_params_custom").Where("is_deleted = 'N' and mold_code = ?", vo.MoldCode).Updates(&model.MoldParamsCustom{
		IsDeleted: "Y",
		UpdatedBy: userId,
	})

	// 新增
	var batchSave = make([]model.MoldParamsCustom, len(customs))
	for i := 0; i < len(customs); i++ {
		custom := customs[i]
		batchSave[i] = model.MoldParamsCustom{
			MoldCode:  moldCode,
			Key:       custom.Key,
			Value:     custom.Value,
			CreatedBy: userId,
			UpdatedBy: userId,
		}

		// var one model.MoldParamsCustom
		// if err := tx.Table("mold_params_custom").Where("is_deleted = 'N' and mold_code = ? and `key` = ?", moldCode, custom.Key).First(&one).Error; err != nil {
		// 	// key 不能重复
		// 	var count int64
		// 	tx.Table("mold_params_custom").Where("is_deleted = 'N' and mold_code = ? and `key` = ?", moldCode, custom.Key).Count(&count)
		// 	if count > 0 {
		// 		panic(base.ResponseEnum[base.MOLD_MOLDING_CUSTOM_KEY_REPEAT])
		// 	}

		// 	// 不存在，新增
		// 	tx.Table("mold_params_custom").Where("is_deleted = 'N' and mold_code = ?", moldCode).Create(&model.MoldParamsCustom{
		// 		MoldCode: moldCode,
		// 		Key:      custom.Key,
		// 		Value:    custom.Value,
		// 	})
		// } else { // 存在，更新
		// 	// key 不能重复
		// 	var count int64
		// 	tx.Table("mold_params_custom").Where("is_deleted = 'N' and mold_code = ? and id != ? and `key` = ?", moldCode, one.ID, custom.Key).Count(&count)
		// 	if count > 0 {
		// 		panic(base.ResponseEnum[base.MOLD_MOLDING_CUSTOM_KEY_REPEAT])
		// 	}

		// 	tx.Table("mold_params_custom").Where("is_deleted = 'N' and mold_code = ?", moldCode).Updates(&model.MoldParamsCustom{
		// 		Key:   custom.Key,
		// 		Value: custom.Value,
		// 	})
		// }
	}

	tx.Table("mold_params_custom").Where("is_deleted = 'N' and mold_code = ?", moldCode).CreateInBatches(batchSave, len(batchSave))

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具成型参数
// @Summary 删除模具成型参数
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/molding/customs [delete]
func DeleteMoldingCustoms(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	ids := vo.IDS
	if len(ids) == 0 {
		panic(base.ParamsErrorN())
	}

	dao.GetConn().Table("mold_params_custom").Where("id in ?", ids).Updates(&model.MoldInfo{
		IsDeleted: "Y",
		UpdatedBy: c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}
