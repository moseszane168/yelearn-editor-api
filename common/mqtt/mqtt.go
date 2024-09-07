package mqtt

import (
	"crf-mold/base"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	CLIENT_ID          = "crf-mold-client"
	PUBLISH_TOPIC_NAME = "mda/Stamping/%s/shift/report"
)

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var client mqtt.Client

func GetMqttClient(fn mqtt.MessageHandler, topic string) mqtt.Client {
	broker := viper.GetString("mqtt.broker")
	port := viper.GetInt("mqtt.port")
	userName := viper.GetString("mqtt.userName")
	password := viper.GetString("mqtt.password")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(CLIENT_ID + base.UUID())
	opts.SetUsername(userName)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(fn)
	opts.ResumeSubs = true
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected")
		if token := c.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
			logrus.Error(token.Error())
		}
	}
	opts.OnConnectionLost = connectLostHandler
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client = c
	return c
}

func Publish(lineLevel string, data string) {
	topic := fmt.Sprintf(viper.GetString("mqtt.publishTopic"), lineLevel)
	token := client.Publish(topic, 0, false, data)
	token.Wait()
	time.Sleep(time.Second)
}

func Subscribe(fn mqtt.MessageHandler, topic string) {
	_ = GetMqttClient(fn, topic)
	logrus.Infof("订阅Topic:%s", topic)
}
