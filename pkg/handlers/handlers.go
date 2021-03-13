package handlers

import (
	"net/http"

	"github.com/sztyui/bookings/pkg/config"
	"github.com/sztyui/bookings/pkg/models"
	"github.com/sztyui/bookings/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Template(w, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.Template(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation for handle Reservations
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "make-reservation.page.html", &models.TemplateData{})
}

// Generals for handle Generals and book
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "generals.page.html", &models.TemplateData{})
}

// Majors for handle Majors and book
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "majors.page.html", &models.TemplateData{})
}

// Availability for handle Availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "search-availability.page.html", &models.TemplateData{})
}

// Contact for handle Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.html", &models.TemplateData{})
}
