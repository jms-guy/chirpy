package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
	"encoding/json"
	"github.com/jms-guy/chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	db *database.Queries
	fileserverHits atomic.Int32
}

type User struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

///// Config handle methods

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, req *http.Request) {	//Creates a new user in database
	type httpRequest struct {
		Email string `json:"email"`
	}

	request := httpRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)	//Gets request data
	if err != nil {
		fmt.Printf("Error decoding request body: %s", err)
		respondWithError(w, 400, "Invalid JSON payload")
		return
	}

	userParams := database.CreateUserParams{	//Create new user parameters
		ID: uuid.New(),
		Email: request.Email,
	}

	newUser, err := cfg.db.CreateUser(req.Context(), userParams)	//Create new user in db
	if err != nil {
		fmt.Printf("Error creating new user: %s", err)
		respondWithError(w, 400, "Could not create user")
		return
	}

	user := User{	//Map new user to config user struct
		ID: newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email: newUser.Email,
	}
	respondWithJSON(w, 201, user)	//Send response data
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {	//Wraps server handles in a function that increases server visit count
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, req *http.Request) {	
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<html>
  						<body>
							<h1>Welcome, Chirpy Admin</h1>
							<p>Chirpy has been visited %d times!</p>
						</body>
					</html>`, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {	
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}