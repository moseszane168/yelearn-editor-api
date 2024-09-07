package rfid

import (
	"sync"
	"time"
)

type RfidDevice struct {
	Host              string    // 主机IP
	Port              int32     // 端口
	Proto             int32     // 协议
	Fd                int       // 连接fd
	Id                string    // RFID设备唯一标识
	LastTagReportTid  string    // 最后一次上报的标签的tid
	LastTagReportTime time.Time // 最后一次上报的时间
}

// 已连接设备集合
var RfidDeviceSet map[string]*RfidDevice

// RfidDeviceSet读写锁
var RwLock sync.RWMutex

type RfidInfo struct {
	RfidDeviceID string
	TID          string
}
