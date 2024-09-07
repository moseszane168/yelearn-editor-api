package constant

// 备件更换类型
const (
	REPAIR      = "repair"
	MAINTENANCE = "maintenance"
	REMODEL     = "remodel"
)

// 保养计划类型
const (
	PLAN_TIMING_TYPE   = "timing"
	PLAN_METERING_TYPE = "metering"
	PLAN_MANUAL_TYPE   = "manual"
)

// 标准类型
const (
	GENERAL = "general"
	SPECIAL = "special"
)

const USERID = "userId"

// 保养状态
const (
	WAIT     = "wait"
	PAUSE    = "pause"
	COMPLETE = "complete"
	RUNNING  = "running"
	TIMEOUT  = "timeout"
	APPROVAL = "approval"
)

// 任务状态
const (
	TASKCOMPLETE = "TaskComplete"
	TASKPAUSE    = "TaskPause"
	TASKCONTINUE = "TaskContinue"
	TASKGENERATE = "TaskGenerate"
	TASKTIMEOUT  = "TaskTimeout"
)

// PD 输送链bincode
const PdChainBinCode1 = "CS-LMSSJ-OUT-002"

// websocket业务系统ID
const (
	WebsocketPlcStackerCoordinate = "1"              // 堆垛机坐标推送
	WebsocketInboundStatus        = "1001"           // 模具入库状态推送
	WebsocketOutboundStatus       = "1002"           // 模具出库状态推送
	WebsocketInboundTasks         = "inbound-tasks"  // 模具入库任务推送
	WebsocketOutboundTasks        = "outbound-tasks" // 模具出库任务推送
	WebsocketPdInOutTasks         = "pd-tasks"       // Pd产线出入库任务推送
)

// 出入库任务状态
const (
	INOUTBOUND_STARTING = "starting"
	INOUTBOUND_TIMEOUT  = "timeout"
	INOUTBOUND_FAILURE  = "failure"
	INOUTBOUND_RUNNING  = "running"
	INOUTBOUND_PENDING  = "pending"
	INOUTBOUND_CANCELED = "canceled"
	INOUTBOUND_COMPLETE = "complete"
)

// 出入库任务排序
const (
	OrderTaskRunning = 10
	OrderTaskPending = 20
	OrderTaskDone    = 30
)

//任务阶段编码
const (
	TaskStageInboundAgvPick         = "100"
	TaskStageInboundStationRoll     = "200"
	TaskStageInboundPlcPick         = "300"
	TaskStageOutboundPlcPickAndRoll = "600"
	TaskStageOutboundAgvPick        = "700"
	TaskStageFinished               = "1000"
	TaskStageFailed                 = "1100"
	TaskStageCanceled               = "1200"
)

const (
	WMS_REQUEST_INBOUND_TASK                 = 1 //入库请求TaskType
	WMS_REQUEST_FREE_INBOUND_STATION         = 1 //入库站点请求TaskType
	WMS_AGV_ARRIVE_INBOUND_STATION           = 2 // AGV到达入站口
	WMS_CHECK_INBOUND_STATION_READY_FOR_ROLL = 3 //确认asrs是否可接收模具
	WMS_CHECK_INBOUND_STATION_RECEIVE_DONE   = 5 //确认asrs是否已接收完成

	WMS_REQUEST_OUTBOUND_TASK = 2 //出库请求TaskType

	WMS_TASK_DELETE = 101
)

