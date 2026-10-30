package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//declscope:shared // every endpoint writes through it
func writeJSONResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}

//declscope:shared // every endpoint writes through it
func writeErrorResponse(w http.ResponseWriter, status int, msg string) {
	writeJSONResponse(w, status, map[string]string{"error": msg})
}
