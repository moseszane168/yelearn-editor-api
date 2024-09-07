package message

import (
	"crf-mold/base"
	"crf-mold/common"
	"crf-mold/common/constant"
	"crf-mold/dao"
	"crf-mold/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Tags 消息中心
// @Summary 消息分页
// @Accept json
// @Produce json
// @Param query query PageMessageInputVO true "PageMessageInputVO"
// @Param AuthToken header string false "Token"
// @Success 200 {object} PageMessageOututVO
// @Router /message/page [get]
func PageMessage(c *gin.Context) {
	var vo PageMessageInputVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		panic(base.ParamsErrorN())
	}

	current, size := base.GetPageParams(c)
	userId := c.GetHeader(constant.USERID)

	tx := dao.GetConn().Table("message").Where("content like concat('%',?,'%') and operator = ? and is_deleted = 'N'",
		vo.CodeOrName, userId).Order("`gmt_created` desc")
	var result []model.Message
	page := base.Page(tx, &result, current, size)

	if len(result) == 0 {
		page.List = []interface{}{}
	} else {
		res := make([]PageMessageOututVO, len(result))
		for i := 0; i < len(result); i++ {
			var item PageMessageOututVO
			base.CopyProperties(&item, result[i])
			res[i] = item
		}
		page.List = res
	}

	c.JSON(http.StatusOK, base.Success(page))
}

// @Tags 消息中心
// @Summary 删除消息
// @Accept json
// @Produce json
// @Param Body body common.DeleteVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"null"}"
// @Router /message [delete]
func DeleteMessage(c *gin.Context) {
	var vo common.DeleteVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	ids := vo.IDS
	userId := c.GetHeader(constant.USERID)

	dao.GetConn().Table("message").Where("id in ? and operator = ?", ids, userId).Updates(&model.Message{
		IsDeleted: "Y",
		UpdatedBy: userId,
	})

	// 消息推送
	client, ok := ClientMap[userId]
	if ok {
		client.Write(fmt.Sprintf(`{"count": %d}`, GetUnreadMessageCount(userId)))
	}

	c.JSON(http.StatusOK, base.Success(true))
}

// @Tags 消息中心
// @Summary 消息已读
// @Accept json
// @Produce json
// @Param Body body common.IDSVO true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"true"}"
// @Router /message/read [put]
func ReadMessage(c *gin.Context) {
	var vo common.IDSVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsErrorN())
	}

	ids := vo.IDS
	userId := c.GetHeader(constant.USERID)

	dao.GetConn().Table("message").Where("id in ? and operator = ?", ids, userId).Updates(&model.Message{
		Status:    "read",
		UpdatedBy: userId,
	})

	// 消息推送
	client, ok := ClientMap[userId]
	if ok {
		client.Write(fmt.Sprintf(`{"count": %d}`, GetUnreadMessageCount(userId)))
	}

	c.JSON(http.StatusOK, base.Success(true))
}
