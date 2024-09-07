//go:build linux

package rfid

/*
#include <stdlib.h>
#include "./lib/rfid.h"

// 函数声明
void scanCallBack_cgo(rfid_device in);
*/
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// RfidRwLock RfidHandlerSet 读写锁
var RfidRwLock sync.RWMutex

var RfidHandlerSet map[string]*RfidHandler

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

type RfidConfig struct {
	Host        string // 主机IP
	Port        int32  // 端口
	Proto       int32  // 协议
	DeviceType  int32  // 设备型号
	Id          string // RFID设备唯一标识
	StationCode string //站口唯一标识
}

func getErrorDescription(code C.int) string {
	var msg [512]byte
	C.error_description(code, (*C.char)(unsafe.Pointer((&msg[0]))))
	return string(msg[:])
}

// GetRfidDeviceHandle 获取RFID设备句柄
func GetRfidDeviceHandle(config *RfidConfig) (*RfidHandler, error) {
	RfidRwLock.RLock()
	if RfidHandlerSet == nil {
		RfidHandlerSet = make(map[string]*RfidHandler)
	}
	if handler, ok := RfidHandlerSet[config.Id]; ok {
		// 存在
		RfidRwLock.RUnlock()
		fmt.Printf("从Set中获取RFID操作句柄, IP: %s Fd: %d", handler.Host, handler.Fd)
		return handler, nil
	}
	RfidRwLock.RUnlock()
	var handler *RfidHandler
	var err error
	if config.DeviceType == 65536 {
		handler, err = NewRfidP218DeviceHandle(config)
	} else {
		handler, err = NewRfidDeviceHandle(config)
	}
	if err != nil {
		fmt.Println(errors.New(config.StationCode + "站台RFID句柄初始化失败"))
		fmt.Println(err.Error())
		return handler, err
	} else {
		// 添加到rfid操作句柄集合中
		RfidRwLock.RLock()
		fmt.Printf("将RFID操作句柄添加到Set中, IP: %s Fd: %d", handler.Host, handler.Fd)
		RfidHandlerSet[handler.Id] = handler
		RfidRwLock.RUnlock()
		return handler, nil
	}
}

// NewRfidDeviceHandle 初始化RFID设备句柄
func NewRfidDeviceHandle(config *RfidConfig) (*RfidHandler, error) {
	// 建立连接
	var fd C.int
	cHostStr := C.CString(config.Host)
	defer C.free(unsafe.Pointer(cHostStr)) // must call

	handler := &RfidHandler{
		Host:  config.Host,
		Port:  config.Port,
		Proto: config.Proto,
		Id:    config.Id,
	}

	ret := C.connect_net(cHostStr, C.int(config.Port), C.transport_protocol(config.Proto), &fd)
	if ret != 0 {
		fmt.Printf("RFID设备连接失败, IP: %s Error: %s", config.Host, getErrorDescription(ret))
		time.Sleep(time.Second)
		return handler, errors.New(getErrorDescription(ret))
	}

	// 设置连接fd
	handler.Fd = int(fd)
	// 获取工作模式
	var mode C.work_mode
	if ret := C.get_work_mode(C.int(int(fd)), &mode); ret != 0 {
		return handler, errors.New(fmt.Sprintf("RFID设备获取工作模式失败, IP: %s Error: %s", config.Host, getErrorDescription(ret)))
	}
	// 设置工作模式(设置为命令模式)
	if ret := C.set_work_mode(C.int(int(fd)), C.command_mode); ret != 0 {
		return handler, errors.New(fmt.Sprintf("RFID设备设置工作模式失败, IP: %s Error: %s", config.Host, getErrorDescription(ret)))
	}
	return handler, nil
}

// NewRfidP218DeviceHandle 初始化RFID型号为P218设备的句柄
func NewRfidP218DeviceHandle(config *RfidConfig) (*RfidHandler, error) {
	// 建立连接
	var fd C.int
	cHostStr := C.CString(config.Host)
	defer C.free(unsafe.Pointer(cHostStr)) // must call

	//C.set_log_level(6)

	handler := &RfidHandler{
		Host:       config.Host,
		Port:       config.Port,
		Proto:      config.Proto,
		DeviceType: config.DeviceType,
		Id:         config.Id,
	}

	ret := C.connect_p218_udp(cHostStr, C.int(config.Port), C.int(config.DeviceType), &fd)
	// UDP 实际上是无连接的, 所以这里ret都会是0
	fmt.Printf("Connect Ret %d", int(ret))
	if ret != 0 {
		fmt.Printf("RFID设备连接失败, IP: %s Error: %s", config.Host, getErrorDescription(ret))
		time.Sleep(time.Second)
		return handler, errors.New(getErrorDescription(ret))
	}

	// 设置连接fd
	handler.Fd = int(fd)
	// 获取工作模式
	var mode C.work_mode
	if ret := C.get_work_mode(C.int(int(fd)), &mode); ret != 0 {
		return handler, errors.New(fmt.Sprintf("RFID设备获取工作模式失败, IP: %s Error: %s", config.Host, getErrorDescription(ret)))
	}
	// 设置工作模式(设置为命令模式)
	if ret := C.set_work_mode(C.int(int(fd)), C.command_mode); ret != 0 {
		return handler, errors.New(fmt.Sprintf("RFID设备设置工作模式失败, IP: %s Error: %s", config.Host, getErrorDescription(ret)))
	}
	return handler, nil
}

