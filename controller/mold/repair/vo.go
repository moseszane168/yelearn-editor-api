package mold

import "crf-mold/base"

//
// 模具维修VO开始
//
type PageRepairVO struct {
	base.PageVO

	PartCodes  string `form:"partCodes"`  // 模具零件号
	CodeOrName string `form:"codeOrName"` // 模具编号和模具名称模糊查找，该参数有值时忽略其它条件参数

	ID              int64  `form:"id"`
	Code            string `form:"code"`            // 维修编号
	LineLevel       string `form:"lineLevel"`       // 线别
	FaultDesc       string `form:"faultDesc"`       // 故障描述
	MoldCode        string `form:"moldCode"`        // 模具编号
	RepairContent   string `form:"repairContent"`   // 维修内容
	ReportTimeBegin string `form:"reportTimeBegin"` // 报告时间开始
	ReportTimeEnd   string `form:"reportTimeEnd"`   // 报告时间结束
	ArriveTimeBegin string `form:"arriveTimeBegin"` // 到达时间开始
	ArriveTimeEnd   string `form:"arriveTimeEnd"`   // 到达时间结束
	FinishTimeBegin string `form:"finishTimeBegin"` // 完成时间开始
	FinishTimeEnd   string `form:"finishTimeEnd"`   // 完成时间结束
	Operator        string `form:"operator"`        // 操作者
	Repairtor       string `form:"repairtor"`       // 维修者
	AuditorProduct  string `form:"auditorProduct"`  // 审核-生产组长
	AuditorEngine   string `form:"auditorEngine"`   // 审核-维修工程师
	ImproveContent  string `form:"improveContent"`  // 永久改善内容
	Lockor          string `form:"lockor"`          // 锁定人
	Unlockor        string `form:"unlockor"`        // 解锁人
	Confirmor       string `form:"confirmor"`       // 确定人
	RepairLevel     string `form:"repairLevel"`     // 维修等级
	RepairStations  string `form:"repairStations"`  // 维修工站

	CodeLike string `form:"codeLike"` // 维修编码模糊
} // @name PageRepairVO

type PageRepairOutVO struct {
	ID               int64         `json:"id"`
	Code             string        `json:"code"`           // 维修编号
	LineLevel        string        `json:"lineLevel"`      // 线别
	FaultDesc        string        `json:"faultDesc"`      // 故障描述
	MoldCode         string        `json:"moldCode"`       // 模具编号
	RepairContent    string        `json:"repairContent"`  // 维修内容
	MoldRepairSpares []RepairSpare `json:"repairSpares"`   // 维修备件
	ReportTime       *base.Time    `json:"reportTime"`     // 报告时间开始
	ArriveTime       *base.Time    `json:"arriveTime"`     // 到达时间开始
	FinishTime       *base.Time    `json:"finishTime"`     // 完成时间开始
	LastTime         int           `json:"lastTime"`       // 维修时长(分钟)
	SumLastTime      int64         `json:"sumLastTime"`    // 总维修时长(分钟)
	Operator         string        `json:"operator"`       // 操作者
	Repairtor        string        `json:"repairtor"`      // 维修者
	AuditorProduct   string        `json:"auditorProduct"` // 审核-生产组长
	AuditorEngine    string        `json:"auditorEngine"`  // 审核-维修工程师
	ImproveContent   string        `json:"improveContent"` // 永久改善内容
	Lockor           string        `json:"lockor"`         // 锁定人
	Unlockor         string        `json:"unlockor"`       // 解锁人
	Confirmor        string        `json:"confirmor"`      // 确定人
	RepairLevel      string        `json:"repairLevel"`    // 维修等级
	RepairStation    string        `json:"repairStation"`  // 维修工站
	RepairStations   []string      `json:"repairStations"` // 维修工站
	PartCodes        []string      `json:"partCodes"`      // 模具零件号
	PartCodeStr      string        `json:"partCodeStr"`    // 模具零件号, '/'拼接
} // @name PageRepairOutVO

type RepairSpare struct {
	Code  string `json:"code"`
	Count int    `json:"count"`
} // @name RepairSpare

type CreateRepairVO struct {
	LineLevel      string        `json:"lineLevel"`                 // 线别
	FaultDesc      string        `json:"faultDesc"`                 // 故障描述
	MoldCode       string        `json:"moldCode"`                  // 模具编号
	RepairContent  string        `json:"repairContent"`             // 维修内容
	RepairSpares   []RepairSpare `json:"repairSpares"`              // 维修备件更换
	ReportTime     *base.Time    `json:"reportTime"`                // 报告时间开始
	ArriveTime     *base.Time    `json:"arriveTime"`                // 到达时间开始
	FinishTime     *base.Time    `json:"finishTime"`                // 完成时间开始
	LastTime       int           `json:"lastTime"  binding:"min=0"` // 维修时长(分钟)
	Operator       string        `json:"operator"`                  // 操作者
	Repairtor      string        `json:"repairtor"`                 // 维修者
	AuditorProduct string        `json:"auditorProduct"`            // 审核-生产组长
	AuditorEngine  string        `json:"auditorEngine"`             // 审核-维修工程师
	ImproveContent string        `json:"improveContent"`            // 永久改善内容
	Lockor         string        `json:"lockor"`                    // 锁定人
	Unlockor       string        `json:"unlockor"`                  // 解锁人
	Confirmor      string        `json:"confirmor"`                 // 确定人
	RepairLevel    string        `json:"repairLevel"`               // 维修等级
	RepairStations []string      `json:"repairStations"`            // 维修工站
} // @name CreateRepairVO

type UpdateRepairVO struct {
	ID             int64         `json:"id"`
	LineLevel      string        `json:"lineLevel"`                // 线别
	FaultDesc      string        `json:"faultDesc"`                // 故障描述
	MoldCode       string        `json:"moldCode"`                 // 模具编号
	RepairContent  string        `json:"repairContent"`            // 维修内容
	RepairSpares   []RepairSpare `json:"repairSpares"`             // 维修备件更换
	ReportTime     *base.Time    `json:"reportTime"`               // 报告时间开始
	ArriveTime     *base.Time    `json:"arriveTime"`               // 到达时间开始
	FinishTime     *base.Time    `json:"finishTime"`               // 完成时间开始
	LastTime       int           `json:"lastTime" binding:"min=0"` // 维修时长(分钟)
	Operator       string        `json:"operator"`                 // 操作者
	Repairtor      string        `json:"repairtor"`                // 维修者
	AuditorProduct string        `json:"auditorProduct"`           // 审核-生产组长
	AuditorEngine  string        `json:"auditorEngine"`            // 审核-维修工程师
	ImproveContent string        `json:"improveContent"`           // 永久改善内容
	Lockor         string        `json:"lockor"`                   // 锁定人
	Unlockor       string        `json:"unlockor"`                 // 解锁人
	Confirmor      string        `json:"confirmor"`                // 确定人
	RepairLevel    string        `json:"repairLevel"`              // 维修等级
	RepairStations []string      `json:"repairStations"`           // 维修工站
} // @name UpdateRepairVO
