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
	mux := http.NewServeMux() // default mux

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

	simpleHandler := MiddleWareChainTest(http.HandlerFunc(SimpleHttpTest), MiddleWaretest, MiddleWareRecoverTest, WriteToConsole) // MiddleWare Chain Example
	mux.Handle("/simple", simpleHandler)
	srv := http.Server{
		Addr:    PORT,
		Handler: mux, // default mux
	}
	fmt.Printf("Sever listening to the port %s\n", PORT)
	srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func SimpleHttpTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from SimpleHttpTest")
}
