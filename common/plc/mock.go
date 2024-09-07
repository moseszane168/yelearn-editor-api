package plc

import (
	"github.com/robinson/gos7"
)

type MockClient struct {
	client gos7.Client
}

func NewMockPlcClient() PlcClientInf {
	return &MockClient{}
}

// 获取字节
func (p *MockClient) ReadPlcByte(db int, offset int) (byte, error) {
	return 1, nil
}

// 写入字节
func (p *MockClient) WritePlcByte(db int, offset int, value byte) error {
	return nil
}

// 获取字
func (p *MockClient) ReadPlcWord(db int, offset int) (uint16, error) {
	return 1, nil
}

// 写入字
func (p *MockClient) WritePlcWord(db int, offset int, value uint16) error {
	return nil
}

// 获取双字
func (p *MockClient) ReadPlcDword(db int, offset int) (uint32, error) {
	return 1, nil
}

// 写入双字
func (p *MockClient) WritePlcDword(db int, offset int, value uint32) error {
	return nil
}

// 读位
func (p *MockClient) ReadPlcBit(db, offset, pos int) (bool, error) {
	return true, nil
}

// 写位
func (p *MockClient) WritePlcBit(db, offset, pos int, value bool) (byte, error) {
	return 0, nil
}
