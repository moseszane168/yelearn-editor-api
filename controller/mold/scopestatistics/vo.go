package scopestatistics

type MoldStatusGroupOutVO struct {
	All       int64 `json:"all"`        // 所有模具
	Normal    int64 `json:"zhengchang"` // 正常模具
	Sequester int64 `json:"fengcun"`    // 封存模具
	Fault     int64 `json:"baofei"`     // 报废模具
} // @name MoldStatusGroupOutVO

type MoldMaintenanceTaskOutVO struct {
	Completed int64 `json:"complete"` // 今天已完成保养
	Wait      int64 `json:"wait"`     // 未完成保养
	TimeOut   int64 `json:"timeout"`  // 超时未保养
	Pause     int64 `json:"pause"`    // 挂起
} // @name MoldMaintenanceTaskOutVO

type MoldExceptionOutVO struct {
	FaultDesc string `json:"faultDesc"` // 故障描述
	Cnt       int64  `json:"cnt"`       // 累计次数
	Cost      int64  `json:"cost"`      // 累计损失工时，单位分钟
} // @name MoldExceptionVO

type AgvTaskVO struct {
	Type string `json:"type"` // 统计类型
} // @name AgvTaskVO

type X struct {
	Year  int64 `json:"year"`
	Month int64 `json:"month"`
	Week  int64 `json:"week"`
	Day   int64 `json:"day"`
}

type MoldMaintenanceTask struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type MoldStatus struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type AgvTaskScopeVO struct {
	X []X     `json:"x"` // x轴
	Y []int64 `json:"y"` // Y轴
} // @name AgvTaskScopeVO

type InOutBoundTaskVO struct {
	Type      string `json:"type"`                                             // 出入库类型, inbound-入库; outbound-出库; 不传或为空-全部类型
	Platform  string `json:"platform"`                                         //平台; 不传或为空-全部平台
	Status    string `json:"status"`                                           // 任务状态, complete-成功；failure-失败; 不传或为空-全部状态
	StartDate string `json:"startDate" binding:"required,datetime=2006-01-02"` //开始日期
	EndDate   string `json:"endDate" binding:"required,datetime=2006-01-02"`   //结束日期
} // @name InOutBoundTaskVO

type InOutBoundTaskOutVO struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
} // @name InOutBoundTaskOutVO
