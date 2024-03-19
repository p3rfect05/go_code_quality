package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/p3rfect05/go_proj/internal/config"
	"github.com/p3rfect05/go_proj/internal/driver"
	"github.com/p3rfect05/go_proj/internal/forms"
	"github.com/p3rfect05/go_proj/internal/helpers"
	"github.com/p3rfect05/go_proj/internal/models"
	"github.com/p3rfect05/go_proj/internal/render"
	"github.com/p3rfect05/go_proj/internal/repository"
	"github.com/p3rfect05/go_proj/internal/repository/db_repo"
)

// TemplateData holds data sent from handlers to templates

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  db_repo.NewPostgresRepo(db.SQL, a),
	}
}

func newTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  db_repo.NewTestDBRepo(a),
	}
}

// NewHandlers sets the repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})

}

// `Reservation` renders "make reservation" page and displays it
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.SearchRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room = room
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	m.App.Session.Put(r.Context(), "reservation", res)
	data := make(map[string]interface{})
	data["reservation"] = res

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "make-reservations.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		err_msg := "can't get reservation from session"
		m.App.Session.Put(r.Context(), "error", err_msg)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		helpers.ServerError(w, errors.New(err_msg))
		return
	}
	err := r.ParseForm()
	if err != nil {
		err_msg := "error, can't parse form"
		m.App.Session.Put(r.Context(), "error", err_msg)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		helpers.ServerError(w, errors.New(err_msg))
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone", "start_date", "end_date")
	form.MinLength("first_name", 3)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservations.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		err_msg := "can't insert reservation into database"
		m.App.Session.Put(r.Context(), "error", err_msg)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		helpers.ServerError(w, errors.New(err_msg))
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		ReservationID: newReservationID,
		RoomID:        reservation.RoomID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		err_msg := "can't insert room restriction into database"
		m.App.Session.Put(r.Context(), "error", err_msg)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		helpers.ServerError(w, errors.New(err_msg))
		return
	}

	//send notifications - first to guess
	htmlMessage := fmt.Sprintf(`
		<strong>Reservation confirmation</string><br>
		Dear %s, <br>
		You have booked reservation from %s to %s
		`, reservation.FirstName, reservation.StartDate.Format(time.DateOnly), reservation.EndDate.Format(time.DateOnly))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.ru",
		Subject:  "Reservation confirmation",
		Content:  htmlMessage,
		Template: "drip.html",
	}
	m.App.MailChannel <- msg

	htmlMessage = fmt.Sprintf(`
		<strong>Reservation confirmation</string><br>
		A reservation has been made for %s from %s to %s
		`, reservation.FirstName, reservation.StartDate.Format(time.DateOnly), reservation.EndDate.Format(time.DateOnly))

	msg = models.MailData{
		To:      reservation.Email,
		From:    "me@here.ru",
		Subject: "Reservation confirmation",
		Content: htmlMessage,
	}
	m.App.MailChannel <- msg
	fmt.Print("email sent")
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders "generals" page and displays it
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders "majors" page and displays it
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders "search availability" page and displays it
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handlers POST request from make-reservation page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start_date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end_date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse room_id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// JSONAvailability handles request for availability and sends JSON response
func (m *Repository) JSONAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Can't parse start date",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Can't parse end date",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Can't get room_id",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	available, err := m.DB.SearchAvailabilityByDateByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Can't connect to the database",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders "contact" page and displays it
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//SETUP LATER: helpers.ServerError(w, errors.New("can't get error from session"))
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL params, builds a session variables and redirects user to make-reservation page
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()

	roomID, err := strconv.Atoi(urlParams.Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := urlParams.Get("s")
	ed := urlParams.Get("e")

	var res models.Reservation
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room, err := m.DB.SearchRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room = room

	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		err_msg := "error, can't parse form"
		m.App.Session.Put(r.Context(), "error", err_msg)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		helpers.ServerError(w, errors.New(err_msg))
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.AuthenticateUser(email, password)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogOut the user out
func (m *Repository) LogOut(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}
