package main

import (
	"api-test/internal/data"
	"api-test/internal/logging"
	"api-test/internal/rest"
	"errors"
	"fmt"
	"net/http"
)

func sendError(w http.ResponseWriter, r *http.Request, data any, err error, status int) {
	switch status {
	case http.StatusNotFound:
		rest.NotFoundResponse(w, r, data, err)
	case http.StatusBadRequest:
		rest.BadRequestResponse(w, r, data, err)
	case http.StatusInternalServerError:
		rest.ServerErrorResponse(w, r, err)
	}
}

func getPayload(w http.ResponseWriter, r *http.Request) (Thing, error) {
	// Get body
	var payload ThingNoID
	err := rest.ReadJSON(r, &payload)

	if err != nil {
		rest.BadRequestResponse(w, r, "unable to read JSON", err)
		return Thing{}, err
	}

	newPayload := Thing{ID: "0", Name: payload.Name, Description: payload.Description}

	return newPayload, nil
}

func (app *application) GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.InfoContext(ctx, "returning items")

	things, err, status := data.SelectAll(app.db)

	if err != nil {
		sendError(w, r, "unable to get data", err, status)
		return
	}

	rest.WriteJSON(w, r, things, http.StatusOK)
}

func (app *application) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	// Get ID from URL
	id := r.URL.Path[1:]

	thing, err, status := data.Select(app.db, id)

	if err != nil {
		sendError(w, r, "unable to get data", err, status)
		return
	}

	// Item found
	logger.InfoContext(ctx, "returning item")
	rest.WriteJSON(w, r, thing, http.StatusOK)

	// Item not found
	rest.NotFoundResponse(w, r, "no item found with ID: "+id, errors.New("item not found"))
}

func (app *application) PostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	// Get body
	payload, err := getPayload(w, r)

	if err != nil {
		return
	}

	var outputID int

	if r.Header.Get("Test") != "true" { // If testing, don't change things
		var status int
		err, status, outputID = data.Insert(app.db, data.Thing(payload))
		if err != nil {
			sendError(w, r, "unable to post data", err, status)
			return
		}
	}

	// Item was posted
	logger.InfoContext(ctx, "item posted")
	rest.WriteJSON(w, r, outputID, http.StatusCreated)
}

func (app *application) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	// Get ID from URL
	id := r.URL.Path[1:]

	// Item not found, create new item
	if r.Header.Get("Test") != "true" { // If testing, don't change things
		err, status := data.Delete(app.db, id)
		if err != nil {
			sendError(w, r, "unable to delete data", err, status)
			return
		}
	}

	logger.InfoContext(ctx, "item deleted")
}

func (app *application) PutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT")

	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)
	// Get ID from URL
	id := r.URL.Path[1:]

	// Get body
	payload, err := getPayload(w, r)

	if err != nil {
		return
	}

	if r.Header.Get("Test") != "true" { // If testing, don't change things
		err, status := data.Update(app.db, id, data.Thing(payload))
		if err != nil {
			sendError(w, r, "unable to put data", err, status)
		}

		if status == http.StatusNotFound { // Item not found, create new item
			data.Insert(app.db, data.Thing(payload))
		}

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

	thing, err, status := data.Select(app.db, id)
	if err != nil {
		sendError(w, r, "unable to find item", err, status)
		return
	}

	// Found item
	var missedChanges [][]string

	for _, change := range payload { // Loop through changes

		if r.Header.Get("Test") != "true" { // If testing, don't change things
			switch change[0] {
			case "Description":
				thing.Description = change[1]
			case "Name":
				thing.Name = change[1]
			default:
				// Change not found
				missedChanges = append(missedChanges, change)
			}
		}
	}

	err, status = data.Update(app.db, id, *thing)
	if err != nil {
		sendError(w, r, "unable to update item", err, status)
		return
	}

	if len(missedChanges) > 0 {
		rest.NotFoundResponse(w, r, missedChanges, errors.New("some keywords were not found"))
		return
	}
	// Item patched
	logger.InfoContext(ctx, "item updated")

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
