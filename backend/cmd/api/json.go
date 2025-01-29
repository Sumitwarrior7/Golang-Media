package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	/* Denial-of-Service (DoS/DDoS) Attacks */
	// If no limit is set, an attacker could send excessively large requests to consume server resources, such as memory and CPU.
	// This could lead to denial of service for legitimate users as the server struggles to handle these requests.
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data) // Here the payload variable which is 'data'in this case, is populated with the json-data from decoder through r.Body
}

// The writeJSONError function is a utility function designed to streamline sending error responses in JSON format to the client
func writeJsonError(w http.ResponseWriter, status int, message string) error {
	// This ensures that the error message is always sent as a key-value pair with "error" as the key : { "error": "your error message here" }
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJSON(w, status, &envelope{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJSON(w, status, &envelope{Data: data})
}