func (r *RfidHandler) ReadRfid() ([]string, error) {
	fmt.Printf("RFID设备IP: %s, Fd: %d\n", r.Host, r.Fd)
	tids := []string{}
	var tags [10]C.tag
	var count C.int
	var ret C.int
	password := []byte{0, 0, 0, 0}
	ret = C.list6c(C.int(r.Fd), C.memory_bank_tid, C.int(0), C.int(12), (*C.uchar)(unsafe.Pointer(&password[0])), C.int(4), (*C.tag)(unsafe.Pointer(&tags[0])), &count)
	if ret == 0 {
		fmt.Printf("设备ID: %s\n", r.Id)
		fmt.Println("列出附近的标签成功", count)
		for j := 0; j < int(count); j++ {
			var epc_hex_str [1024]byte
			C.to_hex_str((*C.uchar)(unsafe.Pointer((&((tags[j].id)[0])))), tags[j].len, (*C.char)(unsafe.Pointer(&epc_hex_str)))
			tid := string(epc_hex_str[:])
			fmt.Println("RFID读取到TID: ", tid[:24])
			tids = append(tids, tid[:24])
		}
		return tids, nil
	} else if ret != 24 {
		r.CloseRfidHandle()
		RfidRwLock.RLock()
		fmt.Printf("读写异常，从Set中删除RFID操作句柄, IP: %s Fd: %d", r.Host, r.Fd)
		delete(RfidHandlerSet, r.Id)
		RfidRwLock.RUnlock()
		return nil, errors.New(fmt.Sprintf("RFID读取异常list6c err, IP: %s ", r.Host))
	}
	return nil, nil
}

func (r *RfidHandler) CloseRfidHandle() {
	if r.Fd != 0 {
		C.disconnect(C.int(r.Fd))
	}
}

