package message

import (
	"crf-mold/dao"
	"crf-mold/model"

	"gorm.io/gorm"
)

type Message struct {
	Content   string
	Operator  []string
	JobId     int64
	RemodelId int64
	Type      string
}

func CreateMessage(tx *gorm.DB, msg Message) {
	batch := make([]model.Message, len(msg.Operator))
	for i := 0; i < len(msg.Operator); i++ {
		operator := msg.Operator[i]

		var entity model.Message
		entity.Content = msg.Content
		entity.Operator = operator
		entity.Status = "unread"
		entity.CreatedBy = "admin"
		entity.UpdatedBy = "admin"
		entity.Type = msg.Type
		entity.JobId = msg.JobId
		entity.RemodelId = msg.RemodelId
		batch[i] = entity
	}

	tx.Table("message").CreateInBatches(&batch, len(batch))
}

func GetUnreadMessageCount(operator string) int64 {
	var count int64
	dao.GetConn().Table("message").Where("status = 'unread' and is_deleted = 'N' and operator = ?", operator).Count(&count)
	return count
}
