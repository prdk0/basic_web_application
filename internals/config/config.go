package config

import (
	"bookings/internals/models"
	"fmt"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Env           EnvRn
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}

type EnvRn struct {
	Dev  bool
	Test bool
}

func (e *EnvRn) SetEviroment(s string) {
	switch s {
	case "dev":
		e.Dev = true
		e.Test = false
	case "test":
		e.Test = true
		e.Dev = false
	default:
		fmt.Println("Wrong entry")
		return
	}
}
