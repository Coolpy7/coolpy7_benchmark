package main

import (
	"client"
	"log"
	"testing"
	"time"
)

func BenchmarkPub(b *testing.B) {
	cl := client.New()
	cf, err := cl.Connect(&client.Config{
		BrokerURL:    *urlString,
		CleanSession: *clearsession,
		KeepAlive:    *pingtime,
		ValidateSubs: true,
		ClientID:     "bench",
	})
	if err != nil {
		log.Println("conn", err)
	}

	err = cf.Wait(time.Second)
	if err != nil {
		log.Println("conn wait", err)
	}
	bigData := []byte(`{"hello":"world"}`)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cb, err := cl.Publish("/bench/pub", bigData, uint8(*qos), false)
			if err != nil {
				b.Fail()
			}
			err = cb.Wait(50 * time.Millisecond)
			if err != nil {
				b.Fail()
			}
		}
	})
}
