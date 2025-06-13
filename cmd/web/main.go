package main

import (
	"bookings/internals/config"
	"bookings/internals/handlers"
	"bookings/internals/render"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const PORT = ":8080"

var app config.AppConfig // make variable app to global so that it availbles to all in main packages
var session *scs.SessionManager

func main() {

	// setting app enviroment
	app.InProduction = false

	// Session settings
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// create template cache from main -> render through config.app, This is doing because it will run only once instead of
	// running multiple times if you call from render package
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	// repository pattern which helps to implement interfaces
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)
	srv := http.Server{
		Addr:    PORT,
		Handler: router(&app),
	}
	fmt.Printf("Sever listening to the port %s\n", PORT)
	srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
