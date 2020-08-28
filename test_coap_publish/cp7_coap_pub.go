package main

import (
	"github.com/jacoblai/go-coap"
	"log"
	"math"
	"math/rand"
)

func main() {
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
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Printf("Response: code:%s, payload:%s", rv.Code, rv.Payload)
	}
}
