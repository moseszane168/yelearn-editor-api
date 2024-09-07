package email

import "crf-mold/base"

type EmailSetVO struct {
	Type     string `json:"type"`     // 类型：maintenance保养超时邮件发送/remodel模具改造邮件发送
	Receives string `json:"receives"` // 邮件接收人，使用逗号分隔
} // @name EmailSetVO

type EmailGetVO struct {
	MaintenanceReceives string `json:"maintenanceReceives"` // 保养超时邮件发送人
	RemodelReceives     string `json:"remodelReceives"`     // 模具改造邮件发送人
} // @name EmailGetVO

type MaintenanceTaskEmailVO struct {
	ID               int64     `json:"id"`               // 任务编码
	TaskCode         string    `json:"taskCode"`         // 模具名称
	MoldCode         string    `json:"moldCode"`         // 模具编码
	MoldName         string    `json:"moldName"`         // 模具名称
	Type             string    `json:"type"`             // 模具类型
	StandardName     string    `json:"standardName"`     // 标准名称
	TaskType         string    `json:"taskType"`         // 任务类型 auto-自动; manual-手动
	ProjectName      string    `json:"projectName"`      // 项目名称
	MaintenanceLevel string    `json:"maintenanceLevel"` // 保养等级
	GmtCreated       base.Time `json:"gmtCreated"`       // 创建时间
	GmtUpdated       base.Time `json:"gmtUpdated"`       // 更新时间
	Status           string    `json:"status"`           // 状态
	PartCodes        string    `json:"partCodes"`        // 零件号
	CreatedBy        string    `json:"createdBy"`        // 创建人
	Operator         string    `json:"operator"`         // 保养人
} // @name MaintenanceTaskEmailVO

type RemodelEmailVO struct {
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
} // @name RemodelEmailVO