//入库任务流转节点
const (
	IN_TASKFLOW_PD_TASK_REQUEST                = 1000 // PD入库信息下发
	IN_TASKFLOW_WAIT_WMS_STATION_REQUEST       = 1010 // 轮询等待WMS发起入库站点请求
	IN_TASKFLOW_CHECK_ASRS_FREE_STATION        = 1020 // 检查可入站口，选择一个站口
	IN_TASKFLOW_CALL_WMS_ALLOWED_STATION       = 1030 // Callback WMS系统，告知可入站口的BinCode，允许入库
	IN_TASKFLOW_WAIT_WMS_RECEIVE_READY_REQUEST = 1040 // 轮询等待WMS发起`确认asrs是否可接收模具`的请求
	IN_TASKFLOW_CHECK_ASRS_ROLL_READY          = 1050 // ASRS检测站口是否跟AGV对接好+入站口可入
	IN_TASKFLOW_ASRS_START_ROLL                = 1060 // ASRS发送命令到PLC入库口开始滚动
	IN_TASKFLOW_CALL_WMS_START_ROLL            = 1070 // Callback WMS系统，已开始转动链条
	IN_TASKFLOW_ASRS_READ_RFID                 = 1080 // ASRS读取RFID信号（在读触发的情况下，读取RFID，读触发消失就需要更新该节点）
	IN_TASKFLOW_WAIT_WMS_RECEIVE_DONE_REQUEST  = 1090 // 轮询等待WMS发起`确认asrs是否已接收完成`的请求
	IN_TASKFLOW_CHECK_ASRS_RECEIVE_DONE        = 1100 // ASRS检测站口是否已接收完成模具
	IN_TASKFLOW_CALL_WMS_RECEIVE_DONE          = 1110 // Callback WMS系统，已接收完成
	IN_TASKFLOW_ASRS_PLC_PRE_CHECK             = 1120 // ASRS系统入库之前PLC前置条件检测
	IN_TASKFLOW_ASRS_ASSIGN_LOCATION           = 1130 // ASRS系统分配立库库位
	IN_TASKFLOW_ASRS_PLC_WRITE_CMD             = 1140 // ASRS系统PLC写命令(包含重置通知和命令下达位、写命令前置条件检测等待写命令就绪、写入库PLC命令)
	IN_TASKFLOW_CHECK_ASRS_PLC_CMD_DONE        = 1150 // ASRS系统检测等待PLC命令执行结果
	IN_TASKFLOW_TASK_DONE                      = 1160 // 入库任务完成，ASRS系统更新储位信息，标记任务结束状态
)

var PdInboundTaskFlowOrders = []int{
	IN_TASKFLOW_WAIT_WMS_STATION_REQUEST,
	IN_TASKFLOW_CHECK_ASRS_FREE_STATION,
	IN_TASKFLOW_CALL_WMS_ALLOWED_STATION,
	IN_TASKFLOW_WAIT_WMS_RECEIVE_READY_REQUEST,
	IN_TASKFLOW_CHECK_ASRS_ROLL_READY,
	IN_TASKFLOW_ASRS_START_ROLL,
	IN_TASKFLOW_CALL_WMS_START_ROLL,
	IN_TASKFLOW_ASRS_READ_RFID,
	IN_TASKFLOW_WAIT_WMS_RECEIVE_DONE_REQUEST,
	IN_TASKFLOW_CHECK_ASRS_RECEIVE_DONE,
	IN_TASKFLOW_CALL_WMS_RECEIVE_DONE,
	IN_TASKFLOW_ASRS_PLC_PRE_CHECK,
	IN_TASKFLOW_ASRS_ASSIGN_LOCATION,
	IN_TASKFLOW_ASRS_PLC_WRITE_CMD,
	IN_TASKFLOW_CHECK_ASRS_PLC_CMD_DONE,
	IN_TASKFLOW_TASK_DONE,
}

var ManualInboundTaskFlowOrders = []int{
	IN_TASKFLOW_ASRS_READ_RFID,
	IN_TASKFLOW_ASRS_PLC_PRE_CHECK,
	IN_TASKFLOW_ASRS_ASSIGN_LOCATION,
	IN_TASKFLOW_ASRS_PLC_WRITE_CMD,
	IN_TASKFLOW_CHECK_ASRS_PLC_CMD_DONE,
	IN_TASKFLOW_TASK_DONE,
}

