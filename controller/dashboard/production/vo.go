package dashboard

type PartQtyVo struct {
	PartNo string `json:"partNo"` //零件号
	Qty    int    `json:"qty"`    //数量
} // @name partQtyVo

type PartQtyProjectVo struct {
	PartNo      string `json:"partNo"`      //零件号
	Qty         int    `json:"qty"`         //数量
	ProjectName string `json:"projectName"` //项目名称
} // @name PartQtyProjectVo

type ShiftProductionVo struct {
	ShiftNo             int16 `json:"shiftNo"`             //班次
	ShiftMin            int   `json:"shiftMin"`            //班次时长
	MaintenanceDuration int64 `json:"maintenanceDuration"` //班次维修时长
	TotalQty            int   `json:"totalQty"`            //总数量
	PartList            []PartQtyVo
} // @name shiftProductionVo

type ShiftProductionProjectVo struct {
	ShiftNo             int16 `json:"shiftNo"`             //班次
	ShiftMin            int   `json:"shiftMin"`            //班次时长
	MaintenanceDuration int64 `json:"maintenanceDuration"` //班次维修时长
	TotalQty            int   `json:"totalQty"`            //总数量
	PartList            []PartQtyProjectVo
} // @name ShiftProductionProjectVo

type LineMaintenanceDurationVo struct {
	LineLevel           string `json:"lineLevel"`           //产线
	ShiftNo             int16  `json:"shiftNo"`             //班次
	MaintenanceDuration int64  `json:"maintenanceDuration"` //班次维修时长
} // @name LineMaintenanceDurationVo

type LineProductionVo struct {
	LineLevel string `json:"lineLevel"` //产线
	ShiftNo   int16  `json:"shiftNo"`   //班次
	ShiftMin  int    `json:"shiftMin"`  //班次时长
	PartNo    string `json:"partNo"`    //零件号
	TotalQty  int    `json:"totalQty"`  //总数量
} // @name LineProductionVo

type LineProductionProjectVo struct {
	LineLevel   string `json:"lineLevel"`   //产线
	ShiftNo     int16  `json:"shiftNo"`     //班次
	ShiftMin    int    `json:"shiftMin"`    //班次时长
	PartNo      string `json:"partNo"`      //零件号
	TotalQty    int    `json:"totalQty"`    //总数量
	ProjectName string `json:"projectName"` //项目名称
} // @name LineProductionProjectVo

type LineProductionResponseVo struct {
	LineLevel  string            `json:"lineLevel"`  //产线
	DayShift   ShiftProductionVo `json:"dayShift"`   //白班
	NightShift ShiftProductionVo `json:"nightShift"` //晚班
} // @name LineProductionResponseVo

type LineProductionProjectResponseVo struct {
	LineLevel  string                   `json:"lineLevel"`  //产线
	DayShift   ShiftProductionProjectVo `json:"dayShift"`   //白班
	NightShift ShiftProductionProjectVo `json:"nightShift"` //晚班
} // @name LineProductionProjectResponseVo

type LineProductionRequestVO struct {
	StartDate string `json:"startDate" binding:"required,datetime=2006-01-02"` //开始日期
	EndDate   string `json:"endDate" binding:"required,datetime=2006-01-02"`   //结束日期
} // @name LineProductionRequestVO

type LineProductionReportRequestVO struct {
	StartDate  string `json:"startDate" binding:"required"`  //开始日期(按周传周数按月传月数)
	EndDate    string `json:"endDate" binding:"required"`    //结束日期(按周传周数按月传月数)
	ReportCode string `json:"reportCode" binding:"required"` //报告编码 统计生产数量:ProductTotleByWeek/ProductTotleByMonth;停机工时:StopLastTimeByWeek/StopLastTimeByMonth;故障次数:ErrorTimesByWeek/ErrorTimesByMonth;维修停机率:OutageRateByWeek/OutageRateByMonth;备件消耗趋势:UserSparePartByWeek/UserSparePartByMonth;
	OrderField string `json:"orderField" binding:"required"` //排序字段 次数：count 时间：time
	LineLevel  string `json:"lineLevel"`                     //产线
} // @name 周报月报趋势图
