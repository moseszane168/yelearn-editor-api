package scrollmsg

import (
	"crf-mold/base"
	redisClient "crf-mold/common/redis"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// @Tags 滚动消息
// @Summary 保存滚动消息
// @Accept json
// @Produce json
// @Param Body body ScrollMessageVO true "ScrollMessageVO"
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"保存成功","result":"null"}"
// @Router /scrollmsg/save [post]
func SaveScrollMessage(c *gin.Context) {
	var vo ScrollMessageVO
	if err := c.ShouldBindBodyWith(&vo, binding.JSON); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	if err := redisClient.RedisClient.SetKey("scrollMsg", vo.Message, -1); err != nil {
		panic(base.ParamsError(err.Error()))
	}
	c.JSON(http.StatusOK, base.SuccessN())
}

// @Tags 滚动消息
// @Summary 查询滚动消息
// @Produce json
// @Param AuthToken header string false "Token"
// @Success 200 {string} json "{"code":200,"message":"请求成功","result":"消息内容"}"
// @Router /scrollmsg/get [get]
func GetScrollMessage(c *gin.Context) {
	val, err := redisClient.RedisClient.GetKey("scrollMsg")
	if err != nil && err.Error() != "redis: nil" {
		panic(base.ParamsError(err.Error()))
	}
	c.JSON(http.StatusOK, base.Success(val))
}
