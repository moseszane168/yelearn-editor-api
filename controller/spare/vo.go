package spare

import (
	"crf-mold/base"
)

//
// 备件基本信息
//

type CreateSpareVO struct {
	Code     string `json:"code" binding:"required"` // 备件编号
	Name     string `json:"name"`                    // 备件名称
	Flavor   string `json:"flavor"`                  // 备件规格
	Material string `json:"material"`                // 材质
} // @name CreateSpareVO

type UpdateSpareVO struct {
	ID       int64  `json:"id" binding:"required"`   // ID
	Code     string `json:"code" binding:"required"` // 备件编号
	Name     string `json:"name"`                    // 备件名称
	Flavor   string `json:"flavor"`                  // 备件规格
	Material string `json:"material"`                // 材质
} // @name UpdateSpareVO

type PageSpareVO struct {
	base.PageVO

	CodeOrName string `form:"codeOrName"` // 备件编号和名称模糊查找，该参数有值时忽略其它条件参数

	Code     string `form:"code"`     // 备件编号
	Name     string `form:"name"`     // 备件名称
	Flavor   string `form:"flavor"`   // 备件规格
	Material string `form:"material"` // 材质
} // @name PageSpareVO

type PageSpareOutVO struct {
	ID         int64     `json:"id"`         // ID
	Code       string    `json:"code"`       // 备件编号
	Name       string    `json:"name"`       // 备件名称
	Flavor     string    `json:"flavor"`     // 备件规格
	Material   string    `json:"material"`   // 材质
	GmtCreated base.Time `json:"gmtCreated"` // 创建时间
	GmtUpdated base.Time `json:"gmtUpdated"` // 修改时间
} // @name PageSpareOutVO

type SpareImportVO struct {
	Code     string `json:"code" excel:"备件编号" validate:"required,unique"` // 备件编号
	Name     string `json:"name" excel:"备件名称"`                            // 备件名称
	Flavor   string `json:"flavor" excel:"备件规格"`                          // 备件规格
	Material string `json:"material" excel:"材质"`                          // 材质
} // @name SpareImportVO

type SpareImportParsedVOS struct {
	Data []SpareImportVO `json:"data" binding:"min=1"` // 导入数据
} // @name SpareImportParsedVOS

type SpareExportVO struct {
	Code     string `excel:"备件编号"` // 备件编号
	Name     string `excel:"备件名称"` // 备件名称
	Flavor   string `excel:"备件规格"` // 备件规格
	Material string `excel:"材质"`   // 材质
} // @name SpareExportVO

//
// 备件库存
//
type InboundSpareRequestVO struct {
	SpareCode string `json:"spareCode" binding:"required"` // 备件编号
	Type      string `json:"type" binding:"required"`      // 入库类型
	Location  string `json:"location" binding:"required"`  // 入库位置
	Count     int    `json:"count" binding:"gte=1"`        // 入库数量
	Remark    string `json:"remark"`                       // 备注
} // @name InboundSpareRequestVO

type OutboundSpareRequestVO struct {
	SpareCode string `json:"spareCode" binding:"required"` // 备件编号
	Type      string `json:"type" binding:"required"`      // 出库类型
	Count     int    `json:"count" binding:"min=1"`        // 出库数量
	Remark    string `json:"remark"`                       // 备注
} // @name OutboundSpareRequestVO

type PageSpareRequestVO struct {
	base.PageVO

	CodeOrName string `form:"codeOrName"` // 备件编号和名称模糊查找，该参数有值时忽略其它条件参数

	SpareCode string `form:"spareCode"` // 备件编号
	Location  string `form:"location"`  // 出入库位置
	Count     int    `form:"count"`     // 出入库数量
} // @name PageSpareRequestVO

type PageSpareRequestOutVO struct {
	SpareCode string `json:"spareCode" form:"spareCode"` // 备件编号
	Name      string `json:"name"`                       // 备件名称
	Flavor    string `json:"flavor"`                     // 备件规格
	Material  string `json:"material"`                   // 材质
	Location  string `json:"location" form:"location"`   // 出入库位置
	Count     int    `json:"count" form:"count"`         // 出入库数量

	InboundTime int64 // 入库时间
} // @name PageSpareRequestOutVO

