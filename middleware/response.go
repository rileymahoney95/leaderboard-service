package middleware

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response for the API
type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

// RespondWithError sends an error response to the client
func RespondWithError(w http.ResponseWriter, code int, message string, err error) {
	var errMsg interface{}
	if err != nil {
		errMsg = err.Error()
	}

	response := ErrorResponse{
		Status:  code,
		Message: message,
		Error:   errMsg,
	}

	RespondWithJSON(w, code, response)
}

// RespondWithJSON sends a JSON response to the client
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
