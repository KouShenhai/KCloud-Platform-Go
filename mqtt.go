package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"strings"
	"time"
)

type MQTT struct {
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
	// 客户端ID
	ClientId string `json:"clientId"`
	// 主机
	Host string `json:"host"`
	// 端口
	Port string `json:"port"`
	// 主题
	Topic string `json:"topic"`
	// Qos【0/1/2】 => 值越大延迟越高
	Qos byte `json:"qos"`
}

func GetMqttClient(config MQTT, callback mqtt.MessageHandler) mqtt.Client {
	broker := fmt.Sprintf("tcp://%s:%s", config.Host, config.Port)
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID(config.ClientId)
	options.SetUsername(config.Username)
	options.SetPassword(config.Password)
	options.OnConnect = onConnect
	// 开启自动重连
	options.SetAutoReconnect(true)
	options.OnConnectionLost = onConnectLost
	options.SetDefaultPublishHandler(defaultPublishHandler)
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("connect failed", zap.Error(token.Error()))
		time.Sleep(5 * time.Second)
		return GetMqttClient(config, callback)
	}
	// 订阅topic
	topics := strings.Split(config.Topic, ",")
	var filters = make(map[string]byte)
	for _, topic := range topics {
		filters[topic] = config.Qos
	}
	client.SubscribeMultiple(filters, callback)
	return client
}

func DisConnectMQTT(client mqtt.Client) {
	client.Disconnect(30 * 1000)
}

func defaultPublishHandler(mqtt.Client, mqtt.Message) {}

func onConnect(mqtt.Client) {}

func PublishMQTT(client mqtt.Client, topic string, qos byte, retained bool, payload interface{}) {
	client.Publish(topic, qos, retained, payload)
}

func DefaultPublishMQTT(client mqtt.Client, topic string, payload interface{}) {
	PublishMQTT(client, topic, 0, false, payload)
}

func onConnectLost(client mqtt.Client, err error) {
	fmt.Println("mqtt connect lost", zap.Error(err))
	// 重连
	for i := 0; i < 5; i++ {
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			fmt.Println("reconnect successfully")
			break
		} else {
			fmt.Println("reconnect failed", zap.Any("attempt", i+1), zap.Error(token.Error()))
			time.Sleep(10)
		}
	}
}
