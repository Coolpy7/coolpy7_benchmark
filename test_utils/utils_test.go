package utils

import (
	"crypto/tls"
	"log"
	"net"
	"packet"
	"testing"
	"time"
)

var ip = "192.168.100.2:1883"

//测试空连接攻击，当连接建立后两秒内不进行MQTT身份验证即被Coolpy7主动断开连接
func TestDdos(t *testing.T) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	time.Sleep(3 * time.Second)
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	pkt := packet.NewConnectPacket()
	pkt.CleanSession = true
	pkt.ClientID = "test1"
	ibuf := make([]byte, pkt.Len())
	pkt.Encode(ibuf)
	n, err := conn.Write(ibuf)
	if err != nil {
		t.Error(n, err)
	}
	buf := make([]byte, 128)
	n, err = conn.Read(buf)
	if err == nil {
		t.Error(n, err)
	}
}

//Tls连接测试，需连接到Coolpy7 TLS Poxy代理服务
func TestTls(t *testing.T) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", ip, conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	pkt := packet.NewConnectPacket()
	pkt.CleanSession = true
	pkt.ClientID = "test1"
	ibuf := make([]byte, pkt.Len())
	pkt.Encode(ibuf)
	_, err = conn.Write(ibuf)
	if err != nil {
		t.Error(err)
	}
	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		t.Error(err)
	}
	cb := packet.NewConnackPacket()
	_, err = cb.Decode(buf)
	if err != nil {
		t.Error()
	}
	if cb.ReturnCode != 0 {
		t.Error("result err")
	}
	t.Log(cb.ReturnCode)
}
