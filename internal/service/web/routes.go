package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)


type Page struct {
	Title string
	Body  []byte
}

type contextKey string

const TemplatesDirKey contextKey = "templatesDir"

var defaultTemplatesDir = filepath.Join("internal", "routes", "templates")

func renderTemplate(w http.ResponseWriter, dir, tmpl string, p *Page) {
	t, err := template.ParseFiles(filepath.Join(dir, tmpl+".html"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, p); err != nil {
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}



type Handler struct {

}

func NewHandler() *Handler{
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux){
	mux.HandleFunc("/",h.handleRoot)
}

func (h *Handler) handleRoot(w http.ResponseWriter, r *http.Request){

}
