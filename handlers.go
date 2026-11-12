package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")

	ctx := r.Context()
	logger := app.LoggerFromContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	logger.InfoContext(ctx, "returning items")

	json.NewEncoder(w).Encode(things)

}

func (app *application) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")

	ctx := r.Context()
	logger := app.LoggerFromContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL
	id := r.URL.Path[1:]

	for _, thing := range things {
		if thing.ID == id {

			// Item found
			logger.InfoContext(ctx, "returning item")

			json.NewEncoder(w).Encode(thing)

			return
		}
	}

	// Item not found
	logger.ErrorContext(ctx, "unable to find item", "error", errors.New("no item found with given ID"))

	w.WriteHeader(http.StatusNotFound)
}

func (app *application) PostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST")

	ctx := r.Context()
	logger := app.LoggerFromContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get body
	var payload thing
	err = json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		logger.ErrorContext(ctx, "unable to read JSON", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Test") != "true" { // If testing, don't change things
		things = append(things, payload)
	}

	// Item was posted
	w.WriteHeader(http.StatusCreated)

	logger.InfoContext(ctx, "item posted")

}

func (app *application) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE")

	ctx := r.Context()
	logger := app.LoggerFromContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL
	id := r.URL.Path[1:]

	for i, thing := range things {
		if thing.ID == id {

			if r.Header.Get("Test") != "true" { // If testing, don't change things
				things = append(things[:i], things[i+1:]...)
			}

			// Item deleted
			logger.InfoContext(ctx, "item deleted")

			return
		}
	}

	// Item not found
	logger.ErrorContext(ctx, "unable to find item", "error", errors.New("no item found with given ID"))

	w.WriteHeader(http.StatusNotFound)

	fmt.Fprintln(w, "No item found with ID: "+id)
}

func (app *application) PutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT")

	ctx := r.Context()
	logger := app.LoggerFromContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL
	id := r.URL.Path[1:]

	// Get body
	var payload thing
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		logger.ErrorContext(ctx, "unable to read JSON", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, thing := range things {
		if thing.ID == id {

			if r.Header.Get("Test") != "true" { // If testing, don't change things
				things = append(append(things[:i], payload), things[i+1:]...)
			}

			// Item replaced
			logger.InfoContext(ctx, "item replaced")

			return
		}
	}

	// Item not found, create new item
	if r.Header.Get("Test") != "true" { // If testing, don't change things
		things = append(things, payload)
	}

	// Item was created
	w.WriteHeader(http.StatusCreated)
	logger.InfoContext(ctx, "item created")
}

func (app *application) PatchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PATCH")

	ctx := r.Context()
	logger := app.LoggerFromContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL
	id := r.URL.Path[1:]

	// Get body
	var payload [][]string
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		logger.ErrorContext(ctx, "unable to read JSON", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, thing := range things {
		if thing.ID == id {

			// Found item

			var missedChanges [][]string

			for _, change := range payload { // Loop through changes

				if r.Header.Get("Test") != "true" { // If testing, don't change things
					switch change[0] {
					case "Name":
						things[i].Name = change[1]
					case "ID":
						things[i].ID = change[1]
					default:
						// Change not found
						missedChanges = append(missedChanges, change)
					}
				}
			}

			if len(missedChanges) > 0 {
				logger.ErrorContext(ctx, "some keywords were not found", "error", errors.New("some keywords were not found"))

				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(missedChanges)
				return
			}
			// Item patched
			logger.InfoContext(ctx, "item updated")
			fmt.Fprintln(w)

			return
		}
	}

	// Item not found
	logger.ErrorContext(ctx, "unable to find item", "error", errors.New("no item found with given ID"))

	w.WriteHeader(http.StatusNotFound)

	fmt.Fprintln(w, "No item found with ID: "+id)
}

func (app *application) OptionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		return
	}
	fmt.Println("OPTIONS")

	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, origin, accept, X-API-Key")
	w.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight for 1 day

	// Respond with 200 OK and no body
	w.WriteHeader(http.StatusOK)
}
