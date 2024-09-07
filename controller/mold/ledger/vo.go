package mold

import (
	"crf-mold/base"
)

// 模具台账
type MoldCustomVO struct {
	Key   string `json:"k"` // 字段名
	Value string `json:"v"` // 字段值g
} // @name MoldCustomVO

type CreateMoldVO struct {
	Customs   []MoldCustomVO `json:"customs"`   // 自定义参数
	PartCodes []string       `json:"partCodes"` // 模具零件号

	Code           string     `json:"code" binding:"required"`      // 模具编号
	Name           string     `json:"name" binding:"required"`      // 模具名称
	ClientName     string     `json:"clientName"`                   // 客户名称
	Weight         float64    `json:"weight"`                       // 模具重量
	ProjectName    string     `json:"projectName"`                  // 项目名称
	SizeLong       float64    `json:"sizeLong" binding:"required"`  // 模具长
	SizeWidth      float64    `json:"sizeWidth" binding:"required"` // 模具宽
	SizeHeigh      float64    `json:"sizeHeigh" binding:"required"` // 模具高
	Platform       string     `json:"platform" binding:"required"`  // 平台
	Process        string     `json:"process" binding:"required"`   // 工序
	PropertyNumber string     `json:"propertyNumber"`               // 资产编号
	Type           string     `json:"type" binding:"required"`      // 模具类型
	Status         string     `json:"status" binding:"required"`    // 模具状态
	MakeDate       *base.Time `json:"makeDate"`                     // 制造日期
	FirstUseDate   *base.Time `json:"firstUseDate"`                 // 第一次使用日期
	YearLimit      float64    `json:"yearLimit"`                    // 使用年限
	Rfid           string     `json:"rfid"`                         // RFID
	LineLevel      string     `json:"lineLevel" binding:"required"` // 使用线别
	Provide        string     `json:"provide" binding:"required"`   // 模具供应商
	Category       string     `json:"category"`                     // 模具类别
} // @name CreateMoldVO

type UpdateMoldVO struct {
	Customs   []MoldCustomVO `json:"customs"`   // 自定义参数
	PartCodes []string       `json:"partCodes"` // 模具零件号

	ID             int64      `json:"id" binding:"required"`        // ID
	Code           string     `json:"code" binding:"required"`      // 模具编号
	Name           string     `json:"name" binding:"required"`      // 模具名称
	ClientName     string     `json:"clientName"`                   // 客户名称
	Weight         float64    `json:"weight"`                       // 模具重量
	ProjectName    string     `json:"projectName"`                  // 项目名称
	SizeLong       float64    `json:"sizeLong" binding:"required"`  // 模具长
	SizeWidth      float64    `json:"sizeWidth" binding:"required"` // 模具宽
	SizeHeigh      float64    `json:"sizeHeigh" binding:"required"` // 模具高
	Platform       string     `json:"platform" binding:"required"`  // 平台
	Process        string     `json:"process" binding:"required"`   // 工序
	PropertyNumber string     `json:"propertyNumber"`               // 资产编号
	Type           string     `json:"type" binding:"required"`      // 模具类型
	Status         string     `json:"status" binding:"required"`    // 模具状态
	MakeDate       *base.Time `json:"makeDate"`                     // 制造日期
	FirstUseDate   *base.Time `json:"firstUseDate"`                 // 第一次使用日期
	YearLimit      float64    `json:"yearLimit"`                    // 使用年限
	Rfid           string     `json:"rfid"`                         // RFID
	LineLevel      string     `json:"lineLevel" binding:"required"` // 使用线别
	Provide        string     `json:"provide" binding:"required"`   // 模具供应商
	Category       string     `json:"category"`                     // 模具类别
} // @name UpdateMoldVO

