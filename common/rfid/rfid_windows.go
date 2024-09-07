//go:build windows

package rfid

import (
	"sync"
)

type RfidConfig struct {
	Host        string // 主机IP
	Port        int32  // 端口
	Proto       int32  // 协议
	DeviceType  int32  // 设备型号
	Id          string // RFID设备唯一标识
	StationCode string //站口唯一标识
}

type RfidHandler struct {
	Host       string // 主机IP
	Port       int32  // 端口
	Proto      int32  // 协议
	DeviceType int32  // 设备型号
	Fd         int    // 连接fd
	Id         string // RFID设备唯一标识
	//LastTagReportTid  string    // 最后一次上报的标签的tid
	//LastTagReportTime time.Time // 最后一次上报的时间
}

// windows下做个空实现
func Exec() {
	RfidDeviceSet = make(map[string]*RfidDevice)
	RwLock = sync.RWMutex{}
	// do nothing
}

// GetRfidDeviceHandle 获取RFID设备句柄
func GetRfidDeviceHandle(config *RfidConfig) (*RfidHandler, error) {
	handler := &RfidHandler{
		Host:  config.Host,
		Port:  config.Port,
		Proto: config.Proto,
		Id:    config.Id,
	}

	return handler, nil
}

func (r *RfidHandler) ReadRfid() ([]string, error) {
	return nil, nil
}
