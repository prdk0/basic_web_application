package handlers

import (
	"bookings/internals/models"
	"context"
	"log"
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
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"generals", "/generals-quaters", "GET", http.StatusOK},
	{"majors", "/majors-suite", "GET", http.StatusOK},
	{"search availability", "/seach-availability", "GET", http.StatusOK},
	{"404 page", "/error_page", "GET", http.StatusNotFound},
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
	// reservation test with session
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

	// reservation test with out session
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Error(err)
	}
	ctx = getCtx(req)

	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	reservation.RoomID = 100

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// reservation test with invalid room id
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Error(err)
	}
	ctx = getCtx(req)

	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
