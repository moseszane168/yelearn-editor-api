package rtsp

import (
	"bufio"
	"crf-mold/base"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"net/http"
)

// @Tags RTSP流转码
// @Summary RTSP流转码
// @Accept json
// @Produce json
// @Param Body body RTSPTransSrv true "body"
// @Param AuthToken header string false "Token"
// @Success 200 {object} string
// @Router /rtsp/play [post]
func PlayRTSP(c *gin.Context) {
	srv := RTSPTransSrv{}
	if err := c.ShouldBindBodyWith(&srv, binding.JSON); err != nil {
		panic(base.ParamsError("不是有效的 RTSP 地址"))
	}

	ret := srv.Service()
	c.JSON(http.StatusOK, base.Success(ret))
}

// Mpeg1Video 接收 mpeg1vido 数据流
func Mpeg1Video(c *gin.Context) {
	logrus.Info("接收推送数据，channel" + c.Param("channel"))
	bodyReader := bufio.NewReader(c.Request.Body)

	for {
		data, err := bodyReader.ReadBytes('\n')
		if err != nil {
			logrus.Error(err.Error())
			break
		}

		WsManager.Groupbroadcast(c.Param("channel"), data)
	}
}

// Wsplay 通过 websocket 播放 mpegts 数据
func Wsplay(c *gin.Context) {
	WsManager.RegisterClient(c)
}
