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
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)
	fmt.Printf("Sever listening to the port %s\n", PORT)
	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
	}
}
