package httpie

import (
	"encoding/json"
	"net/http"
)

// Used for operations that resulted in a failure, returns a JSON error with the specified status code
func WriteErrJson(w http.ResponseWriter, status int, message string) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// Used for operations that resulted in a failure, returns a JSON error
// Determines the status code from the error if possible, defaults to 500
func WriteErr(w http.ResponseWriter, err error) error {
	if httpErr, ok := err.(HttpError); ok {
		return WriteErrJson(w, httpErr.StatusCode(), httpErr.Error())
	} else {
		return WriteErrJson(w, http.StatusInternalServerError, "internal server error")
	}
}

// Used for successful operations that return a body
func WriteOk[T any](w http.ResponseWriter, data T) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}

// Used for successful operations that dont return a body
func WriteAccepted(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Helper to write either an error or a successful response
func WriteOkOrErr[T any](w http.ResponseWriter, data T, err error) {
	if err != nil {
		WriteErr(w, err)
		return
	}
	WriteOk(w, data)
}
