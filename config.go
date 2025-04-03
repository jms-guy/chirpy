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
	platform string
	fileserverHits atomic.Int32
}

type User struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

///// Config handle methods

func (cfg *apiConfig) getSingleChirp(w http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("chirpId")

	id, err := uuid.Parse(chirpId)
	if err != nil {
		fmt.Printf("Failed to parse chirp ID: %s\n", err)
		respondWithError(w, 404, "Chirp not found")
		return
	}
	
	chirp, err := cfg.db.GetSingleChirp(req.Context(), id)
	if err != nil {
		fmt.Printf("Error getting chirp from database: %s\n", err)
		respondWithError(w, 404, "Chirp not found")
		return
	}

	new := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(w, 200, new)
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		fmt.Printf("Error retrieving chirps from database: %s", err)
		respondWithError(w, 400, "Error retrieving data from server")
		return
	}

	returnChirps := []Chirp{}

	for _, chirp := range chirps {
		new := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		}
		returnChirps = append(returnChirps, new)
	}

	respondWithJSON(w, 200, returnChirps)
}

func (cfg *apiConfig) chirpsHandler(w http.ResponseWriter, req *http.Request) {	//Creates a chirp in db
	type chirpRequest struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	request := chirpRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)	//Decodes request data
	if err != nil {
		fmt.Printf("Error decoding response body: %s", err)
		respondWithError(w, 400, "Invalid JSON payload")
		return
	}

	chirpParams := database.CreateChirpParams{	//Create chirp parameters
		ID: uuid.New(),
		Body: request.Body,
		UserID: request.UserID,
	}

	newChirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)	//Create new chirp
	if err != nil {
		fmt.Printf("Error creating new chirp: %s", err)
		respondWithError(w, 400, "Could not create chirp")
		return
	}

	chirp := Chirp{	//Set chirp struct to respond with
		ID: newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body: newChirp.Body,
		UserID: newChirp.UserID,
	}
	respondWithJSON(w, 201, chirp)
}

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
	if cfg.platform != "dev" {
		respondWithError(w, 403, "User does not have access to this page")
		return
	}
	err := cfg.db.ClearUsers(req.Context())
	if err != nil {
		fmt.Printf("Error clearing users table: %s", err)
		respondWithError(w, 500, "Error clearing database")
		return
	}
	respondWithJSON(w, 200, map[string]string{"status": "Users table cleared successfully"})
}