//出库任务流转节点
const (
	OUT_TASKFLOW_PD_TASK_REQUEST               = 2000 // PD入库信息下发,代替产线发送模具出库任务
	OUT_TASKFLOW_WAIT_WMS_OUTBOUND_REQUEST     = 2010 // 轮询等待WMS发起通知asrs出库请求
	OUT_TASKFLOW_ASRS_PLC_PRE_CHECK            = 2020 // ASRS系统出库之前PLC前置条件检测
	OUT_TASKFLOW_ASRS_PLC_WRITE_CMD            = 2030 // ASRS系统PLC写命令(包含重置通知和命令下达位、写命令前置条件检测等待写命令就绪、写出库PLC命令)
	OUT_TASKFLOW_CHECK_ASRS_PLC_CMD_DONE       = 2040 // ASRS系统检测等待PLC命令执行结果
	OUT_TASKFLOW_ASRS_UPDATE_LOCATION          = 2050 // ASRS系统模具出库完成，更新立库库位
	OUT_TASKFLOW_CALL_WMS_OUTBOUND_DONE        = 2060 // Callback WMS系统，告知模具出库完成
	OUT_TASKFLOW_WAIT_WMS_START_ROLL_REQUEST   = 2070 // 轮询等待WMS发起`通知asrs开始传输`的请求
	OUT_TASKFLOW_CHECK_ASRS_ROLL_READY         = 2080 // ASRS检测出库口是否可以开始滚动
	OUT_TASKFLOW_ASRS_START_ROLL               = 2090 // ASRS发送命令到PLC出库口开始滚动
	OUT_TASKFLOW_CALL_WMS_START_ROLL           = 2100 // Callback WMS系统，已开始转动链条
	OUT_TASKFLOW_WAIT_WMS_RECEIVE_DONE_REQUEST = 2110 // 轮询等待WMS发起`已传输完成，可停止链条`的请求
	OUT_TASKFLOW_ASRS_STOP_ROLL                = 2120 // ASRS发送命令到PLC出库口停止滚动
	OUT_TASKFLOW_CHECK_ASRS_STOP_ROLL_DONE     = 2130 // ASRS发送命令检测出站口是否已停止滚动链条中
	OUT_TASKFLOW_CALL_WMS_RECEIVE_DONE         = 2140 // Callback WMS系统，传输完成，agv可离开
	OUT_TASKFLOW_TASK_DONE                     = 2150 // 出库任务完成，标记任务结束状态
)

//带AGV对接
var WithAgvOutboundTaskFlowOrders = []int{
	OUT_TASKFLOW_WAIT_WMS_OUTBOUND_REQUEST,
	OUT_TASKFLOW_ASRS_PLC_PRE_CHECK,
	OUT_TASKFLOW_ASRS_PLC_WRITE_CMD,
	OUT_TASKFLOW_CHECK_ASRS_PLC_CMD_DONE,
	OUT_TASKFLOW_ASRS_UPDATE_LOCATION,
	OUT_TASKFLOW_CALL_WMS_OUTBOUND_DONE,
	OUT_TASKFLOW_WAIT_WMS_START_ROLL_REQUEST,
	OUT_TASKFLOW_CHECK_ASRS_ROLL_READY,
	OUT_TASKFLOW_ASRS_START_ROLL,
	OUT_TASKFLOW_CALL_WMS_START_ROLL,
	OUT_TASKFLOW_WAIT_WMS_RECEIVE_DONE_REQUEST,
	OUT_TASKFLOW_ASRS_STOP_ROLL,
	OUT_TASKFLOW_CHECK_ASRS_STOP_ROLL_DONE,
	OUT_TASKFLOW_CALL_WMS_RECEIVE_DONE,
	OUT_TASKFLOW_TASK_DONE,
}

//无AGV对接
var WithoutAgvOutboundTaskFlowOrders = []int{
	OUT_TASKFLOW_ASRS_PLC_PRE_CHECK,
	OUT_TASKFLOW_ASRS_PLC_WRITE_CMD,
	OUT_TASKFLOW_CHECK_ASRS_PLC_CMD_DONE,
	OUT_TASKFLOW_ASRS_UPDATE_LOCATION,
	OUT_TASKFLOW_TASK_DONE,
}
