package recreation

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	RecreationGovURI = "https://www.recreation.gov"

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ErrCloudFlare struct {
	Status   int
	Contents []byte
}

func (e ErrCloudFlare) Error() string {
	return fmt.Sprintf("Cloudflare Error %d: %s", e.Status, e.Contents)
}

// Obfuscator is what will set all the headers required to avoid detection by cloudflare.
// It also rate limits you to ensure you don't go over the cloudflare limit.
// TODO: make the rate limiting dynamic, ie increase limiting if it's no good
type Obfuscator struct {
	client      *http.Client
	ctx         context.Context
	rateLimiter chan struct{}
}

// TODO: do this better. Should use a real library to make user agent headers
func (c *Obfuscator) Do(req *http.Request) (*http.Response, error) {
	c.rateLimiter <- struct{}{}
	req.Header.Set("User-Agent", UserAgent)
	timer := time.NewTimer(1 * time.Second)
	res, err := c.client.Do(req)
	// wait  at least one second between each call
	<-timer.C
	<-c.rateLimiter
	return res, err
}

// set sensible defaults for http client
func initObfuscator(ctx context.Context) *Obfuscator {
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
		rateLimiter: make(chan struct{}),
		ctx:         ctx,
	}

}
