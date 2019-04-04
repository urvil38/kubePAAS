package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

func NewHTTPClient(timeout *time.Duration) *http.Client {
	if timeout != nil {
		httpTransport := &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConns:        5,
			IdleConnTimeout:     10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
		}
		return &http.Client{
			Transport: httpTransport,
			Timeout:   *timeout,
		}
	}
	return &http.Client{}
}