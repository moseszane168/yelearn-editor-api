package plc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/robinson/gos7"
	"sync"
	"time"
)

var helper gos7.Helper

type PlcClientInf interface {
	ReadPlcByte(db int, offset int) (byte, error)
	WritePlcByte(db int, offset int, value byte) error
	ReadPlcWord(db int, offset int) (uint16, error)
	WritePlcWord(db int, offset int, value uint16) error
	ReadPlcDword(db int, offset int) (uint32, error)
	WritePlcDword(db int, offset int, value uint32) error
	ReadPlcBit(db, offset, pos int) (bool, error)
	WritePlcBit(db, offset, pos int, value bool) (byte, error)
}

type PlcClient struct {
	client gos7.Client
}

// PLCRwLock PLCHandlerSet 读写锁
var PLCRwLock sync.RWMutex

var PLCHandlerSet map[string]*PLCHandler

type PLCHandler struct {
	PlcClientInf
	Host string // 主机IP
	Rack int    // rack
	Slot int    // slot
	Db   int
	Id   string // PLC设备唯一标识
}

type PLCConfig struct {
	Host string // 主机IP
	Rack int    // rack
	Slot int    // slot
	Db   int    //db块
	Id   string // PLC设备唯一标识
}

func NewPlcClient(tcpDevice string, rack, slot int) (PlcClientInf, error) {
	// 建立连接
	handler := gos7.NewTCPClientHandler(tcpDevice, rack, slot)

	// 设置连接属性
	handler.Timeout = 200 * time.Second
	handler.IdleTimeout = 200 * time.Second
	//handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)

	// 创建连接，可以基于这个新建多个会话（client）
	err := handler.Connect()
	p := &PlcClient{}
	if err != nil {
		return p, err
	}
	// 新建一个会话
	p.client = gos7.NewClient(handler)
	return p, nil
}

// NewPLCDeviceHandle 初始化PLC设备句柄
func NewPLCDeviceHandle(config *PLCConfig) (*PLCHandler, error) {
	clientInterface, err := NewPlcClient(config.Host, config.Rack, config.Slot)
	if err != nil {
		fmt.Printf("PLC设备连接失败, IP: %s Error: %s\n", config.Host, err.Error())
		time.Sleep(time.Second)
	}
	handler := &PLCHandler{
		PlcClientInf: clientInterface,
		Host:         config.Host,
		Rack:         config.Rack,
		Slot:         config.Slot,
		Db:           config.Db,
		Id:           config.Id,
	}
	return handler, err
}

// GetPLCDeviceHandle 获取PLC设备句柄
func GetPLCDeviceHandle(config *PLCConfig) (*PLCHandler, error) {
	PLCRwLock.RLock()
	if PLCHandlerSet == nil {
		PLCHandlerSet = make(map[string]*PLCHandler)
	}
	if handler, ok := PLCHandlerSet[config.Id]; ok {
		// 存在
		PLCRwLock.RUnlock()
		return handler, nil
	}
	PLCRwLock.RUnlock()
	handler, err := NewPLCDeviceHandle(config)
	if err != nil {
		fmt.Println(errors.New(config.Id + "PLC设备句柄初始化失败"))
		fmt.Println(err.Error())
		return handler, err
	} else {
		// 添加到rfid操作句柄集合中
		PLCRwLock.RLock()
		fmt.Printf("将PLC设备操作句柄添加到Set中, IP: %s ID: %s", handler.Host, handler.Id)
		PLCHandlerSet[handler.Id] = handler
		PLCRwLock.RUnlock()
		return handler, nil
	}
}

// PLCReadWord PLC读取字
func (p *PLCHandler) PLCReadWord(db int, offset int) (uint16, error) {
	r, err := p.ReadPlcWord(db, offset)
	if err != nil {
		PLCRwLock.RLock()
		fmt.Printf("读写异常，从Set中删除PLC操作句柄, IP: %s Id: %s", p.Host, p.Id)
		delete(PLCHandlerSet, p.Id)
		PLCRwLock.RUnlock()
	}
	return r, err
}

// PLCWriteWord PLC写字
func (p *PLCHandler) PLCWriteWord(db int, offset int, value uint16) error {
	err := p.WritePlcWord(db, offset, value)
	if err != nil {
		PLCRwLock.RLock()
		fmt.Printf("读写异常，从Set中删除PLC操作句柄, IP: %s Id: %s", p.Host, p.Id)
		delete(PLCHandlerSet, p.Id)
		PLCRwLock.RUnlock()
	}
	return err
}

// PLCReadDword PLC读取双字
func (p *PLCHandler) PLCReadDword(db int, offset int) (uint32, error) {
	r, err := p.ReadPlcDword(db, offset)
	if err != nil {
		PLCRwLock.RLock()
		fmt.Printf("读写异常，从Set中删除PLC操作句柄, IP: %s Id: %s", p.Host, p.Id)
		delete(PLCHandlerSet, p.Id)
		PLCRwLock.RUnlock()
	}
	return r, err
}

