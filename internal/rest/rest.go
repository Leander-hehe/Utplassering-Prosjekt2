package rest

import (
	"api-test/internal/logging"
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, data any) error {

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, r *http.Request, data any, status int) error {
	logger := logging.LoggerFromContext(r.Context())

	logger.Info("marshalling data")
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	logger.Info("writing response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	const serverErrorMsg string = "server error"
	WriteError(w, r, serverErrorMsg, err, http.StatusInternalServerError)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, data any, err error) {
	WriteError(w, r, data, err, http.StatusBadRequest)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, data any, err error) {
	WriteError(w, r, data, err, http.StatusNotFound)
}

func WriteError(w http.ResponseWriter, r *http.Request, data any, err error, status int) {
	LogError(r, data, err)

	writeErr := WriteJSON(w, r, data, status)
	if writeErr != nil {
		LogError(r, "error writing respons", writeErr)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func LogError(r *http.Request, data any, err error) {
	logger := logging.LoggerFromContext(r.Context())

	logger.Error(
		"an error occurred",
		"request_method", r.Method,
		"request_url", r.URL.String(),
		"error", err,
		"message", data,
	)
}
