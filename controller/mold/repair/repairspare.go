/**
 * 模具自定义字段
 */

package mold

import (
	"crf-mold/common/constant"
	"crf-mold/model"

	"gorm.io/gorm"
)

/**
 * 保存模具维修更换备件字段，支持新增和修改操作
 */
func SaveRepairSpare(tx *gorm.DB, repairSpares []RepairSpare, id int64, moldId int64, loginName string) error {
	// 删除原有数据
	if err := tx.Table("mold_replace_spare_rel").Where("repair_id = ?", id).Updates(&model.MoldReplaceSpareRel{
		IsDeleted: "Y",
		UpdatedBy: loginName,
	}).Error; err != nil {
		return err
	}

	// 新增
	var batchSave []model.MoldReplaceSpareRel
	for _, v := range repairSpares {

		// 查询上一次冲次和当前模具冲次
		var lastFlushCount int64
		var flushCount int64
		var lastFlushCountRecord model.MoldReplaceSpareRel
		if err := tx.Table("mold_replace_spare_rel").Where("mold_id = ? and spare_code = ? and is_deleted = 'N'", moldId, v.Code).Order("gmt_created desc").First(&lastFlushCountRecord).Error; err == nil {
			lastFlushCount = lastFlushCountRecord.FlushCount
		}

		var moldInfo model.MoldInfo
		if err := tx.Table("mold_info").Where("id = ? and is_deleted = 'N'", moldId).First(&moldInfo).Error; err == nil {
			flushCount = moldInfo.FlushCount
		}

		batchSave = append(batchSave, model.MoldReplaceSpareRel{
			SpareCode:      v.Code,
			Count:          v.Count,
			RepairID:       id,
			MoldID:         moldId,
			Type:           constant.REPAIR,
			CreatedBy:      loginName,
			UpdatedBy:      loginName,
			LastFlushCount: lastFlushCount,
			FlushCount:     flushCount,
		})
	}

	if err := tx.Table("mold_replace_spare_rel").CreateInBatches(batchSave, len(batchSave)).Error; err != nil {
		return err
	}

	return nil
}
