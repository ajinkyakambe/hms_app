package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponseJSON ...
func ResponseJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// ResponseError ...
func ResponseError(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		ResponseJSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ResponseJSON(w, http.StatusBadRequest, nil)
}

// func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
// 	response, err := json.Marshal(payload)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	w.Write([]byte(response))
// }

// // respondError makes the error response with payload as json format
// func respondError(w http.ResponseWriter, code int, message string) {
// 	respondJSON(w, code, map[string]string{"error": message})
// }
