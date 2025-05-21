package handlers

import (
	"bwa/pkg/config"
	"bwa/pkg/render"
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
	render.RenderTemplates(w, "home.page.html")
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "about.page.html")
}
