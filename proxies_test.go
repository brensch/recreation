package recreation

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func TestProxyRequest(t *testing.T) {

// 	req := httptest.NewRequest(http.MethodGet, "https://www.sendo.gov/api/sendo", nil)

// 	newHost := "https://www.recreation.gov"
// 	req2 := ProxyRequest(req, newHost)

// 	if req2.URL.Host != newHost {
// 		t.Errorf("got wrong host: %s - expecting %s", req.URL.Host, newHost)
// 	}

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	_ = res

// }

func TestProxyRequest(t *testing.T) {

	// mux := http.NewServeMux()
	// mux.HandleFunc("/HandleProxyRequest", HandleProxyRequest)
	// server := httptest.NewServer(mux)
	// defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "http://asdfasdfasdfffdd.com/HandleProxyRequest/api/camps/availability/campground/232449/month?start_date=2022-06-01T00%3A00%3A00.000Z", nil)
	w := httptest.NewRecorder()

	HandleProxyRequest(w, req)

	res := w.Result()

	resBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(resBytes))
	// res, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	t.Error(err)
	// }

	_ = res

}
