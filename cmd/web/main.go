package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/p3rfect05/go_proj/internal/config"
	"github.com/p3rfect05/go_proj/internal/handlers"
	"github.com/p3rfect05/go_proj/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {

	//change to true when in production
	appConfig.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}
	fmt.Printf("Server runs on %s port\n", portNumber)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
