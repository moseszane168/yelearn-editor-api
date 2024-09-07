package model

import (
	"crf-mold/base"
	"time"
)

// GetAutoMigrateTables 获取需要自动迁移的数据库对象
func GetAutoMigrateTables() []interface{} {
	AutoMigrateTables := []interface{}{
		&MoldQuality{},
	}
	return AutoMigrateTables
}

// DictGroup [...]
type DictGroup struct {
	ID         int64      `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Code       string     `gorm:"column:code;type:varchar(64);not null" json:"code"`               // 组编码
	Name       string     `gorm:"column:name;type:varchar(255);not null" json:"name"`              // 组名称
	IsDeleted  string     `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	GmtCreated *base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated *base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy  string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string     `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *DictGroup) TableName() string {
	return "dict_group"
}

// DictProperty [...]
type DictProperty struct {
	ID         int64      `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Key        string     `gorm:"column:key;type:varchar(255);not null" json:"key"`                // key
	ValueCn    string     `gorm:"column:value_cn;type:varchar(255)" json:"valueCn"`                // value_cn
	ValueEn    string     `gorm:"column:value_en;type:varchar(255)" json:"valueEn"`                // value_en
	Order      uint32     `gorm:"column:order;type:int(255) unsigned;not null" json:"order"`       // 排序
	GroupCode  string     `gorm:"column:group_code;type:varchar(64);not null" json:"groupCode"`    // 组编码
	IsDeleted  string     `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	GmtCreated *base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated *base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy  string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string     `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *DictProperty) TableName() string {
	return "dict_property"
}

// MoldPartRel [...]
type MoldPartRel struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	MoldCode   string    `gorm:"column:mold_code;type:varchar(64);not null" json:"moldCode"`      // 模具ID
	PartCode   string    `gorm:"column:part_code;type:varchar(255);not null" json:"partCode"`     // 零件号
	IsDeleted  string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
	GmtCreated base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldPartRel) TableName() string {
	return "mold_part_rel"
}

// MoldCustomInfo [...]
type MoldCustomInfo struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	MoldCode   string    `gorm:"column:mold_code;type:varchar(64);not null" json:"moldCode"`      // 模具ID
	Key        string    `gorm:"column:key;type:varchar(255);not null" json:"key"`                // 字段名
	Value      string    `gorm:"column:value;type:varchar(255);not null" json:"value"`            // 字段值
	IsDeleted  string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
	GmtCreated base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldCustomInfo) TableName() string {
	return "mold_custom_info"
}

// MoldInfo [...]
type MoldInfo struct {
	ID             int64      `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Code           string     `gorm:"column:code;type:varchar(64);not null" json:"code"`                       // 模具编号
	Name           string     `gorm:"column:name;type:varchar(255);not null" json:"name"`                      // 模具名称
	ClientName     string     `gorm:"column:client_name;type:varchar(255);not null" json:"clientName"`         // 客户名称
	Weight         float64    `gorm:"column:weight;type:double;not null" json:"weight"`                        // 模具重量
	ProjectName    string     `gorm:"column:project_name;type:varchar(255);not null" json:"projectName"`       // 项目名称
	SizeLong       float64    `gorm:"column:size_long;type:double" json:"sizeLong"`                            // 模具长
	SizeWidth      float64    `gorm:"column:size_width;type:double" json:"sizeWidth"`                          // 模具宽
	SizeHeigh      float64    `gorm:"column:size_heigh;type:double" json:"sizeHeigh"`                          // 模具高
	Platform       string     `gorm:"column:platform;type:varchar(64);not null" json:"platform"`               // 平台
	Process        string     `gorm:"column:process;type:varchar(64)" json:"process"`                          // 工序
	PropertyNumber string     `gorm:"column:property_number;type:varchar(255);not null" json:"propertyNumber"` // 资产编号
	Type           string     `gorm:"column:type;type:varchar(64);not null" json:"type"`                       // 模具类型
	Status         string     `gorm:"column:status;type:varchar(64);not null" json:"status"`                   // 模具状态
	MakeDate       *base.Time `gorm:"column:make_date;type:datetime" json:"makeDate"`                          // 制造日期
	FirstUseDate   *base.Time `gorm:"column:first_use_date;type:datetime" json:"firstUseDate"`                 // 第一次使用日期
	YearLimit      float64    `gorm:"column:year_limit;type:double" json:"yearLimit"`                          // 使用年限
	Rfid           string     `gorm:"column:rfid;type:varchar(64)" json:"rfid"`                                // RFID
	LineLevel      string     `gorm:"column:line_level;type:varchar(64);not null" json:"lineLevel"`            // 使用线别
	Provide        string     `gorm:"column:provide;type:varchar(255);not null" json:"provide"`                // 模具供应商
	FlushCount     int64      `gorm:"column:flush_count;type:bigint(64)" json:"flushCount"`                    // 总冲次
	CalcFlushCount int64      `gorm:"column:calc_flush_count;type:bigint(64)" json:"calcFlushCount"`           // 计算冲次
	Category       string     `gorm:"column:category;type:varchar(2);not null" json:"category"`                // 模具类别
	IsDeleted      string     `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`            // 是否删除
	CreatedBy      string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                     // 创建人
	UpdatedBy      string     `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                     // 修改人
	GmtCreated     *base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`         // 创建时间
	GmtUpdated     *base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`         // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldInfo) TableName() string {
	return "mold_info"
}

// MoldKnowledge [...]
type MoldKnowledge struct {
	ID         int64         `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Name       string        `gorm:"column:name;type:varchar(255);not null" json:"name"`              // 知识名称
	Group      string        `gorm:"column:group;type:bigint(64);not null" json:"group"`              // 知识分类
	Content    string        `gorm:"column:content;type:mediumtext;not null" json:"content"`          // 知识内容
	Files      base.FileList `gorm:"column:files;type:json;default:null" json:"files"`                // 上传的文件
	Style      string        `gorm:"column:style;type:varchar(255);not null" json:"style"`            // 书皮样式
	IsDeleted  string        `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	GmtCreated *base.Time    `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated *base.Time    `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy  string        `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string        `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *MoldKnowledge) TableName() string {
	return "mold_knowledge"
}

