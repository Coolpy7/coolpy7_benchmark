package main

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	pres := []byte(`{"ok":true}`)
	u := `http://localhost:8081/pub?clientid=jacoblai&topic=testtopic/a&qos=0&retain=false`
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(u)
	req.Header.SetMethodBytes([]byte("POST"))
	req.SetBody([]byte(`{"hello":"world"}`))
	//1 cp7启动参数 -as=1 -jsk=coolpy7
	//2 https://jwt.io/ 使用coolpy7作为your-256-bit-secret
	//3 把生成的jwt放到下一行代码的Basic后，注意Basic后跟一个空格
	req.Header.Add("Authorization", "Basic eyJhbGciOiJIUzI1NiJ9.e30.k1PZfshORXyxbck0bv95juNEBvbPNd2L47bqVsy4ix8")

	fasthttp.Do(req, resp)

	bodyBytes := resp.Body()
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	if bytes.Compare(bodyBytes, pres) != 0 {
		log.Fatal(string(bodyBytes))
	}
	log.Println("test ok")

}
