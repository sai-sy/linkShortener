package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

type contextKey string

const TemplatesDirKey contextKey = "templatesDir"

//go:embed templates/*.html
var templatesFS embed.FS

func loadPage(title string) (*Page, error) {
	filename := "templates/" + title + ".html"
	fmt.Println("loadPage filename:", filename)
	body, err := templatesFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFS(templatesFS, "templates/"+tmpl+".html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	title := "index"
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("cannot find")}
	}
	fmt.Println("rendering page", title)
	renderTemplate(w, title, p)

}
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	title := "login"
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("cannot find")}
	}
	fmt.Println("rendering page", title)
	renderTemplate(w, title, p)

}