type PageMoldVO struct {
	base.PageVO
	PartCodes string `form:"partCodes" json:"partCodes"` // 模具零件号

	CodeOrName string `form:"codeOrName"` // 模具编号和名称模糊查找，该参数有值时忽略其它条件参数

	Code              string     `form:"code"`              // 模具编号
	Name              string     `form:"name"`              // 模具名称
	ClientName        string     `form:"clientName"`        // 客户名称
	Weight            float64    `form:"weight"`            // 模具重量
	ProjectName       string     `form:"projectName"`       // 项目名称
	SizeLong          float64    `form:"sizeLong"`          // 模具长
	SizeWidth         float64    `form:"sizeWidth"`         // 模具宽
	SizeHeigh         float64    `form:"sizeHeigh"`         // 模具高
	Platform          string     `form:"platform"`          // 平台
	Process           string     `form:"process"`           // 工序
	PropertyNumber    string     `form:"propertyNumber"`    // 资产编号
	Type              string     `form:"type"`              // 模具类型
	Status            string     `form:"status"`            // 模具状态
	MakeDateBegin     *base.Time `form:"makeDateBegin"`     // 制造日期
	MakeDateEnd       *base.Time `form:"makeDateEnd"`       // 制造日期
	FirstUseDateBegin *base.Time `form:"firstUseDateBegin"` // 第一次使用日期
	FirstUseDateEnd   *base.Time `form:"firstUseDateEnd"`   // 第一次使用日期
	YearLimit         *float64   `form:"yearLimit"`         // 使用年限
	Rfid              string     `form:"rfid"`              // RFID TID
	LineLevel         string     `form:"lineLevel"`         // 使用线别
	Provide           string     `form:"provide"`           // 模具供应商
	Category          string     `form:"category"`          // 模具类别
	ExcludeType       string     `form:"excludeType"`       // 不包含的模具类别 chongkong,chengxing
} // @name PageMoldVO

type PageMoldOutVO struct {
	PartCodes   []string    `form:"partCodes" json:"partCodes"` // 模具零件号
	PartCodeStr string      `json:"partCodeStr"`                // 模具零件号, '/'拼接
	Customs     interface{} `json:"customs"`                    // 自定义参数

	ID             int64      `form:"id" json:"id"`                         // 模具编号
	Code           string     `form:"code" json:"code"`                     // 模具编号
	Name           string     `form:"name" json:"name"`                     // 模具名称
	ClientName     string     `form:"clientName" json:"clientName"`         // 客户名称
	Weight         float64    `form:"weight" json:"weight"`                 // 模具重量
	ProjectName    string     `form:"projectName" json:"projectName"`       // 项目名称
	SizeLong       float64    `form:"sizeLong" json:"sizeLong"`             // 模具长
	SizeWidth      float64    `form:"sizeWidth" json:"sizeWidth"`           // 模具宽
	SizeHeigh      float64    `form:"sizeHeigh" json:"sizeHeigh"`           // 模具高
	Platform       string     `form:"platform" json:"platform"`             // 平台
	Process        string     `form:"process" json:"process"`               // 工序
	PropertyNumber string     `form:"propertyNumber" json:"propertyNumber"` // 资产编号
	Type           string     `form:"type" json:"type"`                     // 模具类型
	Status         string     `form:"status" json:"status"`                 // 模具状态
	MakeDate       *base.Time `form:"makeDate" json:"makeDate"`             // 制造日期
	FirstUseDate   *base.Time `form:"firstUseDate" json:"firstUseDate"`     // 第一次使用日期
	YearLimit      float64    `form:"yearLimit" json:"yearLimit"`           // 使用年限
	Rfid           string     `form:"rfid" json:"rfid"`                     // RFID
	LineLevel      string     `form:"lineLevel" json:"lineLevel"`           // 使用线别
	Provide        string     `form:"provide" json:"provide"`               // 模具供应商
	FlushCount     int64      `json:"flushCount"`                           // 模具冲次
	Category       string     `json:"category"`                             // 模具类别
} // @name PageMoldOutVO

type PageMoldOutV1 struct {
	PageMoldOutVO
	HavePlan bool `json:"have_plan"` // 是否含有维保计划 在分页查询模具列表API中使用
} // @name PageMoldOutV1

