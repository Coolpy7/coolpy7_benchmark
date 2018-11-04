package transport

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"log"
	"strings"
	"errors"
)

// The Dialer handles connecting to a server and creating a connection.
type Dialer struct {
	TLSConfig     *tls.Config
	RequestHeader http.Header

	DefaultTCPPort string
	DefaultTLSPort string
	DefaultWSPort  string
	DefaultWSSPort string

	webSocketDialer *websocket.Dialer

	Ips   map[int]net.IP
	IpIdx int
}

// NewDialer returns a new Dialer.
func NewDialer() *Dialer {
	return &Dialer{
		DefaultTCPPort: "1883",
		DefaultTLSPort: "8883",
		DefaultWSPort:  "80",
		DefaultWSSPort: "443",
		webSocketDialer: &websocket.Dialer{
			Proxy:        http.ProxyFromEnvironment,
			Subprotocols: []string{"mqtt"},
		},
		Ips:   make(map[int]net.IP),
		IpIdx: 0,
	}
}

var sharedDialer *Dialer

func init() {
	sharedDialer = NewDialer()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("init ", err)
	}
	idx := 0
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), "192.168.") {
				sharedDialer.Ips[idx] = ipnet.IP
				idx++
				log.Println(idx, ipnet.IP.String())
			}
		}
	}
}

// Dial is a shorthand function.
func Dial(urlString string) (Conn, error) {
	return sharedDialer.Dial(urlString)
}

// Dial initiates a connection based in information extracted from an URL.
func (d *Dialer) Dial(urlString string) (Conn, error) {
	urlParts, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(urlParts.Host)
	if err != nil {
		host = urlParts.Host
		port = ""
	}

	switch urlParts.Scheme {
	case "tcp", "mqtt":
		if port == "" {
			port = d.DefaultTCPPort
		}
	RELOAD:
		if len(d.Ips) == 0 {
			return nil, errors.New("no ip cat use")
		}
		localaddr := &net.TCPAddr{IP: d.Ips[d.IpIdx]}
		dl := net.Dialer{LocalAddr: localaddr}
		conn, err := dl.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			if err.Error() == "cannot assign requested address" {
				d.IpIdx++
				log.Println(d.IpIdx, "change local address")
				goto RELOAD
			}
			return nil, err
		}

		return NewNetConn(conn), nil
	case "tls", "mqtts":
		if port == "" {
			port = d.DefaultTLSPort
		}

		conn, err := tls.Dial("tcp", net.JoinHostPort(host, port), d.TLSConfig)
		if err != nil {
			return nil, err
		}

		return NewNetConn(conn), nil
	case "ws":
		if port == "" {
			port = d.DefaultWSPort
		}

		wsURL := fmt.Sprintf("ws://%s:%s%s", host, port, urlParts.Path)

		conn, _, err := d.webSocketDialer.Dial(wsURL, d.RequestHeader)
		if err != nil {
			return nil, err
		}

		return NewWebSocketConn(conn), nil
	case "wss":
		if port == "" {
			port = d.DefaultWSSPort
		}

		wsURL := fmt.Sprintf("wss://%s:%s%s", host, port, urlParts.Path)

		d.webSocketDialer.TLSClientConfig = d.TLSConfig
		conn, _, err := d.webSocketDialer.Dial(wsURL, d.RequestHeader)
		if err != nil {
			return nil, err
		}

		return NewWebSocketConn(conn), nil
	}

	return nil, ErrUnsupportedProtocol
}
