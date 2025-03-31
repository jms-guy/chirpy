package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func validateHandler(w http.ResponseWriter, req *http.Request) {
	type chirpyPost struct {
		Body string `json:"body"`
	}
	type validated struct {
		Valid bool `json:"valid"`
		Error string `json:"error,omitempty"`
	}

	respBody := validated{}

	decoder := json.NewDecoder(req.Body)
	request := chirpyPost{}
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respBody.Error = "Invalid JSON payload"
		respondWithJSON(w, 400, respBody) // 400 for bad client input
		return
	}

	if len(request.Body) > 140 {
		respBody.Error = "Chirp is too long"
		respondWithJSON(w, 400, respBody)
	} else {
		respBody.Valid = true
		respondWithJSON(w, 200, respBody)
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, body any) {
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