//
////export scanCallBack
//func scanCallBack(device C.rfid_device) {
//	id := C.GoString(device.id)
//
//	RwLock.RLock()
//	if _, ok := RfidDeviceSet[id]; ok {
//		// 存在
//		return
//	}
//	RwLock.RUnlock()
//
//	RwLock.Lock()
//	host := C.GoString(device.ip)
//	port := int32(device.port)
//	proto := int32(device.proto)
//
//	item := &RfidDevice{
//		Host:  host,
//		Port:  port,
//		Proto: proto,
//		Id:    id,
//	}
//	RfidDeviceSet[id] = item
//
//	go func(item *RfidDevice) {
//		defer func() {
//			if e := recover(); e != nil {
//				// 断开连接
//				debug.PrintStack()
//				fmt.Println(e)
//				RwLock.Lock()
//				delete(RfidDeviceSet, item.Id)
//				RwLock.Unlock()
//			}
//		}()
//
//		var fd C.int
//		cHostStr := C.CString(item.Host)
//		defer C.free(unsafe.Pointer(cHostStr)) // must cal
//		ret := C.connect_net(cHostStr, C.int(item.Port), C.transport_protocol(item.Proto), &fd)
//		if ret != 0 {
//			fmt.Println("连接失败:", getErrorDescription(ret))
//			panic("connect_net err")
//		}
//		// 设置连接fd
//		item.Fd = int(fd)
//		fmt.Printf("RFID host: %s, fd: %d , proto: %d", item.Host, int(fd), int(item.Proto))
//
//		// 获取工作模式
//		var mode C.work_mode
//		ret = C.get_work_mode(C.int(item.Fd), &mode)
//		if ret != 0 {
//			fmt.Println("get_work_mode失败", getErrorDescription(ret))
//			panic("get_work_mode err")
//		}
//
//		// 设置工作模式
//		ret = C.set_work_mode(C.int(item.Fd), C.command_mode)
//		if ret != 0 {
//			fmt.Println("set_work_mode 失败", getErrorDescription(ret))
//			panic("set_work_mode err")
//		}
//
//		// 列出附件标签
//		var tags [10]C.tag
//		var count C.int
//		password := []byte{0, 0, 0, 0}
//		for {
//			ret = C.list6c(C.int(item.Fd), C.memory_bank_tid, C.int(0), C.int(12), (*C.uchar)(unsafe.Pointer(&password[0])), C.int(4), (*C.tag)(unsafe.Pointer(&tags[0])), &count)
//			if ret != 0 {
//				if ret != 24 {
//					panic("list6c err")
//				}
//			} else {
//				fmt.Printf("设备ID: %s", item.Id)
//				fmt.Println("列出附近的标签成功", count)
//
//				for j := 0; j < int(count); j++ {
//					var epc_hex_str [1024]byte
//					C.to_hex_str((*C.uchar)(unsafe.Pointer((&((tags[j].id)[0])))), tags[j].len, (*C.char)(unsafe.Pointer(&epc_hex_str)))
//					tid := string(epc_hex_str[:])
//					fmt.Println("TID: ", tid[:24])
//					item.LastTagReportTid = tid[:24]
//					item.LastTagReportTime = time.Now()
//					rfidDevice := RfidInfo{
//						RfidDeviceID: item.Id,
//						TID:          item.LastTagReportTid,
//					}
//					_ = SetRfidToRedis(rfidDevice)
//				}
//			}
//			time.Sleep(time.Second)
//		}
//	}(item)
//
//	RwLock.Unlock()
//}
//
//// InitRfidDeviceHandle 初始化RFID设备句柄
//func (r *RfidDevice) InitRfidDeviceHandle() (*RfidDevice, error) {
//	// 建立连接
//	var fd C.int
//	cHostStr := C.CString(r.Host)
//	defer C.free(unsafe.Pointer(cHostStr)) // must call
//	if ret := C.connect_net(cHostStr, C.int(r.Port), C.transport_protocol(r.Proto), &fd); ret != 0 {
//		return r, errors.New(fmt.Sprintf("RFID设备连接失败, IP: %s Error: %s", r.Host, getErrorDescription(ret)))
//	}
//	// 设置连接fd
//	r.Fd = int(fd)
//	// 获取工作模式
//	var mode C.work_mode
//	if ret := C.get_work_mode(C.int(int(fd)), &mode); ret != 0 {
//		return nil, errors.New(fmt.Sprintf("RFID设备获取工作模式失败, IP: %s Error: %s", r.Host, getErrorDescription(ret)))
//	}
//	// 设置工作模式(设置为命令模式)
//	if ret := C.set_work_mode(C.int(int(fd)), C.command_mode); ret != 0 {
//		return nil, errors.New(fmt.Sprintf("RFID设备设置工作模式失败, IP: %s Error: %s", r.Host, getErrorDescription(ret)))
//	}
//	return r, nil
//}
//func (r *RfidDevice) ReadRfid() (string, error) {
//	fmt.Printf("RFID设备IP: %s, Fd: %d", r.Host, r.Fd)
//	var tags [10]C.tag
//	var count C.int
//	var ret C.int
//	password := []byte{0, 0, 0, 0}
//	ret = C.list6c(C.int(r.Fd), C.memory_bank_tid, C.int(0), C.int(12), (*C.uchar)(unsafe.Pointer(&password[0])), C.int(4), (*C.tag)(unsafe.Pointer(&tags[0])), &count)
//	if ret == 0 {
//		fmt.Printf("设备ID: %s", r.Id)
//		fmt.Println("列出附近的标签成功", count)
//		for j := 0; j < int(count); j++ {
//			var epc_hex_str [1024]byte
//			C.to_hex_str((*C.uchar)(unsafe.Pointer((&((tags[j].id)[0])))), tags[j].len, (*C.char)(unsafe.Pointer(&epc_hex_str)))
//			tid := string(epc_hex_str[:])
//			fmt.Println("RFID读取到TID: ", tid[:24])
//			return tid[:24], nil
//		}
//	} else if ret != 24 {
//		return "", errors.New(fmt.Sprintf("RFID读取异常list6c err, IP: %s ", r.Host))
//	}
//	return "", nil
//}
//
//func Exec() {
//	// 已连接设备集合
//	RfidDeviceSet = make(map[string]*RfidDevice, 5)
//	RwLock = sync.RWMutex{}
//
//	// 开启任务扫描当前网段的设备
//	r := C.start_device_scan(C.rfid_device_scan_callback(C.scanCallBack_cgo))
//	if r != 0 {
//		fmt.Println(getErrorDescription(r))
//		panic("start_device_scan error:")
//	}
//}
//
//func getErrorDescription(code C.int) string {
//	var msg [512]byte
//	C.error_description(code, (*C.char)(unsafe.Pointer((&msg[0]))))
//	return string(msg[:])
//}

//func getCoreVersion() string {
//	return C.GoString(C.core_version())
//}
