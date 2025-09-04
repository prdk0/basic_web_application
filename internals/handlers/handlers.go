package handlers

import (
	"bookings/internals/config"
	"bookings/internals/forms"
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

type templateData = models.TemplateData

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
	render.RenderTemplates(w, r, "home.page.tmpl", &templateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	// stringMap["test"] = "Hello to About Page!"
	remoteIP := m.App.Session.GetString(r.Context(), "remoteip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplates(w, r, "about.page.tmpl", &templateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "contact.page.tmpl", &templateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "generals.page.tmpl", &templateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "majors.page.tmpl", &templateData{})
}

// Search Availability

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "search-availability.page.tmpl", &templateData{})
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
	var emptyReservation models.Reservation
	data := make(map[string]any)
	data["reservation"] = emptyReservation
	render.RenderTemplates(w, r, "make-reservation.page.tmpl", &templateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsValidEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation
		render.RenderTemplates(w, r, "make-reservation.page.tmpl", &templateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("cannot get item from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]any)
	data["reservation"] = reservation
	render.RenderTemplates(w, r, "reservation-summary.page.tmpl", &templateData{
		Data: data,
	})
}

func (m *Repository) PageNotFound(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplates(w, r, "404.page.tmpl", &models.TemplateData{})
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
