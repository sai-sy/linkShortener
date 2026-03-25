package handlers

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/sai-sy/linkShortener/internal/web/templates"
)

type Page struct {
	Title string
	Body  []byte
}

type contextKey string

const TemplatesDirKey contextKey = "templatesDir"

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	fmt.Println("loadPage filename:", filename)
	body, err := fs.ReadFile(templates.FS(), filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFS(templates.FS(), tmpl+".html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
