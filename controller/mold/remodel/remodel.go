/**
 * 模具改造
 */

package mold

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/controller/email"
	"crf-mold/controller/message"
	"crf-mold/dao"
	"crf-mold/model"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type StatusVO struct {
	ID     int64  `json:"id"`     // ID
	Reason string `json:"reason"` // 撤单原因
}

// @Tags 模具改造
// @Summary 模具改造分页
// @Accept json
// @Produce json
// @Param query query PageRemodelVO true "PageRemodelVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageRemodelOutVO
// @Router /mold/remodel/page [get]
func PageRemodel(c *gin.Context) {
	var vo PageRemodelVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	var result []PageRemodelOutVO
	var page *base.BasePage
	if vo.CodeOrName != "" {
		rawSql := `
		select
			mr.id,
			mr.code,
			mr.mold_code,
			mr.remodel_start_time,
			mr.remodel_end_time,
			mr.finish_time,
			COALESCE(uop.name, mr.director) AS director,
			mr.type,
			mi.project_name,
			mr.location,
			mr.content,
			mr.status,
			mr.withdraw_reason,
			mr.is_delay,
			mr.delay_day,
			ifnull(GROUP_CONCAT(DISTINCT mpr.part_code SEPARATOR '/'),'') as part_codes
		from mold_remodel mr 
		left join mold_info mi on mr.mold_code = mi.code and mi.is_deleted = 'N'
		left join mold_part_rel mpr on mpr.mold_code = mi.code and mpr.is_deleted = 'N'
		left join user_info uop on uop.login_name = mr.director
		where mr.is_deleted = 'N' and (mi.code like concat('%',?,'%') or mi.name like concat('%',?,'%'))
		group by mr.id,mr.code,mr.mold_code,mi.code,mi.project_name,uop.name
		order by mr.finish_time desc, mr.gmt_created desc
		`
		page = base.PageWithRawSQL(dao.GetConn(), &result, vo.GetCurrentPage(), vo.GetSize(), rawSql, vo.CodeOrName, vo.CodeOrName)
	} else {
		var entity model.MoldRemodel
		base.CopyProperties(&entity, vo)
		tx := dao.BuildWhereCondition(dao.GetConn().Table("mold_remodel"), entity)
		// 编码模糊
		if vo.CodeLike != "" {
			tx.Where("(code like concat('%',?,'%') or director like concat('%',?,'%'))", vo.CodeLike, vo.CodeLike)
		}
		tx.Where("is_deleted = 'N'").Order("gmt_created desc")
		page = base.Page(tx, &result, vo.GetCurrentPage(), vo.GetSize())
	}

	pageList := page.List.(*[]PageRemodelOutVO)
	if len(*pageList) == 0 {
		page.List = []interface{}{}
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 模具改造
// @Summary 模具改造查看
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param AuthToken header string false "Token"
// @Success 200 {object} mold.PageRemodelOutVO
// @Router /mold/remodel/one [get]
func OneRemodel(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		panic(base.ParamsErrorN())
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(base.ParamsErrorN())
	}

	var moldRemodel []model.MoldRemodel
	dao.GetConn().Table("mold_remodel").Where("id = ? and is_deleted = 'N'", id).Find(&moldRemodel)

	if len(moldRemodel) > 0 {
		m := moldRemodel[0]
		var vo PageRemodelOutVO
		base.CopyProperties(&vo, m)

		c.JSON(http.StatusOK, base.Success(vo))
	} else {
		c.JSON(http.StatusOK, base.SuccessN())
	}
}

// @Tags 模具改造
// @Summary 删除模具改造
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/remodel [delete]
func DeleteRemodel(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	ids := vo.IDS
	if len(ids) == 0 {
		panic(base.ParamsErrorN())
	}

	dao.GetConn().Table("mold_remodel").Where("id in ?  and is_deleted = 'N'", ids).Updates(&model.MoldRemodel{
		IsDeleted: "Y",
		UpdatedBy: c.GetHeader(constant.USERID),
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具改造
// @Summary 添加模具改造
// @Accept json
// @Produce json
// @Param Body body CreateRemodelVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/remodel [post]
func CreateRemodel(c *gin.Context) {
	var vo CreateRemodelVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}

	// 校验结束时间必须在当前时间之后
	t, err := time.Parse(base.TimeFormart, vo.RemodelEndTime.String())
	if err != nil {
		panic(base.ParamsErrorN())
	}

	if base.DiffDaySince(t) > 0 {
		panic(base.ResponseEnum[base.REMODEL_END_TIME_NOT_BEFORE_NOW])
	}

	var moldRemodel model.MoldRemodel
	base.CopyProperties(&moldRemodel, vo)

	moldRemodel.Code = common.GenerateCode("GZ")

	userId := c.GetHeader(constant.USERID)
	moldRemodel.CreatedBy = userId
	moldRemodel.UpdatedBy = userId
	moldRemodel.Status = "wait"
	moldRemodel.FinishTime = nil

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("mold_remodel").Create(&moldRemodel).Error; err != nil {
		panic(err)
	}

	// 消息发送
	message.CreateMessage(tx, message.Message{
		Content:   email.FormatMessage(email.REMODEL_CREATE, moldRemodel.Code).Content,
		Operator:  []string{moldRemodel.Director},
		RemodelId: moldRemodel.ID,
		Type:      constant.REMODEL,
	})

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(moldRemodel.ID))
}

// @Tags 模具改造
// @Summary 更新模具改造
// @Accept json
// @Produce json
// @Param Body body UpdateRemodelVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/remodel [put]
func UpdateRemodel(c *gin.Context) {
	var vo UpdateRemodelVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	// 参数校验,TODO

	if vo.ID == 0 {
		panic(base.ParamsErrorN())
	}

	// 结束时间不能在当前时间之前
	t, err := time.Parse(base.TimeFormart, vo.RemodelEndTime.String())
	if err != nil {
		panic(base.ParamsErrorN())
	}

	if base.DiffDaySince(t) > 0 {
		panic(base.ResponseEnum[base.REMODEL_END_TIME_NOT_BEFORE_NOW])
	}

	// 存在
	var one model.MoldRemodel
	if err := dao.GetConn().Table("mold_remodel").Where("id = ? and is_deleted = 'N'", vo.ID).First(&one).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 只能编辑撤回的单
	if one.Status != "withdraw" {
		panic(base.ParamsErrorN())
	}

	// 更新
	userId := c.GetHeader(constant.USERID)
	var moldRemodel model.MoldRemodel
	base.CopyProperties(&moldRemodel, vo)
	moldRemodel.UpdatedBy = userId
	moldRemodel.Status = "wait"

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	if err := tx.Table("mold_remodel").Where("id = ? and is_deleted = 'N'", moldRemodel.ID).Updates(&moldRemodel).Error; err != nil {
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具改造
// @Summary 完工
// @Accept json
// @Produce json
// @Param Body body common.IDVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/remodel/completed [post]
func UpdateRemodelStatusCompleted(c *gin.Context) {
	var vo common.IDVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	if vo.ID == 0 {
		panic(base.ParamsErrorN())
	}

	// 存在
	var entity model.MoldRemodel
	if err := dao.GetConn().Table("mold_remodel").Where("id = ? and is_deleted = 'N'", vo.ID).First(&entity).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 只能从待提交状态撤回
	if entity.Status != "wait" {
		panic(base.ParamsErrorN())
	}

	t := base.Now()
	dao.GetConn().Table("mold_remodel").Where("id = ? and is_deleted = 'N'", vo.ID).Updates(&model.MoldRemodel{
		UpdatedBy:  c.GetHeader(constant.USERID),
		Status:     "complete",
		FinishTime: &t,
	})

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 模具改造
// @Summary 撤单
// @Accept json
// @Produce json
// @Param Body body StatusVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /mold/remodel/withdraw [post]
func UpdateRemodelStatusWithdraw(c *gin.Context) {
	var vo StatusVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	if vo.ID == 0 {
		panic(base.ParamsErrorN())
	}

	// 撤销原因
	if vo.Reason == "" {
		panic(base.ResponseEnum[base.WITHDRAW_CAN_NOT_EMPTY])
	}

	// 存在
	var entity model.MoldRemodel
	if err := dao.GetConn().Table("mold_remodel").Where("id = ? and is_deleted = 'N'", vo.ID).First(&entity).Error; err != nil {
		panic(base.ParamsErrorN())
	}

	// 只能从待提交状态撤回
	if entity.Status != "wait" {
		panic(base.ParamsErrorN())
	}

	dao.GetConn().Table("mold_remodel").Where("id = ? and is_deleted = 'N'", vo.ID).Updates(&model.MoldRemodel{
		UpdatedBy:      c.GetHeader(constant.USERID),
		Status:         "withdraw",
		WithdrawReason: vo.Reason,
	})

	c.JSON(http.StatusOK, base.Success(true))
}

func HandleExpiredRemodelOrderFunc() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error("HandleExpiredRemodelOrderFunc error Panic info is: %v", err)
			debug.PrintStack()
			logrus.WithField("stack", string(debug.Stack())).Errorf("panic: %v\n", err)
		}
	}()

	logrus.Info("HandleExpiredRemodelOrder处理过期的改造订单执行中:", base.Now().String())
	// 获取当前时间
	now := base.Now().DateString()

	// 找到所以已经延期的改造订单
	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	var delayOrder []model.MoldRemodel
	tx.Table("mold_remodel").Where("status = 'wait' and remodel_end_time < ? and is_deleted = 'N'", now).Find(&delayOrder)

	// 更新每一条延期改造订单的延期天数
	for i := 0; i < len(delayOrder); i++ {
		id := delayOrder[i].ID
		timeStr := delayOrder[i].RemodelEndTime
		t, _ := base.FormatTime(timeStr.String())
		delayDay := base.DiffDaySince(t.Time())

		tx.Table("mold_remodel").Where("id = ?", id).Updates(&model.MoldRemodel{
			IsDelay:   "Y",
			UpdatedBy: "admin",
			DelayDay:  delayDay,
		})
	}
	tx.Commit()

	var result []email.RemodelEmailVO
	sql := `
		select
			mr.id,
			mr.code,
			mr.mold_code,
			mr.remodel_start_time,
			mr.remodel_end_time,
			mr.finish_time,
			mr.director,
			mr.type,
			mi.project_name,
			mr.location,
			mr.content,
			mr.status,
			mr.withdraw_reason,
			mr.is_delay,
			mr.delay_day,
			ifnull(GROUP_CONCAT(DISTINCT mpr.part_code SEPARATOR '/'),'') as part_codes
		from mold_remodel mr 
		left join mold_info mi on mr.mold_code = mi.code and mi.is_deleted = 'N'
		left join mold_part_rel mpr on mpr.mold_code = mi.code and mpr.is_deleted = 'N'
		where mr.is_deleted = 'N' and delay_day >= 7 and delay_email = 'N'
		group by mr.id,mr.code,mr.mold_code,mi.code,mi.project_name
		order by mr.finish_time desc, mr.gmt_created desc
	`
	// 查找需要发送邮件的改造单
	dao.GetConn().Raw(sql, now).Scan(&result)

	// 给订单的责任人人发消息
	for i := 0; i < len(result); i++ {
		message.CreateMessage(tx, message.Message{
			Content:   email.FormatMessage(email.REMODEL_TIMEOUT, result[i].Code).Content,
			Operator:  []string{result[i].Director},
			RemodelId: result[i].ID,
			Type:      constant.REMODEL,
		})
	}

	// 对于每一条已经延期七天的改造订单发送一封邮件
	receivers, err := email.GetMoldTaskTimeoutEmailReceiver(email.REMODEL_TIMEOUT)
	if err != nil {
		logrus.Error(err)
	} else {
		if len(result) > 0 {
			for _, receiver := range receivers {
				if err := email.SendEmail(receiver, email.REMODEL_TIMEOUT, result); err != nil {
					logrus.WithField("address", receiver).Error(err)
				}
			}
		}
	}

	// 发送邮件
	for i := 0; i < len(result); i++ {
		order := result[i]
		// 更新对应订单的状态,TODO:待对接后进行
		tx.Table("mold_remodel").Where("id = ?", order.ID).Updates(&model.MoldRemodel{
			UpdatedBy:  "admin",
			DelayEmail: "Y",
		})
	}
}

//
// 处理过期的改造订单
//
func HandleExpiredRemodelOrder() {
	f := HandleExpiredRemodelOrderFunc

	// 每天凌晨12:10点执行
	c := cron.New()
	c.AddFunc("0 10 0 * * ?", f)
	c.Start()

	logrus.Info("CRON启动,等待执行")
}

//
// 处理过期的改造订单
//
func ExposeTimingFunc(c *gin.Context) {
	HandleExpiredRemodelOrderFunc()
	c.JSON(http.StatusOK, base.Success(true))
}
