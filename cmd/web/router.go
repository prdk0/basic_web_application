package main

import (
	"bookings/internals/config"
	"bookings/internals/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func router(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/generals-quaters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	//Search Availability
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	// Reservation
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	// login
	mux.Route("/user", func(mux chi.Router) {
		// mux.Use(redirectIfAuthenticated)
		mux.Get("/login", handlers.Repo.ShowLogin)
		mux.Post("/login", handlers.Repo.PostShowLogin)
	})
	mux.Get("/user/logout", handlers.Repo.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/ls-reservation-new", handlers.Repo.AdminListNewReservations)
		mux.Get("/ls-reservation-all", handlers.Repo.AdminsListAllReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
	})

	// 404
	mux.NotFound(handlers.Repo.PageNotFound)

	staticFileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))
	return mux
}
