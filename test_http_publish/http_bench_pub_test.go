package test_http_publish

import (
	"bytes"
	"encoding/base64"
	"github.com/valyala/fasthttp"
	"testing"
)

func BenchmarkPub(b *testing.B) {
	b.ReportAllocs()
	pres := []byte(`{"ok":true}`)
	u := `http://localhost:8081/pub?clientid=jacoblai&topic=testtopic/a&qos=0&retain=false`
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()

			req.SetRequestURI(u)
			req.Header.SetMethodBytes([]byte("POST"))
			req.SetBody([]byte(`{"hello":"world"}`))
			req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))

			fasthttp.Do(req, resp)

			bodyBytes := resp.Body()
			if bytes.Compare(bodyBytes, pres) != 0 {
				b.Error("payload")
			}
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
		}
	})
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
