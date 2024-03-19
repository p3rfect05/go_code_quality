package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/p3rfect05/go_proj/internal/models"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
	MailChannel   chan models.MailData
}