// MoldRemodel 模具改造
type MoldRemodel struct {
	ID               int64      `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"-"`
	Code             string     `gorm:"column:code;type:varchar(64);not null" json:"code"`                        // 改造编号
	MoldCode         string     `gorm:"column:mold_code;type:varchar(64);not null" json:"moldCode"`               // 模具编码
	RemodelStartTime *base.Time `gorm:"column:remodel_start_time;type:datetime;not null" json:"remodelStartTime"` // 改造开始时间
	RemodelEndTime   *base.Time `gorm:"column:remodel_end_time;type:datetime;not null" json:"remodelEndTime"`     // 改造结束时间
	FinishTime       *base.Time `gorm:"column:finish_time;type:datetime" json:"finishTime"`                       // 改造完成时间
	Type             string     `gorm:"column:type;type:varchar(64);not null" json:"type"`                        // 改造类别
	Location         string     `gorm:"column:location;type:varchar(64);not null" json:"location"`                // 改造地点
	Content          string     `gorm:"column:content;type:varchar(255);not null" json:"content"`                 // 改造内容
	IsDelay          string     `gorm:"column:is_delay;type:varchar(1);not null;default:N" json:"isDelay"`        // 是否延期
	DelayEmail       string     `gorm:"column:delay_email;type:varchar(1);not null;default:N" json:"delayEmail"`  // 是否发送延期通知邮件
	Director         string     `gorm:"column:director;type:varchar(64);not null" json:"director"`                // 责任人
	DelayDay         int        `gorm:"column:delay_day;type:int(1)" json:"delayDay"`                             // 延期天数，范围：1-7
	Status           string     `gorm:"column:status;type:varchar(64);default:wait" json:"status"`                // 单据状态：wait待完成、complete完工、withdraw撤单
	WithdrawReason   string     `gorm:"column:withdraw_reason;type:varchar(255)" json:"withdrawReason"`           // 撤单原因
	IsDeleted        string     `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`             // 是否删除
	GmtCreated       time.Time  `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`          // 创建时间
	GmtUpdated       time.Time  `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`          // 修改时间
	CreatedBy        string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                      // 创建人
	UpdatedBy        string     `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                      // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *MoldRemodel) TableName() string {
	return "mold_remodel"
}

// MoldRepair [...]
type MoldRepair struct {
	ID             int64      `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Code           string     `gorm:"column:code;type:varchar(64);not null" json:"code"`                       // 维修编号
	LineLevel      string     `gorm:"column:line_level;type:varchar(64);not null" json:"lineLevel"`            // 线别
	FaultDesc      string     `gorm:"column:fault_desc;type:varchar(255);not null" json:"faultDesc"`           // 故障描述
	MoldCode       string     `gorm:"column:mold_code;type:varchar(64);not null" json:"moldCode"`              // 模具编号
	RepairContent  string     `gorm:"column:repair_content;type:varchar(255)" json:"repairContent"`            // 维修内容
	ReportTime     *base.Time `gorm:"column:report_time;type:datetime" json:"reportTime"`                      // 报告时间
	ArriveTime     *base.Time `gorm:"column:arrive_time;type:datetime" json:"arriveTime"`                      // 到达时间
	FinishTime     *base.Time `gorm:"column:finish_time;type:datetime" json:"finishTime"`                      // 完成时间
	LastTime       int        `gorm:"column:last_time;type:int(11)" json:"lastTime"`                           // 维修时长(分钟)
	Operator       string     `gorm:"column:operator;type:varchar(64);not null" json:"operator"`               // 操作者
	Repairtor      string     `gorm:"column:repairtor;type:varchar(64);not null" json:"repairtor"`             // 维修者
	AuditorProduct string     `gorm:"column:auditor_product;type:varchar(64);not null" json:"auditorProduct"`  // 审核-生产组长
	AuditorEngine  string     `gorm:"column:auditor_engine;type:varchar(64);not null" json:"auditorEngine"`    // 审核-维修工程师
	ImproveContent string     `gorm:"column:improve_content;type:varchar(255);not null" json:"improveContent"` // 永久改善内容
	Lockor         string     `gorm:"column:lockor;type:varchar(64);not null" json:"lockor"`                   // 锁定人
	Unlockor       string     `gorm:"column:unlockor;type:varchar(64);not null" json:"unlockor"`               // 解锁人
	Confirmor      string     `gorm:"column:confirmor;type:varchar(64);not null" json:"confirmor"`             // 确定人
	RepairLevel    string     `gorm:"column:repair_level;type:varchar(64);not null" json:"repairLevel"`        // 维修等级
	RepairStation  string     `gorm:"column:repair_station;type:varchar(255)" json:"repairStation"`            // 维修工站
	IsDeleted      string     `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`            // 是否删除
	CreatedBy      string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                     // 创建人
	UpdatedBy      string     `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                     // 修改人
	GmtCreated     base.Time  `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`         // 创建时间
	GmtUpdated     base.Time  `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`         // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldRepair) TableName() string {
	return "mold_repair"
}

// MoldQuality [...]
type MoldQuality struct {
	ID              int64      `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Code            string     `gorm:"column:code;type:varchar(64);not null" json:"code"`                    // 质量编号
	LineLevel       string     `gorm:"column:line_level;type:varchar(64);not null" json:"lineLevel"`         // 线别
	MoldId          int64      `gorm:"column:mold_id;type:bigint(64)" json:"moldId"`                         // 模具ID
	OrderCode       string     `gorm:"column:order_code;type:varchar(64);not null" json:"orderCode"`         // 工单号
	ProductionCount int        `gorm:"column:production_count;type:int(11);not null" json:"productionCount"` // 生产数量
	DefectCount     int        `gorm:"column:defect_count;type:int(11);not null" json:"defectCount"`         // 不良数量
	DefectDesc      string     `gorm:"column:defect_desc;type:varchar(255);not null" json:"defectDesc"`      // 不良描述
	QualityContent  string     `gorm:"column:quality_content;type:varchar(255)" json:"qualityContent"`       // 质量内容
	IsDeleted       string     `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`         // 是否删除
	CreatedBy       string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                  // 录入人
	UpdatedBy       string     `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                  // 修改人
	GmtCreated      *base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`      // 录入时间
	GmtUpdated      *base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`      // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldQuality) TableName() string {
	return "mold_quality"
}

// MoldReplaceSpareRel 模具维修备件更换
type MoldReplaceSpareRel struct {
	ID                int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"-"`
	RepairID          int64     `gorm:"column:repair_id;type:bigint(64)" json:"repairId"`                    // 维修ID
	MaintenanceTaskID int64     `gorm:"column:maintenance_task_id;type:bigint(64)" json:"maintenanceTaskId"` // 保养任务ID
	MoldID            int64     `gorm:"column:mold_id;type:bigint(64)" json:"moldId"`                        // 模具ID
	Type              string    `gorm:"column:type;type:varchar(50)" json:"type"`                            // 更换类型
	SpareCode         string    `gorm:"column:spare_code;type:varchar(64);not null" json:"spareCode"`        // 备件编号
	Count             int       `gorm:"column:count;type:int(11);not null" json:"count"`                     // 备件数量
	Remark            string    `gorm:"column:remark;type:varchar(50)" json:"remark"`                        // 备注
	IsDeleted         string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`        // 是否删除
	CreatedBy         string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                 // 创建人
	UpdatedBy         string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                 // 修改人
	GmtCreated        time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`     // 创建时间
	GmtUpdated        time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`     // 修改时间
	FlushCount        int64     `gorm:"column:flush_count;type:bigint(64)" json:"flushCount"`                // 冲次
	LastFlushCount    int64     `gorm:"column:last_flush_count;type:bigint(64)" json:"lastFlushCount"`       // 上次冲次
}

// TableName get sql table name.获取数据库表名
func (m *MoldReplaceSpareRel) TableName() string {
	return "mold_replace_spare_rel"
}

// Properties [...]
type Properties struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Key        string    `gorm:"column:key;type:varchar(64);not null" json:"key"`                 // key
	Value      string    `gorm:"column:value;type:varchar(2000);not null" json:"value"`           // value
	IsDeleted  string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	GmtCreated base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *Properties) TableName() string {
	return "properties"
}

// SpareInfo [...]
type SpareInfo struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Code       string    `gorm:"column:code;type:varchar(64);not null" json:"code"`               // 备件编号
	Name       string    `gorm:"column:name;type:varchar(64);not null" json:"name"`               // 备件名称
	Flavor     string    `gorm:"column:flavor;type:varchar(64);not null" json:"flavor"`           // 备件规格
	Material   string    `gorm:"column:material;type:varchar(64);not null" json:"material"`       // 材质
	IsDeleted  string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	GmtCreated base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *SpareInfo) TableName() string {
	return "spare_info"
}

// SpareRequest [...]
type SpareRequest struct {
	ID          int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	SpareCode   string    `gorm:"column:spare_code;type:varchar(64);not null" json:"spareCode"`    // 备件编号
	Type        string    `gorm:"column:type;type:varchar(64);not null" json:"type"`               // 出入库类型
	Location    string    `gorm:"column:location;type:varchar(64);not null" json:"location"`       // 出入库位置
	Count       int       `gorm:"column:count;type:int(11);not null" json:"count"`                 // 出入库数量
	Remark      string    `gorm:"column:remark;type:varchar(255)" json:"remark"`                   // 备注
	IsDeleted   string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	InboundTime int64     `gorm:"column:inbound_time;type:bigint(64)" json:"inboundTime"`          // 入库时间
	GmtCreated  base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated  base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy   string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy   string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *SpareRequest) TableName() string {
	return "spare_request"
}

// UserAuthority [...]
type UserAuthority struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Code       string    `gorm:"column:code;type:varchar(64);not null" json:"code"`                // 编码
	Key        string    `gorm:"column:key;type:varchar(64);not null" json:"key"`                  // 资源标识，唯一
	Name       string    `gorm:"column:name;type:varchar(255);not null" json:"name"`               // 资源名称
	URI        string    `gorm:"column:uri;type:varchar(255);not null" json:"uri"`                 // 资源url
	Method     string    `gorm:"column:method;type:varchar(10);not null" json:"method"`            // 请求方式
	GroupName  string    `gorm:"column:group_name;type:varchar(255);not null" json:"groupName"`    // 资源组名称
	Display    string    `gorm:"column:display;type:varchar(1);not null;default:Y" json:"display"` // 是否返回前端显示
	IsDeleted  string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`     // 是否删除
	GmtCreated base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`  // 创建时间
	GmtUpdated base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`  // 修改时间
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`              // 创建人
	UpdatedBy  string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`              // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *UserAuthority) TableName() string {
	return "user_authority"
}

