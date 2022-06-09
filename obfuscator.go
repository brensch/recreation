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
)

// Obfuscator is what will set all the headers required to avoid detection by cloudflare.
type Obfuscator interface {
	// Made this differ from the standard http.Client interface in order to
	// ensure people are aware they must sneak.
	DoSneakily(req *http.Request) (*http.Response, error)
}

type ErrCloudFlare struct {
	Status   int
	Contents []byte
}

func (e ErrCloudFlare) Error() string {
	return fmt.Sprintf("Cloudflare Error %d: %s", e.Status, e.Contents)
}

type AgentRandomiser struct {
	client *http.Client
	ctx    context.Context
}

func (c AgentRandomiser) DoSneakily(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", randomUserAgent())
	req = req.WithContext(c.ctx)
	res, err := c.client.Do(req)
	return res, err
}

// set sensible defaults for http client.
// returning by value to try and avoid the heap
func InitAgentRandomiser(ctx context.Context) AgentRandomiser {
	return AgentRandomiser{
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
	}

}
