package mold

import "crf-mold/base"

//
// 模具质量VO开始
//
type PageQualityVO struct {
	base.PageVO

	PartCodes string `form:"partCodes"` // 模具零件号
	LikeQuery string `form:"likeQuery"` // 线别、模具编号和模具名称模糊查找，该参数有值时忽略其它条件参数

	ID                 int64  `form:"id"`
	Code               string `form:"code"`               // 质量编号
	LineLevel          string `form:"lineLevel"`          // 线别
	MoldCode           string `form:"moldCode"`           // 模具编号
	OrderCode          string `form:"orderCode"`          // 工单号
	ProductionCountMin int    `form:"productionCountMin"` // 生产数量最小值
	ProductionCountMax int    `form:"productionCountMax"` // 生产数量最大值
	DefectCountMin     int    `form:"defectCountMin"`     // 不良数量最小值
	DefectCountMax     int    `form:"defectCountMax"`     // 不良数量最大值
	DefectDesc         string `form:"defectDesc"`         // 不良描述
	QualityContent     string `form:"qualityContent"`     // 质量内容
	CreatedBy          string `form:"createdBy"`          // 录入人
} // @name PageQualityVO

type PageQualityOutVO struct {
	ID              int64      `json:"id"`
	Code            string     `json:"code"`            // 质量编号
	LineLevel       string     `json:"lineLevel"`       // 线别
	MoldId          int64      `json:"moldId"`          // 模具ID
	MoldCode        string     `json:"moldCode"`        // 模具编号
	OrderCode       string     `json:"orderCode"`       // 工单号
	ProductionCount int        `json:"productionCount"` // 生产数量
	DefectCount     int        `json:"defectCount"`     // 不良数量
	DefectDesc      string     `json:"defectDesc"`      // 不良描述
	QualityContent  string     `json:"qualityContent"`  // 质量内容
	CreatedBy       string     `json:"createdBy"`       // 录入人
	UpdatedBy       string     `json:"updatedBy"`       // 更新人
	GmtCreated      *base.Time `json:"gmtCreated"`      // 录入时间
	GmtUpdated      *base.Time `json:"gmtUpdated"`      // 更新时间
	PartCodes       string     `json:"partCodes"`       // 模具零件号, '/'拼接
} // @name PageQualityOutVO

type CreateQualityVO struct {
	LineLevel       string `json:"lineLevel"`       // 线别
	MoldId          int64  `json:"moldId"`          // 模具ID
	OrderCode       string `json:"orderCode"`       // 工单号
	ProductionCount int    `json:"productionCount"` // 生产数量
	DefectCount     int    `json:"defectCount"`     // 不良数量
	DefectDesc      string `json:"defectDesc"`      // 不良描述
	QualityContent  string `json:"qualityContent"`  // 质量内容
} // @name CreateQualityVO

type UpdateQualityVO struct {
	ID              int64  `json:"id"`
	LineLevel       string `json:"lineLevel"`       // 线别
	MoldId          int64  `json:"moldId"`          // 模具ID
	OrderCode       string `json:"orderCode"`       // 工单号
	ProductionCount int    `json:"productionCount"` // 生产数量
	DefectCount     int    `json:"defectCount"`     // 不良数量
	DefectDesc      string `json:"defectDesc"`      // 不良描述
	QualityContent  string `json:"qualityContent"`  // 质量内容
} // @name UpdateQualityVO
