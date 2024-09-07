/**
 * 统一响应体
 */

package base

type ResponseCode string

type Response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Result  interface{}  `json:"result"`
}

const (
	SUCCESS             ResponseCode = "0"
	BUSY                ResponseCode = "100"
	TIMEOUT             ResponseCode = "300"
	PARAMS_ERROR        ResponseCode = "400"
	UNAUTHORIZED        ResponseCode = "401"
	NO_PERSSION         ResponseCode = "403"
	NETWORK_EXCEPTION   ResponseCode = "404"
	DATA_NOT_EXIST      ResponseCode = "600"
	UNKNOW_ERROR        ResponseCode = "999"
	JSON_VALIDATE_ERROR ResponseCode = "777"

	RESOURCE_EXIST         ResponseCode = "1000"
	USER_OR_PASSWORD_ERROR ResponseCode = "1001"
	USER_ID_REPEAT         ResponseCode = "1002"
	PASSWORD_FORMAT_ERROR  ResponseCode = "1003"
	USER_NOT_EXIST         ResponseCode = "1004"
	USER_LOCKED            ResponseCode = "1005"
	DICT_VALUE_EXIST       ResponseCode = "1006"

	MOLD_CODE_EXIST                  ResponseCode = "2000"
	MOLD_FLUSH_COUNT_ZERO_ONLY_ADMIN ResponseCode = "2001"
	MOLD_EXIST_IN_ASRS               ResponseCode = "2002"

	SPARE_CODE_EXIST         ResponseCode = "3000"
	SPARE_REQUEST_NOT_ENOUGH ResponseCode = "3001"

	MOLD_DOC_EXIST ResponseCode = "4000"

	KNOWLEDGE_EXIST           ResponseCode = "5000"
	ELASTICSEARCH_SAVE_FAULRE ResponseCode = "5001"

	MOLD_BOM_EXIST     ResponseCode = "6000"
	MOLD_BOM_NOT_EXIST ResponseCode = "6001"

	MOLD_MOLDING_CUSTOM_KEY_REPEAT ResponseCode = "7000"

	WITHDRAW_CAN_NOT_EMPTY          ResponseCode = "8000"
	REMODEL_END_TIME_NOT_BEFORE_NOW ResponseCode = "8001"

	WS_ONLINE_MESSAGE_CODE      ResponseCode = "9001"
	WS_OFFLINE_MESSAGE_CODE     ResponseCode = "9002"
	WS_BUSINESS_SYSTEM_ID_ERROR ResponseCode = "9003"

	ADMIN_CAN_NOT_DELETE ResponseCode = "10000"

	MAINTENANCE_NAME_EXIST                       ResponseCode = "11000"
	MAINTENANCE_PLAN_NAME_EXIST                  ResponseCode = "11001"
	MAINTENANCE_PLAN_TYPE_EXIST                  ResponseCode = "11002"
	MAINTENANCE_STANDARD_GENERAL_EXIST           ResponseCode = "11003"
	MAINTENANCE_TASK_HUANG_UP_CAN_NOT_THAN_THREE ResponseCode = "11004"
	MAINTENANCE_TASK_CONTENT_UNCOMPLETE          ResponseCode = "11005"
	MAINTENANCE_STANDARD_CONTENT_CAN_NOT_EMPTY   ResponseCode = "11006"
	MAINTENANCE_STANDARD_REL_MOLD_CAN_NOT_EMPTY  ResponseCode = "11007"
	MAINTENANCE_PLAN_NOT_RUNNING                 ResponseCode = "11008"
	MAINTENANCE_PLAN_NOT_EXIST                   ResponseCode = "11009"
	MAINTENANCE_CHARGE_FAILED                    ResponseCode = "11010"
	MAINTENANCE_MANUAL_ADD_FAILED                ResponseCode = "11011"
	MAINTENANCE_PLAN_OPERATE_NOT_ALLOWED         ResponseCode = "11012"
	MAINTENANCE_PLAN_CRON_INVALID                ResponseCode = "11013"

	EMAIL_UNOKNOW_TYPE ResponseCode = "12000"
	EMAIL_FORMAT_ERROR ResponseCode = "12001"

	STEREOSCOPIC_EXIST_RUNNING_INOUTBOUND_TASK ResponseCode = "13000"
	STEREOSCOPIC_MOLD_NOT_EXIST                ResponseCode = "13001"
	MOLD_STATION_NOT_EXIST                     ResponseCode = "13002"
	STEREOSCOPIC_NO_EMPTY_LOCATION             ResponseCode = "13003"
	PLEASE_NOT_REPEAT_EXECUTION                ResponseCode = "13004"
	JOB_NOT_EXIST                              ResponseCode = "13005"
	OUTBOUND_PLC_PRE_CHECK_FAILED              ResponseCode = "13006"
	OUTBOUND_PENDING_TASK_COUNT_LIMITED        ResponseCode = "13007"
	OUTBOUND_CANCEL_PENDING_TASK_FAILED        ResponseCode = "13008"
	OUTBOUND_STOP_RUNNING_TASK_FAILED          ResponseCode = "13009"
	OUTBOUND_DUPLICATED_TASK                   ResponseCode = "13010"
	INBOUND_PLC_PRE_CHECK_FAILED               ResponseCode = "13011"
	WEB_INBOUND_FEATURE_NOT_ALLOWED            ResponseCode = "13012"
	OUTBOUND_AGV_NOT_ALLOWED                   ResponseCode = "13013"
	WEB_OUTBOUND_FEATURE_NOT_ALLOWED           ResponseCode = "13014"
	WEB_INOUTBOUND_FEATURE_NOT_ALLOWED         ResponseCode = "13015"

	MOLD_WAREHOUSE_MANAGE_CATEGORY_CONFLICT ResponseCode = "14000"
	MOLD_WAREHOUSE_MOLD_CODE_DUPLICATED     ResponseCode = "14001"
)

