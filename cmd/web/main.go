package main

import (
	"bwa/pkg/config"
	"bwa/pkg/handlers"
	"bwa/pkg/render"
	"fmt"
	"log"
	"net/http"
)

const PORT = ":8080"

func main() {

	var app config.AppConfig
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
