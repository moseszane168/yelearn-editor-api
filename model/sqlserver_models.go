package model

import "crf-mold/base"

// CY_WMS_Request 产线呼叫WMS
type CY_WMS_Request struct {
	TaskId       string `gorm:"primaryKey;column:TaskId;type:nvarchar(20);not null" json:"-"`
	TaskNo       string `gorm:"column:TaskNo;type:nvarchar(20)" json:"TaskNo"`
	TaskType     int    `gorm:"column:TaskType;type:int" json:"TaskType"`
	TaskStatus   int    `gorm:"column:TaskStatus;type:int" json:"TaskStatus"`
	CreateTime   string `gorm:"column:CreateTime;type:datetime" json:"CreateTime"`
	UpdateTime   string `gorm:"column:UpdateTime;type:datetime" json:"UpdateTime"`
	ErrorCode    string `gorm:"column:ErrorCode;type:nvarchar(10)" json:"ErrorCode"`
	ErrorMsg     string `gorm:"column:ErrorMsg;type:nvarchar(200)" json:"ErrorMsg"`
	Flag         int    `gorm:"column:Flag;type:int" json:"Flag"`
	PartNo       string `gorm:"column:PartNo;type:nvarchar(30)" json:"PartNo"`
	DeviceStatus int    `gorm:"column:DeviceStatus;type:int" json:"DeviceStatus"`
	BInCode      string `gorm:"column:BInCode;type:nvarchar(30)" json:"BinCode"`
}

// LM_WMS_CallBack 立库呼叫WMS
type LM_WMS_CallBack struct {
	TaskId     string `gorm:"primaryKey;column:TaskId;type:nvarchar(20);not null" json:"-"`
	TaskNo     string `gorm:"column:TaskNo;type:nvarchar(20)" json:"TaskNo"`
	TaskType   int    `gorm:"column:TaskType;type:int" json:"TaskType"`
	TaskStatus int    `gorm:"column:TaskStatus;type:int" json:"TaskStatus"`
	PartNo     string `gorm:"column:PartNo;type:nvarchar(30)" json:"PartNo"`
	BinCode    string `gorm:"column:BinCode;type:nvarchar(30)" json:"BinCode"`

	CreateTime string `gorm:"column:CreateTime;type:datetime" json:"CreateTime"`
	UpdateTime string `gorm:"column:UpdateTime;type:datetime" json:"UpdateTime"`
	ErrorCode  string `gorm:"column:ErrorCode;type:nvarchar(10)" json:"ErrorCode"`
	ErrorMsg   string `gorm:"column:ErrorMsg;type:nvarchar(200)" json:"ErrorMsg"`
	Flag       int    `gorm:"column:Flag;type:int" json:"Flag"`
}

// PD_WMS_Request 产线呼叫WMS中间表
type PD_WMS_Request struct {
	TaskId       string     `gorm:"primaryKey;column:TaskId;type:nvarchar(50);not null" json:"-"`
	TaskNo       string     `gorm:"column:TaskNo;type:nvarchar(50)" json:"TaskNo"`
	TaskType     int        `gorm:"column:TaskType;type:int" json:"TaskType"`
	TaskStatus   *int       `gorm:"column:TaskStatus;type:int" json:"TaskStatus"`
	CreateTime   base.Time  `gorm:"column:CreateTime;type:datetime" json:"CreateTime"`
	UpdateTime   *base.Time `gorm:"column:UpdateTime;type:datetime" json:"UpdateTime"`
	ErrorCode    string     `gorm:"column:ErrorCode;type:nvarchar(10)" json:"ErrorCode"`
	ErrorMsg     string     `gorm:"column:ErrorMsg;type:nvarchar(200)" json:"ErrorMsg"`
	Flag         int        `gorm:"column:Flag;type:int" json:"Flag"`
	PartNo       string     `gorm:"column:PartNo;type:nvarchar(20)" json:"PartNo"`
	DeviceStatus int        `gorm:"column:DeviceStatus;type:int;default:1" json:"DeviceStatus"`
	BinCode      string     `gorm:"column:BinCode;type:nvarchar(30)" json:"BinCode"`
}

// WMS_ASRS_Request WMS-ASRS中间表
type WMS_ASRS_Request struct {
	TaskId     string     `gorm:"primaryKey;column:TaskId;type:nvarchar(50);not null" json:"-"`
	TaskNo     string     `gorm:"column:TaskNo;type:nvarchar(50)" json:"TaskNo"`
	TaskType   int        `gorm:"column:TaskType;type:int;not null" json:"TaskType"`
	TaskStatus *int       `gorm:"column:TaskStatus;type:int" json:"TaskStatus"`
	CreateTime base.Time  `gorm:"column:CreateTime;type:datetime" json:"CreateTime"`
	UpdateTime *base.Time `gorm:"column:UpdateTime;type:datetime" json:"UpdateTime"`
	ErrorCode  string     `gorm:"column:ErrorCode;type:nvarchar(10)" json:"ErrorCode"`
	ErrorMsg   string     `gorm:"column:ErrorMsg;type:nvarchar(200)" json:"ErrorMsg"`
	Flag       int        `gorm:"column:Flag;type:int" json:"Flag"`
	PartNo     string     `gorm:"column:PartNo;type:nvarchar(50)" json:"PartNo"`
	BinCode    string     `gorm:"column:BinCode;type:nvarchar(20)" json:"BinCode"`
}

// ASRS_WMS_CallBack ASRS_WMS回调中间表
type ASRS_WMS_CallBack struct {
	TaskId       string     `gorm:"primaryKey;column:TaskId;type:nvarchar(50);not null" json:"-"`
	TaskNo       string     `gorm:"column:TaskNo;type:nvarchar(50)" json:"TaskNo"`
	TaskType     int        `gorm:"column:TaskType;type:int;not null" json:"TaskType"`
	SrcLocation  string     `gorm:"column:SrcLocation;type:nvarchar(20)" json:"SrcLocation"`
	DestLocation string     `gorm:"column:DestLocation;type:nvarchar(20)" json:"DestLocation"`
	PartNo       string     `gorm:"column:PartNo;type:nvarchar(50)" json:"PartNo"`
	TaskStatus   *int       `gorm:"column:TaskStatus;type:int" json:"TaskStatus"`
	CreateTime   base.Time  `gorm:"column:CreateTime;type:datetime" json:"CreateTime"`
	UpdateTime   *base.Time `gorm:"column:UpdateTime;type:datetime" json:"UpdateTime"`
	ErrorCode    string     `gorm:"column:ErrorCode;type:nvarchar(10)" json:"ErrorCode"`
	ErrorMsg     string     `gorm:"column:ErrorMsg;type:nvarchar(200)" json:"ErrorMsg"`
	Flag         int        `gorm:"column:Flag;type:int" json:"Flag"`
}