// UserAuthorityRel [...]
type UserAuthorityRel struct {
	ID            int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	LoginName     string    `gorm:"column:login_name;type:varchar(64);not null" json:"loginName"`         // 用户登录名
	AuthorityCode string    `gorm:"column:authority_code;type:varchar(64);not null" json:"authorityCode"` // 权限code
	IsDeleted     string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`         // 是否删除
	GmtCreated    base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`      // 创建时间
	GmtUpdated    base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`      // 修改时间
	CreatedBy     string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                  // 创建人
	UpdatedBy     string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                  // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *UserAuthorityRel) TableName() string {
	return "user_authority_rel"
}

// UserInfo [...]
type UserInfo struct {
	ID         int64      `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	LoginName  string     `gorm:"column:login_name;type:varchar(64);not null" json:"loginName"`    // 工号
	Name       string     `gorm:"column:name;type:varchar(64);not null" json:"name"`               // 姓名
	Department string     `gorm:"column:department;type:varchar(64);not null" json:"department"`   // 部门
	Password   string     `gorm:"column:password;type:varchar(255)" json:"password"`               // 密码，rsa加密
	IsRoot     string     `gorm:"column:is_root;type:varchar(1);default:N" json:"isRoot"`          // 是否超级管理员
	IsLocked   string     `gorm:"column:is_locked;type:varchar(1);default:N" json:"isLocked"`      // 是否锁定
	IsDeleted  string     `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	UnlockTime *time.Time `gorm:"column:unlock_time;type:datetime" json:"unlock_time"`             // 解锁时间
	GmtCreated base.Time  `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated base.Time  `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy  string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string     `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *UserInfo) TableName() string {
	return "user_info"
}

