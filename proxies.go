package recreation

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type GetAvailabilityReq struct {
	CampgroundID string    `json:"campground_id,omitempty"`
	TargetTime   time.Time `json:"target_time,omitempty"`
}

type ProxyClient struct {
	client http.Client
}

func (c *ProxyClient) Do(req *http.Request) (*http.Response, error) {

	return c.client.Do(req)
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// since this is going to be run as a cloud function, globals seem to be the only way to pass things in
// since the signature is fixed
func HandleProxyRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("proxying")

	proxy.ServeHTTP(w, r)
	fmt.Println("proxied", r.URL)

}

// func Director(r *http.Request) {
// 	fmt.Println("directed")
// 	r.URL.Host = "https://www.recreation.gov"
// 	r.Header.Set("User-Agent", api.RandomUserAgent())
// }

// func ProxyRequest(r *http.Request, newHost string) *http.Request {

// 	// i don't see any reason to change the context from what it was originally, so keeping it the same
// 	req := r.Clone(r.Context())
// 	req.URL.Host = newHost

// 	return req
// }

// // HandleGetAvailability is intended to be run as a cloud function in GCP.
// // I am spreading these across all the regions of GCP so that they will all have different IP addresses.
// // Cloud functions are extremely cheap so the theory is that this will actually lead to less resource usage.
// func HandleGetAvailability(w http.ResponseWriter, r *http.Request) {
// 	uuid := uuid.New().String()
// 	log = log.With(zap.String("request_id", uuid))

// 	var req GetAvailabilityReq
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		log.Error("failed to decode request", zap.Error(err))
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	avail, err := api.GetAvailability(context.Background(), log, "https://us-west1-campr-app.cloudfunctions.net/HandleProxyRequest", req.CampgroundID, req.TargetTime)
// 	if err != nil {
// 		log.Error("failed to get availability", zap.Error(err))
// 		w.WriteHeader(http.StatusTeapot)
// 		return
// 	}

// 	w.Header().Add("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(avail)
// 	if err != nil {
// 		log.Error("failed to encode avail into response", zap.Error(err))
// 		w.WriteHeader(http.StatusTeapot)
// 		return
// 	}
// }
