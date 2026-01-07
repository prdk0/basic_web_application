package handlers

import (
	"bookings/internals/config"
	"bookings/internals/driver"
	"bookings/internals/forms"
	"bookings/internals/helpers"
	"bookings/internals/models"
	"bookings/internals/render"
	"bookings/internals/repository"
	"bookings/internals/repository/dbrepo"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

type templateData = models.TemplateData

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NeTestwRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &templateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &templateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &templateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &templateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &templateData{})
}

// Search Availability

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &templateData{})
}

func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllrooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "database error")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]any)
	data["rooms"] = rooms

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)
	render.Template(w, r, "choose-room.page.tmpl", &templateData{
		Data: data,
	})
}

type JsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		resp := JsonResponse{
			Ok:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "	")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	availableRoom, err := m.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomId)
	if err != nil {
		resp := JsonResponse{
			Ok:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "	")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := JsonResponse{
		Ok:        availableRoom,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomId),
	}

	msg, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(msg)
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	roomId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find the id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, r.URL.Query().Get("s"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find the start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, r.URL.Query().Get("e"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find the end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var reservation models.Reservation
	room, err := m.DB.GetRoomById(roomId)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find the room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.Room.RoomName = room.RoomName
	reservation.RoomID = roomId
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Reservation

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restrictions, err := m.DB.GetRestrictions()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find the restrictions")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomById(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find the room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]any)
	data["reservation"] = res
	data["restrictions"] = restrictions
	render.Template(w, r, "make-reservation.page.tmpl", &templateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse the form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsValidEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation
		if m.App.Env.Test {
			http.Error(w, "Invalid input format", http.StatusSeeOther)
		}
		sd := reservation.StartDate.Format("2006-01-02")
		ed := reservation.EndDate.Format("2006-01-02")

		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed
		render.Template(w, r, "make-reservation.page.tmpl", &templateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	newReservationId, err := m.DB.InsertReservation(reservation)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation in to database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction_selected_value := r.FormValue("restrictions")

	selected_restriction_value, err := strconv.Atoi(restriction_selected_value)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error in selected value")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		RestrictionID: selected_restriction_value,
		ReservationID: sql.NullInt32{Int32: 0, Valid: false},
	}

	if selected_restriction_value == 1 {
		restriction.ReservationID = sql.NullInt32{Int32: int32(newReservationId), Valid: true}
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restrictions")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	htmlMessage := fmt.Sprintf(`
		<string>Reservation Confirmation</strong><pre>
		Dear %s:, <br>
		This is to confirm your reservation from %s to %s.
		`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "pradeek.k@gmail.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "email.html",
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	urlSlitter := strings.Split(r.RequestURI, "/")
	roomId, err := strconv.Atoi(urlSlitter[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get the roomId")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get the session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomId

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("cannot get item from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]any)
	data["reservation"] = reservation
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	render.Template(w, r, "reservation-summary.page.tmpl", &templateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// Login
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user := models.User{}

	user.Email = email
	user.Password = password

	m.App.Session.Put(r.Context(), "loginDetails", user)

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsValidEmail("email")
	if !form.Valid() {
		loginDetails := m.App.Session.Get(r.Context(), "loginDetails").(models.User)
		data := make(map[string]any)
		data["loginDetails"] = loginDetails
		err := render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		return
	}
	user_id, _, err := m.DB.Authenticate(user.Email, user.Password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "wrong username or password")
		loginDetails := m.App.Session.Get(r.Context(), "loginDetails").(models.User)
		data := make(map[string]any)
		data["loginDetails"] = loginDetails
		err := render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		return
	}
	user.ID = user_id
	m.App.Session.Put(r.Context(), "loginDetails", user)
	m.App.Session.Put(r.Context(), "user_id", user.ID)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (m *Repository) AdminsListAllReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]any)
	data["reservations"] = reservations

	err = render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (m *Repository) AdminsListNewReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]any)
	data["reservations"] = reservations

	err = render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	explode := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(explode[4])
	if err != nil {
		helpers.ServerError(w, err)
	}
	src := explode[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src
	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
	}
	data := make(map[string]any)
	data["reservation"] = res
	err = render.Template(w, r, "admin-show-reservations.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	explode := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(explode[4])
	if err != nil {
		helpers.ServerError(w, err)
	}
	src := explode[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src
	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")
	email := r.Form.Get("email")
	phone := r.Form.Get("phone")

	res.FirstName = first_name
	res.LastName = last_name
	res.Email = email
	res.Phone = phone

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.MinLength("last_name", 3)
	form.IsValidEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = res
		if m.App.Env.Test {
			http.Error(w, "Invalid input format", http.StatusSeeOther)
		}
		render.Template(w, r, "admin-show-reservations.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	err = m.DB.UpdateReservationById(res)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "flash", "successfully updated")
	redirectUrl := fmt.Sprintf("/admin/ls-reservation-%s", src)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)

}

func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "src")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	err = m.DB.UpdateProcessedForReservation(id, 1)

	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation successfully processed")
	http.Redirect(w, r, fmt.Sprintf("/admin/ls-reservation-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "src")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	err = m.DB.DeleteReservation(id)

	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "flash", "Deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/ls-reservation-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {

	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]any)
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
		}

		for _, y := range restrictions {
			if y.ReservationID.Valid {
				// reservation
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = int(y.ReservationID.Int32)
				}
			} else {
				// block
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	err = render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (m *Repository) PageNotFound(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "404.page.tmpl", &models.TemplateData{})
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
