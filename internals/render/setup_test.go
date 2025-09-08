package render

import (
	"bookings/internals/config"
	"bookings/internals/models"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	// setting app enviroment
	testApp.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// Session settings
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
	testApp.Session = session

	app = &testApp
	os.Exit(m.Run())
}

type myResponseWriter struct{}

func (my *myResponseWriter) Header() http.Header {
	var h http.Header
	return h
}

func (my *myResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (my *myResponseWriter) WriteHeader(statusCode int) {}
