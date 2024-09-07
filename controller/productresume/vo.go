package productresume

import (
	"crf-mold/base"
	"crf-mold/model"
)

type PageProductResumeInputVO struct {
	base.PageVO

	CodeOrName string `form:"codeOrName"` // 模具编号/零件号/工单号模糊查找
	LineLevel  string `form:"lineLevel"`  // 线别
	RfidTid    string `form:"rfidTid"`    // RFID TID
	StartTime  string `form:"startTime"`  // 开始时间
	EndTime    string `form:"endTime"`    // 结束时间
} // @name PageProductResumeInputVO

type PageProductResumeOutputVO struct {
	ID           int64     `json:"id"`
	MoldCode     string    `json:"moldCode"`     // 模具编码
	MoldName     string    `json:"moldName"`     // 模具名称
	MoldType     string    `json:"moldType"`     // 模具类型
	OrderCode    string    `json:"orderCode"`    // 工单号
	PartCode     string    `json:"partCode"`     // 零件号
	LineLevel    string    `json:"lineLevel"`    // 线别
	SumCount     int64     `json:"sumCount"`     // 总计数
	CompleteTime base.Time `json:"completeTime"` // 生产完成时间
	Count        int64     `json:"count"`        // 工单数量
} // @name PageProductResumeOutputVO

type ProductResumeExportVO struct {
	MoldCode     string    `excel:"模具编码"`           // 模具编码
	MoldName     string    `excel:"模具名称"`           // 模具名称
	MoldType     string    `excel:"模具类型"`           // 模具类型
	OrderCode    string    `excel:"工单号"`            // 工单号
	PartCode     string    `excel:"零件号"`            // 零件号
	LineLevel    string    `excel:"线别" dict:"line"` // 线别
	CompleteTime base.Time `excel:"生产完成时间"`         // 生产完成时间
	Count        int64     `excel:"工单数量"`           // 工单数量
} // @name ProductResumeExportVO

type MoldMaintenanceInfo struct {
	model.MoldInfo
	PartCode string `json:"partCode"` // 零件号
}

type PendingCheckProductResumeMoldInfo struct {
	model.MoldInfo
	ResumeID    int64  `json:"resumeId"`    //履历ID
	OrderCode   string `json:"orderCode"`   // 工单号
	ResumeCount int64  `json:"resumeCount"` // 工单数量
}

type MaintenancePlanRel struct {
	model.MoldMaintenancePlan
	MoldId int64 `json:"moldId"` //模具ID
}
