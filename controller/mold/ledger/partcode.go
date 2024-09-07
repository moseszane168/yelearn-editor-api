/**
 * 模具台账零件号
 */

package mold

import (
	"crf-mold/dao"
	"crf-mold/model"

	"gorm.io/gorm"
)

func SavePartCodes(tx *gorm.DB, partCodes []string, oldCode, newCode string, loginName string) error {

	// 删除原有数据
	if err := tx.Table("mold_part_rel").Where("mold_code = ?", oldCode).Updates(&model.MoldPartRel{
		IsDeleted: "Y",
		UpdatedBy: loginName,
	}).Error; err != nil {
		return err
	}

	// 新增
	var batchSave []model.MoldPartRel
	for _, v := range partCodes {
		batchSave = append(batchSave, model.MoldPartRel{
			MoldCode:  newCode,
			PartCode:  v,
			CreatedBy: loginName,
			UpdatedBy: loginName,
		})
	}

	if err := tx.Table("mold_part_rel").CreateInBatches(batchSave, len(batchSave)).Error; err != nil {
		return err
	}

	return nil
}

func GetMoldPartCodes(code string) []string {
	var res []model.MoldPartRel
	dao.GetConn().Table("mold_part_rel").Where("mold_code = ? and is_deleted = 'N'", code).Find(&res)

	var result []string
	for _, v := range res {
		result = append(result, v.PartCode)
	}

	return result
}
