package auth

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetBearerToken(headers http.Header) (string, error) {	//Gets bearer authorization token from request header
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("no token found")
	}
	tokenString = tokenString[len("Bearer "):]
	return tokenString, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {	//Generates a JWT authorization token
	claims := jwt.RegisteredClaims{	//Creates claims payload for token
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: hex.EncodeToString(userID[:]),	//Encode uuid into string form
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)	//Create a token with claims payload 

	s, err := token.SignedString([]byte(tokenSecret))	//Create a verification signature of the token based off the tokenSecret string
	if err != nil {
		return "", err
	}

	return s, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {	//Takes a token string, and validates it based off a stored secret string
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {	//Parses the token string into a claims token
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("token is invalid or expired")
	}
	idString, err := parsedToken.Claims.GetSubject()	//Gets userID string from token payload
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error getting id from token: %w", err)
	}
	
	idBytes, err := hex.DecodeString(idString)	//Converts token string into bytes
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error decoding id string: %w", err)
	}
	
	id, err := uuid.FromBytes(idBytes)	//Converts id token bytes into usable uuid variable
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error getting id from byte value: %w", err)
	}

	return id, nil
}