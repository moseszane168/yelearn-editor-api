package mqtt_test

import (
	"fmt"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/require"
)

const (
	BROKER               = "172.16.1.171"
	PORT                 = 1883
	USER_NAME            = "admin"
	PASSWORD             = "public"
	CLIENT_ID            = "go_mqtt_client"
	PUBLISH_TOPIC_NAME   = "mda/Stamping/1/shift/report"
	SUBSCRIBE_TOPIC_NAME = "mda/Stamping/+/shift/report"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func getMqttClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", BROKER, PORT))
	opts.SetClientID(CLIENT_ID)
	opts.SetUsername(USER_NAME)
	opts.SetPassword(PASSWORD)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func TestMqtt(t *testing.T) {
	require.NotNil(t, getMqttClient())
	time.Sleep(time.Second * 1000)
}

func TestPublicsh(t *testing.T) {
	client := getMqttClient()
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish(PUBLISH_TOPIC_NAME, 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
	time.Sleep(time.Second * 100)
}

func TestSubscribe(t *testing.T) {
	client := getMqttClient()
	token := client.Subscribe(SUBSCRIBE_TOPIC_NAME, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", SUBSCRIBE_TOPIC_NAME)
	time.Sleep(time.Second * 100)
}
