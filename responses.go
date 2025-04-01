package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	body := map[string]string{"error": msg}
	data, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshalling JSON data: %s", err)
        fallback := map[string]string{"error": "Something went wrong"}
        fallbackData, _ := json.Marshal(fallback) // Safe fallback, ignoring error here
        w.Write(fallbackData)
        return
	}
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, body any) {	//Sets the response 
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    data, err := json.Marshal(body)
    if err != nil {
        log.Printf("Error marshalling JSON data: %s", err)
        fallback := map[string]string{"error": "Something went wrong"}
        fallbackData, _ := json.Marshal(fallback) // Safe fallback, ignoring error here
        w.Write(fallbackData)
        return
    }
    w.Write(data)
}