package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateHandler(w http.ResponseWriter, req *http.Request) {	//Validates that a post is not longer than the 140 char limit
	type chirpyPost struct {
		Body string `json:"body"`
	}
	type validated struct {
		CleanedBody string `json:"cleaned_body"`
		Valid bool `json:"valid"`
		Error string `json:"error,omitempty"`
	}

	respBody := validated{}
	decoder := json.NewDecoder(req.Body)
	request := chirpyPost{}
	err := decoder.Decode(&request)
	respBody.CleanedBody = badWordReplacement(request.Body)	//Set response "body", replace any profanity
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

func badWordReplacement(text string) string {	//Replaces profanity in post with "****"
	wordsInPost := strings.Fields(text)
	filtered := make([]string, len(wordsInPost))
	for i, word := range wordsInPost {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			word = "****"
		}
		filtered[i] = word
	}
	return strings.Join(filtered, " ")
}