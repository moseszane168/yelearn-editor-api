package mold

import (
	"crf-mold/base"
)

//
// 保养标准
//
type MaintenanceStandardContentVO struct {
	Cycle               int    `json:"cycle"`               // 周期
	Item                string `json:"item"`                // 维护项
	JudgeMethod         string `json:"judgeMethod"`         // 维护判断方法
	MaintenanceMethod   string `json:"maintenanceMethod"`   // 维护方法
	EligibilityCriteria string `json:"eligibilityCriteria"` // 合格标准
} // @name MaintenanceStandardContentVO

type CreateMoldMaintenanceStandardVO struct {
	Time         string                         `json:"time" binding:"required"`         // 创建时间
	Name         string                         `json:"name" binding:"required"`         // 保养名称
	Code         string                         `json:"code" binding:"required"`         // 标准编号
	Version      string                         `json:"version" binding:"required"`      // 版本
	Type         string                         `json:"type" binding:"required"`         // 模具类型
	Level        string                         `json:"level" binding:"required"`        // 保养级别
	StandardType string                         `json:"standardType" binding:"required"` // 标准类型
	Content      []MaintenanceStandardContentVO `json:"content"`                         // 保养标准内容
	MoldIds      []int64                        `json:"moldIds"`                         // 关联模具，only 专用标准
} // @name CreateMoldMaintenanceStandardVO

type UpdateMoldMaintenanceStandardVO struct {
	ID           int64                          `json:"id" binding:"required"`           // ID
	Name         string                         `json:"name" binding:"required"`         // 保养名称
	Code         string                         `json:"code" binding:"required"`         // 标准编号
	Version      string                         `json:"version" binding:"required"`      // 版本
	Type         string                         `json:"type" binding:"required"`         // 模具类型
	Level        string                         `json:"level" binding:"required"`        // 保养级别
	StandardType string                         `json:"standardType" binding:"required"` // 标准类型
	Content      []MaintenanceStandardContentVO `json:"content"`                         // 保养标准内容
	MoldIds      []int64                        `json:"moldIds"`                         // 关联模具，only 专用标准
} // @name UpdateMoldMaintenanceStandardVO

type PageMaintenanceStandardInputVO struct {
	base.PageVO
	CodeOrName string `form:"codeOrName"`

	GmtCreatedBegin string `form:"gmtCreatedBegin"` // 创建时间开始
	GmtCreatedEnd   string `form:"gmtCreatedEnd"`   // 创建时间结束
	Name            string `form:"name"`            // 保养名称
	Code            string `form:"code"`            // 标准编号
	Level           string `form:"level"`           // 保养级别
	CreatedBy       string `form:"createdBy"`       // 创建人
} // @name PageMaintenanceStandardInputVO

type PageMaintenanceStandardOutputVO struct {
	ID           int64     `json:"id"`           // ID
	Name         string    `json:"name"`         // 保养名称
	Code         string    `json:"code"`         // 标准编号
	Version      string    `json:"version"`      // 版本
	Type         string    `json:"type"`         // 模具类型
	Level        string    `json:"level"`        // 保养级别
	StandardType string    `json:"standardType"` // 标准类型
	CreatedBy    string    `json:"createdBy"`    // 创建人
	GmtCreated   base.Time `json:"gmtCreated"`   // 创建时间
} // @name PageMaintenanceStandardOutputVO

type MoldRelPageOutVO struct {
	ID        int64  `json:"id"`        // ID
	Name      string `json:"name"`      // 模具名称
	Code      string `json:"code"`      // 模具编码
	LineLevel string `json:"lineLevel"` // 线别
	Process   string `json:"process"`   // 工序
	Level     string `json:"level"`     // 保养级别
	PartCodes string `json:"partCodes"` // 零件号
} // @name MoldRelPageOutVO

type MoldMaintenanceStandardOutVO struct {
	ID           int64     `json:"id"`           // ID
	Name         string    `json:"name"`         // 保养名称
	Code         string    `json:"code"`         // 标准编号
	Version      string    `json:"version"`      // 版本
	Type         string    `json:"type"`         // 模具类型
	Level        string    `json:"level"`        // 保养级别
	StandardType string    `json:"standardType"` // 标准类型
	CreatedBy    string    `json:"createdBy"`    // 创建人
	GmtCreated   base.Time `json:"gmtCreated"`   // 创建时间

	Molds   []MoldRelPageOutVO             `json:"moldRel"` // 关联模具
	Content []MaintenanceStandardContentVO `json:"content"` // 标准详情
} // @name MoldMaintenanceStandardOutVO

