package userservice

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/urvil38/kubepaas/types"
)

const (
	userserviceEndpoint = "https://kubepaas.appspot.com/v1/%s"
)

func NewHTTPClient(timeout *time.Duration) *types.HttpClient {
	if timeout != nil {
		httpTransport := &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 20 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConns:        5,
			IdleConnTimeout:     20 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
		}
		return &types.HttpClient{
			Client: &http.Client{
				Transport: httpTransport,
				Timeout:   *timeout,
			},
		}
	}
	return &types.HttpClient{
		Client: &http.Client{},
	}
}
