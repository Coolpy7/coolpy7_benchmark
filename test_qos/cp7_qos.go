package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//测试逻辑说明
//首先以subtestqos1作为clientid的客户端连接到服务器端
//然后订阅一个qos1主题testqos1topic
//然后断开连接
//新建立另一个客户端连接到服务器，然后以同相的qos1推送个消息到相同主题
//断开推送消息客户端
//以subtestqos1为client id客户端重新登陆，不用做任何动作应该连接成功后会收到pub推送123456即为测试成功
func main() {

	//1.连接成功后按qos订阅主题，然后连接关闭
	//2.新建一个不同clientid的连接，然后按qos推消息到相同主题
	//3.然后以第1步相同clientid并以clearsession为false的状态连接服务器，一切正常会收接到消息
	//测试不同qos修改以下参数
	//go firstsub(1,true)
	//time.Sleep(1*time.Second)
	//go firstpub(1)
	//time.Sleep(1*time.Second)
	go firstsub(1, false)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func firstpub(qos int) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetAutoReconnect(false).SetCleanSession(false).SetClientID("pubber")
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Println("connect local mqtt err:", token.Error())
	}
	if token := c.Publish("foo", uint8(qos), false, []byte("bar")); token.Wait() && token.Error() != nil {
		log.Println("connect local mqtt err:", token.Error())
	}
}

func firstsub(qos int, disconn bool) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetAutoReconnect(false).SetCleanSession(false).SetClientID("subber")
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Println("connect local mqtt err:", token.Error())
	}

	if token := c.Subscribe("foo", uint8(qos), brokerLoadHandler); token.Wait() && token.Error() != nil {
		log.Println("connect local mqtt err:", token.Error())
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		if disconn {
			c.Disconnect(20)
		}
	}()

}

func brokerLoadHandler(client MQTT.Client, msg MQTT.Message) {
	log.Println(msg.Topic(), string(msg.Payload()))
}
