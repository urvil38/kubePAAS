package authservice

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/urvil38/kubepaas/types"
)

const (
	userserviceAPI = "https://kubepaas.appspot.com/v1/%s"
)

func newHTTPClient(timeout *time.Duration) *types.HttpClient {
	if timeout != nil {
		httpTransport := &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 20 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 7 * time.Second,
			MaxIdleConns:        5,
			IdleConnTimeout:     1 * time.Second,
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
