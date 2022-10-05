package render

import (
	"bytes"
	"experimenting_with_sessions/pkg/config"
	"experimenting_with_sessions/pkg/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buff := new(bytes.Buffer)

	td = AddDefaultData(td) // Add default data

	_ = t.Execute(buff, td)

	_, err := buff.WriteTo((w))
	if err != nil {
		fmt.Println("Error writing template to browser: ", err)
	}

}

// Creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	// Create a Go map
	// Create a template cache variable that holds all of my templates and creates it at the start of the application
	myCache := map[string]*template.Template{}

	// Tell Go to search for all the files at a certain location
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// Create a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Does this templates match any layouts?
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil

}
