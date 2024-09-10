package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Poojasadgir/room-reservation/internal/config"
)

var app *config.AppConfig

// NewHelpers initializes the app configuration for helpers package
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError logs a client error with a given status and returns an HTTP error response with the same status code.
// It takes in a http.ResponseWriter and an integer status code as parameters.
// It logs the error using the app.InfoLog and returns an HTTP error response with the given status code.
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of ", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError logs the error and sends a 500 Internal Server Error response to the client.
// It takes a http.ResponseWriter and an error as input parameters.
// The error is logged along with a stack trace and a 500 status code is sent to the client.
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// IsAuthenticated checks if the user is authenticated by checking if the "user_id" key exists in the session.
// It takes a pointer to an http.Request as a parameter and returns a boolean value.
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}
