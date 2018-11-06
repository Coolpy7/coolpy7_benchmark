package main

import (
	"flag"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

// 大文件传输测试工具
// 本工具建立一个订单及通过一个发布测试一个大文件通信质量

func main() {
	var ip = flag.String("ip", "192.168.100.2", "ip")
	flag.Parse()

	opts := MQTT.NewClientOptions().AddBroker("tcp://" + *ip + ":1883")
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Println("connect local mqtt err:", token.Error())
	}

	if token := c.Subscribe("bigfile", 0, brokerLoadHandler); token.Wait() && token.Error() != nil {
		log.Println("connect local mqtt err:", token.Error())
	}

	bts, _ := ioutil.ReadFile("E:\\迅雷下载\\HBuilder.9.0.6.windows.zip")
	log.Println("send size:", len(bts))
	c.Publish("bigfile", 0, false, bts)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func brokerLoadHandler(client MQTT.Client, msg MQTT.Message) {
	log.Println(msg.Topic(), len(msg.Payload()))
}