type MoldExportVO struct {
	Code           string     `excel:"模具编号"`                          // 模具编号
	Name           string     `excel:"模具名称"`                          // 模具名称
	ClientName     string     `excel:"客户名称" dict:"clientName"`        // 客户名称
	Weight         float64    `excel:"模具重量"`                          // 模具重量
	ProjectName    string     `excel:"项目名称"`                          // 项目名称
	SizeLong       float64    `excel:"模具长"`                           // 模具长
	SizeWidth      float64    `excel:"模具宽"`                           // 模具宽
	SizeHeigh      float64    `excel:"模具高"`                           // 模具高
	Platform       string     `excel:"平台" dict:"platform"`            // 平台
	Process        string     `excel:"工序" dict:"operation"`           // 工序
	PropertyNumber string     `excel:"资产编号"`                          // 资产编号
	Type           string     `excel:"模具类型" dict:"moldType"`          // 模具类型
	Status         string     `excel:"模具状态" dict:"moldStatus"`        // 模具状态
	MakeDate       *base.Time `excel:"制造日期"`                          // 制造日期
	FirstUseDate   *base.Time `excel:"第一次使用日期"`                       // 第一次使用日期
	YearLimit      float64    `excel:"使用年限"`                          // 使用年限
	Rfid           string     `excel:"RFID"`                          // RFID
	LineLevel      string     `excel:"使用线别" dict:"line"`              // 使用线别
	Provide        string     `excel:"模具供应商" dict:"moldManufacturer"` // 模具供应商
	Category       string     `excel:"模具属性"`                          // 模具属性

	PartCodes string `excel:"模具零件号"` // 模具零件号
} // @name MoldExportVO

type MoldImportVO struct {
	PartCodes      string     `json:"partCodes" excel:"模具零件号"`                                                       // 模具零件号
	Code           string     `json:"code" excel:"模具编号" validate:"required,unique"`                                  // 模具编号
	Name           string     `json:"name" excel:"模具名称" validate:"required"`                                         // 模具名称
	ClientName     string     `json:"clientName" excel:"客户名称" validate:"required,enum=clientName" dict:"clientName"` // 客户名称
	Weight         *float64   `json:"weight" excel:"模具重量"`                                                           // 模具重量
	ProjectName    string     `json:"projectName" excel:"项目名称"`                                                      // 项目名称
	SizeLong       float64    `json:"sizeLong" excel:"模具长" validate:"required"`                                      // 模具长
	SizeWidth      float64    `json:"sizeWidth" excel:"模具宽" validate:"required"`                                     // 模具宽
	SizeHeigh      float64    `json:"sizeHeigh" excel:"模具高" validate:"required"`                                     // 模具高
	Platform       string     `json:"platform" excel:"平台" validate:"required,enum=platform" dict:"platform"`         // 平台
	Process        string     `json:"process" excel:"工序" validate:"required,enum=operation" dict:"operation"`        // 工序
	PropertyNumber string     `json:"propertyNumber" excel:"资产编号"`                                                   // 资产编号
	Type           string     `json:"type" excel:"模具类型" validate:"required,enum=moldType" dict:"moldType"`           // 模具类型
	Status         string     `json:"status" excel:"模具状态" validate:"required,enum=moldStatus" dict:"moldStatus"`     // 模具状态
	MakeDate       *base.Time `json:"makeDate" excel:"制造日期"`                                                         // 制造日期
	FirstUseDate   *base.Time `json:"firstUseDate" excel:"第一次使用日期"`                                                  // 第一次使用日期
	// YearLimit      float64   `json:"yearLimit" excel:"使用年限" validate:"required"`                                              // 使用年限
	Rfid      string `json:"rfid" excel:"RFID"`                                                                       // RFID
	LineLevel string `json:"lineLevel" excel:"使用线别" validate:"required,enum=line" dict:"line"`                        // 使用线别
	Provide   string `json:"provide" excel:"模具供应商" validate:"required,enum=moldManufacturer" dict:"moldManufacturer"` // 模具供应商
	Category  string `json:"category" excel:"模具属性" validate:"required,enum=moldProperties" dict:"moldProperties"`     // 模具属性
} // @name MoldImportVO

