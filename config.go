package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
	"encoding/json"
	"github.com/jms-guy/chirpy/internal/database"
	"github.com/jms-guy/chirpy/internal/auth"
	"github.com/google/uuid"
)

type apiConfig struct {
	db *database.Queries
	platform string
	tokenSecret string
	fileserverHits atomic.Int32
}

type User struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

///// Config handle methods

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {	//Returns a user struct from the database based on a given email and password string
	type httpRequest struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds *int `json:"expires_in_seconds,omitempty"`
	}

	request := httpRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)	//Gets request data
	if err != nil {
		fmt.Printf("Error decoding request body: %s", err)
		respondWithError(w, 400, "Invalid JSON payload")
		return
	}

	newUser, err := cfg.db.GetUserFromEmail(req.Context(), request.Email)	//Gets user struct
	if err != nil {
		fmt.Printf("Error getting user from db: %s", err)
		respondWithError(w, 401, "Incorrect email")
		return
	}

	if err := auth.CheckPasswordHash(newUser.HashedPassword, request.Password); err != nil {	//Authenticates user from password string, against hash in user struct
		respondWithError(w, 401, "Incorrect password")
		return
	}

	var expiresIn time.Duration		//Sets the expiration time for the authorization token at 1 hour if expiration date not set by client
	if request.ExpiresInSeconds != nil {
		if *request.ExpiresInSeconds > 3600 {
			expiresIn = 3600 * time.Second
		} else {
			expiresIn = time.Duration(*request.ExpiresInSeconds) * time.Second
		}
	} else {
		expiresIn = 3600 * time.Second
	}

	token, err := auth.MakeJWT(newUser.ID, cfg.tokenSecret, expiresIn)	//Creates token for user
	if err != nil {
		respondWithError(w, 500, "Error creating token")
		return
	}

	user := User{	//Structures user data in json payload
		ID: newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email: newUser.Email,
		Token: token,
	}
	respondWithJSON(w, 200, user)
}

func (cfg *apiConfig) getSingleChirp(w http.ResponseWriter, req *http.Request) {	//Gets a single chirp from database
	chirpId := req.PathValue("chirpId")	//Gets chirpid string from path

	id, err := uuid.Parse(chirpId)	//Decodes id string into uuid to find in database
	if err != nil {
		fmt.Printf("Failed to parse chirp ID: %s\n", err)
		respondWithError(w, 404, "Chirp not found")
		return
	}
	
	chirp, err := cfg.db.GetSingleChirp(req.Context(), id)	//Fetches chirp from database
	if err != nil {
		fmt.Printf("Error getting chirp from database: %s\n", err)
		respondWithError(w, 404, "Chirp not found")
		return
	}

	new := Chirp{	//Structures chirp in json payload
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(w, 200, new)
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, req *http.Request) {	//Returns all chirps from database
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
	}

	token, tokenErr := auth.GetBearerToken(req.Header)	//Gets authorization token from request
	if tokenErr != nil {
		respondWithError(w, 401, "Missing authorization header")
		return
	}
	id, valErr := auth.ValidateJWT(token, cfg.tokenSecret)	//Validates authorization
	if valErr != nil {
		respondWithError(w, 401, "Unauthorized")
		return
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
		UserID: id,
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
		Password string `json:"password"`
		Email string `json:"email"`
	}

	request := httpRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)	//Gets request data
	if err != nil {
		fmt.Printf("Error decoding request body: %s", err)
		respondWithError(w, 400, "Invalid JSON payload")
		return
	}

	hash, err := auth.HashPassword(request.Password)	//Hashes password
	if err != nil {
		fmt.Printf("Error hashing password: %s", err)
		respondWithError(w, 400, "Invalid password string")
		return
	}

	userParams := database.CreateUserParams{	//Create new user parameters
		ID: uuid.New(),
		Email: request.Email,
		HashedPassword: hash,
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