package message

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var ClientMap = make(map[string]Client)

// 读写锁
var WsLock sync.RWMutex

// 已连接的Client
type Client struct {
	conn *websocket.Conn
}

// 发送消息
func (c *Client) Write(msg string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

// http升级websocket校验
var upgrader = websocket.Upgrader{
	// 这个是校验请求来源
	// 在这里我们不做校验，直接return true
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketUpgrade websocket接口,从http升级到websocket协议
func WebSocketUpgrade(c *gin.Context) {
	loginName := c.Query("loginName")
	logrus.Info(loginName)

	// 将普通的http GET请求升级为websocket请求
	client, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	// 加入连接客户列表
	WsLock.Lock()
	if old, ok := ClientMap[loginName]; ok {
		old.conn.Close()
	}

	cc := Client{
		conn: client,
	}
	ClientMap[loginName] = cc

	WsLock.Unlock()

	// 推送数量
	count := GetUnreadMessageCount(loginName)
	cc.Write(fmt.Sprintf(`{"count": %d}`, count))
}