type MoldImportParsedVOS struct {
	Data []MoldImportVO `json:"data"` // 导入数据
} // @name MoldImportParsedVOS

type MoldFlushCountVO struct {
	ID       int64  `json:"id" binding:"required"` // 模具id
	PassWord string `json:"password"`              // 管理员密码，加密
}

// 模具文档库
type SaveMoldDocVO struct {
	FileKey  string `json:"fileKey" binding:"required"`  // 文件key
	MoldCode string `json:"moldCode" binding:"required"` // 模具code
	Name     string `json:"name" binding:"required"`     // 文档名称
	Version  string `json:"version" binding:"required"`  // 版本
	Content  string `json:"content" binding:"required"`  // 变更内容
	Remark   string `json:"remark"`                      // 备注
} // @name SaveMoldDocVO

// 查询零件号
type PartCodeQueryMoldVO struct {
	PartCode string `json:"partCode" binding:"required"` // 零件号
} // @name PartCodeQueryMoldVO

// 模具BOM
type CreateBomVO struct {
	StandardComponent string `json:"standardComponent"` // 标准件:Y 非标准件：N
	MoldCode          string `json:"moldCode"`          // 模具编号
	PartName          string `json:"partName"`          // 零件名称
	PartCode          string `json:"partCode"`          // 零件号
	Count             int    `json:"count"`             // 数量
	Material          string `json:"material"`          // 材料
	Flavor            string `json:"flavor"`            // 规格
	Standard          string `json:"standard"`          // 标准
	HeatTreating      string `json:"heatTreating"`      // 热处理
	AssemblyRelation  string `json:"assemblyRelation"`  // 装配关系
	Remark            string `json:"remark"`            // 备注
} // @name CreateBomVO

type UpdateBomVO struct {
	ID               int64  `json:"id"`
	MoldCode         string `json:"moldCode"`         // 模具编号
	PartName         string `json:"partName"`         // 零件名称
	PartCode         string `json:"partCode"`         // 零件号
	Count            int    `json:"count"`            // 数量
	Material         string `json:"material"`         // 材料
	Flavor           string `json:"flavor"`           // 规格
	Standard         string `json:"standard"`         // 标准
	HeatTreating     string `json:"heatTreating"`     // 热处理
	AssemblyRelation string `json:"assemblyRelation"` // 装配关系
	Remark           string `json:"remark"`           // 备注
} // @name UpdateBomVO

type PageBomVO struct {
	base.PageVO
	CodeOrName string `form:"codeOrName"` // 模具编号和名称模糊查找，该参数有值时忽略其它条件参数

	StandardComponent string `form:"standardComponent"` // 标准件:Y 非标准件：N
	MoldCode          string `form:"moldCode"`          // 模具编号
	PartName          string `form:"partName"`          // 零件名称
	PartCode          string `form:"partCode"`          // 零件号
	Count             int    `form:"count"`             // 数量
	Material          string `form:"material"`          // 材料
	Flavor            string `form:"flavor"`            // 规格
	Standard          string `form:"standard"`          // 标准
	HeatTreating      string `form:"heatTreating"`      // 热处理
	AssemblyRelation  string `form:"assemblyRelation"`  // 装配关系
	Remark            string `form:"remark"`            // 备注
} // @name PageBomVO

type PageBomOutVO struct {
	ID                int64  `json:"id"`
	StandardComponent string `json:"standardComponent"` // 标准件:Y 非标准件：N
	MoldCode          string `json:"moldCode"`          // 模具编号
	PartName          string `json:"partName"`          // 零件名称
	PartCode          string `json:"partCode"`          // 零件号
	Count             int    `json:"count"`             // 数量
	Material          string `json:"material"`          // 材料
	Flavor            string `json:"flavor"`            // 规格
	Standard          string `json:"standard"`          // 标准
	HeatTreating      string `json:"heatTreating"`      // 热处理
	AssemblyRelation  string `json:"assemblyRelation"`  // 装配关系
	Remark            string `json:"remark"`            // 备注
} // @name PageBomOutVO

