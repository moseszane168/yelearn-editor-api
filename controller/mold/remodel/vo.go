package mold

import "crf-mold/base"

//
// 模具改造
//
type PageRemodelVO struct {
	base.PageVO

	CodeOrName string `form:"codeOrName"` // 模具编号和模具名称模糊查找，该参数有值时忽略其它条件参数

	ID                    int64     `form:"id"`
	Code                  string    `form:"code"`                  // 改造编码
	MoldCode              string    `form:"moldCode"`              // 模具编码
	RemodelStartTimeBegin string    `form:"remodelStartTimeBegin"` // 改造开始时间Begin
	RemodelStartTimeEnd   string    `form:"remodelStartTimeEnd"`   // 改造开始时间End
	RemodelEndTimeBegin   string    `form:"remodelEndTimeBegin"`   // 改造结束时间Begin
	RemodelEndTimeEnd     string    `form:"remodelEndTimeEnd"`     // 改造结束时间End
	FinishTime            base.Time `form:"finishTime"`            // 改造完成时间
	Director              string    `form:"director"`              // 责任人
	Type                  string    `form:"type"`                  // 改造类别
	Location              string    `form:"location"`              // 改造地点
	Content               string    `form:"content"`               // 改造内容
	Status                string    `form:"status"`                // 状态
	WithdrawReason        string    `form:"withdrawReason"`        // 撤销原因
	IsDelay               string    `form:"isDelay"`               // 是否延期

	CodeLike string `form:"codeLike"` // 改造编码模糊
} // @name PageRemodelVO

type PageRemodelOutVO struct {
	ID               int64      `json:"id"`
	Code             string     `json:"code"`             // 改造编码
	MoldCode         string     `json:"moldCode"`         // 模具编码
	PartCodes        string     `json:"partCodes"`        // 零件号
	ProjectName      string     `json:"projectName"`      // 项目名称
	RemodelStartTime *base.Time `json:"remodelStartTime"` // 改造开始时间
	RemodelEndTime   *base.Time `json:"remodelEndTime"`   // 改造结束时间
	FinishTime       *base.Time `json:"finishTime"`       // 改造完成时间
	Director         string     `json:"director"`         // 责任人
	Type             string     `json:"type"`             // 改造类别
	Location         string     `json:"location"`         // 改造地点
	Content          string     `json:"content"`          // 改造内容
	Status           string     `json:"status"`           // 状态
	WithdrawReason   string     `json:"withdrawReason"`   // 撤销原因
	IsDelay          string     `json:"isDelay"`          // 是否延期
	DelayDay         int        `json:"delayDay"`         // 延期天数，范围：1-7
} // @name PageRemodelOutVO

type CreateRemodelVO struct {
	MoldCode         string     `json:"moldCode" binding:"required"`         // 模具编码
	RemodelStartTime *base.Time `json:"remodelStartTime" binding:"required"` // 改造开始时间
	RemodelEndTime   *base.Time `json:"remodelEndTime" binding:"required"`   // 改造结束时间
	Director         string     `json:"director" binding:"required"`         // 责任人
	Type             string     `json:"type" binding:"required"`             // 改造类别
	Location         string     `json:"location" binding:"required"`         // 改造地点
	Content          string     `json:"content" binding:"max=100"`           // 改造内容
} // @name CreateRemodelVO

type UpdateRemodelVO struct {
	ID               int64      `json:"id"`
	MoldCode         string     `json:"moldCode"`         // 模具编码
	RemodelStartTime *base.Time `json:"remodelStartTime"` // 改造开始时间
	RemodelEndTime   *base.Time `json:"remodelEndTime"`   // 改造结束时间
	Director         string     `json:"director"`         // 责任人
	Type             string     `json:"type"`             // 改造类别
	Location         string     `json:"location"`         // 改造地点
	Content          string     `json:"content"`          // 改造内容
} // @name UpdateRemodelVO
