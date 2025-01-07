package main

import (
	"forum/cmd/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

// templateData is the structure that holds data that is passed to the HTML templates.
type templateData struct {
	CurrentYear     int
	Thread          *models.Thread
	Threads         []*models.Thread
	User            *models.User
	Form            any
	Flash           string
	IsAuthenticated bool
}

// newTemplate initializes a templateData struct with the current year and a flash message.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

// newTemplateCache creates a cache of parsed HTML templates.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
