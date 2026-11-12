package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	standard := alice.New()

	mux.HandleFunc("GET /", app.GetHandler)
	mux.HandleFunc("GET /healthcheck", app.HealthcheckHandler)
	mux.HandleFunc("GET /{id}", app.GetByIDHandler)
	mux.HandleFunc("POST /", app.PostHandler)
	mux.HandleFunc("DELETE /{id}", app.DeleteHandler)
	mux.HandleFunc("PUT /{id}", app.PutHandler)
	mux.HandleFunc("PATCH /{id}", app.PatchHandler)
	mux.HandleFunc("GET /test", app.TestAllHandlers)
	mux.HandleFunc("OPTIONS /", app.OptionsHandler)

	handler := standard.Then(mux)
	return handler
}
