package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/gopheramit/Simple-Form/internal/pkg"
	"github.com/gopheramit/Simple-Form/internal/pkg/forms"
)

type templateData struct {
	//	CSRFToken   string
	//	CurrentYear int
	Form *forms.Form
	//	Flash       string
	Snippet *pkg.EmailInput
	//Snippets        []*models.Snippet
	//IsAuthenticated bool
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
