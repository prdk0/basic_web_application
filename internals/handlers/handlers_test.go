package handlers

import (
	"bookings/internals/models"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	// {"home", "/", "GET", []postData{}, http.StatusOK},
	// {"about", "/about", "GET", []postData{}, http.StatusOK},
	// {"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// {"generals", "/generals-quaters", "GET", []postData{}, http.StatusOK},
	// {"majors", "/majors-suite", "GET", []postData{}, http.StatusOK},
	// {"search availability", "/seach-availability", "GET", []postData{}, http.StatusOK},
	// {"make reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"404 page", "/error_page", "GET", []postData{}, http.StatusNotFound},
	// {"post search availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2025-08-21"},
	// 	{key: "end", value: "2025-08-23"},
	// }, http.StatusOK},
	// {"post search availability json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2025-08-21"},
	// 	{key: "end", value: "2025-08-23"},
	// }, http.StatusOK},
	// {"make reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Doe"},
	// 	{key: "email", value: "me@here.com"},
	// 	{key: "phone", value: "2342344"},
	// }, http.StatusOK},
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
			values := url.Values{}
			for _, v := range e.params {
				values.Add(v.key, v.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 4,
		Room: models.Room{
			ID:       4,
			RoomName: "General's Quaters",
		},
	}

	req, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := getCtx(req)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
