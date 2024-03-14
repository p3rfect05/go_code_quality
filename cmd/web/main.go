package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/p3rfect05/go_proj/internal/config"
	"github.com/p3rfect05/go_proj/internal/handlers"
	"github.com/p3rfect05/go_proj/internal/helpers"
	"github.com/p3rfect05/go_proj/internal/models"
	"github.com/p3rfect05/go_proj/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}
	fmt.Printf("Server runs on %s port\n", portNumber)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// what will be in the session
	gob.Register(models.Reservation{})
	//change to true when in production
	appConfig.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache", err)
		return err
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)
	helpers.NewHelpers(&appConfig)

	return nil
}
