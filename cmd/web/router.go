package main

import (
	"bwa/pkg/config"
	"bwa/pkg/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func router(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	// mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	// mux.Use(WriteToConsole)
	// mux.Use(MiddleWaretest)
	// mux.Use(MiddleWareRecoverTest)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	return mux
}
