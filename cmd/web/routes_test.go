package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/p3rfect05/go_proj/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	// do nothing
	default:

		t.Errorf("Type is not *chi.Mux, type is %T", v)
	}
}
