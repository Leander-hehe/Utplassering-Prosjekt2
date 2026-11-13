package main

import (
	"api-test/internal/data"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

// !!! Fix needed in test-handlers.go

type Thing struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ThingNoID struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type application struct {
	db *sql.DB
}

var app application

func main() {
	os.Setenv("APP_API_KEY", "test123")

	// Connect to SQLite database (creates file if it doesn't exist)
	db, err := sql.Open("sqlite", "data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	data.Setup(db)

	// Store database pointer in app.db
	app.db = db

	handler := app.routes()

	http.ListenAndServe(":8080", app.CheckAuthorization(handler))

}

func (app *application) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
