package main

import (
	"github.com/p3rfect05/go_proj/config"
	"github.com/p3rfect05/go_proj/pkg/handlers"
	"github.com/p3rfect05/go_proj/pkg/render"
	"log"
	"net/http"
)

func main() {
	var appConfig config.AppConfig
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = false
	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)
	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)
	_ = http.ListenAndServe(":8080", nil)
}
