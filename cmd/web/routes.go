package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/p3rfect05/go_proj/internal/config"
	"github.com/p3rfect05/go_proj/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	//mux.Use(Auth)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)
	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Get("/user/logout", handlers.Repo.LogOut)

	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.JSONAvailability)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
	})

	return mux
}
