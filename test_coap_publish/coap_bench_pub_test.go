package main

import (
	"bytes"
	"github.com/jacoblai/go-coap"
	"math"
	"math/rand"
	"testing"
)

func BenchmarkPub(b *testing.B) {
	min := 0
	max := math.MaxUint16
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: uint16(rand.Intn(max-min) + min),
		Payload:   []byte(`{"data":"hello mqtt from coap"}`),
	}

	req.SetPathString("/pub")
	req.SetOption(coap.LocationQuery, "Basic amFjb2I6cGFzcw==")
	req.SetOption(coap.URIQuery, "clientid=jacoblai&topic=testtopic/+/a/#&qos=0&retain=false")

	c, err := coap.Dial("udp", "localhost:1883")
	if err != nil {
		b.Error(err)
	}
	pres := []byte(`{"ok":true}`)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rv, err := c.Send(req)
			if err != nil {
				b.Error(err)
			}
			if bytes.Compare(rv.Payload, pres) != 0 {
				b.Error("payload")
			}
		}
	})
}
