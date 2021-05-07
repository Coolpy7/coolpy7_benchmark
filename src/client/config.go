package client

import (
	"coolpy7_benchmark/src/packet"
	"coolpy7_benchmark/src/transport"
)

// A Config holds information about establishing a connection to a broker.
type Config struct {
	Dialer       *transport.Dialer
	BrokerURL    string
	ClientID     string
	CleanSession bool
	KeepAlive    string
	WillMessage  *packet.Message
	ValidateSubs bool
}

// NewConfig creates a new Config using the specified URL.
func NewConfig(url string) *Config {
	return &Config{
		BrokerURL:    url,
		CleanSession: true,
		KeepAlive:    "30s",
		ValidateSubs: true,
	}
}

// NewConfigWithClientID creates a new Config using the specified URL and client ID.
func NewConfigWithClientID(url, id string) *Config {
	config := NewConfig(url)
	config.ClientID = id
	return config
}
