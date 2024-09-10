package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Poojasadgir/room-reservation/internal/config"
	"github.com/Poojasadgir/room-reservation/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

// A dummy writer to satisfy response headers for testing
type dummyWriter struct{}

func (tw *dummyWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *dummyWriter) WriteHeader(i int) {

}

func (tw *dummyWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func TestMain(m *testing.M) {
	// What is going to be put into the session
	gob.Register(models.Reservation{})

	// Change this to true when in production
	testApp.InProduction = false

	// Logging and error handling
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}
