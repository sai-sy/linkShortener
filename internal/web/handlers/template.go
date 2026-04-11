package handlers

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/web/app/templates"
)

type Page struct {
	Title         string
	Body          []byte
	Authenticated bool
	ProfileURL    string
	Routemaps     []db.Routemap
}

type contextKey string

const TemplatesDirKey contextKey = "templatesDir"

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := fs.ReadFile(templates.FS(), filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFS(templates.FS(), "base.html", tmpl+".html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "base", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
