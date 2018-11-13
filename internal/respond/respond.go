package respond

import (
	"encoding/json"
	"net/http"
)

// JSON writes the status code, the error if it's not null, and otherwise the
// data.
func JSON(w http.ResponseWriter, data interface{}, err error, code int) {
	var payload interface{}
	if err != nil {
		payload = map[string]string{"error": err.Error()}
	} else {
		payload = data
	}

	w.WriteHeader(code)

	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
