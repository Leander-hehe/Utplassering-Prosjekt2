package main

import (
	"api-test/internal/logging"
	"api-test/internal/rest"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.InfoContext(ctx, "returning items")
	rest.WriteJSON(w, r, things, http.StatusOK)
}

func (app *application) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	// Get ID from URL
	id := r.URL.Path[1:]

	for _, thing := range things {
		if thing.ID == id {

			// Item found
			logger.InfoContext(ctx, "returning item")
			rest.WriteJSON(w, r, thing, http.StatusOK)

			return
		}
	}

	// Item not found
	rest.NotFoundResponse(w, r, "no item found with ID: "+id, errors.New("item not found"))
}

func (app *application) PostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	// Get body
	var payload thing
	err := rest.ReadJSON(r, &payload)

	if err != nil {
		rest.BadRequestResponse(w, r, "unable to read JSON", err)
		return
	}

	if r.Header.Get("Test") != "true" { // If testing, don't change things
		things = append(things, payload)
	}

	// Item was posted
	logger.InfoContext(ctx, "item posted")
	rest.WriteJSON(w, r, nil, http.StatusCreated)
}

func (app *application) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

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
	rest.NotFoundResponse(w, r, "no item found with ID: "+id, errors.New("item not found"))
}

func (app *application) PutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)
	// Get ID from URL
	id := r.URL.Path[1:]

	// Get body
	var payload thing
	err := rest.ReadJSON(r, &payload)

	if err != nil {
		rest.BadRequestResponse(w, r, "unable to read JSON", err)
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
	logger.InfoContext(ctx, "item created")
	rest.WriteJSON(w, r, nil, http.StatusCreated)
}

func (app *application) PatchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PATCH")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	// Get ID from URL
	id := r.URL.Path[1:]

	// Get body
	var payload [][]string
	err := rest.ReadJSON(r, &payload)

	if err != nil {
		rest.BadRequestResponse(w, r, "unable to read JSON", err)
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
				rest.NotFoundResponse(w, r, missedChanges, errors.New("some keywords were not found"))
				return
			}
			// Item patched
			logger.InfoContext(ctx, "item updated")

			return
		}
	}

	// Item not found
	rest.NotFoundResponse(w, r, "no item found with ID: "+id, errors.New("item not found"))
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
	rest.WriteJSON(w, r, nil, http.StatusOK)
}
