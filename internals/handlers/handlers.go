package handlers

import (
	"bookings/internals/config"
	"bookings/internals/models"
	"bookings/internals/render"
	"encoding/json"
	"fmt"
	"log"
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
	render.RenderTemplates(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	// stringMap["test"] = "Hello to About Page!"
	remoteIP := m.App.Session.GetString(r.Context(), "remoteip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplates(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Search Availability

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	dateData := fmt.Sprintf("start date is %s and end date is %s", start, end)
	w.Write([]byte(dateData))
}

type JsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := JsonResponse{
		Ok:      true,
		Message: "Hello from Json Response",
	}
	msg, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(msg)
}

// Reservation

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}