type MoldStandardSelectVO struct {
	ID               int64  `json:"id"`               // ID
	MaintenanceName  string `json:"maintenanceName"`  // 保养名称
	MaintenanceLevel string `json:"maintenanceLevel"` // 保养级别
} // @name MoldStandardSelectVO

//
// 模具保养计划
//
type CronRuntimeOutVO struct {
	RunTimeList []string `json:"runTimeList"` //Cron表达式未来运行时间
} // @name CronRuntimeOutVO

type CronExpressionVo struct {
	CronExpression string //Cron表达式
} // @name CronExpressionVo

type CreateMoldMaintenancePlanVO struct {
	Name         string  `json:"name" binding:"required"`     // 模具保养计划名称
	PlanType     string  `json:"planType" binding:"required"` // 模具保养计划类型:timing-计时; metering-计量
	MoldType     string  `json:"moldType"`                    // 模具类型
	TaskStart    int64   `json:"taskStart"`                   // 任务生成区间最小值
	TaskEnd      int64   `json:"taskEnd"`                     // 任务生成区间最大值，大于
	TaskStandard int64   `json:"taskStandard"`                // 任务生成标准值，只做显示用
	Operate      string  `json:"operate"`                     // 任务生成区间最小值操作符：小于等于lte、小于lt、等于ge
	MoldIds      []int64 `json:"moldIds"`                     // 计时类型计划关联的模具ID列表
	PlanCron     string  `json:"planCron"`                    // 计时cron表达式
	TimeoutHours int     `json:"timeoutHours"`                // 计时类型计划超时时间(小时)
} // @name CreateMoldMaintenancePlanVO

type UpdateMoldMaintenancePlanVO struct {
	ID           int64   `json:"id" binding:"required"`
	Name         string  `json:"name" binding:"required"`     // 模具保养计划名称
	PlanType     string  `json:"planType" binding:"required"` // 模具保养计划类型:timing-计时; metering-计量
	MoldType     string  `json:"moldType"`                    // 模具类型
	TaskStart    int64   `json:"taskStart"`                   // 任务生成区间最小值
	TaskEnd      int64   `json:"taskEnd"`                     // 任务生成区间最大值，大于
	TaskStandard int64   `json:"taskStandard"`                // 保养冲刺，任务生成标准值，只做显示用
	Operate      string  `json:"operate"`                     // 任务生成区间最小值操作符：小于等于lte、小于lt、等于ge
	MoldIds      []int64 `json:"moldIds"`                     // 计时类型计划关联的模具ID列表
	PlanCron     string  `json:"planCron"`                    // 计时cron表达式
	TimeoutHours int     `json:"timeoutHours"`                // 计时类型计划超时时间(小时)
} // @name UpdateMoldMaintenancePlanVO

type PageMaintenancePlanInputVO struct {
	base.PageVO
	CodeOrName string `form:"codeOrName"`

	GmtCreatedBegin string `form:"gmtCreatedBegin" json:"gmtCreatedBegin"` // 创建时间开始
	GmtCreatedEnd   string `form:"gmtCreatedEnd" json:"gmtCreatedEnd"`     // 创建时间结束
} // @name PageMaintenancePlanInputVO

type PageMaintenancePlanOutputVO struct {
	ID           int64      `json:"id"`
	Code         string     `json:"code"`         // 模具编号
	Name         string     `json:"name"`         // 模具保养计划名称
	PlanType     string     `json:"planType"`     // 模具保养计划类型
	MoldType     string     `json:"moldType"`     // 模具类型
	TaskStart    int64      `json:"taskStart"`    // 任务生成区间最小值
	TaskEnd      int64      `json:"taskEnd"`      // 任务生成区间最大值，大于
	TaskStandard int64      `json:"taskStandard"` // 保养冲刺，任务生成标准值，只做显示用
	Operate      string     `json:"operate"`      // 任务生成区间最小值操作符：小于等于lte、小于lt、等于ge
	Status       string     `json:"status"`       // 状态：running开启，pause暂停
	CreatedBy    string     `json:"createdBy"`    // 创建人
	GmtCreated   *base.Time `json:"gmtCreated"`   // 创建时间
	Department   string     `json:"department"`   // 部门
} // @name PageMaintenancePlanOutputVO

