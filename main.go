package main

import (
	"fmt"
	"net/http"
	"os"
)

type thing struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type application struct{}

var app application

var things = []thing{
	{ID: "1", Name: "Ok"},
}

func main() {
	os.Setenv("APP_API_KEY", "test123")

	handler := app.routes()

	http.ListenAndServe(":8080", app.CheckAuthorization(handler))
}

func (app *application) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w)
}