type SpareRequestImportVO struct {
	SpareCode string `json:"spareCode" excel:"备件编号" validate:"required"`                     // 备件编号
	Type      string `json:"type" excel:"入库类型" validate:"required,enum=storageType"`         // 出入库类型
	Location  string `json:"location" excel:"入库储位" validate:"required,enum=placeForStorage"` // 出入库位置
	Count     int    `json:"count" excel:"入库数量" validate:"required,min=1"`                   // 出入库数量
	Remark    string `json:"remark" excel:"备注"`                                              // 备注
} // @name SpareRequestImportVO

type SpareRequestImportParsedVOS struct {
	Data []SpareRequestImportVO `json:"data" binding:"min=1"` // 导入数据
} // @name SpareRequestImportParsedVOS

type SpareRequestExportVO struct {
	SpareCode string `excel:"备件编号"` // 备件编号
	Name      string `excel:"备件名称"` // 备件名称
	Flavor    string `excel:"备件规格"` // 备件规格
	Material  string `excel:"材质"`   // 材质
	Location  string `excel:"入库储位"` // 入库储位
	Count     int    `excel:"库存数量"` // 库存数量
} // @name SpareRequestExportVO

type SpareRequestRecordOutVO struct {
	SpareCode  string `json:"spareCode"`  // 备件编号
	Location   string `json:"location"`   // 出入库位置
	Type       string `json:"type"`       // 出入库类型
	Count      int    `json:"count"`      // 出入库数量
	GmtCreated string `json:"gmtCreated"` // 出入库时间
	CreatedBy  string `json:"createdBy"`  // 出入库人
	Remark     string `json:"remark"`     // 备注
} // @name SpareRequestRecordOutVO

type SpareRequestRecordVO struct {
	Code            string `form:"code" binding:"required"` // 备件编号
	GmtCreatedBegin string `form:"gmtCreatedBegin"`         // 出入库时间开始
	GmtCreatedEnd   string `form:"gmtCreatedEnd"`           // 出入库时间结束
} // @name SpareRequestRecordVO

type SpareRequestRecordPageVO struct {
	base.PageVO

	Code            string `form:"code" binding:"required"` // 备件编号
	GmtCreatedBegin string `form:"gmtCreatedBegin"`         // 出入库时间开始
	GmtCreatedEnd   string `form:"gmtCreatedEnd"`           // 出入库时间结束
} // @name SpareRequestRecordVO

//
// 备件履历
//
type SpareResumePageQueryVO struct {
	base.PageVO
	SpareCode string   `json:"spareCode"`
	SpareName string   `json:"spareName"`
	MoldCode  string   `json:"moldCode"`
	RfidTid   string   `json:"rfidTid"` // RFID TID
	Type      string   `json:"type"`
	Time      []string `json:"time"`
} // @name SpareResumePageQueryVO

type SpareResumePageVO struct {
	SpareCode  string    `json:"spareCode"`  // 备件编号
	SpareName  string    `json:"spareName"`  // 备件名称
	MoldCode   string    `json:"moldCode"`   // 模具编号
	RelCode    string    `json:"relCode"`    // 对应单号
	CreatedBy  string    `json:"createdBy"`  // 使用人
	GmtCreated base.Time `json:"gmtCreated"` // 使用时间
	Count      int64     `json:"count"`      // 使用数量
	FlushCount int64     `json:"flushCount"` // 备件使用冲次
} // @name SpareResumePageVO

type SpareDetailsExportVO struct {
	SpareCode  string    `json:"spareCode" excel:"备件编号"`    // 备件编号
	SpareName  string    `json:"spareName" excel:"备件名称"`    // 备件名称
	MoldCode   string    `json:"moldCode" excel:"模具编号"`     // 模具编号
	RelCode    string    `json:"relCode" excel:"对应单号"`      // 对应单号
	CreatedBy  string    `json:"createdBy" excel:"使用人"`     // 使用人
	GmtCreated base.Time `json:"gmtCreated" excel:"使用时间"`   // 使用时间
	Count      int64     `json:"count" excel:"使用数量"`        // 使用数量
	FlushCount int64     `json:"flushCount" excel:"备件使用冲次"` // 备件使用冲次
} // @name SpareDetailsExportVO