type MaintenancePlanMoldRelOutVO struct {
	ID          int64  `json:"id"`          // ID
	Name        string `json:"name"`        // 模具名称
	Code        string `json:"code"`        // 模具编码
	ProjectName string `json:"projectName"` // 项目名称
	Type        string `json:"type"`        // 模具类型
	LineLevel   string `json:"lineLevel"`   // 线别
	Process     string `json:"process"`     // 工序
	PartCodes   string `json:"partCodes"`   // 零件号
} // @name MaintenancePlanMoldRelOutVO

type MaintenancePlanOutVO struct {
	ID           int64                         `json:"id"`           // ID
	Code         string                        `json:"code"`         // 模具编号
	Name         string                        `json:"name"`         // 模具保养计划名称
	PlanType     string                        `json:"planType"`     // 模具保养计划类型
	MoldType     string                        `json:"moldType"`     // 模具类型
	TaskStart    int64                         `json:"taskStart"`    // 任务生成区间最小值
	TaskEnd      int64                         `json:"taskEnd"`      // 任务生成区间最大值，大于
	TaskStandard int64                         `json:"taskStandard"` // 保养冲刺，任务生成标准值，只做显示用
	Operate      string                        `json:"operate"`      // 任务生成区间最小值操作符：小于等于lte、小于lt、等于ge
	PlanCron     string                        `json:"planCron"`     // 计时cron表达式
	TimeoutHours int                           `json:"timeoutHours"` // 计时类型计划超时时间(小时)
	Molds        []MaintenancePlanMoldRelOutVO `json:"moldRel"`      // 计时类型计划关联模具
	Status       string                        `json:"status"`       // 状态：running开启，pause暂停
	CreatedBy    string                        `json:"createdBy"`    // 创建人
	GmtCreated   *base.Time                    `json:"gmtCreated"`   // 创建时间
} // @name MaintenancePlanOutVO

type MaintenanceTimingPlanOutVO struct {
	ID           int64   // ID
	Name         string  // 计时保养计划名称
	PlanCron     string  // 计时cron表达式
	TimeoutHours int     // 计时类型计划超时时间(小时)
	Molds        []int64 // 计时类型计划关联模具
} // @name MaintenanceTimingPlanOutVO

//
// 模具保养任务
//
type MaintenanceTaskContentVO struct {
	Cycle               int    `json:"cycle"`               // 周期
	Item                string `json:"item"`                // 维护项
	JudgeMethod         string `json:"judgeMethod"`         // 维护判断方法
	MaintenanceMethod   string `json:"maintenanceMethod"`   // 维护方法
	EligibilityCriteria string `json:"eligibilityCriteria"` // 合格标准
	Status              string `json:"isMaintenance"`       // 是否保养,Y/N
	Remark              string `json:"remark"`              // 备注
} // @name MaintenanceTaskContentVO

type ReplaceSpareVO struct {
	SpareCode string `json:"spareCode"`
	Count     int    `json:"count"`
	Remark    string `json:"remark"`
} // @name ReplaceSpareVO

type AddMoldMaintenanceTaskVO struct {
	MoldID       int64  `json:"moldId" binding:"required"`   // 模具ID
	MoldType     string `json:"moldType" binding:"required"` // 模具类型 chongkong-冲孔，chengxing-成型
	PartCode     string `json:"partCode"`                    // 零件号
	TimeoutHours int    `json:"timeoutHours"`                // 任务超时时间(小时)
	ApprovalId   string `json:"approvalId"`                  // 审核人ID
	IsApproval   int64  `json:"isApproval"`                  // 是否需审批：0-需 1-无需
} // @name AddMoldMaintenanceTaskVO

type SubmitMoldMaintenanceTaskVO struct {
	ID                        int64                      `json:"id" binding:"required"`                        // ID
	Operator                  string                     `json:"operator" binding:"required"`                  // 保养人
	Time                      string                     `json:"time" binding:"required"`                      // 保养时间
	MaintenanceLevel          string                     `json:"maintenanceLevel" binding:"required"`          // 保养级别
	StandardName              string                     `json:"standardName" binding:"required"`              // 标准名称
	MoldMaintenanceStandardId int64                      `json:"moldMaintenanceStandardId" binding:"required"` // 标准ID
	Remark                    string                     `json:"remark" `                                      // 任务备注
	Content                   []MaintenanceTaskContentVO `json:"content"`                                      // 保养任务内容
	ReplaceSpare              []ReplaceSpareVO           `json:"replaceSpare"`                                 // 更换备件
} // @name SubmitMoldMaintenanceTaskVO

