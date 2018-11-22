package main

import (
	"client"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"packet"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var urlString = flag.String("url", "tcp://127.0.0.1:1883", "broker url")
var topic = flag.String("topic", "cp7sub%i", "the used topic")
var workers = flag.Int("workers", 100, "number of workers")
var cs = flag.String("cid", "testclient", "client id start with")
var qos = flag.Uint("qos", 0, "sub qos level")
var ping = flag.Int("ping", 0, "disable ping")

func main() {
	flag.Parse()

	for i := 0; i < *workers; i++ {
		id := strconv.Itoa(i)
		if i%1000 == 0 {
			log.Println(id)
		}

		cl := client.New()
		if *ping == 1 {
			cl.Ping = false
		}
		cl.Callback = func(msg *packet.Message, err error) error {
			if err != nil {
				log.Println("callback", err)
			}
			log.Println(msg.Topic, msg.QOS, len(msg.Payload))
			return nil
		}

		cf, err := cl.Connect(&client.Config{
			BrokerURL:    *urlString,
			CleanSession: false,
			KeepAlive:    "30s",
			ValidateSubs: true,
			ClientID:     *cs + id,
		})
		if err != nil {
			log.Println("conn", err)
		}

		err = cf.Wait(10 * time.Second)
		if err != nil {
			log.Println("conn wait", err)
		}

		if *topic != "null" {
			sf, err := cl.Subscribe(strings.Replace(*topic, "%i", id, 1), uint8(*qos))
			if err != nil {
				log.Println("sub", err)
			}

			err = sf.Wait(10 * time.Second)
			if err != nil {
				log.Println("sub wait", err)
			}
		}
	}

	fmt.Println("finish", *workers)

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