// UserInfo [...]
type LoginInfo struct {
	ID            int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	LoginName     string    `gorm:"column:login_name;type:varchar(64);not null" json:"login_name"`          // 工号
	LastLoginDate time.Time `gorm:"column:last_login_date;type:datetime;default:null" json:"lastLoginDate"` // 上次登录时间
	FaultCount    int       `gorm:"column:fault_count;type:int(11)" json:"faultCount"`                      // 登录失败次数
	IsDeleted     string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`           // 是否删除
	GmtCreated    base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`        // 创建时间
	GmtUpdated    base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`        // 修改时间
	CreatedBy     string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                    // 创建人
	UpdatedBy     string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                    // 修改人
}

// TableName get sql table name.获取数据库表名
func (m *LoginInfo) TableName() string {
	return "login_info"
}

// MoldDoc 模具文档
type MoldDoc struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	FileKey    string    `gorm:"column:file_key;type:varchar(255);not null" json:"fileKey"`       // 文件Minio key
	MoldCode   string    `gorm:"column:mold_code;type:varchar(64);not null" json:"moldCode"`      // 模具编号
	Name       string    `gorm:"column:name;type:varchar(255);not null" json:"name"`              // 文档名称
	Version    string    `gorm:"column:version;type:varchar(64);not null" json:"version"`         // 版本
	Content    string    `gorm:"column:content;type:varchar(255);not null" json:"content"`        // 变更内容
	Remark     string    `gorm:"column:remark;type:varchar(255)" json:"remark"`                   // 备注
	IsDeleted  string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
	GmtCreated base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
}

// MoldBom 模具bom
type MoldBom struct {
	ID                int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	StandardComponent string    `gorm:"column:standard_component;type:varchar(1);not null;default:Y" json:"standardComponent"` // 标准件
	MoldCode          string    `gorm:"column:mold_code;type:varchar(64);not null" json:"moldCode"`                            // 模具编号
	PartName          string    `gorm:"column:part_name;type:varchar(255);not null" json:"partName"`                           // 零件代号
	PartCode          string    `gorm:"column:part_code;type:varchar(255);not null" json:"partCode"`                           // 零件名称
	Count             int       `gorm:"column:count;type:int(11);not null" json:"count"`                                       // 数量
	Material          string    `gorm:"column:material;type:varchar(255)" json:"material"`                                     // 材料
	Flavor            string    `gorm:"column:flavor;type:varchar(255);not null" json:"flavor"`                                // 规格
	Standard          string    `gorm:"column:standard;type:varchar(255)" json:"standard"`                                     // 标准
	HeatTreating      string    `gorm:"column:heat_treating;type:varchar(255)" json:"heatTreating"`                            // 热处理
	AssemblyRelation  string    `gorm:"column:assembly_relation;type:varchar(255)" json:"assemblyRelation"`                    // 装配关系
	Remark            string    `gorm:"column:remark;type:varchar(255)" json:"remark"`                                         // 备注
	IsDeleted         string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`                          // 是否删除
	CreatedBy         string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                                   // 创建人
	UpdatedBy         string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                                   // 修改人
	GmtCreated        base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`                       // 创建时间
	GmtUpdated        base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`                       // 修改时间
}

// MoldParams 模具成型参数
type MoldParams struct {
	ID              int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	MoldCode        string    `gorm:"column:mold_code;type:varchar(64);not null" json:"moldCode"`              // 模具编号
	Height          float64   `gorm:"column:height;type:double;not null" json:"height"`                        // 模具高度
	DefensiveDevice string    `gorm:"column:defensive_device;type:varchar(1);not null" json:"defensiveDevice"` // 防错装置
	Desc            string    `gorm:"column:desc;type:varchar(255);not null" json:"desc"`                      // 描述
	IsDeleted       string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`            // 是否删除
	CreatedBy       string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                     // 创建人
	UpdatedBy       string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                     // 修改人
	GmtCreated      time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`         // 创建时间
	GmtUpdated      time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`         // 修改时间
}

// MoldParamsCustom 模具成型参数自定义字段
type MoldParamsCustom struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	MoldCode   string    `gorm:"column:mold_code;type:varchar(64)" json:"moldCode"`               // 模具ID
	Key        string    `gorm:"column:key;type:varchar(255);not null" json:"key"`                // 字段名
	Value      string    `gorm:"column:value;type:varchar(255);not null" json:"value"`            // 字段值
	IsDeleted  string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
	GmtCreated time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
}

// MoldMaintenancePlan 模具保养计划
type MoldMaintenancePlan struct {
	ID           int64      `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Code         string     `gorm:"column:code;type:varchar(64);not null" json:"code"`               // 模具编号
	Name         string     `gorm:"column:name;type:varchar(50);not null" json:"name"`               // 模具保养计划名称
	PlanType     string     `gorm:"column:plan_type;type:varchar(50);not null" json:"plan_type"`     // 计划类型: timing-计时; metering-计量
	MoldType     string     `gorm:"column:mold_type;type:varchar(50);not null" json:"moldType"`      // 模具类型
	TaskStart    int64      `gorm:"column:task_start;type:int(11);not null" json:"taskStart"`        // 任务生成区间最小值，小于等于
	TaskEnd      int64      `gorm:"column:task_end;type:int(11);not null" json:"taskEnd"`            // 任务生成区间最大值，大于
	TaskStandard int64      `gorm:"column:task_standard;type:int(11);not null" json:"taskStandard"`  // 任务生成标准值，只做显示用
	Operate      string     `gorm:"column:operate;type:varchar(64);not null" json:"operate"`         // 任务生成区间最小值操作符：小于等于lte、小于lt、等于ge
	PlanCron     string     `gorm:"column:plan_cron;type:varchar(100)" json:"plan_cron"`             // 计时cron表达式
	TimeoutHours int        `gorm:"column:timeout_hours;type:int(11)" json:"timeout_hours"`          // 计时类型计划超时时间(小时)
	Status       string     `gorm:"column:status;type:varchar(50)" json:"status"`                    // 状态
	IsDeleted    string     `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	CreatedBy    string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy    string     `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
	GmtCreated   *base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated   *base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldMaintenancePlan) TableName() string {
	return "mold_maintenance_plan"
}

