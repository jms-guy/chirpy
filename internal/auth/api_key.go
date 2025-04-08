package auth

import (
	"fmt"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("no api key found")
	}
	apiKey = apiKey[len("ApiKey "):]
	return apiKey, nil
}