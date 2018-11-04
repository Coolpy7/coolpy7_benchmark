package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"os"
	"os/signal"
	"packet"
	"syscall"
	"transport"
)

var urlString = flag.String("url", "tcp://192.168.200.238:1883", "broker url")
var workers = flag.Int("workers", 100, "number of workers")

func main() {
	flag.Parse()

	fmt.Printf("Start benchmark of %s using %d workers.\n", *urlString, *workers)

	clients := make(map[string]transport.Conn)

	for i := 0; i < *workers; i++ {
		id := strconv.Itoa(i)
		if i%1000 == 0 {
			log.Println(i)
		}
		name, conn := consumer(id)
		clients[name] = conn
	}

	fmt.Println("finish", len(clients))

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

func connection(id string) transport.Conn {
	conn, err := transport.Dial(*urlString)
	if err != nil {
		panic(err)
	}

	mqttURL, err := url.Parse(*urlString)
	if err != nil {
		panic(err)
	}

	connect := packet.NewConnectPacket()
	connect.ClientID = "connect/" + id

	if mqttURL.User != nil {
		connect.Username = mqttURL.User.Username()
		pw, _ := mqttURL.User.Password()
		connect.Password = pw
	}

	err = conn.Send(connect)
	if err != nil {
		panic(err)
	}

	return conn
}

func consumer(id string) (string, transport.Conn) {
	name := "consumer/" + id
	return name, connection(name)
}
