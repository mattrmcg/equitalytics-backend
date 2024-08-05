package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// writes JSON byte slice with given 'v' (should be map or struct)
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// writes given error into json for return
func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}
