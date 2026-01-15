package routes

import (
    "context"
    "fmt"
    "html/template"
    "net/http"
    "os"
    "path/filepath"

    "github.com/sai-sy/linkShortener/internal/api"
    "github.com/sai-sy/linkShortener/internal/db"
)

type contextKey string

const TemplatesDirKey contextKey = "templatesDir"

var defaultTemplatesDir = filepath.Join("internal", "routes", "templates")

type Page struct {
    Title string
    Body  []byte
}

func templatesDirFromContext(ctx context.Context) string {
    if ctx == nil {
        return defaultTemplatesDir
    }
    if dir, ok := ctx.Value(TemplatesDirKey).(string); ok && dir != "" {
        return dir
    }
    return defaultTemplatesDir
}

func loadPage(dir, title string) (*Page, error) {
    filename := filepath.Join(dir, title+".html")
    fmt.Println("loadPage filename:", filename)
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

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

func NewRouter(ctx context.Context, queries *db.Queries) http.Handler {
    mux := http.NewServeMux()

    templateDir := templatesDirFromContext(ctx)

    api.RegisterRoutes(mux, queries)

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        title := "index"
        p, err := loadPage(templateDir, title)
        if err != nil {
            p = &Page{Title: title, Body: []byte("cannot find")}
        }
        fmt.Println("rendering page", title)
        renderTemplate(w, templateDir, "index", p)
    })

    return mux
}
