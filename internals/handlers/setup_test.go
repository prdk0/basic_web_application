package handlers

import (
	"bookings/internals/config"
	"bookings/internals/driver"
	"bookings/internals/models"
	"bookings/internals/render"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var pathToTemplates = "./../../templates"

var functions = template.FuncMap{}

var app config.AppConfig // make variable app to global so that it availbles to all in main packages
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

func getroutes() http.Handler {
	gob.Register(models.Reservation{})

	// setting app enviroment
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Session settings
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	db, err := driver.ConnectSQL("")
	if err != nil {
		app.ErrorLog.Println(err)
	}
	// create template cache from main -> render through config.app, This is doing because it will run only once instead of
	// running multiple times if you call from render package
	tc, err := createTestTemplateCache()
	if err != nil {
		log.Fatal("CreateTemplateCache failed")
	}
	app.TemplateCache = tc
	app.UseCache = true

	// repository pattern which helps to implement interfaces
	repo := NewRepo(&app, db)
	NewHandlers(repo)
	render.NewTemplate(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/generals-quaters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	//Search Availability
	mux.Get("/seach-availability", Repo.SearchAvailability)
	mux.Post("/search-availability", Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	// Reservation
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	staticFileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))
	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session to all request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func createTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
