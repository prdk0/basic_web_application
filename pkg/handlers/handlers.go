package handlers

import (
	"bookings/pkg/config"
	"bookings/pkg/models"
	"bookings/pkg/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remoteip", remoteIP)
	render.RenderTemplates(w, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	// stringMap["test"] = "Hello to About Page!"
	remoteIP := m.App.Session.GetString(r.Context(), "remoteip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplates(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
