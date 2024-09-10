package render

import (
	"net/http"
	"testing"

	"github.com/Poojasadgir/room-reservation/internal/models"
)

func TestDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")

	result := DefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("Flash value of '123' not found in session")
	}
}

func TestTemplate(t *testing.T) {
	templatePath = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	var ww dummyWriter
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser")
	}
	err = Template(&ww, r, "non-existant.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Tried to render non-existant template")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/test-url", nil)

	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewRenderer(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	templatePath = "./../../templates"
	_, err := CreateTemplateCache()

	if err != nil {
		t.Error(err)
	}
}