type PageMaintenanceTaskInputVO struct {
	base.PageVO
	CodeOrName      string `form:"codeOrName"`      // 模具编号/零件号
	GmtCreatedBegin string `form:"gmtCreatedBegin"` // 创建时间开始
	GmtCreatedEnd   string `form:"gmtCreatedEnd"`   // 创建时间结束

	TaskCode         string `form:"taskCode"`
	MoldCode         string `form:"moldCode"`
	RfidTid          string `form:"rfidTid"` // RFID TID
	Type             string `form:"type"`
	ProjectName      string `form:"projectName"`
	Operator         string `form:"operator"` // 保养人
	StandardName     string `form:"standardName"`
	TaskType         string `form:"taskType"` // 任务类型 auto-自动; manual-手动
	MaintenanceLevel string `form:"maintenanceLevel"`
	Status           string `form:"status"`
	PartCodes        string `form:"partCodes"`
} // @name PageMaintenanceTaskInputVO

type PageMaintenanceTaskOutVO struct {
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
} // @name PageMaintenanceTaskOutVO

type StatusVO struct {
	ID     int64  `json:"id" validate:"required"` // id
	Reason string `json:"reason"`                 // 挂起原因，如果是挂起则该字段有内容
} // @name StatusVO

type MaintenanceTaskOneOutVO struct {
	ID                        int64  `json:"id"`
	Code                      string `json:"code"`                      // 任务编码
	MoldName                  string `json:"moldName"`                  // 模具名称
	MoldCode                  string `json:"moldCode"`                  // 模具编码
	Type                      string `json:"type"`                      // 模具类型
	StandardName              string `json:"standardName"`              // 标准名称
	MoldMaintenanceStandardId *int64 `json:"moldMaintenanceStandardId"` // 标准ID
	TaskType                  string `json:"taskType"`                  // 任务类型 auto-自动; manual-手动
	Remark                    string `json:"remark"`                    // 任务备注
	MaintenanceLevel          string `json:"maintenanceLevel"`          // 保养等级
	CreatedBy                 string `json:"createdBy"`                 // 创建人

	Operator string `json:"operator"` // 保养人
	Time     string `json:"time"`     // 保养时间

	Content      []MaintenanceTaskContentVO `json:"content"`      // 保养任务内容
	ReplaceSpare []ReplaceSpareOneVO        `json:"replaceSpare"` // 更换备件
} // @name MaintenanceTaskOneOutVO

type MaintenanceTaskChargeVO struct {
	Id         int64  `json:"id" binding:"required"`         // 保养任务ID
	ChargeName string `json:"chargeName" binding:"required"` // 负责人名称
} // @name MaintenanceTaskChargeVO

type MaintenanceTaskApprovalVO struct {
	Id              int64  `json:"id" binding:"required"`              // 保养任务ID
	ApprovalStatus  int64  `json:"approvalStatus" binding:"required"`  // 审批状态：0-待审批 1-通过 2-驳回
	ApprovalComment string `json:"approvalComment" binding:"required"` // 审批意见
} // @name MaintenanceTaskApprovalVO

type MaintenanceTaskDeleteVO struct {
	Id     int64  `json:"id" binding:"required"`     // 保养任务ID
	Status string `json:"status" binding:"required"` // 任务状态, complete-已保养任务; timeout-超时任务; wait-待保养任务; timeout-超时任务
} // @name MaintenanceTaskDeleteVO

type ReplaceSpareOneVO struct {
	SpareCode string `json:"spareCode"` // 备件编号
	SpareName string `json:"spareName"` // 备件名称
	Flavor    string `json:"flavor"`    // 规格型号
	Count     int    `json:"count"`     // 更换数量
	Remark    string `json:"remark"`    // 备注
} // @name ReplaceSpareOneVO

//
// 模具生命周期
//
type MoldMaintenanceTaskLifecycleVO struct {
	ID                    int64     `json:"id"`
	MoldMaintenanceTaskID int64     `json:"moldMaintenanceTaskId"` // 模具标准任务ID
	Title                 string    `json:"title"`                 // 标题,字典key
	Operator              string    `json:"operator"`              // 操作人
	Time                  base.Time `json:"time"`                  // 时间
	Content               string    `json:"content"`               // 内容
} // @name MoldMaintenanceTaskLifecycleVO
