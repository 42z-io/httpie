package httpie

import (
	"encoding/json"
	"io"
	"net/http"
)

// Read a JSON document from the request body and unmarshal it into the provided data object
func ReadJson[T any](r *http.Request, data *T) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrBadRequest
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return ErrBadRequest
	}
	return nil
}