type BomExportVO struct {
	MoldCode         string `json:"moldCode" excel:"模具编号"`         // 模具编号
	PartName         string `json:"partName" excel:"零件名称"`         // 零件名称
	PartCode         string `json:"partCode" excel:"零件号"`          // 零件号
	Count            int    `json:"count" excel:"数量"`              // 数量
	Material         string `json:"material" excel:"材料"`           // 材料
	Flavor           string `json:"flavor" excel:"规格"`             // 规格
	Standard         string `json:"standard" excel:"标准"`           // 标准
	HeatTreating     string `json:"heatTreating" excel:"热处理"`      // 热处理
	AssemblyRelation string `json:"assemblyRelation" excel:"装配关系"` // 装配关系
	Remark           string `json:"remark" excel:"备注"`             // 备注
} // @name BomExportVO

type BomImportVO struct {
	MoldCode         string `json:"moldCode" validate:"required" excel:"模具编号"` // 模具编号
	PartName         string `json:"partName" excel:"零件名称" validate:"required"` // 零件名称
	PartCode         string `json:"partCode" excel:"零件号" validate:"required"`  // 零件号
	Count            int    `json:"count" excel:"数量" validate:"required"`      // 数量
	Material         string `json:"material" excel:"材料"`                       // 材料
	Flavor           string `json:"flavor" excel:"规格" validate:"required"`     // 规格
	Standard         string `json:"standard" excel:"标准"`                       // 标准
	HeatTreating     string `json:"heatTreating" excel:"热处理"`                  // 热处理
	AssemblyRelation string `json:"assemblyRelation" excel:"装配关系"`             // 装配关系
	Remark           string `json:"remark" excel:"备注"`                         // 备注
} // @name BomImportVO

type BomImportParsedVOS struct {
	Standard string        `json:"standard"` // 是否标准件
	Data     []BomImportVO `json:"data"`     // 导入数据
} // @name BomImportParsedVOS

//
// 模具成型参数
//

type MoldingCustomVO struct {
	ID    int64  `json:"id"`
	Key   string `json:"k"` // 字段名
	Value string `json:"v"` // 字段值g
} // @name MoldCustomVO

type SaveMoldingParamsVO struct {
	MoldCode        string  `json:"moldCode"`        // 模具编号
	Height          float64 `json:"height"`          // 模具高度
	DefensiveDevice string  `json:"defensiveDevice"` // 防错装置
	Desc            string  `json:"desc"`            // 描述
} // @name SaveMoldingParamsVO

type PageMoldingParamsOutVO struct {
	MoldCode        string            `json:"moldCode"`        // 模具编号
	Height          float64           `json:"height"`          // 模具高度
	DefensiveDevice string            `json:"defensiveDevice"` // 防错装置
	Desc            string            `json:"desc"`            // 描述
	Customs         []MoldingCustomVO `json:"customs"`         // 自定义字段
} // @name PageMoldingParamsOutVO

type CreateMoldingParamsCustomsVO struct {
	MoldCode string            `json:"moldCode"` // 模具编码
	Customs  []MoldingCustomVO `json:"customs"`  // 数据
} // @name CreateMoldingParamsCustomsVO

// 保养履历
type PageMoldMaintenanceOutVO struct {
	Id               int64  `json:"id"`               // 任务ID
	Code             string `json:"code"`             // 任务编号
	MaintenanceLevel string `json:"maintenanceLevel"` // 保养级别
	StandardName     string `json:"standardName"`     // 标准名称
	Time             string `json:"time"`             // 保养时间
	UpdatedBy        string `json:"updatedBy"`        // 保养人
} // @name PageMoldMaintenanceOutVO

// 备件更换履历
type PageReplaceSpareVO struct {
	SpareCode string     `json:"spareCode"` // 备件编码
	SpareName string     `json:"spareName"` // 备件名称
	Reason    string     `json:"reason"`    // 更换原因
	Time      *base.Time `json:"time"`      // 更换时间
	Count     int64      `json:"count"`     // 数量
	Code      string     `json:"code"`      // 维修/保养单号
} // @name PageReplaceSpareVO
