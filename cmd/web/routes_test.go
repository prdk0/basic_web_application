package main

import (
	"bookings/internals/config"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := router(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Errorf("type is not *chi.Mux, type is %T", v)
	}
}
