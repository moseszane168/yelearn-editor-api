package websocket

import (
	"crf-mold/base"
	"crf-mold/common/constant"
	"crf-mold/common/plc"
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	_ "net/http"
	"time"
)

// ToClientChan 客户端信息通道
var ToClientChan chan clientInfo

// clientInfo 通道结构体
type clientInfo struct {
	ClientId   string
	SendUserId string
	MessageId  string
	Code       string
	Msg        string
	Data       *interface{}
}

type RetData struct {
	MessageId  string      `json:"messageId"`
	SendUserId string      `json:"sendUserId"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
}

// 心跳间隔
var heartbeatInterval = 10 * time.Second

func init() {
	ToClientChan = make(chan clientInfo, 1000)
}

var Manager = NewClientManager() // 管理者

// SendMessage2Client 发送信息到指定客户端
func SendMessage2Client(clientId string, sendUserId string, code string, msg string, data *interface{}) (messageId string) {
	messageId = base.GenUUID()
	//发送到客户端
	SendToClient(messageId, clientId, sendUserId, code, msg, data)

	return
}

// CloseClient 关闭客户端
func CloseClient(clientId, systemId string) {
	if conn, err := Manager.GetByClientId(clientId); err == nil && conn != nil {
		if conn.SystemId != systemId {
			return
		}
		Manager.DisConnect <- conn
		//log.Printf("客户端 %s 被关闭连接", clientId)
	}

	return
}

// AddClient2Group 添加客户端到分组
func AddClient2Group(systemId string, groupName string, clientId string, userId string, extend string) {
	if client, err := Manager.GetByClientId(clientId); err == nil {
		//添加到分组group
		Manager.AddClient2Group(groupName, client, userId, extend)
	}
}

// SendMessage2Group 发送信息到指定分组
func SendMessage2Group(systemId, sendUserId, groupName string, code string, msg string, data *interface{}) (messageId string) {
	messageId = base.UUID()
	Manager.Send2Group(systemId, messageId, sendUserId, groupName, code, msg, data)
	return
}

// SendMessage2System 发送信息到指定业务系统
func SendMessage2System(systemId, sendUserId string, code string, msg string, data interface{}) {
	messageId := base.UUID()
	Manager.Send2BusinessSystem(systemId, messageId, sendUserId, code, msg, &data)
}

// GetSystemGroupClientList 获取分组列表
func GetSystemGroupClientList(systemId *string, groupName *string) map[string]interface{} {
	var clientList []string
	retList := Manager.GetGroupClientList(GenGroupKey(*systemId, *groupName))
	clientList = append(clientList, retList...)

	return map[string]interface{}{
		"count": len(clientList),
		"list":  clientList,
	}
}

// SendToClient 通过本服务器发送信息
func SendToClient(messageId, clientId string, sendUserId string, code string, msg string, data *interface{}) {
	//log.Printf("%s 发送到客户端 %s", sendUserId, clientId)
	ToClientChan <- clientInfo{ClientId: clientId, MessageId: messageId, SendUserId: sendUserId, Code: code, Msg: msg, Data: data}
	return
}

// ListenSendMessage 监听并发送给客户端信息
func ListenSendMessage() {
	for {
		clientInfo := <-ToClientChan
		//log.WithFields(log.Fields{
		//	"clientId":   clientInfo.ClientId,
		//	"messageId":  clientInfo.MessageId,
		//	"sendUserId": clientInfo.SendUserId,
		//	"code":       clientInfo.Code,
		//	"msg":        clientInfo.Msg,
		//	"data":       clientInfo.Data,
		//}).Info("发送到客户端....")
		if conn, err := Manager.GetByClientId(clientInfo.ClientId); err == nil && conn != nil {
			if err := Render(conn.Socket, clientInfo.MessageId, clientInfo.SendUserId, clientInfo.Code, clientInfo.Msg, clientInfo.Data); err != nil {
				Manager.DisConnect <- conn
				//log.WithFields(log.Fields{
				//	"clientId": clientInfo.ClientId,
				//	"msg":      clientInfo.Msg,
				//}).Error("客户端异常离线：" + err.Error())
			}
		}
	}
}

func Render(conn *websocket.Conn, messageId string, sendUserId string, code string, message string, data interface{}) error {
	if message == "rtsp" {
		binaryData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return conn.WriteMessage(websocket.BinaryMessage, binaryData)
	} else {
		return conn.WriteJSON(RetData{
			Code:       code,
			MessageId:  messageId,
			SendUserId: sendUserId,
			Message:    message,
			Result:     data,
		})
	}

}

// PingTimer 启动定时器进行心跳检测
func PingTimer() {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			<-ticker.C
			//发送心跳
			for clientId, conn := range Manager.AllClient() {
				if err := conn.Socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					Manager.DisConnect <- conn
					log.Errorf("发送心跳失败: %s 总连接数：%d", clientId, Manager.Count())
				}
			}
		}

	}()
}

type PLCStacker struct {
	StackerType int //堆垛机类型 1小库/2大库
	XCoordinate int // x轴坐标
}

// SimulatePLCWebsocket 模拟PLC堆垛机Websocket请求
func SimulatePLCWebsocket() {
	// 小库模拟
	go func() {
		var smallXStart, smallX = 15505, 15505
		var smallXEnd = 685
		var gap = 5
		for ; smallX >= smallXEnd; smallX -= gap {
			time.Sleep(50 * time.Millisecond)
			payload := PLCStacker{StackerType: 1, XCoordinate: smallX}
			// 发送到对应业务系统中
			SendMessage2System(constant.WebsocketPlcStackerCoordinate, "", string(base.SUCCESS), "堆垛机运行中", payload)
			if (smallX - gap) < smallXEnd {
				smallX = smallXStart
			}
		}
	}()
	// 大库模拟
	go func() {
		var bigXStart, bigX = 16502, 16502
		var bigXEnd = 689
		var gap = 5
		for ; bigX >= bigXEnd; bigX -= gap {
			time.Sleep(50 * time.Millisecond)
			payload := PLCStacker{StackerType: 2, XCoordinate: bigX}
			// 发送到对应业务系统中
			SendMessage2System(constant.WebsocketPlcStackerCoordinate, "", string(base.SUCCESS), "堆垛机运行中", payload)
			if (bigX - gap) < bigXEnd {
				bigX = bigXStart
			}
		}
	}()
}

// InoutBoundStatusReadPoints 模具出入库状态读取结构体
type InoutBoundStatusReadPoints struct {
	PlcConfig         plc.PLCConfig
	StationUniqueCode string
	ReadOffset1       int
	ReadPoint1        int
	ReadOffset2       int
	ReadPoint2        int
}

// InoutBoundStatusResponse 模具出入库状态响应结构体
type InoutBoundStatusResponse struct {
	Type              string `json:"type"`              // 出入库类型 inbound-入库; outbound-出库
	StationUniqueCode string `json:"stationUniqueCode"` //站口唯一标识
	IsRoll            bool   `json:"isRoll"`            //是否滚动中
	Occupy            bool   `json:"occupy"`            //是否占位中
	MoldStatus        int    `json:"moldStatus"`        // 模具状态, 1-正常 2-待保养 3-报废或封存 4-未知
}

// SimulateInboundStatusWebsocket 模拟模具入库状态推送
func SimulateInboundStatusWebsocket() {
	inboundStations := []string{"101", "201", "203"}
	for _, stationCode := range inboundStations {
		// 定时读取各站口模具入库状态
		go func(code string) {
			var payload InoutBoundStatusResponse
			var count int
			for {
				payload.Type = "inbound"
				payload.StationUniqueCode = code
				payload.MoldStatus = 1 // 入库默认为正常模具
				count++
				if count < 12 {
					payload.IsRoll = true
					payload.Occupy = false
				} else if count < 18 {
					payload.IsRoll = false
					payload.Occupy = true
				} else if count < 23 {
					payload.IsRoll = false
					payload.Occupy = false
				} else {
					count = 0
					continue
				}
				// 发送到对应业务系统中
				SendMessage2System(constant.WebsocketInboundStatus, "", string(base.SUCCESS), "模具入库状态推送", payload)
				// 间隔1s读取推送一次
				time.Sleep(time.Second)
			}
		}(stationCode)
	}
}

// SimulateOutboundStatusWebsocket 模拟模具出库状态推送
func SimulateOutboundStatusWebsocket() {
	inboundStations := []string{"102", "202", "204"}
	for _, stationCode := range inboundStations {
		// 定时读取各站口模具出库状态
		go func(code string) {
			var payload InoutBoundStatusResponse
			var count int
			for {
				payload.Type = "outbound"
				payload.StationUniqueCode = code
				payload.MoldStatus = 1 // 出库默认为正常模具
				count++
				if count < 12 {
					payload.IsRoll = true
					payload.Occupy = false
				} else if count < 18 {
					payload.IsRoll = false
					payload.Occupy = true
				} else if count < 23 {
					payload.IsRoll = false
					payload.Occupy = false
				} else {
					count = 0
					time.Sleep(time.Second * 3)
					continue
				}
				// 发送到对应业务系统中
				SendMessage2System(constant.WebsocketOutboundStatus, "", string(base.SUCCESS), "模具出库状态推送", payload)
				// 间隔1s读取推送一次
				time.Sleep(time.Second)
			}
		}(stationCode)
	}
}
