package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Poojasadgir/room-reservation/internal/config"
	"github.com/Poojasadgir/room-reservation/internal/models"
	"github.com/justinas/nosurf"
)

var templatePath = "./templates"
var app *config.AppConfig
var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

// NewRenderer creates a new renderer with the given AppConfig.
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate takes a time.Time object and returns a string representation of the date in the format "2006-01-02".
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate formats a given time.Time object to a string using the provided format string.
// The format string should be in the same format as the standard library's time package.
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// DefaultData populates the default data for templates.
// It takes a pointer to a TemplateData struct and a pointer to an http.Request struct.
// It returns a pointer to a TemplateData struct.
func DefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// PopString retrieves and removes a string value from the session.
	// Token returns a CSRF token for the given request.
	// Exists returns true if the given key exists in the session.
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Template executes a given template with the provided TemplateData and writes the output to the http.ResponseWriter.
// If app.UseCache is true, it uses the template cache from app.TemplateCache, otherwise it creates a new cache using CreateTemplateCache().
// It returns an error if the template cannot be found in the cache or if there is an error writing the template to the response writer.
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// Get the template cache from AppConfig
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Pulls a template out of the template map(tc)
	t, ok := tc[tmpl]
	if !ok {
		log.Println("Could not get template from template cache")
		return errors.New("can't get template from cache")
	}

	// Byte buffer
	buf := new(bytes.Buffer)
	td = DefaultData(td, r)
	_ = t.Execute(buf, td)

	//Writes the template from the byte buffer to the response writer
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}
	return nil
}

// CreateTemplateCache creates a map of pre-parsed templates by parsing all the page templates
// and their corresponding layout templates. It returns a map of template names to their
// corresponding *template.Template pointers. It also returns an error if any error occurs
// while parsing the templates.
func CreateTemplateCache() (map[string]*template.Template, error) {
	layoutPath := fmt.Sprintf("%s/*.layout.tmpl", templatePath)
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", templatePath))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(layoutPath)
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(layoutPath)
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

// Iterate generates a slice of integers from 0 to count-1.
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}