// MoldMaintenancePlanRel 模具保养计划关联模具表
type MoldMaintenancePlanRel struct {
	ID                    int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	MoldMaintenancePlanID int64     `gorm:"column:mold_maintenance_plan_id;type:bigint(64);not null" json:"moldMaintenancePlanId"` // 模具保养计划ID
	MoldId                int64     `gorm:"column:mold_id;type:bigint(64);not null" json:"moldId"`                                 // 模具ID
	IsDeleted             string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`                          // 是否删除
	CreatedBy             string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                                   // 创建人
	UpdatedBy             string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                                   // 修改人
	GmtCreated            time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`                       // 创建时间
	GmtUpdated            time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`                       // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldMaintenancePlanRel) TableName() string {
	return "mold_maintenance_plan_rel"
}

// MoldMaintenanceStandard 模具保养标准
type MoldMaintenanceStandard struct {
	ID           int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Name         string    `gorm:"column:name;type:varchar(50);not null" json:"name"`                  // 保养名称
	Code         string    `gorm:"column:code;type:varchar(64)" json:"code"`                           // 标准编号
	Version      string    `gorm:"column:version;type:varchar(50);not null" json:"version"`            // 版本
	Type         string    `gorm:"column:type;type:varchar(50);not null" json:"type"`                  // 模具类型
	Level        string    `gorm:"column:level;type:varchar(50);not null" json:"level"`                // 保养级别
	StandardType string    `gorm:"column:standard_type;type:varchar(50);not null" json:"standardType"` // 标准类型
	IsDeleted    string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`       // 是否删除
	CreatedBy    string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                // 创建人
	UpdatedBy    string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                // 修改人
	GmtCreated   string    `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`    // 创建时间
	GmtUpdated   time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`    // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldMaintenanceStandard) TableName() string {
	return "mold_maintenance_standard"
}

// MoldMaintenanceStandardContent 模具保养标准内容
type MoldMaintenanceStandardContent struct {
	ID                        int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	MoldMaintenanceStandardID int64     `gorm:"column:mold_maintenance_standard_id;type:bigint(64);not null" json:"moldMaintenanceStandardId"` // 模具保养标准表ID
	Cycle                     int       `gorm:"column:cycle;type:int(11);not null" json:"cycle"`                                               // 周期
	Item                      string    `gorm:"column:item;type:varchar(50);not null" json:"item"`                                             // 维护项
	JudgeMethod               string    `gorm:"column:judge_method;type:varchar(300);not null" json:"judgeMethod"`                             // 维护判断方法
	MaintenanceMethod         string    `gorm:"column:maintenance_method;type:varchar(300);not null" json:"maintenanceMethod"`                 // 维护方法
	EligibilityCriteria       string    `gorm:"column:eligibility_criteria;type:varchar(300);not null" json:"eligibilityCriteria"`             // 合格标准
	Remark                    string    `gorm:"column:remark;type:varchar(300)" json:"remark"`                                                 // 备注
	Order                     int       `gorm:"column:order;type:int(11)" json:"order"`                                                        // 顺序
	IsDeleted                 string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`                                  // 是否删除
	CreatedBy                 string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                                           // 创建人
	UpdatedBy                 string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                                           // 修改人
	GmtCreated                time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`                               // 创建时间
	GmtUpdated                time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`                               // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldMaintenanceStandardContent) TableName() string {
	return "mold_maintenance_standard_content"
}

// MoldMaintenanceStandardRel 模具保养标准关联表
type MoldMaintenanceStandardRel struct {
	ID                        int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	MoldMaintenanceStandardID int64     `gorm:"column:mold_maintenance_standard_id;type:bigint(64);not null" json:"moldMaintenanceStandardId"` // 模具保养标准
	MoldId                    int64     `gorm:"column:mold_id;type:bigint(64);not null" json:"moldId"`                                         // 模具ID
	IsDeleted                 string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`                                  // 是否删除
	CreatedBy                 string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                                           // 创建人
	UpdatedBy                 string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                                           // 修改人
	GmtCreated                time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`                               // 创建时间
	GmtUpdated                time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`                               // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldMaintenanceStandardRel) TableName() string {
	return "mold_maintenance_standard_rel"
}

// MoldMaintenanceTask 模具保养任务
type MoldMaintenanceTask struct {
	ID                        int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	Code                      string    `gorm:"column:code;type:varchar(64);not null" json:"code"`                                    // 任务编号
	MoldID                    int64     `gorm:"column:mold_id;type:bigint(64);not null" json:"moldId"`                                // 模具ID
	MoldMaintenancePlainID    int64     `gorm:"column:mold_maintenance_plain_id;type:bigint(64)" json:"moldMaintenancePlainId"`       // 模具标准计划ID
	MoldMaintenanceStandardID int64     `gorm:"column:mold_maintenance_standard_id;type:bigint(64)" json:"moldMaintenanceStandardId"` // 模具标准ID
	Status                    string    `gorm:"column:status;type:varchar(64);not null" json:"status"`                                // 状态
	TaskType                  string    `gorm:"column:task_type;type:varchar(64);default:auto" json:"task_type"`                      // 任务类型 manual-手动 auto-自动
	PlanType                  string    `gorm:"column:plan_type;type:varchar(50);not null" json:"plan_type"`                          // 计划类型: timing-计时; metering-计量
	StandardName              string    `gorm:"column:standard_name;type:varchar(50)" json:"standardName"`                            // 标准名称
	Remark                    string    `gorm:"column:remark;type:varchar(300)" json:"remark"`                                        // 任务备注
	MaintenanceLevel          string    `gorm:"column:maintenance_level;type:varchar(50)" json:"maintenanceLevel"`                    // 保养级别
	Operator                  string    `gorm:"column:operator;type:varchar(64)" json:"operator"`                                     // 保养人
	Time                      string    `gorm:"column:time;type:varchar(50)" json:"time"`                                             // 保养时间
	Reason                    string    `gorm:"column:reason;type:varchar(50)" json:"reason"`                                         // 挂起原因
	TaskGenInterval           int64     `gorm:"column:task_gen_interval;type:bigint(64)" json:"taskGenInterval"`                      // 任务生成空间
	HangUpCount               int64     `gorm:"column:hang_up_count;type:int(11)" json:"hangUpCount"`                                 // 挂起次数
	TimeoutCount              int64     `gorm:"column:timeout_count;type:int(11)" json:"timeoutCount"`                                // 超时次数
	TimeoutTime               time.Time `gorm:"column:timeout_time;type:datetime;default:null" json:"timeoutTime"`                    // 计时超时时间
	IsDeleted                 string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`                         // 是否删除
	PartCode                  string    `gorm:"column:part_code;type:varchar(64)" json:"partCode"`                                    // 零件号
	CreatedBy                 string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                                  // 创建人
	UpdatedBy                 string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                                  // 修改人
	GmtCreated                time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`                      // 创建时间
	GmtUpdated                time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`                      // 修改时间
	ApprovalId                string    `gorm:"column:approval_id;type:varchar(64);default:null" json:"approvalId"`                   // 审核人ID
	IsApproval                int64     `gorm:"column:is_approval;type:tinyint(1);default:0" json:"isApproval"`                       // 是否需审批：0-需 1-无需
	ApprovalStatus            int64     `gorm:"column:approval_status;type:tinyint(1);default:0" json:"approvalStatus"`               // 审批状态：0-待审批 1-通过 2-驳回
	ApprovalComment           string    `gorm:"column:approval_comment;type:varchar(50);default:null" json:"approvalComment"`         // 审批意见
}

