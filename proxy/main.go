package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// url, err := url.Parse("https://www.recreation.gov")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// proxy := httputil.ReverseProxy{
	// 	Transport: roundTripper(rt),
	// 	Director: func(req *http.Request) {
	// 		req.URL.Scheme = "https"
	// 		req.URL.Host = "www.recreation.gov"
	// 		req.Header.Set("User-Agent", api.RandomUserAgent())
	// 		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/HandleProxyRequest")
	// 	},
	// }

	// proxy := httputil.ReverseProxy{
	// 	Transport: roundTripper(rt),
	// 	Director: func(req *http.Request) {
	// 		target := "www.recreation.gov"
	// 		req.URL.Scheme = "https"
	// 		req.URL.Host = target
	// 		req.Host = target
	// 		req.Header["X-Forwarded-For"] = nil
	// 		req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")

	// 		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/HandleProxyRequest")
	// 	},
	// }

	// proxy.Director = recreation.Director

	// http.HandleFunc("/HandleProxyRequest/", proxy.ServeHTTP)

	http.ListenAndServe(":8080", nil)

}

func rt(req *http.Request) (*http.Response, error) {
	log.Printf("request received. url=%s", req.URL)
	for header, values := range req.Header {
		fmt.Println(header, values)
	}
	// req.Header.Set("Host", "nghttp2.org") // <--- I set it here as well
	defer log.Printf("request complete. url=%s", req.URL)

	return http.DefaultTransport.RoundTrip(req)
}

// roundTripper makes func signature a http.RoundTripper
type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
