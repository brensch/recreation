package recreation

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyRequest(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://asdfasdfasdfffdd.com/HandleProxyRequest/api/camps/availability/campground/232449/month?start_date=2022-06-01T00%3A00%3A00.000Z", nil)
	w := httptest.NewRecorder()

	// TODO: mock this so it doesn't go to the real site
	HandleProxyRequest(w, req)

	res := w.Result()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(resBytes))

}
