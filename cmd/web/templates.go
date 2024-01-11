package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.vadimpush.net/internal/models"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := make(map[string]*template.Template, 10)

	pages, err := filepath.Glob("./ui/html/pages/*.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tmpl, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")

		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseGlob("./ui/html/partials/*.html")

		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseFiles(page)

		cache[name] = tmpl
	}

	return cache, nil
}