// PLCWriteDword PLC写双字
func (p *PLCHandler) PLCWriteDword(db int, offset int, value uint32) error {
	err := p.WritePlcDword(db, offset, value)
	if err != nil {
		PLCRwLock.RLock()
		fmt.Printf("读写异常，从Set中删除PLC操作句柄, IP: %s Id: %s", p.Host, p.Id)
		delete(PLCHandlerSet, p.Id)
		PLCRwLock.RUnlock()
	}
	return err
}

// PLCReadBit PLC读取位
func (p *PLCHandler) PLCReadBit(db, offset, pos int) (bool, error) {
	r, err := p.ReadPlcBit(db, offset, pos)
	if err != nil {
		PLCRwLock.RLock()
		fmt.Printf("读写异常，从Set中删除PLC操作句柄, IP: %s Id: %s", p.Host, p.Id)
		delete(PLCHandlerSet, p.Id)
		PLCRwLock.RUnlock()
	}
	return r, err
}

// PLCWriteBit PLC写位
func (p *PLCHandler) PLCWriteBit(db, offset, pos int, value bool) (byte, error) {
	r, err := p.WritePlcBit(db, offset, pos, value)
	if err != nil {
		PLCRwLock.RLock()
		fmt.Printf("读写异常，从Set中删除PLC操作句柄, IP: %s Id: %s", p.Host, p.Id)
		delete(PLCHandlerSet, p.Id)
		PLCRwLock.RUnlock()
	}
	return r, err
}

// 获取字节
func (p *PlcClient) ReadPlcByte(db int, offset int) (byte, error) {
	if p.client == nil {
		return 0, errors.New("[PLC操作client为nil]")
	}
	buffer := make([]byte, 1)
	err := p.client.AGReadDB(db, offset, 1, buffer)
	if err != nil {
		return 0, err
	}

	return buffer[0], nil
}

// 写入字节
func (p *PlcClient) WritePlcByte(db int, offset int, value byte) error {
	if p.client == nil {
		return errors.New("[PLC操作client为nil]")
	}
	buffer := make([]byte, 1)
	buffer[0] = value
	err := p.client.AGWriteDB(db, offset, 1, buffer)
	if err != nil {
		return err
	}
	return nil
}

// 获取字
func (p *PlcClient) ReadPlcWord(db int, offset int) (uint16, error) {
	if p.client == nil {
		return 0, errors.New("[PLC操作client为nil]")
	}
	buffer := make([]byte, 2)
	err := p.client.AGReadDB(db, offset, 2, buffer)
	if err != nil {
		return 0, err
	}

	// 字节转整数
	x := uint16(0)
	bytesBuffer := bytes.NewBuffer(buffer)
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x, nil
}

// 写入字
func (p *PlcClient) WritePlcWord(db int, offset int, value uint16) error {
	if p.client == nil {
		return errors.New("[PLC操作client为nil]")
	}
	// 整数转字节
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, value)

	buffer := bytesBuffer.Bytes()
	err := p.client.AGWriteDB(db, offset, 2, buffer)
	if err != nil {
		return err
	}

	return nil
}

// 获取双字
func (p *PlcClient) ReadPlcDword(db int, offset int) (uint32, error) {
	if p.client == nil {
		return 0, errors.New("[PLC操作client为nil]")
	}
	buffer := make([]byte, 4)
	err := p.client.AGReadDB(db, offset, 4, buffer)
	if err != nil {
		return 0, err
	}

	// 字节转整数
	x := uint32(0)
	bytesBuffer := bytes.NewBuffer(buffer)
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x, nil
}

// 写入双字
func (p *PlcClient) WritePlcDword(db int, offset int, value uint32) error {
	if p.client == nil {
		return errors.New("[PLC操作client为nil]")
	}
	// 整数转字节
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, value)

	buffer := bytesBuffer.Bytes()
	err := p.client.AGWriteDB(db, offset, 4, buffer)
	if err != nil {
		return err
	}

	return nil
}

// 读位
func (p *PlcClient) ReadPlcBit(db, offset, pos int) (bool, error) {
	if p.client == nil {
		return false, errors.New("[PLC操作client为nil]")
	}
	b, err := p.ReadPlcByte(db, offset)
	if err != nil {
		return false, err
	}

	return helper.GetBoolAt(b, uint(pos)), nil
}

// 写位
func (p *PlcClient) WritePlcBit(db, offset, pos int, value bool) (byte, error) {
	if p.client == nil {
		return 0, errors.New("[PLC操作client为nil]")
	}
	// 读取旧的字节
	b, err := p.ReadPlcByte(db, offset)
	if err != nil {
		return 0, err
	}

	// 修改旧值中的位
	newValue := helper.SetBoolAt(b, uint(pos), value)
	p.WritePlcByte(db, offset, newValue)
	return b, nil
}
