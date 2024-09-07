package message

import (
	"crf-mold/base"
)

type PageMessageInputVO struct {
	base.PageVO
	CodeOrName string `form:"codeOrName"`
} // @name PageMessageInputVO

type PageMessageOututVO struct {
	ID         int64     `json:"id"`
	Status     string    `json:"status"`     // 状态:READ/UNREAD
	Content    string    `json:"content"`    // 消息内容
	JobId      int64     `json:"jobId"`      // 保养任务ID
	RemodelId  int64     `json:"remodelId"`  // 改造任务ID
	Type       string    `json:"type"`       // 消息类型：maintenance保养任务、remodel改造任务
	GmtCreated base.Time `json:"gmtCreated"` // 创建时间
} // @name PageMessageOututVO
