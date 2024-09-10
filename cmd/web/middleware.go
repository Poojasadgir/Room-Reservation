package main

import (
	"net/http"

	"github.com/Poojasadgir/room-reservation/internal/helpers"
	"github.com/justinas/nosurf"
)

// NoSurf adds CSRF protection to all POST requests
// It returns a http.Handler that wraps the input http.Handler
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on each request.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth is a middleware that checks if the user is authenticated before allowing access to the next handler.
// If the user is not authenticated, it redirects them to the login page and sets an error message in the session.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
