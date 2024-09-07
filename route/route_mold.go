package route

import (
	"crf-mold/common/websocket"
	dashboard "crf-mold/controller/dashboard/production"
	"crf-mold/controller/dict"
	"crf-mold/controller/email"
	"crf-mold/controller/excel"
	"crf-mold/controller/file"
	"crf-mold/controller/message"
	knowledge "crf-mold/controller/mold/knowledge"
	mold "crf-mold/controller/mold/ledger"
	maintenance "crf-mold/controller/mold/maintenance"
	quality "crf-mold/controller/mold/quality"
	remodel "crf-mold/controller/mold/remodel"
	repair "crf-mold/controller/mold/repair"
	"crf-mold/controller/mold/scopestatistics"
	"crf-mold/controller/mold/scrollmsg"
	"crf-mold/controller/productresume"
	"crf-mold/controller/rtsp"
	"crf-mold/controller/spare"
	"crf-mold/controller/user"

	"github.com/gin-gonic/gin"
)

func InitMoldRoute(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		// 字典组
		v1.POST("/dict/group", dict.CreateDictGroup)
		v1.DELETE("/dict/group", dict.DeleteDictGroup)
		v1.PUT("/dict/group", dict.UpdateDictGroup)
		v1.GET("/dict/group", dict.ListDictGroup)

		v1.GET("/dict/all", dict.ListAll)

		// 字典属性
		v1.POST("/dict", dict.CreateDictProperty)
		v1.DELETE("/dict", dict.DeleteDictProperty)
		v1.PUT("/dict", dict.UpdateDictProperty)
		v1.GET("/dict", dict.ListDictProperty)
		v1.GET("/dict/page", dict.PageDictProperty)

		// 用户中心
		v1.GET("/user/rsa", user.GetRsaPublicKey)
		v1.POST("/user/encrypt", user.EncryptPassword)
		v1.POST("/user/login", user.LoginIn)
		v1.POST("/user/logout", user.LogOut)

		v1.POST("/user/authority", user.GetUserAuthority)
		v1.POST("/user", user.CreateUser)
		v1.PUT("/user", user.UpdateUser)
		v1.GET("/user", user.ListUser)
		v1.GET("/user/one", user.SingelUserAdmin)
		v1.GET("/user/self", user.LoginUserInfo)
		v1.GET("/user/page", user.PageUser)
		v1.DELETE("/user", user.DeleteUser)
		v1.PUT("/user/reset", user.ResetPassWord)
		v1.POST("/user/pwd", user.ModifyPassword)

		// 权限
		v1.GET("/user/authority", user.GetUserAuthority)
		v1.GET("/user/authority/one", user.GetUserAuthorityById)
		v1.GET("/user/authorities", user.GetAuthoritys)
		v1.PUT("/user/authority", user.UpdateUserAuthority)

		// 台账
		v1.POST("/mold", mold.CreateMold)
		v1.PUT("/mold", mold.UpdateMold)
		v1.DELETE("/mold", mold.DeleteMold)
		v1.GET("/mold/page", mold.PageMold)
		v1.GET("/mold/one", mold.OneMold)
		v1.POST("/mold/parse", mold.ParseMoldExcel)
		v1.GET("/mold/excel", mold.ExportMold)
		v1.POST("/mold/excel", mold.ImportMold)
		v1.PUT("/mold/flush", mold.ZeroMoldFlushCount)

		v1.GET("/mold/part/code/group", mold.MoldPartCodeGroup)
		v1.POST("/mold/by/part/code", mold.MoldQueryByPartCode)

		// Excel模板
		v1.GET("/template", excel.DownloadTemplate)

		// 文件上传
		v1.POST("/file", file.UploadFile)

		// 知识库
		v1.POST("/knowledge", knowledge.CreateMoldKnowledge)
		v1.DELETE("/knowledge", knowledge.DeleteMoldKnowledge)
		v1.PUT("/knowledge", knowledge.UpdateMoldKnowledge)
		v1.GET("/knowledge", knowledge.ListMoldKnowledge)
		v1.GET("/knowledge/page", knowledge.PageMoldKnowledge)
		v1.GET("/knowledge/one", knowledge.OneKnowledge)
		v1.GET("/knowledge/search", knowledge.SearchKnowledge)

		// 模具维修
		v1.POST("/mold/repair", repair.CreateRepair)
		v1.PUT("/mold/repair", repair.UpdateRepair)
		v1.DELETE("/mold/repair", repair.DeleteRepair)
		v1.GET("/mold/repair/page", repair.PageRepair)
		v1.GET("/mold/repair/one", repair.OneRepair)

		// 模具质量
		v1.POST("/mold/quality", quality.CreateQuality)
		v1.PUT("/mold/quality", quality.UpdateQuality)
		v1.DELETE("/mold/quality", quality.DeleteQuality)
		v1.GET("/mold/quality/page", quality.PageQuality)
		v1.GET("/mold/quality/one", quality.OneQuality)

		// 模具改造
		v1.POST("/mold/remodel", remodel.CreateRemodel)
		v1.PUT("/mold/remodel", remodel.UpdateRemodel)
		v1.DELETE("/mold/remodel", remodel.DeleteRemodel)
		v1.GET("/mold/remodel/page", remodel.PageRemodel)
		v1.GET("/mold/remodel/one", remodel.OneRemodel)

		v1.POST("/mold/remodel/completed", remodel.UpdateRemodelStatusCompleted)
		v1.POST("/mold/remodel/withdraw", remodel.UpdateRemodelStatusWithdraw)

		// 改造定时任务
		v1.PUT("/mold/remodel/time", remodel.ExposeTimingFunc)

		// 模具文档
		v1.GET("/mold/doc/page", mold.PageMoldDoc)
		v1.POST("/mold/doc", mold.SaveMoldDoc)

		// 备件
		v1.POST("/spare", spare.CreateSpare)
		v1.PUT("/spare", spare.UpdateSpare)
		v1.DELETE("/spare", spare.DeleteSpare)
		v1.GET("/spare/page", spare.PageSpare)
		v1.GET("/spare/one", spare.OneSpare)
		v1.POST("/spare/parse", spare.ParseSpareExcel)
		v1.GET("/spare/excel", spare.ExportSpare)
		v1.POST("/spare/excel", spare.ImportSpare)

		// 备件库存
		v1.POST("/spare/inbound", spare.InboundSpareRequest)
		v1.POST("/spare/outbound", spare.OutboundSpareRequest)

		v1.GET("/spare/request", spare.OneSpareRequest)
		v1.GET("/spare/request/page", spare.PageSpareRequest)
		v1.GET("/spare/inbound", spare.OneSpareRequestInbound)
		v1.GET("/spare/outbound", spare.OneSpareRequestOutbound)
		v1.GET("/spare/inbound/page", spare.PageOneSpareRequestInbound)
		v1.GET("/spare/outbound/page", spare.PageOneSpareRequestOutbound)

		v1.POST("/spare/request/parse", spare.ParseSpareRequestExcel)
		v1.GET("/spare/request/excel", spare.ExportSpareRequest)
		v1.POST("/spare/request/excel", spare.ImportSpareRequest)

		// 模具BOM
		v1.POST("/mold/bom", mold.CreateBom)
		v1.PUT("/mold/bom", mold.UpdateBom)
		v1.DELETE("/mold/bom", mold.DeleteBom)
		v1.GET("/mold/bom/page", mold.PageBom)
		v1.GET("/mold/bom/one", mold.OneBom)
		v1.POST("/mold/bom/parse", mold.ParseBomExcel)
		v1.GET("/mold/bom/excel", mold.ExportBom)
		v1.POST("/mold/bom/excel", mold.ImportBom)

		// 模具成型参数
		v1.PUT("/mold/molding", mold.SaveMoldingParams)
		v1.GET("/mold/molding/one", mold.OneMoldingParams)

		v1.PUT("/mold/molding/customs", mold.SaveMoldingParamsCustoms)

		v1.DELETE("/mold/molding/customs", mold.DeleteMoldingCustoms)

		// 保养履历
		v1.GET("/mold/maintenance/page", mold.PageMoldMaintenance)

		// 备件更换履历
		v1.GET("/spare/replace/page", mold.PageReplaceSpare)

		// 模具保养标准
		v1.PUT("/maintenance/standard", maintenance.UpdateMoldMaintenanceStandard)
		v1.GET("/maintenance/standard/one", maintenance.OneMoldMaintenanceStandard)
		v1.GET("/maintenance/standard/page", maintenance.PageMoldMaintenanceStandard)
		v1.GET("/maintenance/standard/list", maintenance.ListMoldMaintenanceStandard)
		v1.POST("/maintenance/standard", maintenance.CreateMoldMaintenanceStandard)
		v1.DELETE("/maintenance/standard", maintenance.DeleteMoldMaintenanceStandard)
		v1.GET("/maintenance/standard/mold/page", maintenance.MoldRelPage)
		v1.GET("/maintenance/standard/select", maintenance.MoldStandardSelectList)

		// 模具保养计划
		v1.POST("/maintenance/query/cron", maintenance.QueryCronRunTime)
		v1.POST("/maintenance/plan", maintenance.CreateMoldMaintenancePlan)
		v1.DELETE("/maintenance/plan", maintenance.DeleteMoldMaintenancePlan)
		v1.PUT("/maintenance/plan", maintenance.UpdateMoldMaintenancePlan)
		v1.GET("/maintenance/plan/list", maintenance.ListMoldMaintenancePlan)
		v1.GET("/maintenance/plan/page", maintenance.PageMoldMaintenancePlan)
		v1.GET("/maintenance/plan/one", maintenance.OneMoldMaintenancePlan)
		v1.POST("/maintenance/plan/status", maintenance.ToggleMoldMaintenancePlan)

		// 模具保养任务
		v1.POST("/maintenance/task/new", maintenance.AddMoldMaintenanceTask)
		v1.POST("/maintenance/task", maintenance.SubmitMoldMaintenanceTask)
		v1.GET("/maintenance/task/page", maintenance.PageMoldMaintenanceTask)
		v1.GET("/maintenance/task/one", maintenance.OneMoldMaintenanceTask)
		v1.POST("/maintenance/task/status", maintenance.ToggleMoldMaintenanceTask)
		v1.GET("/maintenance/task/lifecycle", maintenance.LifeCycle)

		v1.POST("/maintenance/task/charge", maintenance.ChangeChargeMaintenanceTask)
		v1.POST("/maintenance/task/delete", maintenance.DeleteMaintenanceTask)
		v1.POST("/maintenance/task/approval", maintenance.ApprovalMaintenanceTask)

		// 生产履历
		v1.GET("/mold/productresume/page", productresume.PageProductResume)
		v1.GET("/mold/productresume/excel", productresume.Export)

		// 消息中心
		v1.GET("/message/page", message.PageMessage)
		v1.DELETE("/message", message.DeleteMessage)
		v1.GET("/ws", message.WebSocketUpgrade)
		v1.PUT("message/read", message.ReadMessage)

		// 邮件
		v1.PUT("/email/timeout", email.UpdateMoldTaskTimeoutEmail)
		v1.GET("/email/timeout", email.GetMoldTaskTimeoutEmail)

		// 统计
		v1.GET("/mold/statistics/mold/status", scopestatistics.MoldStatusGroup)
		v1.GET("/mold/statistics/maintenance/task", scopestatistics.MaintenanceStatusGroup)
		v1.POST("/mold/statistics/agv/task", scopestatistics.AgvTask)
		v1.GET("/mold/statistics/repair", scopestatistics.MoldException)
		v1.POST("/mold/statistics/inoutbound/task", scopestatistics.InOutBoundTask)

		//晨会Dashboard
		v1.POST("/dashboard/production", dashboard.LineProduction)
		//周报-维修停机率报告查询
		v1.POST("/dashboard/stopRate", dashboard.ProductStatistics)
		//备件消耗趋势-周月报
		v1.POST("/dashboard/userSparePart", dashboard.UserSparePart)
		//日报-质量
		v1.POST("/dashboard/rejects", dashboard.DefectiveProducts)
		//日报-维修
		v1.POST("/dashboard/maintain", dashboard.MaintainListPlan)
		//故障top10
		v1.POST("/dashboard/errorTopTen", dashboard.ErrorTopTen)
		//趋势图-周/月趋势图
		v1.POST("/dashboard/tendencyChart", dashboard.TendencyChart)
		//趋势图产线列表
		v1.POST("/dashboard/getLineList", dashboard.GetLineList)
		//趋势图数据生成手动触发
		v1.POST("/dashboard/generateProductReport", dashboard.GenerateProductReport)

		// 备件履历
		v1.POST("/spare/resume/page", spare.SpareResumePage)
		v1.GET("/spare/resume/excel", spare.SpareResumeExport)

		// 系统websocket接口
		v1.GET("/system/ws", websocket.Run)

		// 滚动消息
		v1.POST("/scrollmsg/save", scrollmsg.SaveScrollMessage)
		v1.GET("/scrollmsg/get", scrollmsg.GetScrollMessage)

		//rtsp视频流转换
		v1.POST("/rtsp/play", rtsp.PlayRTSP)
		v1.POST("/rtsp/upload/:channel", rtsp.Mpeg1Video)
		v1.GET("/rtsp/live/:channel", rtsp.Wsplay)
	}
}
