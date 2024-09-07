package dashboard

type OutageFactorVo struct {
	LineLevel string  `json:"lineLevel"` //产线编号
	ShiftNo   string  `json:"shiftNo"`   // 1白班2晚班
	LastTime  int     `json:"LastTime"`  //停机时长
	RateNum   float32 `json:"rateNum"`   //维修停机率
} //维修停机率

type OutageFactorPercentVo struct {
	OutageFactorVo
	Rate string `json:"rate"` //维修停机率百分比
} //维修停机率

type RejectsProducts struct {
	LineLevel   string `json:"lineLevel"`   //产线编号
	DefectDesc  string `json:"defectDesc"`  //不良描述
	DefectCount int    `json:"defectCount"` //不良数量
} //质量

type MaintainList struct {
	LineLevel     string `json:"lineLevel"`     //产线编号
	RepairContent string `json:"repairContent"` //故障类别
	ShiftNo       string `json:"shiftNo"`       //1白班2晚班
	LastTime      int    `json:"lastTime"`      //维修时长
} //维修

type MonthOutageFactorVo struct {
	LineLevel  string  `json:"lineLevel"`  //产线编号
	GmtCreated string  `json:"gmtCreated"` //日期
	RateNum    float32 `json:"rateNum"`    //维修停机率
} //维修停机率-月报

type MonthOutageFactorAll struct {
	MofvList []MonthOutageFactorVo `json:"mofvList"` //维修停机率数组
	DateList []string              `json:"dateList"` //日期数组
} //维修停机率-月报返回结构

type ProductTotle struct {
	LineLevel string `json:"lineLevel"` //产线编号
	ShiftDate string `json:"shiftDate"` //班次日期
	TotleNum  int    `json:"totleNum"`  //生产数量
} //生产数量统计列表

type ProductTotleReport struct {
	ProductList []ProductTotle `json:"productList"` //生产数量统计列表
	DateList    []string       `json:"dateList"`    //日期数组
} //生产趋势

type UserSparePartVo struct {
	DateTime string `json:"dateTime"` //日期
	TotleNum int    `json:"totleNum"` //消耗数量
} //备件消耗

type UserSparePartList struct {
	UspvList []UserSparePartVo `json:"uspvList"` //备件消耗
	DateList []string          `json:"dateList"` //日期数组
} //备件消耗结构体

type ErrorTopTens struct {
	FaultDesc string `json:"faultDesc"` //故障描述/类型
	CountNum  int    `json:"countNum"`  //统计次数
	LastTime  int    `json:"lastTime"`  //统计时长
} //故障统计前10

type TendencyChartVo struct {
	LineLevel    string  `json:"lineLevel"`    //产线编号
	ReportDate   string  `json:"reportDate"`   //x轴
	ReportValue1 float32 `json:"reportValue1"` //y1轴
	ReportValue2 float32 `json:"reportValue2"` //y1轴
	Value1Desc   string  `json:"value1Desc"`   //y1轴说明
	Value2Desc   string  `json:"value2Desc"`   //y2轴说明
} //周、月趋势图

type LineList struct {
	LineLevel string `json:"lineLevel"` //产线编号
	LineName  string `json:"LineName"`  //产线编号
} //产线列表
