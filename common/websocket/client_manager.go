package websocket

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// ClientManager 连接管理
type ClientManager struct {
	ClientIdMap     map[string]*Client // 全部的连接
	ClientIdMapLock sync.RWMutex       // 读写锁

	Connect    chan *Client // 连接处理
	DisConnect chan *Client // 断开连接处理

	GroupLock sync.RWMutex
	Groups    map[string][]string

	SystemClientsLock sync.RWMutex
	SystemClients     map[string][]string //业务系统的连接
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		ClientIdMap:   make(map[string]*Client),
		Connect:       make(chan *Client, 10000),
		DisConnect:    make(chan *Client, 10000),
		Groups:        make(map[string][]string, 100),
		SystemClients: make(map[string][]string, 100),
	}

	return
}

// Start 管理处理程序
func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Connect:
			// 建立连接事件
			manager.EventConnect(client)
		case conn := <-manager.DisConnect:
			// 断开连接事件
			manager.EventDisconnect(conn)
		}
	}
}

// EventConnect 建立连接事件
func (manager *ClientManager) EventConnect(client *Client) {
	manager.AddClient(client)

	log.WithFields(log.Fields{
		"clientId": client.ClientId,
		"counts":   Manager.Count(),
	}).Info("客户端已连接")
}

// EventDisconnect 断开连接时间
func (manager *ClientManager) EventDisconnect(client *Client) {
	//关闭连接
	_ = client.Socket.Close()
	manager.DelClient(client)

	//mJson, _ := json.Marshal(map[string]string{
	//	"clientId": client.ClientId,
	//	"userId":   client.UserId,
	//	"extend":   client.Extend,
	//})
	//data := string(mJson)
	//sendUserId := ""
	//
	////发送下线通知
	//if len(client.GroupList) > 0 {
	//	for _, groupName := range client.GroupList {
	//		SendMessage2Group(client.SystemId, sendUserId, groupName, base.WS_ONLINE_MESSAGE_CODE, "客户端下线", &data)
	//	}
	//}

	log.WithFields(log.Fields{
		"clientId": client.ClientId,
		"counts":   Manager.Count(),
		"seconds":  uint64(time.Now().Unix()) - client.ConnectTime,
	}).Info("客户端已断开")

	//标记销毁
	client.IsDeleted = true
	client = nil
}

// AddClient 添加客户端
func (manager *ClientManager) AddClient(client *Client) {
	manager.ClientIdMapLock.Lock()
	defer manager.ClientIdMapLock.Unlock()

	manager.ClientIdMap[client.ClientId] = client
}

// AllClient 获取所有的客户端
func (manager *ClientManager) AllClient() map[string]*Client {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()

	return manager.ClientIdMap
}

// Count 客户端数量
func (manager *ClientManager) Count() int {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()

	return len(manager.ClientIdMap)
}

// DelClient 删除客户端
func (manager *ClientManager) DelClient(client *Client) {
	manager.delClientIdMap(client.ClientId)

	//删除所在的分组
	if len(client.GroupList) > 0 {
		for _, groupName := range client.GroupList {
			manager.delGroupClient(GenGroupKey(client.SystemId, groupName), client.ClientId)
		}
	}

	// 删除系统里的客户端
	manager.delSystemClient(client)
}

// delClientIdMap 删除clientIdMap
func (manager *ClientManager) delClientIdMap(clientId string) {
	manager.ClientIdMapLock.Lock()
	defer manager.ClientIdMapLock.Unlock()

	delete(manager.ClientIdMap, clientId)
}

// GetByClientId 通过clientId获取
func (manager *ClientManager) GetByClientId(clientId string) (*Client, error) {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()

	if client, ok := manager.ClientIdMap[clientId]; !ok {
		return nil, errors.New("客户端不存在")
	} else {
		return client, nil
	}
}