var ResponseEnum = map[ResponseCode]*Response{
	SUCCESS:             {SUCCESS, "请求成功", nil},
	BUSY:                {BUSY, "系统繁忙", nil},
	TIMEOUT:             {TIMEOUT, "请求超时", nil},
	PARAMS_ERROR:        {PARAMS_ERROR, "参数错误", nil},
	UNAUTHORIZED:        {UNAUTHORIZED, "未授权访问,请先登录", nil},
	NO_PERSSION:         {NO_PERSSION, "请求拒绝,无权限访问", nil},
	NETWORK_EXCEPTION:   {NETWORK_EXCEPTION, "网络异常", nil},
	DATA_NOT_EXIST:      {DATA_NOT_EXIST, "数据不存在", nil},
	UNKNOW_ERROR:        {UNKNOW_ERROR, "未知错误", nil},
	JSON_VALIDATE_ERROR: {JSON_VALIDATE_ERROR, "json格式校验失败", nil},

	RESOURCE_EXIST:         {RESOURCE_EXIST, "资源已存在", nil},
	USER_OR_PASSWORD_ERROR: {USER_OR_PASSWORD_ERROR, "账号或密码错误", nil},
	USER_ID_REPEAT:         {USER_ID_REPEAT, "工号重复", nil},
	PASSWORD_FORMAT_ERROR:  {PASSWORD_FORMAT_ERROR, "密码格式错误", nil},
	USER_NOT_EXIST:         {USER_NOT_EXIST, "用户不存在", nil},
	USER_LOCKED:            {USER_LOCKED, "账号已被锁定，请联系管理员重置密码", nil},
	DICT_VALUE_EXIST:       {DICT_VALUE_EXIST, "字典值已存在", nil},

	// 模具
	MOLD_CODE_EXIST:                  {MOLD_CODE_EXIST, "模具编码已存在", nil},
	MOLD_FLUSH_COUNT_ZERO_ONLY_ADMIN: {MOLD_FLUSH_COUNT_ZERO_ONLY_ADMIN, "模具冲次置零仅限管理员操作", nil},

	// 备件
	SPARE_CODE_EXIST:         {SPARE_CODE_EXIST, "备件编码已存在", nil},
	SPARE_REQUEST_NOT_ENOUGH: {SPARE_REQUEST_NOT_ENOUGH, "备件库存不足,无法出库", nil},

	// 模具文档
	MOLD_DOC_EXIST: {MOLD_DOC_EXIST, "无法修改,对应名称和版本的文档已经存在", nil},

	// 知识库
	KNOWLEDGE_EXIST:           {KNOWLEDGE_EXIST, "指定名称的知识库已存在", nil},
	ELASTICSEARCH_SAVE_FAULRE: {ELASTICSEARCH_SAVE_FAULRE, "保存到ES失败", nil},

	// 模具BOM
	MOLD_BOM_EXIST:     {MOLD_BOM_EXIST, "模具BOM已存在", nil},
	MOLD_BOM_NOT_EXIST: {MOLD_BOM_NOT_EXIST, "模具BOM不存在", nil},

	// 模具成型参数
	MOLD_MOLDING_CUSTOM_KEY_REPEAT: {MOLD_MOLDING_CUSTOM_KEY_REPEAT, "模具成型参数自定义字段不能重复", nil},

	// 模具改造
	WITHDRAW_CAN_NOT_EMPTY:          {WITHDRAW_CAN_NOT_EMPTY, "撤销原因不能为空", nil},
	REMODEL_END_TIME_NOT_BEFORE_NOW: {REMODEL_END_TIME_NOT_BEFORE_NOW, "改造周期结束时间不能小于当前时间", nil},

	// 用户
	ADMIN_CAN_NOT_DELETE: {ADMIN_CAN_NOT_DELETE, "不能删除Admin用户", nil},

	// 保养
	MAINTENANCE_NAME_EXIST:      {MAINTENANCE_NAME_EXIST, "标准名称已存在", nil},
	MAINTENANCE_PLAN_NAME_EXIST: {MAINTENANCE_PLAN_NAME_EXIST, "模具计划名称已存在", nil},
	MAINTENANCE_PLAN_TYPE_EXIST: {MAINTENANCE_PLAN_TYPE_EXIST, "对应模具类型的模具计划已存在", nil},

	MAINTENANCE_STANDARD_GENERAL_EXIST: {MAINTENANCE_STANDARD_GENERAL_EXIST, "对应模具类型和线别的通用标准已存在", nil},

	MAINTENANCE_TASK_HUANG_UP_CAN_NOT_THAN_THREE: {MAINTENANCE_TASK_HUANG_UP_CAN_NOT_THAN_THREE, "保养任务挂起次数不能超过三次", nil},
	MAINTENANCE_TASK_CONTENT_UNCOMPLETE:          {MAINTENANCE_TASK_CONTENT_UNCOMPLETE, "存在未完成的保养内容,无法提交", nil},

	MAINTENANCE_STANDARD_CONTENT_CAN_NOT_EMPTY:  {MAINTENANCE_STANDARD_CONTENT_CAN_NOT_EMPTY, "标准详情不能为空", nil},
	MAINTENANCE_STANDARD_REL_MOLD_CAN_NOT_EMPTY: {MAINTENANCE_STANDARD_REL_MOLD_CAN_NOT_EMPTY, "专用标准关联模具不能为空", nil},
	MAINTENANCE_PLAN_NOT_RUNNING:                {MAINTENANCE_PLAN_NOT_RUNNING, "对应模具类型的模具计划未运行", nil},
	MAINTENANCE_PLAN_NOT_EXIST:                  {MAINTENANCE_PLAN_NOT_EXIST, "对应模具类型的模具计划不存在", nil},
	MAINTENANCE_CHARGE_FAILED:                   {MAINTENANCE_CHARGE_FAILED, "普通权限只能指派一次", nil},
	MAINTENANCE_MANUAL_ADD_FAILED:               {MAINTENANCE_MANUAL_ADD_FAILED, "只允许给冲裁/成型模具做手动新增保养单", nil},
	MAINTENANCE_PLAN_OPERATE_NOT_ALLOWED:        {MAINTENANCE_PLAN_OPERATE_NOT_ALLOWED, "不允许在线边服务器上操作保养计划", nil},
	MAINTENANCE_PLAN_CRON_INVALID:               {MAINTENANCE_PLAN_CRON_INVALID, "计时保养计划Cron表达式无效", nil},

	// 邮件
	EMAIL_UNOKNOW_TYPE: {EMAIL_UNOKNOW_TYPE, "未知的类型", nil},
	EMAIL_FORMAT_ERROR: {EMAIL_FORMAT_ERROR, "邮箱格式错误", nil},

	// 立库出入库

	STEREOSCOPIC_EXIST_RUNNING_INOUTBOUND_TASK: {STEREOSCOPIC_EXIST_RUNNING_INOUTBOUND_TASK, "存在正在执行中的出入库任务", nil},
	STEREOSCOPIC_MOLD_NOT_EXIST:                {STEREOSCOPIC_MOLD_NOT_EXIST, "立库中模具不存在", nil},
	MOLD_STATION_NOT_EXIST:                     {MOLD_STATION_NOT_EXIST, "立库站口不存在", nil},
	STEREOSCOPIC_NO_EMPTY_LOCATION:             {STEREOSCOPIC_NO_EMPTY_LOCATION, "立库不存在空闲的库位", nil},
	PLEASE_NOT_REPEAT_EXECUTION:                {PLEASE_NOT_REPEAT_EXECUTION, "请勿重复执行", nil},
	JOB_NOT_EXIST:                              {JOB_NOT_EXIST, "任务不存在", nil},
	OUTBOUND_PLC_PRE_CHECK_FAILED:              {OUTBOUND_PLC_PRE_CHECK_FAILED, "立库PLC前置检测不通过", nil},
	OUTBOUND_PENDING_TASK_COUNT_LIMITED:        {OUTBOUND_PENDING_TASK_COUNT_LIMITED, "超过立库待处理任务数量限制20个", nil},
	OUTBOUND_CANCEL_PENDING_TASK_FAILED:        {OUTBOUND_CANCEL_PENDING_TASK_FAILED, "取消待执行的出库任务失败", nil},
	OUTBOUND_STOP_RUNNING_TASK_FAILED:          {OUTBOUND_STOP_RUNNING_TASK_FAILED, "终止运行中的任务待失败", nil},
	OUTBOUND_DUPLICATED_TASK:                   {OUTBOUND_DUPLICATED_TASK, "已存在重复的出库任务待执行或正在执行", nil},
	WEB_INBOUND_FEATURE_NOT_ALLOWED:            {WEB_INBOUND_FEATURE_NOT_ALLOWED, "Web入库下单功能不可使用", nil},
	OUTBOUND_AGV_NOT_ALLOWED:                   {OUTBOUND_AGV_NOT_ALLOWED, "大模具不能使用AGV", nil},
	WEB_OUTBOUND_FEATURE_NOT_ALLOWED:           {WEB_OUTBOUND_FEATURE_NOT_ALLOWED, "Web出库下单功能不可使用", nil},
	WEB_INOUTBOUND_FEATURE_NOT_ALLOWED:         {WEB_INOUTBOUND_FEATURE_NOT_ALLOWED, "Web出入库功能不可使用", nil},

	//模具仓管理
	MOLD_WAREHOUSE_MANAGE_CATEGORY_CONFLICT: {MOLD_WAREHOUSE_MANAGE_CATEGORY_CONFLICT, "模具属性跟储位属性不匹配", nil},
	MOLD_WAREHOUSE_MOLD_CODE_DUPLICATED:     {MOLD_WAREHOUSE_MOLD_CODE_DUPLICATED, "该模具编号已经存在立库中", nil},
}

func Success(result interface{}) *Response {
	return &Response{
		Code:    "0",
		Message: "请求成功",
		Result:  result,
	}
}

func SuccessN() *Response {
	return ResponseEnum[SUCCESS]
}

func ParamsError(result string) *Response {
	return &Response{
		Code:    "400",
		Message: result,
		Result:  nil,
	}
}

func ParamsErrorN() *Response {
	return ResponseEnum[PARAMS_ERROR]
}

func UnknowError() *Response {
	return ResponseEnum[UNKNOW_ERROR]
}
