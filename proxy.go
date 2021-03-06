package recreation

import (
	"net/http"
	"net/http/httputil"
	"strings"
)

// from https://cloud.google.com/functions/docs/locations

// Tier 1 pricing
// Cloud Functions is available in the following regions with Tier 1 pricing:
// -----
// us-west1 (Oregon) leaf icon Low CO2
// us-central1 (Iowa) leaf icon Low CO2
// us-east1 (South Carolina)
// us-east4 (Northern Virginia)
// europe-west1 (Belgium) leaf icon Low CO2
// europe-west2 (London)
// asia-east1 (Taiwan)
// asia-east2 (Hong Kong)
// asia-northeast1 (Tokyo)
// asia-northeast2 (Osaka)

// Tier 2 pricing
// Cloud Functions is available in the following region with Tier 2 pricing:
// -----
// us-west2 (Los Angeles)
// us-west3 (Salt Lake City)
// us-west4 (Las Vegas)
// northamerica-northeast1 (Montreal) leaf icon Low CO2
// southamerica-east1 (Sao Paulo) leaf icon Low CO2
// europe-west3 (Frankfurt)
// europe-west6 (Zurich) leaf icon Low CO2
// europe-central2 (Warsaw)
// australia-southeast1 (Sydney)
// asia-south1 (Mumbai)
// asia-southeast1 (Singapore)
// asia-southeast2 (Jakarta)
// asia-northeast3 (Seoul)

var (
	proxies = []string{
		// tier 1
		"us-west1",
		"us-central1",
		"us-east1",
		"us-east4",
		"europe-west1",
		"europe-west2",
		"asia-east1",
		"asia-east2",
		"asia-northeast1",
		"asia-northeast2",

		// tier 2
		"us-west2",
		"us-west3",
		"us-west4",
		"northamerica-northeast1",
		"southamerica-east1",
		"europe-west3",
		"europe-west6",
		"europe-central2",
		"australia-southeast1",
		"asia-south1",
		"asia-southeast1",
		"asia-southeast2",
		"asia-northeast3",
	}
)

func MakeProxy() httputil.ReverseProxy {

	return httputil.ReverseProxy{
		// TODO: reenable roundtripper with some logging potentially
		// Transport: roundTripper(rt),
		Director: func(req *http.Request) {
			target := "www.recreation.gov"
			req.URL.Scheme = "https"
			req.URL.Host = target
			req.Host = target
			req.Header["X-Forwarded-For"] = nil
			req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/HandleProxyRequest")
		},
	}
}

func rt(req *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(req)
}

// roundTripper makes func signature a http.RoundTripper
type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