// Send2Group 发送到分组
func (manager *ClientManager) Send2Group(systemId, messageId, sendUserId, groupName string, code string, msg string, data *interface{}) {
	if len(groupName) > 0 {
		clientIds := manager.GetGroupClientList(GenGroupKey(systemId, groupName))
		if len(clientIds) > 0 {
			for _, clientId := range clientIds {
				if _, err := Manager.GetByClientId(clientId); err == nil {
					//发送到分组下的客户端
					SendToClient(messageId, clientId, sendUserId, code, msg, data)
				} else {
					//删除分组
					manager.delGroupClient(GenGroupKey(systemId, groupName), clientId)
				}
			}
		}
	}
}

// Send2BusinessSystem 发送给指定业务系统
func (manager *ClientManager) Send2BusinessSystem(systemId, messageId string, sendUserId string, code string, msg string, data *interface{}) {
	if len(systemId) > 0 {
		clientIds := Manager.GetSystemClientList(systemId)
		if len(clientIds) > 0 {
			for _, clientId := range clientIds {
				SendToClient(messageId, clientId, sendUserId, code, msg, data)
			}
		}
	}
}

// AddClient2Group 添加到指定分组
func (manager *ClientManager) AddClient2Group(groupName string, client *Client, userId string, extend string) {
	//标记当前客户端的userId
	client.UserId = userId
	client.Extend = extend

	//判断之前是否有添加过
	for _, groupValue := range client.GroupList {
		if groupValue == groupName {
			return
		}
	}

	// 为属性添加分组信息
	groupKey := GenGroupKey(client.SystemId, groupName)

	// 添加到管理分组中
	manager.addClient2Group(groupKey, client)

	// 添加客户端所属分组
	client.GroupList = append(client.GroupList, groupName)

	//mJson, _ := json.Marshal(map[string]string{
	//	"clientId": client.ClientId,
	//	"userId":   client.UserId,
	//	"extend":   client.Extend,
	//})
	//data := string(mJson)
	//sendUserId := ""
	//
	////广播发送通知到业务系统指定分组
	//SendMessage2Group(client.SystemId, sendUserId, groupName, base.WS_ONLINE_MESSAGE_CODE, "客户端上线", &data)
}

// addClient2Group 添加到管理分组中
func (manager *ClientManager) addClient2Group(groupKey string, client *Client) {
	manager.GroupLock.Lock()
	defer manager.GroupLock.Unlock()
	manager.Groups[groupKey] = append(manager.Groups[groupKey], client.ClientId)
}

// delGroupClient 删除分组里的客户端
func (manager *ClientManager) delGroupClient(groupKey string, clientId string) {
	manager.GroupLock.Lock()
	defer manager.GroupLock.Unlock()

	for index, groupClientId := range manager.Groups[groupKey] {
		if groupClientId == clientId {
			manager.Groups[groupKey] = append(manager.Groups[groupKey][:index], manager.Groups[groupKey][index+1:]...)
			break
		}
	}
}

// GetGroupClientList 获取管理分组的客户端列表
func (manager *ClientManager) GetGroupClientList(groupKey string) []string {
	manager.GroupLock.RLock()
	defer manager.GroupLock.RUnlock()
	return manager.Groups[groupKey]
}

// AddClient2SystemClient 添加到业务系统客户端列表
func (manager *ClientManager) AddClient2SystemClient(systemId string, client *Client) {
	manager.SystemClientsLock.Lock()
	defer manager.SystemClientsLock.Unlock()
	manager.SystemClients[systemId] = append(manager.SystemClients[systemId], client.ClientId)
}

// delSystemClient 删除业务系统里的客户端
func (manager *ClientManager) delSystemClient(client *Client) {
	manager.SystemClientsLock.Lock()
	defer manager.SystemClientsLock.Unlock()

	for index, clientId := range manager.SystemClients[client.SystemId] {
		if clientId == client.ClientId {
			manager.SystemClients[client.SystemId] = append(manager.SystemClients[client.SystemId][:index], manager.SystemClients[client.SystemId][index+1:]...)
			break
		}
	}
}

// GetSystemClientList 获取指定业务系统的客户端列表
func (manager *ClientManager) GetSystemClientList(systemId string) []string {
	manager.SystemClientsLock.RLock()
	defer manager.SystemClientsLock.RUnlock()
	return manager.SystemClients[systemId]
}
