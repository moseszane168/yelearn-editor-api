package productresume

import (
	"crf-mold/dao"
	"crf-mold/model"
	"encoding/json"
	"runtime/debug"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

const SUBSCRIBE_TOPIC_NAME = "mda/Stamping/+/shift/report"

type MDA_PART struct {
	PartNo string
	QtyOK  int64
	QtyNOK int64
}

type MDA struct {
	StartTime string     `json:"StartTime"` //班次开始
	EndTime   string     `json:"EndTime"`   //班次结束
	ShiftId   int64      `json:"ShiftId"`   //班次编号
	ShiftNo   int64      `json:"ShiftNo"`   //班次类型
	ShiftMin  int64      `json:"ShiftMin"`
	ShiftDate string     `json:"ShiftDate"` //班次日期
	PartList  []MDA_PART `json:"PartList"`  //详细产量
	TimeStamp string     `json:"TimeStamp"`
}

// 生产履历数据处理
func ProductionResumeMessagePubHandler(_ mqtt.Client, msg mqtt.Message) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error("ProductionResumeMessagePubHandler error Panic info is: %v", err)
			debug.PrintStack()
			logrus.WithField("stack", string(debug.Stack())).Errorf("panic: %v\n", err)
		}
	}()

	payLoad := msg.Payload()
	data := string(payLoad)

	logrus.WithField("data", data).Info("recv mqtt data")

	// 保存记录生产数据
	var mqttHistory model.MqttHistory
	mqttHistory.Data = data
	dao.GetConn().Table("mqtt_history").Save(&mqttHistory)

	// 反序列化
	var mda MDA
	err := json.Unmarshal([]byte(payLoad), &mda)
	if err != nil {
		logrus.WithField("data", data).Error("Mqtt message format error")
		return
	}

	// 接收到生产数据
	topic := msg.Topic()
	logrus.WithField("topic", topic).WithField("data", data).Info("接收到生产数据")

	go HandleProductionResume(topic, mda)
}

func HandleProductionResume(topic string, mda MDA) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error("HandleProductionResume error Panic info is: %v", err)
			debug.PrintStack()
			logrus.WithField("stack", string(debug.Stack())).Errorf("panic: %v\n", err)
		}
	}()

	pr := NewProductResume(topic, mda)

	tx := dao.GetConn().Begin()
	defer dao.TransactionRollback(tx)

	// 3、保存生产履历
	pr.SaveProductionResume(tx)

	// 添加产线生产履历记录
	pr.SaveLineProductionResume(tx)

	tx.Commit()

}
