/**
 * 模具自定义字段
 */

package mold

import (
	"crf-mold/model"

	"gorm.io/gorm"
)

/**
 * 保存台账自定义字段，支持新增和修改操作
 */
func SaveCustomInfos(tx *gorm.DB, vos []MoldCustomVO, oldCode, newCode string, loginName string) error {
	// 删除原有数据
	if err := tx.Table("mold_custom_info").Where("mold_code = ?", oldCode).Updates(&model.MoldCustomInfo{
		IsDeleted: "Y",
		UpdatedBy: loginName,
	}).Error; err != nil {
		return err
	}

	// 新增
	var batchSave []model.MoldCustomInfo
	for _, v := range vos {
		batchSave = append(batchSave, model.MoldCustomInfo{
			Key:       v.Key,
			Value:     v.Value,
			MoldCode:  newCode,
			CreatedBy: loginName,
			UpdatedBy: loginName,
		})
	}

	if err := tx.Table("mold_custom_info").CreateInBatches(batchSave, len(batchSave)).Error; err != nil {
		return err
	}

	return nil
}
