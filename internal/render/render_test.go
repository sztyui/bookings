package render

import (
	"net/http"
	"testing"

	"github.com/sztyui/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("failed")
	}
}

func TestTemplate(t *testing.T) {
	tc, err := CreateTemplateCache(templateCacheFolder)
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = Template(&ww, r, "about.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	tc, err := CreateTemplateCache(templateCacheFolder)
	if err != nil {
		t.Error(err)
	}
	if len(tc) == 0 {
		t.Error("Could not find the templates.")
	}
}
