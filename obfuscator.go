package recreation

import (
	"context"
	"net"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Obfuscator is what will set all the headers required to avoid detection by cloudflare
type Obfuscator struct {
	client *http.Client
	ctx    context.Context
}

// TODO: do this better. Should use a real library to make user agent headers
func (c *Obfuscator) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	return c.client.Do(req)
}

// set sensible defaults for http client
func InitObfuscator(ctx context.Context) *Obfuscator {
	return &Obfuscator{
		client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
		ctx: ctx,
	}

}
