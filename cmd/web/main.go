package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/p3rfect05/go_proj/internal/config"
	"github.com/p3rfect05/go_proj/internal/driver"
	"github.com/p3rfect05/go_proj/internal/handlers"
	"github.com/p3rfect05/go_proj/internal/helpers"
	"github.com/p3rfect05/go_proj/internal/models"
	"github.com/p3rfect05/go_proj/internal/render"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(appConfig.MailChannel)
	listenToMail()
	fmt.Println("Started mail listener...")

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}
	fmt.Printf("Server runs on %s port\n", portNumber)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what will be in the session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	mailChan := make(chan models.MailData)
	appConfig.MailChannel = mailChan

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

	//connect to database
	log.Println("Connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 database=go_proj user=postgres password=1234")
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	log.Println("Connected to db!")
	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache", err)
		return nil, err
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&appConfig)
	helpers.NewHelpers(&appConfig)

	return db, nil
}
