package main

import (
	"net/http"
	"os"
)

// middleware.go
func (app *application) CheckAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get environment secret from a .env file or something equivalent
		key := r.Header.Get("X-API-Key")
		// if key doesn't match env variable then return unauthorized
		if key != os.Getenv("APP_API_KEY") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// if nothing is caught by the if statement then run the request through
		next.ServeHTTP(w, r)
	})
}
