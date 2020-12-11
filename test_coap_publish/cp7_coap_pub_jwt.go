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
	//1 cp7启动参数 -as=1 -jsk=coolpy7
	//2 https://jwt.io/ 使用coolpy7作为your-256-bit-secret
	//3 把生成的jwt放到下一行代码的Basic后，注意Basic后跟一个空格
	req.SetOption(coap.LocationQuery, "Basic eyJhbGciOiJIUzI1NiJ9.e30.k1PZfshORXyxbck0bv95juNEBvbPNd2L47bqVsy4ix8")
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