// TableName get sql table name.获取数据库表名
func (m *MoldMaintenanceTask) TableName() string {
	return "mold_maintenance_task"
}

// MoldMaintenanceTaskContent 模具保养标准内容
type MoldMaintenanceTaskContent struct {
	ID                    int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"id"`
	MoldMaintenanceTaskID int64     `gorm:"column:mold_maintenance_task_id;type:bigint(64);not null" json:"moldMaintenanceTaskId"` // 模具保养任务ID
	Cycle                 int       `gorm:"column:cycle;type:int(11);not null" json:"cycle"`                                       // 周期
	Item                  string    `gorm:"column:item;type:varchar(50);not null" json:"item"`                                     // 维护项
	JudgeMethod           string    `gorm:"column:judge_method;type:varchar(300);not null" json:"judgeMethod"`                     // 维护判断方法
	MaintenanceMethod     string    `gorm:"column:maintenance_method;type:varchar(300);not null" json:"maintenanceMethod"`         // 维护方法
	EligibilityCriteria   string    `gorm:"column:eligibility_criteria;type:varchar(300);not null" json:"eligibilityCriteria"`     // 合格标准
	Remark                string    `gorm:"column:remark;type:varchar(300)" json:"remark"`                                         // 备注
	Order                 int       `gorm:"column:order;type:int(11)" json:"order"`                                                // 顺序
	Status                string    `gorm:"column:status;type:varchar(1);not null" json:"status"`                                  // 状态
	IsDeleted             string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`                          // 是否删除
	CreatedBy             string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                                   // 创建人
	UpdatedBy             string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                                   // 修改人
	GmtCreated            time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`                       // 创建时间
	GmtUpdated            time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`                       // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *MoldMaintenanceTaskContent) TableName() string {
	return "mold_maintenance_task_content"
}

// MoldMaintenanceTaskLifecycle 模具保养任务生命周期（任务时间线）
type MoldMaintenanceTaskLifecycle struct {
	ID                    int64     `gorm:"primaryKey;column:id;type:bigint(64) unsigned;not null" json:"id"`
	MoldMaintenanceTaskID int64     `gorm:"column:mold_maintenance_task_id;type:bigint(64);not null" json:"moldMaintenanceTaskId"` // 模具标准任务ID
	Title                 string    `gorm:"column:title;type:varchar(50);not null" json:"title"`                                   // 标题,字典key
	Operator              string    `gorm:"column:operator;type:varchar(64)" json:"operator"`                                      // 操作人
	Time                  string    `gorm:"column:time;type:datetime;not null" json:"time"`                                        // 时间
	Content               string    `gorm:"column:content;type:varchar(64);not null" json:"content"`                               // 内容
	IsDeleted             string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`                          // 是否删除
	CreatedBy             string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                                   // 创建人
	UpdatedBy             string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`                                   // 修改人
	GmtCreated            base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`                       // 创建时间
	GmtUpdated            base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"`                       // 修改时间
	ApprovalId            string    `gorm:"column:approval_id;type:varchar(64);default:null" json:"approvalId"`                    // 审核人ID
	IsApproval            int64     `gorm:"column:is_approval;type:tinyint(1);default:0" json:"isApproval"`                        // 是否需审批：0-需 1-无需
	ApprovalStatus        int64     `gorm:"column:approval_status;type:tinyint(1);default:0" json:"approvalStatus"`                // 审批状态：0-待审批 1-通过 2-驳回
	ApprovalComment       string    `gorm:"column:approval_comment;type:varchar(50);default:null" json:"approvalComment"`          // 审批意见
}

// TableName get sql table name.获取数据库表名
func (m *MoldMaintenanceTaskLifecycle) TableName() string {
	return "mold_maintenance_task_lifecycle"
}

// MoldProductResume 模具生产履历
type MoldProductResume struct {
	ID       int64  `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"-"`
	MoldCode string `gorm:"column:mold_code;type:varchar(64);not null" json:"moldCode"`  // 模具编码，冗余存储
	MoldName string `gorm:"column:mold_name;type:varchar(255);not null" json:"moldName"` // 模具名称，冗余存储
	MoldType string `gorm:"column:mold_type;type:varchar(64);not null" json:"moldType"`  // 模具类型，冗余存储
	MoldId   int64  `gorm:"column:mold_id;type:bigint(64);not null" json:"moldId"`       // 模具ID

	OrderCode    string    `gorm:"column:order_code;type:varchar(64);not null" json:"orderCode"`    // 工单号
	PartCode     string    `gorm:"column:part_code;type:varchar(64);not null" json:"partCode"`      // 零件号
	LineLevel    string    `gorm:"column:lineLevel;type:varchar(64);not null" json:"lineLevel"`     // 线别
	CompleteTime string    `gorm:"column:complete_time;type:datetime;not null" json:"completeTime"` // 生产完成时间
	Count        int64     `gorm:"column:count;type:bigint(64);not null" json:"count"`              // 工单数量
	IsDeleted    string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	IsChecked    string    `gorm:"column:is_checked;type:varchar(1);default:Y" json:"isChecked"`    // 是否已经检查过
	CreatedBy    string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy    string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
	GmtCreated   time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated   time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
}

