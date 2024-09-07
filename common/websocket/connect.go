package websocket

import (
	"crf-mold/base"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type renderData struct {
	ClientId string `json:"clientId"`
}

func ConnRender(conn *websocket.Conn, data interface{}) (err error) {
	err = conn.WriteJSON(base.Response{
		Code:    base.SUCCESS,
		Message: "success",
		Result:  data,
	})

	return
}

func Run(cxt *gin.Context) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(cxt.Writer, cxt.Request, nil)
	if err != nil {
		log.Errorf("websocket upgrade error: %v", err)
		http.NotFound(cxt.Writer, cxt.Request)
		return
	}

	//设置读取消息大小上限制
	conn.SetReadLimit(maxMessageSize)

	//解析参数, 获取业务系统ID
	systemId := cxt.Request.FormValue("systemId")
	if len(systemId) == 0 {
		_ = Render(conn, "", "", string(base.WS_BUSINESS_SYSTEM_ID_ERROR), "系统ID不能为空", []string{})
		_ = conn.Close()
		return
	}

	clientId := base.GenUUID()

	clientSocket := NewClient(clientId, systemId, conn)

	Manager.AddClient2SystemClient(systemId, clientSocket)

	//读取客户端消息
	clientSocket.Read()

	if err = ConnRender(conn, renderData{ClientId: clientId}); err != nil {
		_ = conn.Close()
		return
	}

	// 用户连接事件
	Manager.Connect <- clientSocket
}
