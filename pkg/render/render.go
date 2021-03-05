package render

import (
	"bytes"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"text/template"

	"github.com/sztyui/bookings/pkg/config"
	"github.com/sztyui/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData for adding the default data to my
// TemplateData model
func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// Template is for rendering template files
func Template(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache("./templates")
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Template not found", tmpl)
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to buffer", err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache(tmpl string) (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(path.Join(tmpl, "*.page.html"))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(path.Join(tmpl, "*.layout.html"))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(path.Join(tmpl, "*.layout.html"))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
