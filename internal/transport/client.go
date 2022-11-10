package transport

import (
	"net/http"
	"time"
)

func InitHTTPClient() *http.Client {

	return &http.Client{
		Transport: &loggingRoundTripper{Proxied: http.DefaultTransport},
		Timeout:   1000 * time.Second,
	}
}
