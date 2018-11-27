package main

import (
	"client"
	"flag"
	"fmt"
	"gopkg.in/robfig/cron.v2"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var urlString = flag.String("url", "tcp://username:password@192.168.100.2:1883", "broker url")
var topic = flag.String("topic", "cp7sub%i", "pub topic")
var workers = flag.Int("workers", 200, "number of workers")
var cs = flag.String("cid", "client", "client id start with")
var qos = flag.Uint("qos", 0, "pub qos level")
var clearsession = flag.Bool("clear", true, "clear session")

func main() {
	flag.Parse()

	clients := make(map[string]*client.Client)

	for i := 0; i < *workers; i++ {
		id := strconv.Itoa(i)

		cl := client.New()
		cf, err := cl.Connect(&client.Config{
			BrokerURL:    *urlString,
			CleanSession: *clearsession,
			KeepAlive:    "180s",
			ValidateSubs: true,
			ClientID:     *cs + id,
		})
		if err != nil {
			log.Println("conn", err)
		}

		err = cf.Wait(1 * time.Second)
		if err != nil {
			log.Println("conn wait", err)
		}

		clients[*cs+id] = cl
	}

	c := cron.New()
	c.AddFunc("*/10 * * * * ?", func() {
		go func() {
			for i := 0; i < *workers; i++ {
				id := strconv.Itoa(i)
				for _, v := range clients {
					v.Publish(strings.Replace(*topic, "%i", id, 1), []byte("test"), uint8(*qos), false)
					time.Sleep(5 * time.Millisecond)
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()
	})
	c.Start()

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
