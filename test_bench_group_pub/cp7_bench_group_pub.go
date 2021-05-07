package main

import (
	"coolpy7_benchmark/src/client"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var urlString = flag.String("url", "tcp://username:password@127.0.0.1:1883", "broker url")
var topic = flag.String("topic", "$share/group1", "pub topic")
var workers = flag.Int("workers", 300, "number of workers")
var interval = flag.Int("i", 10, "interval of connecting to the broker(ms)")
var interval_of_msg = flag.Int("I", 1000, "interval of publishing message(ms)")
var size = flag.Int("s", 4096, "payload size")
var cs = flag.String("cid", "client", "client id start with")
var qos = flag.Uint("qos", 0, "pub qos level")
var clearsession = flag.Bool("clear", true, "clear session")
var pingtime = flag.String("keepalive", "300s", "keepalive")

func main() {
	flag.Parse()

	clients := make(map[string]*client.Client)

	for i := 0; i < *workers; i++ {
		id := strconv.Itoa(i)

		cl := client.New()
		cf, err := cl.Connect(&client.Config{
			BrokerURL:    *urlString,
			CleanSession: *clearsession,
			KeepAlive:    *pingtime,
			ValidateSubs: true,
			ClientID:     *cs + id,
		})
		if err != nil {
			log.Println("conn", err)
		}

		err = cf.Wait(2 * time.Second)
		if err != nil {
			log.Println("conn wait", err)
		}

		clients[*cs+id] = cl
		time.Sleep(time.Duration(*interval) * time.Millisecond)
	}

	bigData := make([]byte, *size, *size)

	for i := 0; i < *workers; i++ {
		go func() {
			for _, v := range clients {
				_, _ = v.Publish(*topic, bigData, uint8(*qos), false)
			}
		}()
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
