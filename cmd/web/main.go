package main

import (
	"bookings/internals/config"
	"bookings/internals/driver"
	"bookings/internals/handlers"
	"bookings/internals/helpers"
	"bookings/internals/models"
	"bookings/internals/render"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const PORT = ":8080"

var app config.AppConfig // make variable app to global so that it available to all in main packages
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	listenForMail()

	srv := http.Server{
		Addr:    PORT,
		Handler: router(&app),
	}
	fmt.Printf("Sever listening to the port %s\n", PORT)
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// setting app enviroment
	app.InProduction = false
	app.Env.SetEviroment("dev")
	//logger

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

	// connect to database
	log.Println("connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=pradeek password=Deepakr_123")
	if err != nil {
		log.Fatal("cannot connect to database")
	}

	// create template cache from main -> render through config.app, This is doing because it will run only once instead of
	// running multiple times if you call from render package
	tc, err := render.CreateTemplateCache()
	if err != nil {
		return nil, errors.New("CreateTemplateCache failed")
	}
	app.TemplateCache = tc
	app.UseCache = false

	// repository pattern which helps to implement interfaces
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