// LineProductResume 产线生产履历
type LineProductResume struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"-"`
	LineLevel  string    `gorm:"column:line_level;type:varchar(64);not null" json:"lineLevel"`    // 线别
	ShiftId    int       `gorm:"column:shift_id;type:int(11);not null" json:"shiftId"`            // 班次ID即工单号
	ShiftMin   int       `gorm:"column:shift_min;type:int(11);not null" json:"shiftMin"`          // 班次时长（分钟）
	StartTime  string    `gorm:"column:start_time;type:datetime;not null" json:"startTime"`       // 班次开始时间
	EndTime    string    `gorm:"column:end_time;type:datetime;not null" json:"endTime"`           // 班次结束时间
	ShiftNo    int16     `gorm:"column:shift_no;type:smallint;not null" json:"shiftNo"`           // 班次号；1-早班；2-晚班
	ShiftDate  string    `gorm:"column:shift_date;type:varchar(10);not null" json:"shiftDate"`    // 班次日期
	PartCode   string    `gorm:"column:part_code;type:varchar(64);not null" json:"partCode"`      // 零件号
	QtyOk      int       `gorm:"column:qty_ok;type:int(11);not null" json:"qtyOk"`                // 生产数量
	QtyNOk     int       `gorm:"column:qty_nok;type:int(11);default:0" json:"qtyNOk"`             // 不良数量
	IsDeleted  string    `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	GmtCreated time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
}

// Message 消息
type Message struct {
	ID         int64  `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"-"`
	Status     string `gorm:"column:status;type:varchar(50);not null" json:"status"`           // 状态
	Content    string `gorm:"column:content;type:varchar(255);not null" json:"content"`        // 消息内容
	Operator   string `gorm:"column:operator;type:varchar(64);not null" json:"operator"`       // 邮件接收人
	JobId      int64  `gorm:"column:job_id;type:bigint(64);not null" json:"jobId"`             // 任务ID
	RemodelId  int64  `gorm:"column:remodel_id;type:bigint(64);not null" json:"remodelId"`     // 改造ID
	Type       string `gorm:"column:type;type:varchar(64);not null" json:"type"`               // 邮件类型：1、remodel 改造任务消息 2、maintenance保养任务
	IsDeleted  string `gorm:"column:is_deleted;type:varchar(1);default:N" json:"isDeleted"`    // 是否删除
	GmtCreated string `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated string `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy  string `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// MqttHistory 登录信息
type MqttHistory struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint(64);not null" json:"-"`
	Data       string    `gorm:"column:data;type:varchar(1024);not null" json:"data"`             // mqtt生产数据
	GmtCreated time.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated time.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`             // 创建人
	UpdatedBy  string    `gorm:"column:updated_by;type:varchar(64)" json:"updatedBy"`             // 修改人
}

// MoldStereoscopicWarehouse 立库表
type MoldStereoscopicWarehouse struct {
	ID   int64  `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"-"`
	Name string `gorm:"column:name;type:varchar(50)" json:"name"` // 立库名称
}

// MoldStereoscopicWarehouseLocation 立库储位表
type MoldStereoscopicWarehouseLocation struct {
	ID                           int64     `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"-"`
	Layer                        int       `gorm:"column:layer;type:int(11)" json:"layer"`                                                      // 层
	Col                          int       `gorm:"column:col;type:int(11)" json:"col"`                                                          // 列
	Cell                         int       `gorm:"column:cell;type:int(11)" json:"cell"`                                                        // 格
	InboundStation               string    `gorm:"column:inbound_station;type:varchar(50)" json:"inboundStation"`                               // 出库站点编号
	OutboundStation              string    `gorm:"column:outbound_station;type:varchar(50)" json:"outboundStation"`                             // 入库站点编号
	MoldID                       int64     `gorm:"column:mold_id;type:bigint(20)" json:"moldId"`                                                // 模具ID
	MoldCode                     string    `gorm:"column:mold_code;type:varchar(50)" json:"moldCode"`                                           // 模具code
	MoldSteroscopicWarehouseName string    `gorm:"column:mold_steroscopic_warehouse_name;type:varchar(50)" json:"moldSteroscopicWarehouseName"` // 立库名称
	MoldSteroscopicWarehouseID   int64     `gorm:"column:mold_steroscopic_warehouse_id;type:bigint(20)" json:"moldSteroscopicWarehouseId"`      // 立库ID
	UpdateTime                   base.Time `gorm:"column:update_time;type:datetime" json:"updateTime"`                                          // 更新时间
	Exist                        *bool     `gorm:"column:exist;type:tinyint(1)" json:"exist"`                                                   // 是否存在
	HasRfid                      *bool     `gorm:"column:has_rfid;type:tinyint(1)" json:"hasRfid"`                                              // 是否读取到RFID
	Locked                       *bool     `gorm:"column:locked;type:tinyint(1);" json:"locked"`                                                // 是否被锁定
	Category                     string    `gorm:"column:category;type:varchar(1)" json:"category"`                                             // 储位属性
}

// MoldStereoscopicWarehouseStation 立库站口表
type MoldStereoscopicWarehouseStation struct {
	ID                      int64  `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"-"`
	Layer                   int    `gorm:"column:layer;type:int(11)" json:"layer"`                                          // 层
	Col                     int    `gorm:"column:col;type:int(11)" json:"col"`                                              // 列
	Cell                    int    `gorm:"column:cell;type:int(11)" json:"cell"`                                            // 格
	Offset                  int    `gorm:"column:offset;type:int(11)" json:"offset"`                                        // 站口起始偏移
	Name                    string `gorm:"column:name;type:varchar(255)" json:"name"`                                       // 站口名称
	Code                    string `gorm:"column:code;type:varchar(255)" json:"code"`                                       // 站口物理地址编码
	Type                    string `gorm:"column:type;type:varchar(64)" json:"type"`                                        // 站口类型：出站/入站
	OutboundOffset          int    `gorm:"column:outbound_offset;type:int(11)" json:"outboundOffset"`                       // 出站口命令偏移
	InboundOffset           int    `gorm:"column:inbound_offset;type:int(11)" json:"inboundOffset"`                         // 入站口命令偏移
	StereoscopicWarehouseID int64  `gorm:"column:stereoscopic_warehouse_id;type:bigint(64)" json:"stereoscopicWarehouseId"` // 立库ID，1小库/2大库
	Category                string `gorm:"column:category;type:varchar(2)" json:"category"`                                 // 模具类别：A/B
	UniqueCode              string `gorm:"column:unique_code;type:varchar(5)" json:"uniqueCode"`                            // 站口唯一标识
	RfidDeviceID            string `gorm:"column:rfid_device_id;type:varchar(64)" json:"rfidDeviceId"`                      //RFID 设备ID
}

// MoldInoutBoundJob 出入库任务表
type MoldInoutBoundJob struct {
	ID                      int64      `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"id"`
	SrcLocation             string     `gorm:"column:src_location;type:varchar(50)" json:"srcLocation"`                      // 起始位置编码
	DestLocation            string     `gorm:"column:dest_location;type:varchar(50)" json:"destLocation"`                    // 终点位置编码
	MoldID                  int64      `gorm:"column:mold_id;type:bigint(20)" json:"moldId"`                                 // 模具ID
	MoldCode                string     `gorm:"column:mold_code;type:varchar(50)" json:"moldCode"`                            // 模具code
	PartCode                string     `gorm:"column:part_code;type:varchar(255);" json:"partCode"`                          // 零件号
	Layer                   int        `gorm:"column:layer;type:int(11)" json:"layer"`                                       // 层
	Col                     int        `gorm:"column:col;type:int(11)" json:"col"`                                           // 列
	Cell                    int        `gorm:"column:cell;type:int(11)" json:"cell"`                                         // 格
	FaultReason             string     `gorm:"column:fault_reason;type:varchar(255)" json:"faultReason"`                     // 错误原因
	Type                    string     `gorm:"column:type;type:varchar(50)" json:"type"`                                     // 类型：出库/入库
	Status                  string     `gorm:"column:status;type:varchar(50)" json:"status"`                                 // 状态
	GmtCreated              *base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"`              // 创建时间
	CreatedBy               string     `gorm:"column:created_by;type:varchar(64)" json:"createdBy"`                          // 创建人
	TaskNo                  string     `gorm:"column:task_no;type:varchar(255)" json:"taskNo"`                               // 任务编码，同agv中间表
	TaskStage               string     `gorm:"column:task_stage;type:varchar(5)" json:"taskStage"`                           // 任务阶段编码 100-Agv搬运;200-模具入站;300-模具入库; 600-模具出库; 700-模具出站;800-AGV搬运;1000-任务完成;1100-任务失败;1200-任务取消
	Order                   int        `gorm:"column:order;type:int(11)" json:"order"`                                       // 10-进行中; 20-排队中; 30-失败/完成/取消
	Reason                  string     `gorm:"column:reason;type:varchar(50)" json:"reason"`                                 // 出入库原因
	StereoscopicWarehouseID int        `gorm:"column:stereoscopic_warehouse_id;type:int(11)" json:"stereoscopicWarehouseId"` // 立库ID，1小库/2大库
	GmtCompleted            *base.Time `gorm:"column:gmt_completed;type:datetime" json:"gmtCompleted"`                       // 完成时间
	Agv                     *bool      `gorm:"column:agv;type:tinyint(1)" json:"agv"`                                        // 是否调度agv
	Rfid                    *bool      `gorm:"column:rfid;type:tinyint(1)" json:"rfid"`                                      // 是否需要用户手动传入rfid
	StationUniqueCode       string     `gorm:"column:station_unique_code;type:varchar(5)" json:"uniqueCode"`                 // 出入站口唯一标识
}

