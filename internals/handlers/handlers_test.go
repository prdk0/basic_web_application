package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"generals", "/generals-quaters", "GET", []postData{}, http.StatusOK},
	{"majors", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search availability", "/seach-availability", "GET", []postData{}, http.StatusOK},
	{"make reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
}

func TestHandler(t *testing.T) {
	routes := getroutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			rs, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if rs.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, rs.StatusCode)
			}
		} else {

		}
	}
}