// MoldInoutBoundJobDetail 出入库任务详情表
type MoldInoutBoundJobDetail struct {
	ID          int64      `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"id"`
	JobID       int64      `gorm:"column:job_id;type:bigint(20)" json:"jobId"`                      // 出入库任务表ID
	Step        int        `gorm:"column:step;type:int(11)" json:"stepCode"`                        // 步骤代码
	StepName    string     `gorm:"column:step_name;type:varchar(255)" json:"stepName"`              // 步骤名称
	FaultReason string     `gorm:"column:fault_reason;type:varchar(255)" json:"faultReason"`        // 错误原因
	Status      string     `gorm:"column:status;type:varchar(50)" json:"status"`                    // 状态, success-成功 failure-失败 exception-异常
	GmtCreated  *base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
}

// InoutBoundTaskFlow 出入库任务流转表
type InoutBoundTaskFlow struct {
	ID                int64     `gorm:"primaryKey;column:id;type:bigint(20);not null" json:"id"`
	TaskNo            string    `gorm:"column:task_no;type:varchar(255)" json:"taskNo"`                      // 任务编码
	JobID             int64     `gorm:"column:job_id;type:bigint(20)" json:"jobId"`                          // 出入库任务表ID
	FlowOrder         int       `gorm:"column:flow_order;type:int(11)" json:"flowOrder"`                     // 任务流转的顺序
	SrcLocation       string    `gorm:"column:src_location;type:varchar(50)" json:"srcLocation"`             // 起始位置编码
	DestLocation      string    `gorm:"column:dest_location;type:varchar(50)" json:"destLocation"`           // 终点位置编码
	Type              string    `gorm:"column:type;type:varchar(50)" json:"type"`                            // 类型：出库/入库
	Agv               bool      `gorm:"column:agv;type:tinyint(1)" json:"agv"`                               // 是否调度agv
	Rfid              string    `gorm:"column:rfid;type:varchar(64)" json:"rfid"`                            // RFID
	MoldCode          string    `gorm:"column:mold_code;type:varchar(50)" json:"moldCode"`                   // 模具code
	StationUniqueCode string    `gorm:"column:station_unique_code;type:varchar(5)" json:"stationUniqueCode"` // 出入站口唯一标识
	ErrorCode         string    `gorm:"column:error_code;type:nvarchar(10)" json:"ErrorCode"`
	ErrorMsg          string    `gorm:"column:error_msg;type:nvarchar(200)" json:"ErrorMsg"`
	Flag              bool      `gorm:"column:flag;type:tinyint(1);default:0" json:"flag"`               // 流转处理标志
	GmtCreated        base.Time `gorm:"column:gmt_created;type:datetime;default:null" json:"gmtCreated"` // 创建时间
	GmtUpdated        base.Time `gorm:"column:gmt_updated;type:datetime;default:null" json:"gmtUpdated"` // 修改时间
}
